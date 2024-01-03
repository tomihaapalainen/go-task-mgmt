package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pressly/goose"
	"github.com/tomihaapalainen/go-task-mgmt/assert"
	"github.com/tomihaapalainen/go-task-mgmt/constants"
	"github.com/tomihaapalainen/go-task-mgmt/dotenv"
	"github.com/tomihaapalainen/go-task-mgmt/model"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
	"golang.org/x/crypto/bcrypt"
)

var tDB *sql.DB
var testAdminIn schema.UserIn
var testAdmin model.User
var testProjectManagerIn schema.UserIn
var testProjectManager model.User
var testUserIn schema.UserIn
var testUser model.User
var testUserForRoleIn schema.UserIn
var testUserForRole model.User
var testProject model.Project
var testProjectForDeletion model.Project
var testTaskForDeletion model.Task

func TestMain(m *testing.M) {
	dotenv.ParseDotenv("../.env")
	tDB, _ = sql.Open("sqlite3", "file:.///db.sqlite3?_fk=ON")

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal("err setting dialect: ", err)
	}
	if err := goose.Up(tDB, "../migrations"); err != nil {
		log.Fatal("err running migrations: ", err)
	}

	testAdminIn, testAdmin = createTestUserWithRole("testadmin@example.com", "Testpass1", constants.AdminRoleID)
	testProjectManagerIn, testProjectManager = createTestUserWithRole("testprojectmanager@example.com", "Testpass1", constants.ProjectManagerRoleID)
	testUserIn, testUser = createTestUserWithRole("testuser@example.com", "Testpass1", constants.UserRoleID)
	testUserForRoleIn, testUserForRole = createTestUserWithRole("testuserforrole@example.com", "Testpass1", constants.UserRoleID)
	testProject = createTestProject("Test project", testAdmin.ID)
	testProjectForDeletion = createTestProject("Test project for deletion", testAdmin.ID)
	testTaskForDeletion = createTestTask(testUser.ID, testUser.ID, "Test user task", "Test user task content", constants.Todo)

	code := m.Run()
	if err := goose.Down(tDB, "../migrations"); err != nil {
		log.Fatal("err running migrations: ", err)
	}
	os.Exit(code)
}

func createContext(method, url, jsonStr string) (*httptest.ResponseRecorder, echo.Context) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	return rec, c
}

func createContextWithParams(method, url, jsonStr string, names []string, values []string) (*httptest.ResponseRecorder, echo.Context) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetParamNames(names...)
	c.SetParamValues(values...)
	return rec, c
}

func createTestUserWithRole(email, password string, roleID constants.RoleID) (schema.UserIn, model.User) {
	userIn := schema.UserIn{}
	user := model.User{}
	userIn.Email = email
	userIn.Password = password
	user.Email = email
	user.RoleID = roleID
	b, _ := bcrypt.GenerateFromPassword([]byte(userIn.Email), 4)
	user.PasswordHash = string(b)
	err := user.Create(tDB)
	if err != nil {
		log.Fatal("err creating test user: ", err)
	}
	return userIn, user
}

func createTestProject(name string, userID int) model.Project {
	p := model.Project{}
	p.Name = name
	p.Description = "Test description"
	p.UserID = userID
	err := p.Create(tDB)
	if err != nil {
		log.Fatal("err creating test user ", err)
	}
	return p
}

func createTestTask(assigneeID, creatorID int, title, content string, status constants.TaskStatus) model.Task {
	t := model.Task{
		ProjectID:  testProject.ID,
		AssigneeID: assigneeID,
		CreatorID:  creatorID,
		Title:      title,
		Content:    content,
		Status:     status,
	}
	err := t.Create(tDB)
	if err != nil {
		log.Fatal("err creating task: ", err)
	}
	return t
}

func login(t *testing.T, email, password string) schema.AuthResponse {
	jsonStr := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
	rec, c := createContext("POST", "http://localhost:8080/auth/login", jsonStr)
	err := HandlePostLogIn(tDB)(c)
	assert.AssertEq(t, err, nil)
	res := schema.AuthResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)
	return res
}
