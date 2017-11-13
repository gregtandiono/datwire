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
}

func (suite *UserServiceTestSuite) SetupSuite() {
	suite.userService = &bolt.UserService{}
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
		ID:       uuid.FromStringOrNil("099ef5d7-04d2-43b0-a765-907216f388da"),
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

	u, err := suite.userService.User(uuid.FromStringOrNil("099ef5d7-04d2-43b0-a765-907216f388da"))
	suite.Nil(err)
	suite.Equal("Gregory Tandiono", u.Name, "name should match")
	suite.Equal("gtandiono", u.Username, "username should match")
	suite.Equal("admin", u.Type, "type should match")
}

func (suite *UserServiceTestSuite) TestUserService_SetName() {
	suite.userService.Open()
	defer suite.userService.Close()
	err := suite.userService.SetName(uuid.FromStringOrNil("099ef5d7-04d2-43b0-a765-907216f388da"), "Benjamin")
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestUserService_SetName_VerifySet() {
	suite.userService.Open()
	defer suite.userService.Close()
	u, err := suite.userService.User(uuid.FromStringOrNil("099ef5d7-04d2-43b0-a765-907216f388da"))
	suite.Nil(err)
	suite.Equal("Benjamin", u.Name, "name should be updated")
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
