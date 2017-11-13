package http_test

import (
	"bytes"
	"datwire/pkg/apps/auth"
	dwhttp "datwire/pkg/http"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type authResponseTemplate struct {
	Message string
	Error   string
	Data    *auth.Auth
}

type AuthHandlerTestSuite struct {
	suite.Suite
}

func (suite *AuthHandlerTestSuite) TestAuthHandler_Authorization() {
	request, _ := http.NewRequest("POST", "/auth", bytes.NewBufferString(`{
		"user_id": "e1212aac-35e3-493f-b574-653832214e56",
		"password": "themostawesomepasswordintheworld",
		"hashed_password": "$2a$10$8F56R2MnW4mTwLclwQ4gw.D1v9d8WU6W.bs96ZSjRc69NW2Os2E0e"
	}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	h := dwhttp.NewAuthHandler()
	h.ServeHTTP(response, request)

	var responseBody *authResponseTemplate
	json.Unmarshal(response.Body.Bytes(), &responseBody)
	suite.Equal("", responseBody.Error, "error should be empty")
	suite.Equal("success", responseBody.Message, "message should match")
	suite.NotEmpty(responseBody.Data.Token)
}

func (suite *AuthHandlerTestSuite) TestAuthHandler_Authorization_PassDoesNotMatch() {
	request, _ := http.NewRequest("POST", "/auth", bytes.NewBufferString(`{
		"user_id": "e1212aac-35e3-493f-b574-653832214e56",
		"password": "thispasswordwillfail",
		"hashed_password": "$2a$10$8F56R2MnW4mTwLclwQ4gw.D1v9d8WU6W.bs96ZSjRc69NW2Os2E0e"
	}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	h := dwhttp.NewAuthHandler()
	h.ServeHTTP(response, request)

	var responseBody *authResponseTemplate
	json.Unmarshal(response.Body.Bytes(), &responseBody)
	suite.Equal("password does not match", responseBody.Error, "error should be empty")
	suite.Equal("fail", responseBody.Message, "message should match")
}

func TestAuthHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}
