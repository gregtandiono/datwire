package bolt_test

import (
	"datwire/pkg/bolt"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type AuthServiceTestSuite struct {
	suite.Suite
	authService *bolt.AuthService
	userID      uuid.UUID
}

func (suite *AuthServiceTestSuite) SetupSuite()    {}
func (suite *AuthServiceTestSuite) TearDownSuite() {}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
