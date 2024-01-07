package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/tomihaapalainen/go-task-mgmt/assert"
	"github.com/tomihaapalainen/go-task-mgmt/constants"
	"github.com/tomihaapalainen/go-task-mgmt/model"
	"github.com/tomihaapalainen/go-task-mgmt/mw"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
)

func TestPostRegisterUserShouldPass(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "Testpassword1"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)

	assert.AssertEq(t, rec.Code, http.StatusOK)
	u := model.User{}
	err = json.NewDecoder(rec.Body).Decode(&u)
	assert.AssertEq(t, err, nil)
	assert.AssertNotEq(t, u.ID, 0)
	assert.AssertEq(t, u.RoleID, constants.UserRoleID)
	assert.AssertEq(t, u.Email, "test@example.com")
}

func TestPostRegisterShouldFail(t *testing.T) {
	testCases := []struct {
		requestBody string
	}{
		{`{"email": "", "password": "Testpass1"}`},
		{`{"email": "testuser123@example.com", "password": ""}`},
		{`{"email": "", "password": ""}`},
		{`{"email": "testuser1234@example.com"}`},
		{`{"password": "Testpass1"}`},
		{`{"email": "testuser1234@example.com", "password": "Testpass"}`},
		{`{"email": "testuser1234@example.com", "password": "testpas1"}`},
		{`{"email": "testuser1234@example.com", "password": "TESTPAS1"}`},
		{`{"email": "testuser1234@example.com", "password": "Testpa1"}`},
		{`{}`},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Body: %s", tc.requestBody), func(t *testing.T) {
			rec, c := createContext("POST", "http://localhost:8080/auth/register", tc.requestBody)

			err := HandlePostRegister(tDB)(c)
			assert.AssertEq(t, err, nil)
			assert.AssertEq(t, rec.Code, http.StatusBadRequest)
			res := schema.MessageResponse{}
			err = json.NewDecoder(rec.Body).Decode(&res)
			assert.AssertEq(t, err, nil)
		})
	}
}

func TestPostLogInWithValidCredentialsShouldPass(t *testing.T) {
	jsonStr := `{"email": "testauth@example.com", "password": "Testpassword1"}`

	_, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)

	rec, c := createContext("POST", "http://localhost:8080/auth/login", jsonStr)

	err = HandlePostLogIn(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusOK)
	authResponse := schema.AuthResponse{}
	err = json.NewDecoder(rec.Body).Decode(&authResponse)
	assert.AssertEq(t, err, nil)
}

func TestPostLoginWithInvalidContentTypeShouldFail(t *testing.T) {
	jsonStr := `{"email": "testauth@example.com", "password": "Testpassword1"}`
	req := httptest.NewRequest("POST", "http://localhost:8080/auth/register", bytes.NewBuffer([]byte(jsonStr)))

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := mw.ContentApplicationJSONOnly(HandlePostLogIn(tDB))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
}
