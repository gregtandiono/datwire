package http_test

import (
	"bytes"
	"datwire/pkg/apps/user"
	"datwire/pkg/bolt"
	dwhttp "datwire/pkg/http"
	"datwire/pkg/shared"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	sysbolt "github.com/boltdb/bolt"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type getUserResponseTemplate struct {
	Message string
	Error   string
	Data    *user.User
}

type getUsersResponseTemplate struct {
	Message string
	Error   string
	Data    []user.User
}

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
	mockData := []byte(`{
		"id": "` + suite.userID_1.String() + `",
		"name": "Gregory Tandiono",
		"username": "gtandiono",
		"password": "themostawesomepasswordintheworld",
		"type": "admin"
	}`)

	request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(mockData))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	h := dwhttp.NewUserHandler()
	h.UserService.Open()
	defer h.UserService.Close()
	h.ServeHTTP(response, request)

	var responseBody *shared.ResponseTemplate
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	suite.Equal("", responseBody.Error, "error should be empty")
	suite.Equal("success", responseBody.Message, "message should match")
}

func (suite *UserHandlerTestSuite) TestUserHandler_CreateUser_VerifyCreate() {
	request, _ := http.NewRequest("GET", "/users/"+suite.userID_1.String(), nil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	h := dwhttp.NewUserHandler()
	h.UserService.Open()
	defer h.UserService.Close()
	h.ServeHTTP(response, request)

	var responseBody *getUserResponseTemplate
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	suite.Equal("", responseBody.Error, "error should be empty")
	suite.Equal("success", responseBody.Message, "message should match")
	suite.Equal("Gregory Tandiono", responseBody.Data.Name, "name should match")
	suite.Equal("gtandiono", responseBody.Data.Username, "username should match")
	suite.Equal("admin", responseBody.Data.Type, "type should match")
}

func (suite *UserHandlerTestSuite) TestUserHandler_FetchUser() {
	request, _ := http.NewRequest("GET", "/users/"+suite.userID_2.String(), nil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	h := dwhttp.NewUserHandler()
	h.UserService.Open()
	defer h.UserService.Close()
	h.ServeHTTP(response, request)

	var responseBody *getUserResponseTemplate
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	suite.Equal("", responseBody.Error, "error should be empty")
	suite.Equal("success", responseBody.Message, "message should match")
	suite.Equal("Augustus Kwok", responseBody.Data.Name, "name should match")
	suite.Equal("akwok", responseBody.Data.Username, "username should match")
	suite.Equal("admin", responseBody.Data.Type, "type should match")
}

func (suite *UserHandlerTestSuite) TestUserHandler_FetchUsers() {
	request, _ := http.NewRequest("GET", "/users", nil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	h := dwhttp.NewUserHandler()
	h.UserService.Open()
	defer h.UserService.Close()
	h.ServeHTTP(response, request)

	var responseBody *getUsersResponseTemplate
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	suite.Equal("", responseBody.Error, "error should be empty")
	suite.Equal("success", responseBody.Message, "message should match")
	suite.Equal(2, len(responseBody.Data), "user amount should match")
}

func (suite *UserHandlerTestSuite) TestUserHandler_SetName() {
	mockData := []byte(`{
		"name": "Benjamin Tandiono"
	}`)

	request, _ := http.NewRequest("PUT", "/users/"+suite.userID_1.String(), bytes.NewBuffer(mockData))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	h := dwhttp.NewUserHandler()
	h.UserService.Open()
	defer h.UserService.Close()
	h.ServeHTTP(response, request)

	var responseBody *shared.ResponseTemplate
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	suite.Equal("", responseBody.Error, "error should be empty")
	suite.Equal("success", responseBody.Message, "message should match")
}

func (suite *UserHandlerTestSuite) TestUserHandler_SetName_VerifySet() {
	request, _ := http.NewRequest("GET", "/users/"+suite.userID_1.String(), nil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	h := dwhttp.NewUserHandler()
	h.UserService.Open()
	defer h.UserService.Close()
	h.ServeHTTP(response, request)

	var responseBody *getUserResponseTemplate
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	suite.Equal("", responseBody.Error, "error should be empty")
	suite.Equal("success", responseBody.Message, "message should match")
	suite.Equal("Benjamin Tandiono", responseBody.Data.Name, "name should match")
	suite.Equal("gtandiono", responseBody.Data.Username, "username should match")
	suite.Equal("admin", responseBody.Data.Type, "type should match")
}

func (suite *UserHandlerTestSuite) TestUserHandler_RemoveUser() {
	request, _ := http.NewRequest("DELETE", "/users/"+suite.userID_2.String(), nil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	h := dwhttp.NewUserHandler()
	h.UserService.Open()
	defer h.UserService.Close()
	h.ServeHTTP(response, request)

	var responseBody *shared.ResponseTemplate
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	suite.Equal("", responseBody.Error, "error should be empty")
	suite.Equal("success", responseBody.Message, "message should match")
}

func (suite *UserHandlerTestSuite) TestUserHandler_RemoveUser_VerifyRemoval() {
	request, _ := http.NewRequest("GET", "/users/"+suite.userID_2.String(), nil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	h := dwhttp.NewUserHandler()
	h.UserService.Open()
	defer h.UserService.Close()
	h.ServeHTTP(response, request)

	var responseBody *getUserResponseTemplate
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	suite.Equal("fail", responseBody.Message, "message should match")
	suite.Equal("user does not exist", responseBody.Error, "error message should match")
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
