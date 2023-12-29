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
	req := httptest.NewRequest("POST", "http://localhost:8080/auth/register", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := HandlePostRegister(tDB)(c)
	if err != nil {
		t.Fatalf("register failed: %+v\n", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status code - expected %d - was %d - body: %s\n", http.StatusOK, rec.Code, rec.Body)
	}
	u := schema.UserOut{}
	if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
		t.Fatal("err decoding response body: ", err)
	}
	if u.ID <= 0 {
		t.Fatalf("invalid user ID %d\n", u.ID)
	}
	if u.RoleID != 1 {
		t.Fatalf("role id - expected %d - was %d\n", 1, u.RoleID)
	}
	if u.Email != "test@example.com" {
		t.Fatalf("email - expected 'test@example.com' - was '%s'\n", u.Email)
	}
}

func TestPostRegisterUserWithoutEmailShouldFail(t *testing.T) {
	jsonStr := `{"password": "testpassword"}`
	req := httptest.NewRequest("POST", "http://localhost:8080/auth/register", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := HandlePostRegister(tDB)(c)
	if err != nil {
		t.Fatalf("err handling post register: %+v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status code - expected %d - was %d", http.StatusBadRequest, rec.Code)
	}
	res := schema.ErrorResponse{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatal("err decoding response: ", err)
	}
}

func TestPostRegisterUserWithEmptyEmailShouldFail(t *testing.T) {
	jsonStr := `{"email": "", "password": "testpassword"}`
	req := httptest.NewRequest("POST", "http://localhost:8080/auth/register", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := HandlePostRegister(tDB)(c)
	if err != nil {
		t.Fatalf("err handling post register: %+v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status code - expected %d - was %d", http.StatusBadRequest, rec.Code)
	}
	res := schema.ErrorResponse{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatal("err decoding response: ", err)
	}
}

func TestPostRegisterUserWithoutPasswordShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com"}`
	req := httptest.NewRequest("POST", "http://localhost:8080/auth/register", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := HandlePostRegister(tDB)(c)
	if err != nil {
		t.Fatalf("err handling post register: %+v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status code - expected %d - was %d", http.StatusBadRequest, rec.Code)
	}
	res := schema.ErrorResponse{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatal("err decoding response: ", err)
	}
}

func TestPostRegisterUserWithEmptyPasswordShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": ""}`
	req := httptest.NewRequest("POST", "http://localhost:8080/auth/register", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := HandlePostRegister(tDB)(c)
	if err != nil {
		t.Fatalf("err handling post register: %+v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status code - expected %d - was %d", http.StatusBadRequest, rec.Code)
	}
	res := schema.ErrorResponse{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatal("err decoding response: ", err)
	}
}

func TestPostRegisterUserWithPasswordWithoutDigitShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "Testtest"}`
	req := httptest.NewRequest("POST", "http://localhost:8080/auth/register", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := HandlePostRegister(tDB)(c)
	if err != nil {
		t.Fatalf("err handling post register: %+v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status code - expected %d - was %d", http.StatusBadRequest, rec.Code)
	}
	res := schema.ErrorResponse{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatal("err decoding response: ", err)
	}
}

func TestPostRegisterUserWithPasswordWithoutUpperCaseCharacterShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "testtest1"}`
	req := httptest.NewRequest("POST", "http://localhost:8080/auth/register", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := HandlePostRegister(tDB)(c)
	if err != nil {
		t.Fatalf("err handling post register: %+v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status code - expected %d - was %d", http.StatusBadRequest, rec.Code)
	}
	res := schema.ErrorResponse{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatal("err decoding response: ", err)
	}
}

func TestPostRegisterUserWithPasswordWithoutLowerCaseCharacterShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "TESTTEST1"}`
	req := httptest.NewRequest("POST", "http://localhost:8080/auth/register", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := HandlePostRegister(tDB)(c)
	if err != nil {
		t.Fatalf("err handling post register: %+v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status code - expected %d - was %d", http.StatusBadRequest, rec.Code)
	}
	res := schema.ErrorResponse{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatal("err decoding response: ", err)
	}
}

func TestPostRegisterUserWithShortPasswordShouldFail(t *testing.T) {
	jsonStr := `{"email": "test@example.com", "password": "Test123"}`
	req := httptest.NewRequest("POST", "http://localhost:8080/auth/register", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := HandlePostRegister(tDB)(c)
	if err != nil {
		t.Fatalf("err handling post register: %+v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status code - expected %d - was %d", http.StatusBadRequest, rec.Code)
	}
	res := schema.ErrorResponse{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatal("err decoding response: ", err)
	}
}

func TestPostLogInWithValidCredentialsShouldPass(t *testing.T) {
	jsonStr := `{"email": "testauth@example.com", "password": "Testpassword1"}`
	req := httptest.NewRequest("POST", "http://localhost:8080/auth/register", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	err := HandlePostRegister(tDB)(c)
	if err != nil {
		t.Fatalf("register failed: %+v\n", err)
	}

	req = httptest.NewRequest("POST", "http://localhost:8080/auth/login", bytes.NewBuffer([]byte(jsonStr)))
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = HandlePostLogIn(tDB)(c)
	if err != nil {
		t.Fatalf("log in failed: %+v\n", err)
	}
	if rec.Code != 200 {
		t.Fatalf("status code - expected %d - was %d\n", http.StatusOK, rec.Code)
	}
	authResponse := schema.AuthResponse{}
	if err := json.NewDecoder(rec.Body).Decode(&authResponse); err != nil {
		t.Fatal(err)
	}
}
