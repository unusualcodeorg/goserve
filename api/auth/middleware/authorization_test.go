package middleware

import (
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/unusualcodeorg/goserve/api/auth"
	"github.com/unusualcodeorg/goserve/api/auth/model"
	"github.com/unusualcodeorg/goserve/api/user"
	userModel "github.com/unusualcodeorg/goserve/api/user/model"
	"github.com/unusualcodeorg/goserve/arch/network"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAuthorizationProvider_NoRole(t *testing.T) {
	mockAuthService := new(auth.MockService)
	mockUserService := new(user.MockService)

	token := "Bearer token"
	userId := primitive.NewObjectID()
	keystoreId := primitive.NewObjectID()
	claims := &jwt.RegisteredClaims{ID: "claimId", Subject: userId.Hex()}
	user := &userModel.User{ID: userId}
	keystore := &model.Keystore{ID: keystoreId}

	mockAuthService.On("VerifyToken", "token").Return(claims, nil)
	mockAuthService.On("ValidateClaims", claims).Return(true)
	mockUserService.On("FindUserById", userId).Return(user, nil)
	mockAuthService.On("FindKeystore", user, claims.ID).Return(keystore, nil)

	rr := network.MockTestAuthorizationProvider(t, "",
		NewAuthenticationProvider(mockAuthService, mockUserService),
		NewAuthorizationProvider(),
		network.MockSuccessMsgHandler("success"),
		primitive.E{Key: network.AuthorizationHeader, Value: token},
	)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"permission denied: role missing"`)
}

func TestAuthorizationProvider_WrongRole(t *testing.T) {
	mockAuthService := new(auth.MockService)
	mockUserService := new(user.MockService)

	token := "Bearer token"
	userId := primitive.NewObjectID()
	roleId := primitive.NewObjectID()
	keystoreId := primitive.NewObjectID()
	claims := &jwt.RegisteredClaims{ID: "claimId", Subject: userId.Hex()}
	role := &userModel.Role{ID: roleId, Code: "TEST"}
	user := &userModel.User{ID: userId, RoleDocs: []*userModel.Role{role}}
	keystore := &model.Keystore{ID: keystoreId}

	mockAuthService.On("VerifyToken", "token").Return(claims, nil)
	mockAuthService.On("ValidateClaims", claims).Return(true)
	mockUserService.On("FindUserById", userId).Return(user, nil)
	mockAuthService.On("FindKeystore", user, claims.ID).Return(keystore, nil)

	rr := network.MockTestAuthorizationProvider(t, "WRONG",
		NewAuthenticationProvider(mockAuthService, mockUserService),
		NewAuthorizationProvider(),
		network.MockSuccessMsgHandler("success"),
		primitive.E{Key: network.AuthorizationHeader, Value: token},
	)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"permission denied: does not have suffient role"`)
}

func TestAuthorizationProvider_Success(t *testing.T) {
	mockAuthService := new(auth.MockService)
	mockUserService := new(user.MockService)

	token := "Bearer token"
	userId := primitive.NewObjectID()
	roleId := primitive.NewObjectID()
	keystoreId := primitive.NewObjectID()
	claims := &jwt.RegisteredClaims{ID: "claimId", Subject: userId.Hex()}
	role := &userModel.Role{ID: roleId, Code: "TEST"}
	user := &userModel.User{ID: userId, RoleDocs: []*userModel.Role{role}}
	keystore := &model.Keystore{ID: keystoreId}

	mockAuthService.On("VerifyToken", "token").Return(claims, nil)
	mockAuthService.On("ValidateClaims", claims).Return(true)
	mockUserService.On("FindUserById", userId).Return(user, nil)
	mockAuthService.On("FindKeystore", user, claims.ID).Return(keystore, nil)

	rr := network.MockTestAuthorizationProvider(t, "TEST",
		NewAuthenticationProvider(mockAuthService, mockUserService),
		NewAuthorizationProvider(),
		network.MockSuccessMsgHandler("success"),
		primitive.E{Key: network.AuthorizationHeader, Value: token},
	)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"success"`)
}
