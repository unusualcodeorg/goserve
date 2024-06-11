package auth

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/profile"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/schema"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	IsEmailRegisted(email string) bool
	SignUpBasic(signupDto *dto.SignUpBasic) (*dto.UserAuth, error)
}

type service struct {
	network.BaseService
	keystoreQuery mongo.Query[schema.Keystore]
	userService   user.UserService
	tokenService  core.TokenService
}

func NewAuthService(
	db mongo.Database,
	dbQueryTimeout time.Duration,
	userService user.UserService,
	tokenService core.TokenService,
) AuthService {
	s := service{
		BaseService:   network.NewBaseService(dbQueryTimeout),
		keystoreQuery: mongo.NewQuery[schema.Keystore](db, schema.KeystoreCollectionName),
		userService:   userService,
		tokenService:  tokenService,
	}
	return &s
}

func (s *service) IsEmailRegisted(email string) bool {
	user, _ := s.userService.FindUserByEmail(email)
	return user != nil
}

func (s *service) SignUpBasic(signupDto *dto.SignUpBasic) (*dto.UserAuth, error) {
	role, err := s.userService.FindRoleByCode(schema.RoleCodeLearner)
	if err != nil {
		return nil, err
	}
	roles := make([]schema.Role, 1)
	roles[0] = *role

	hashed, err := bcrypt.GenerateFromPassword([]byte(signupDto.Password), 5)
	if err != nil {
		return nil, err
	}

	user, err := schema.NewUser(signupDto.Email, string(hashed), &signupDto.Name, signupDto.ProfilePicUrl, roles)
	if err != nil {
		return nil, err
	}

	user, err = s.userService.CreateUser(user)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.tokenService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	tokens := dto.NewUserToken(accessToken, refreshToken)
	return dto.NewUserAuth(user, tokens), nil
}
