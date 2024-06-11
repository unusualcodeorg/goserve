package core

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/utils"
)

type TokenService interface {
	GenerateToken(user *schema.User) (string, string, error)
	CreateKeystore(client *schema.User, primaryKey string, secondaryKey string) (*schema.Keystore, error)
	VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error)
	DecodeToken(tokenStr string) (*jwt.RegisteredClaims, error)
	SignToken(claims jwt.RegisteredClaims) (string, error)
}

type tokenService struct {
	network.BaseService
	keystoreQuery mongo.Query[schema.Keystore]
	// token
	rsaPrivateKey        *rsa.PrivateKey
	rsaPublicKey         *rsa.PublicKey
	accessTokenValidity  time.Duration
	refreshTokenValidity time.Duration
	tokenIssuer          string
	tokenAudience        string
}

func NewTokenService(db mongo.Database, dbQueryTimeout time.Duration, env *config.Env) TokenService {
	privatePem, err := utils.LoadPEMFileInto(env.RSAPrivateKeyPath)
	if err != nil {
		panic(err)
	}
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePem)
	if err != nil {
		panic(err)
	}

	publicPem, err := utils.LoadPEMFileInto(env.RSAPublicKeyPath)
	if err != nil {
		panic(err)
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPem)
	if err != nil {
		panic(err)
	}

	s := tokenService{
		BaseService:   network.NewBaseService(dbQueryTimeout),
		keystoreQuery: mongo.NewQuery[schema.Keystore](db, schema.KeystoreCollectionName),
		// token key
		rsaPrivateKey: rsaPrivateKey,
		rsaPublicKey:  rsaPublicKey,
		// token claim
		accessTokenValidity:  time.Duration(env.AccessTokenValiditySec),
		refreshTokenValidity: time.Duration(env.RefreshTokenValiditySec),
		tokenIssuer:          env.TokenIssuer,
		tokenAudience:        env.TokenAudience,
	}
	return &s
}

func (s *tokenService) GenerateToken(user *schema.User) (string, string, error) {
	primaryKey, err := utils.GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}
	secondaryKey, err := utils.GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}

	_, err = s.CreateKeystore(user, primaryKey, secondaryKey)
	if err != nil {
		return "", "", err
	}

	now := jwt.NewNumericDate(time.Now())

	accessTokenClaims := jwt.RegisteredClaims{
		Issuer:    s.tokenIssuer,
		Subject:   user.ID.Hex(),
		Audience:  []string{s.tokenAudience},
		IssuedAt:  now,
		NotBefore: now,
		ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenValidity * time.Second)),
		ID:        primaryKey,
	}

	refreshTokenClaims := jwt.RegisteredClaims{
		Issuer:    s.tokenIssuer,
		Subject:   user.ID.Hex(),
		Audience:  []string{s.tokenAudience},
		IssuedAt:  now,
		NotBefore: now,
		ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenValidity * time.Second)),
		ID:        secondaryKey,
	}

	accessToken, err := s.SignToken(accessTokenClaims)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.SignToken(refreshTokenClaims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *tokenService) CreateKeystore(client *schema.User, primaryKey string, secondaryKey string) (*schema.Keystore, error) {
	ctx, cancel := s.Context()
	defer cancel()

	doc, err := schema.NewKeystore(client.ID, primaryKey, secondaryKey)
	if err != nil {
		return nil, err
	}

	id, err := s.keystoreQuery.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	doc.ID = *id
	return doc, nil
}

func (s *tokenService) SignToken(claims jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(s.rsaPrivateKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (s *tokenService) VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.Parse(tokenStr, func(tkn *jwt.Token) (any, error) {
		return s.rsaPublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		if claims, ok := token.Claims.(jwt.RegisteredClaims); ok {
			return &claims, nil
		}
	}

	return nil, jwt.ErrTokenMalformed
}

func (s *tokenService) DecodeToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.Parse(tokenStr, func(tkn *jwt.Token) (any, error) {
		return s.rsaPublicKey, nil
	})

	if token.Valid {
		if claims, ok := token.Claims.(jwt.RegisteredClaims); ok {
			return &claims, nil
		}
	}

	if err != nil {
		return nil, err
	}

	return nil, jwt.ErrTokenMalformed
}
