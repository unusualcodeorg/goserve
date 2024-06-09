package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type AuthService interface {
	verifyToken(tokenStr string) (*jwt.RegisteredClaims, error)
	decodeToken(tokenStr string) (*jwt.RegisteredClaims, error)
	signToken(claims jwt.RegisteredClaims) (string, error)
}

type service struct {
	network.BaseService
	rsaPrivateKey string
	rsaPublicKey  string
	authQuery     mongo.DatabaseQuery[schema.Role]
}

func NewAuthService(env config.Env, db mongo.Database, dbQueryTimeout time.Duration) AuthService {
	s := service{
		BaseService:   network.NewBaseService(dbQueryTimeout),
		rsaPrivateKey: env.RSAPrivateKey,
		rsaPublicKey:  env.RSAPublicKey,
		authQuery:     mongo.NewDatabaseQuery[schema.Role](db, schema.CollectionName),
	}
	return &s
}

func (s *service) signToken(claims jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	signed, err := token.SignedString(s.rsaPrivateKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (s *service) verifyToken(tokenStr string) (*jwt.RegisteredClaims, error) {
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

func (s *service) decodeToken(tokenStr string) (*jwt.RegisteredClaims, error) {
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
