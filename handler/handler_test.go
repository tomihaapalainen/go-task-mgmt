package handler

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/pressly/goose"
	"github.com/tomihaapalainen/go-task-mgmt/dotenv"
	"github.com/tomihaapalainen/go-task-mgmt/model"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
	"golang.org/x/crypto/bcrypt"
)

var tDB *sql.DB
var testAdminIn schema.UserIn
var testAdmin model.User
var testUserIn schema.UserIn
var testUser model.User
var testProject model.Project

func TestMain(m *testing.M) {
	dotenv.ParseDotenv("../.env")
	tDB, _ = sql.Open("sqlite3", "file:.///db.sqlite3?_fk=ON")

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal("err setting dialect: ", err)
	}
	if err := goose.Up(tDB, "../migrations"); err != nil {
		log.Fatal("err running migrations: ", err)
	}
	createTestAdmin()
	createTestUser()
	createTestProject()

	code := m.Run()
	if err := goose.Down(tDB, "../migrations"); err != nil {
		log.Fatal("err running migrations: ", err)
	}
	os.Exit(code)
}

func createTestAdmin() {
	testAdminIn.Email = "testadmin@example.com"
	testAdminIn.Password = "Testuser1"
	testAdmin.Email = "testadmin@example.com"
	b, _ := bcrypt.GenerateFromPassword([]byte(testAdminIn.Password), 4)
	testAdmin.PasswordHash = string(b)
	testAdmin.RoleID = 1
	err := testAdmin.Create(tDB)
	if err != nil {
		log.Fatal("err creating test admin ", err)
	}
}

func createTestUser() {
	testUserIn.Email = "testuser@example.com"
	testUserIn.Password = "Testuser1"
	testUser.Email = "testuser@example.com"
	b, _ := bcrypt.GenerateFromPassword([]byte(testUserIn.Password), 4)
	testUser.PasswordHash = string(b)
	testUser.RoleID = 3
	err := testUser.Create(tDB)
	if err != nil {
		log.Fatal("err creating test user ", err)
	}
}

func createTestProject() {
	testProject.Name = "Default test project"
	testProject.Description = "Test description"
	testProject.UserID = testAdmin.ID
	err := testProject.Create(tDB)
	if err != nil {
		log.Fatal("err creating test user ", err)
	}
}
