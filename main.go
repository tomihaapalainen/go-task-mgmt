package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	config.ENV = *env
	config.PORT = *port

	db, err := sql.Open("sqlite3", "file:.///db.sqlite3?_fk=ON&_journal=WAL")
	if err != nil {
		log.Fatal("err opening database", err)
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost"},
		AllowHeaders: []string{"Accept", "Content-Type", "Origin"},
		AllowMethods: []string{"DELETE", "GET", "OPTIONS", "PATCH", "POST"},
	}))

	e.Use(mw.ContentTypeApplicationJSONOnly)

	authGroup := e.Group("/auth")
	authGroup.POST("/register", handler.HandlePostRegister(db))
	authGroup.POST("/login", handler.HandlePostLogIn(db))

	roleGroup := e.Group("/role", mw.JwtMiddleware)
	roleGroup.PATCH("/assign", handler.HandlePatchAssignRole(db), mw.PermissionRequired(db, "manage roles"))

	projectGroup := e.Group("/project", mw.JwtMiddleware)
	projectGroup.POST("/create", handler.HandlePostCreateProject(db), mw.PermissionRequired(db, "create project"))
	projectGroup.GET("/:id", handler.HandleGetProjectID(db), mw.PermissionRequired(db, "read project"))
	projectGroup.PATCH("/:id", handler.HandlePatchProjectID(db), mw.PermissionRequired(db, "update project"))
	projectGroup.DELETE("/:id", handler.HandleDeleteProject(db), mw.PermissionRequired(db, "delete project"))

	taskGroup := projectGroup.Group("/:projectID")
	taskGroup.POST("/task/create", handler.HandlePostCreateTask(db), mw.PermissionRequired(db, "create task"))
	taskGroup.GET("/task/:id", handler.HandleGetTaskID(db), mw.PermissionRequired(db, "read task"))
	taskGroup.PATCH("/task/:id", handler.HandlePatchTaskID(db), mw.PermissionRequired(db, "update task"))
	taskGroup.POST("/task/:id", handler.HandleDeleteTask(db), mw.PermissionRequired(db, "delete task"))

	e.Start(config.PORT)
}
