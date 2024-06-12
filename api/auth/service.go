package auth

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user"
	userModel "github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/utils"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	IsEmailRegisted(email string) bool
	SignUpBasic(signUpDto *dto.SignUpBasic) (*dto.UserAuth, error)
	SignInBasic(signInDto *dto.SignInBasic) (*dto.UserAuth, error)
	GenerateToken(user *userModel.User) (string, string, error)
	CreateKeystore(client *userModel.User, primaryKey string, secondaryKey string) (*model.Keystore, error)
	FindKeystore(client *userModel.User, primaryKey string) (*model.Keystore, error)
	VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error)
	DecodeToken(tokenStr string) (*jwt.RegisteredClaims, error)
	SignToken(claims jwt.RegisteredClaims) (string, error)
	ValidateClaims(claims *jwt.RegisteredClaims) bool
	FindApiKey(key string) (*model.ApiKey, error)
}

type service struct {
	network.BaseService
	keystoreQuery mongo.Query[model.Keystore]
	apikeyQuery   mongo.Query[model.ApiKey]
	userService   user.UserService
	// token
	rsaPrivateKey        *rsa.PrivateKey
	rsaPublicKey         *rsa.PublicKey
	accessTokenValidity  time.Duration
	refreshTokenValidity time.Duration
	tokenIssuer          string
	tokenAudience        string
}

func NewAuthService(
	db mongo.Database,
	dbQueryTimeout time.Duration,
	env *config.Env,
	userService user.UserService,
) AuthService {
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

	s := service{
		BaseService:   network.NewBaseService(dbQueryTimeout),
		userService:   userService,
		keystoreQuery: mongo.NewQuery[model.Keystore](db, model.KeystoreCollectionName),
		apikeyQuery:   mongo.NewQuery[model.ApiKey](db, model.CollectionName),
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

func (s *service) IsEmailRegisted(email string) bool {
	user, _ := s.userService.FindUserByEmail(email)
	return user != nil
}

func (s *service) SignUpBasic(signUpDto *dto.SignUpBasic) (*dto.UserAuth, error) {
	role, err := s.userService.FindRoleByCode(userModel.RoleCodeLearner)
	if err != nil {
		return nil, err
	}
	roles := make([]userModel.Role, 1)
	roles[0] = *role

	hashed, err := bcrypt.GenerateFromPassword([]byte(signUpDto.Password), 5)
	if err != nil {
		return nil, err
	}

	user, err := userModel.NewUser(signUpDto.Email, string(hashed), &signUpDto.Name, signUpDto.ProfilePicUrl, roles)
	if err != nil {
		return nil, err
	}

	user, err = s.userService.CreateUser(user)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	tokens := dto.NewUserToken(accessToken, refreshToken)
	return dto.NewUserAuth(user, tokens), nil
}

func (s *service) SignInBasic(signInDto *dto.SignInBasic) (*dto.UserAuth, error) {
	return nil, nil
}

func (s *service) GenerateToken(user *userModel.User) (string, string, error) {
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

func (s *service) CreateKeystore(client *userModel.User, primaryKey string, secondaryKey string) (*model.Keystore, error) {
	ctx, cancel := s.Context()
	defer cancel()

	doc, err := model.NewKeystore(client.ID, primaryKey, secondaryKey)
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

func (s *service) FindKeystore(client *userModel.User, primaryKey string) (*model.Keystore, error) {
	ctx, cancel := s.Context()
	defer cancel()
	filter := bson.M{"client": client.ID, "primaryKey": primaryKey, "status": true}
	return s.keystoreQuery.FindOne(ctx, filter, nil)
}

func (s *service) SignToken(claims jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(s.rsaPrivateKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (s *service) VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(tkn *jwt.Token) (any, error) {
		return s.rsaPublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if token.Valid {
		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
			return claims, nil
		}
	}

	return nil, jwt.ErrTokenMalformed
}

func (s *service) DecodeToken(tokenStr string) (*jwt.RegisteredClaims, error) {
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

func (s *service) ValidateClaims(claims *jwt.RegisteredClaims) bool {
	invalid := claims.Issuer != s.tokenIssuer ||
		claims.Subject == "" ||
		len(claims.Audience) == 0 ||
		claims.Audience[0] != s.tokenAudience ||
		claims.NotBefore == nil ||
		claims.ExpiresAt == nil ||
		claims.ID == ""

	return !invalid
}

func (s *service) FindApiKey(key string) (*model.ApiKey, error) {
	ctx, cancel := s.Context()
	defer cancel()

	filter := bson.M{"key": key, "status": true}

	apikey, err := s.apikeyQuery.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	return apikey, nil
}
