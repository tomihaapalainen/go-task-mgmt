package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pressly/goose"
	"github.com/tomihaapalainen/go-task-mgmt/assert"
	"github.com/tomihaapalainen/go-task-mgmt/dotenv"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
)

var tDB *sql.DB

func TestMain(m *testing.M) {
	dotenv.ParseDotenv("../.env")
	tDB, _ = sql.Open("sqlite3", "file:.///db.sqlite3?_fk=ON")

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal("err setting dialect: ", err)
	}
	if err := goose.Up(tDB, "../migrations"); err != nil {
		log.Fatal("err running migrations: ", err)
	}
	code := m.Run()
	os.Remove("db.sqlite3")
	os.Exit(code)
}

func TestPostRegisterUserShouldPass(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "Testpassword1"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)

	assert.AssertEq(t, rec.Code, http.StatusOK)
	u := schema.UserOut{}
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
	res := schema.ErrorResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithEmptyEmailShouldFail(t *testing.T) {
	jsonStr := `{"email": "", "password": "testpassword"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.ErrorResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithoutPasswordShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.ErrorResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithEmptyPasswordShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": ""}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.ErrorResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithPasswordWithoutDigitShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "Testtest"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.ErrorResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithPasswordWithoutUpperCaseCharacterShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "testtest1"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.ErrorResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithPasswordWithoutLowerCaseCharacterShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "TESTTEST1"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.ErrorResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
}

func TestPostRegisterUserWithShortPasswordShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "Test123"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)

	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusBadRequest)
	res := schema.ErrorResponse{}
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

func createContext(method, url, jsonStr string) (*httptest.ResponseRecorder, echo.Context) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	return rec, c
}
