package http_test

import (
	"datwire/pkg/apps/user"
	"datwire/pkg/bolt"
	"log"
	"testing"

	sysbolt "github.com/boltdb/bolt"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	userService *bolt.UserService
	userID_1    uuid.UUID
	userID_2    uuid.UUID
}

func (suite *UserHandlerTestSuite) SetupSuite() {
	suite.userService = &bolt.UserService{}
	suite.userID_1 = uuid.FromStringOrNil("099ef5d7-04d2-43b0-a765-907216f388da")
	suite.userID_2 = uuid.FromStringOrNil("028b5c04-f91e-4312-990d-33525456d1a3")

	suite.userService.Open()
	defer suite.userService.Close()

	suite.userService.CreateUser(&user.User{
		ID:       suite.userID_2,
		Name:     "Augustus Kwok",
		Username: "akwok",
		Password: "superdupermart",
		Type:     "admin",
	})
}

func (suite *UserHandlerTestSuite) TearDownSuite() {
	db, err := sysbolt.Open("user.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *sysbolt.Tx) error {
		err := tx.DeleteBucket([]byte("users"))
		if err != nil {
			return err
		}
		return nil
	})
}

func (suite *UserHandlerTestSuite) TestUserHandler_CreateUser() {
}

func (suite *UserHandlerTestSuite) TestUserHandler_FetchUser() {
}

func (suite *UserHandlerTestSuite) TestUserHandler_FetchUsers() {
}

func (suite *UserHandlerTestSuite) TestUserHandler_SetName() {
}

func (suite *UserHandlerTestSuite) TestUserHandler_SetName_VerifySet() {
}

func (suite *UserHandlerTestSuite) TestUserHandler_DeleteUser() {
}

func (suite *UserHandlerTestSuite) TestUserHandler_DeleteUser_VerifyDelete() {
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
