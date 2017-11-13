package bolt_test

import (
	"datwire/pkg/bolt"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type AuthServiceTestSuite struct {
	suite.Suite
	authService    *bolt.AuthService
	userID         uuid.UUID
	hashedPassword string
}

func (suite *AuthServiceTestSuite) SetupSuite() {
	suite.authService = &bolt.AuthService{}
	suite.userID = uuid.FromStringOrNil("573c9c34-ff55-4af2-9ac1-985412a8af69")
	suite.hashedPassword = "$2a$10$8F56R2MnW4mTwLclwQ4gw.D1v9d8WU6W.bs96ZSjRc69NW2Os2E0e"
}

func (suite *AuthServiceTestSuite) TestAuthorization() {
	authObj, err := suite.authService.Authorize(
		"themostawesomepasswordintheworld",
		suite.hashedPassword,
		suite.userID,
	)
	suite.Nil(err)
	suite.NotNil(authObj.Token)
	suite.Equal(suite.userID.String(), authObj.ID, "user id should match")
}

func (suite *AuthServiceTestSuite) TestAuthorization_PasswordDoesNotMatch() {
	_, err := suite.authService.Authorize(
		"somesillypasswordthatwilldefinitelyfail",
		suite.hashedPassword,
		suite.userID,
	)
	suite.NotNil(err)
	suite.Equal("password does not match", err.Error(), "error message should match")
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
