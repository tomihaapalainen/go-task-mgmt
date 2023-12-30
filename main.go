package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/tomihaapalainen/go-task-mgmt/config"
	"github.com/tomihaapalainen/go-task-mgmt/dotenv"
	"github.com/tomihaapalainen/go-task-mgmt/handler"
	"github.com/tomihaapalainen/go-task-mgmt/mw"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dotenv.ParseDotenv(".env")

	env := flag.String("env", "dev", "run environment dev|test|prod")
	port := flag.String("port", ":8080", "application port, e.g. ':8080'")
	if env == nil {
		log.Fatal("env flag was nil")

	}
	if port == nil {
		log.Fatal("port flag was nil")
	}
	config.ENV = *env
	config.PORT = *port

	db, err := sql.Open("sqlite3", "file:.///db.sqlite3?_fk=ON&_journal=WAL")
	if err != nil {
		log.Fatal("err opening database", err)
	}

	e := echo.New()

	authGroup := e.Group("/auth")
	authGroup.POST("/register", handler.HandlePostRegister(db))
	authGroup.POST("/login", handler.HandlePostLogIn(db))

	projectGroup := e.Group("/project", mw.JwtMiddleware)
	projectGroup.POST("/create", handler.HandlePostCreateProject(db), mw.PermissionRequired(db, "create project"))

	e.Start(config.PORT)
}
