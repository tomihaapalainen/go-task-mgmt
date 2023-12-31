package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/tomihaapalainen/go-task-mgmt/assert"
	"github.com/tomihaapalainen/go-task-mgmt/model"
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
	assert.AssertEq(t, u.RoleID, 1)
	assert.AssertEq(t, u.Email, "test@example.com")
}

func TestPostRegisterUserWithoutEmailShouldFail(t *testing.T) {
	jsonStr := `{"password": "testpassword"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.MessageResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithEmptyEmailShouldFail(t *testing.T) {
	jsonStr := `{"email": "", "password": "testpassword"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.MessageResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithoutPasswordShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.MessageResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithEmptyPasswordShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": ""}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.MessageResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithPasswordWithoutDigitShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "Testtest"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.MessageResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithPasswordWithoutUpperCaseCharacterShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "testtest1"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.MessageResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithPasswordWithoutLowerCaseCharacterShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "TESTTEST1"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.MessageResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithShortPasswordShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "Test123"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.MessageResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
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
