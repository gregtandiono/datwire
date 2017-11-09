package bolt_test

import (
	"datwire/pkg/bolt"

	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	userService *bolt.UserService
}

func (suite *UserServiceTestSuite) SetupSuite() {
	suite.userService = bolt.NewUserService()
}

func (suite *UserServiceTestSuite) TearDownSuite() {
}
