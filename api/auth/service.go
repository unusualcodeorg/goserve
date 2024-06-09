package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
	FindRoles(roleIds []primitive.ObjectID) ([]schema.Role, error)
	verifyToken(tokenStr string) (*jwt.RegisteredClaims, error)
	decodeToken(tokenStr string) (*jwt.RegisteredClaims, error)
	signToken(claims jwt.RegisteredClaims) (string, error)
}

type service struct {
	network.BaseService
	rsaPrivateKey string
	rsaPublicKey  string
	roleQuery     mongo.Query[schema.Role]
	keystoreQuery mongo.Query[schema.Keystore]
}

func NewAuthService(db mongo.Database, dbQueryTimeout time.Duration, env *config.Env) AuthService {
	s := service{
		BaseService:   network.NewBaseService(dbQueryTimeout),
		rsaPrivateKey: env.RSAPrivateKey,
		rsaPublicKey:  env.RSAPublicKey,
		roleQuery:     mongo.NewQuery[schema.Role](db, schema.RolesCollectionName),
		keystoreQuery: mongo.NewQuery[schema.Keystore](db, schema.KeystoreCollectionName),
	}
	return &s
}

func (s *service) FindRoles(roleIds []primitive.ObjectID) ([]schema.Role, error) {
	ctx, cancel := s.Context()
	cancel()
	filter := bson.M{"_id": bson.M{"$in": roleIds}}
	return s.roleQuery.FindAll(ctx, filter, nil)
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
