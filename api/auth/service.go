package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/schema"
	authschema "github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user"
	userSchema "github.com/unusualcodeorg/go-lang-backend-architecture/api/user/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUpBasic(signupDto *dto.SignUpBasic) (*dto.UserAuth, error)
	generateToken(user userSchema.User) (*dto.UserTokens, error)
	createKeystore(client *userSchema.User, primaryKey string, secondaryKey string) (*schema.Keystore, error)
	verifyToken(tokenStr string) (*jwt.RegisteredClaims, error)
	decodeToken(tokenStr string) (*jwt.RegisteredClaims, error)
	signToken(claims jwt.RegisteredClaims) (string, error)
}

type service struct {
	network.BaseService
	keystoreQuery mongo.Query[authschema.Keystore]
	userService   user.UserService
	// token
	rsaPrivateKey           string
	rsaPublicKey            string
	accessTokenValiditySec  string
	refreshTokenValiditySec string
	tokenIssuer             string
	TokenAudience           string
}

func NewAuthService(db mongo.Database, dbQueryTimeout time.Duration, env *config.Env, userService user.UserService) AuthService {
	s := service{
		BaseService:   network.NewBaseService(dbQueryTimeout),
		keystoreQuery: mongo.NewQuery[authschema.Keystore](db, authschema.KeystoreCollectionName),
		userService:   userService,
		// token
		rsaPrivateKey:           env.RSAPrivateKey,
		rsaPublicKey:            env.RSAPublicKey,
		accessTokenValiditySec:  env.AccessTokenValiditySec,
		refreshTokenValiditySec: env.RefreshTokenValiditySec,
		tokenIssuer:             env.TokenIssuer,
		TokenAudience:           env.TokenAudience,
	}
	return &s
}

func (s *service) SignUpBasic(signupDto *dto.SignUpBasic) (*dto.UserAuth, error) {
	user, _ := s.userService.FindUserByEmail(signupDto.Email)
	if user != nil {
		e := errors.New("user already exists")
		return nil, network.BadRequestError(e.Error(), e)
	}

	role, err := s.userService.FindRoleByCode(userSchema.RoleCodeLearner)
	if err != nil {
		return nil, err
	}
	roles := make([]userSchema.Role, 1)
	roles[0] = *role

	hashed, err := bcrypt.GenerateFromPassword([]byte(signupDto.Password), 5)
	if err != nil {
		return nil, err
	}

	_, err = userSchema.NewUser(signupDto.Email, string(hashed), &signupDto.Name, &signupDto.Password, roles)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *service) createKeystore(client *userSchema.User, primaryKey string, secondaryKey string) (*schema.Keystore, error) {
	ctx, cancel := s.Context()
	defer cancel()

	keystore, err := schema.NewKeystore(client.ID, primaryKey, secondaryKey)
	if err != nil {
		return nil, err
	}

	doc := keystore.GetDocument()

	id, err := s.keystoreQuery.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	doc.ID = *id
	return doc, nil
}

func (s *service) generateToken(user userSchema.User) (*dto.UserTokens, error) {
	// accessTokenKey := utils.GenerateRandomString(32)
	// refreshTokenKey := utils.GenerateRandomString(32)

	return nil, nil
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
