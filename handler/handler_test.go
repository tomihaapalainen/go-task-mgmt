package handler

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/pressly/goose"
	"github.com/tomihaapalainen/go-task-mgmt/dotenv"
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
	if err := goose.Down(tDB, "../migrations"); err != nil {
		log.Println("ERR DOWN")
		log.Fatal("err running migrations: ", err)
	}
	os.Exit(code)
}
