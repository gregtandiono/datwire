package bolt_test

import (
	"datwire/pkg/apps/user"
	"datwire/pkg/bolt"
	"log"
	"testing"

	sysbolt "github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"

	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	userService *bolt.UserService
	userID_1    uuid.UUID
	userID_2    uuid.UUID
}

func (suite *UserServiceTestSuite) SetupSuite() {
	suite.userService = &bolt.UserService{}
	suite.userID_1 = uuid.FromStringOrNil("099ef5d7-04d2-43b0-a765-907216f388da")
	suite.userID_2 = uuid.FromStringOrNil("028b5c04-f91e-4312-990d-33525456d1a3")

	// seed one data
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

func (suite *UserServiceTestSuite) TearDownSuite() {
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

func (suite *UserServiceTestSuite) TestUserService_CreateUser() {
	suite.userService.Open()
	defer suite.userService.Close()

	err := suite.userService.CreateUser(&user.User{
		ID:       suite.userID_1,
		Name:     "Gregory Tandiono",
		Username: "gtandiono",
		Password: "thisisasuperawesomepasswordyo",
		Type:     "admin",
	})

	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestUserService_CreateUser_VerifyCreation() {
	suite.userService.Open()
	defer suite.userService.Close()

	u, err := suite.userService.User(suite.userID_1)
	suite.Nil(err)
	suite.Equal("Gregory Tandiono", u.Name, "name should match")
	suite.Equal("gtandiono", u.Username, "username should match")
	suite.Equal("admin", u.Type, "type should match")
}

func (suite *UserServiceTestSuite) TestUserService_FetchAllUsers() {
	suite.userService.Open()
	defer suite.userService.Close()

	users, err := suite.userService.Users()
	suite.Nil(err)
	suite.Equal(2, len(users), "user amount should match")
}

func (suite *UserServiceTestSuite) TestUserService_SetName() {
	suite.userService.Open()
	defer suite.userService.Close()
	err := suite.userService.SetName(suite.userID_1, "Benjamin")
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestUserService_SetName_VerifySet() {
	suite.userService.Open()
	defer suite.userService.Close()
	u, err := suite.userService.User(suite.userID_1)
	suite.Nil(err)
	suite.Equal("Benjamin", u.Name, "name should be updated")
}

func (suite *UserServiceTestSuite) TestUserService_RemoveUser() {
	suite.userService.Open()
	defer suite.userService.Close()

	err := suite.userService.DeleteUser(suite.userID_2)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestUserService_RemoveUser_VerifyRemoval() {
	suite.userService.Open()
	defer suite.userService.Close()

	u, err := suite.userService.User(suite.userID_2)
	suite.NotNil(err)
	suite.Equal("user does not exist", err.Error(), "error message should match")
	suite.Nil(u)
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
