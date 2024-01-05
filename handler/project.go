package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tomihaapalainen/go-task-mgmt/model"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
)

func HandlePostCreateProject(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		user := c.Get("user").(model.User)

		project := model.Project{}
		if err := json.NewDecoder(c.Request().Body).Decode(&project); err != nil {
			log.Println("err decoding json: ", err)
			return errors.New("invalid request data")
		}
		project.UserID = user.ID

		if project.UserID <= 0 {
			log.Println("invalid user ID", project.UserID)
			return c.JSON(http.StatusBadRequest, schema.MessageResponse{Message: "user ID must be a positive integer"})
		}

		project.Name = strings.TrimSpace(project.Name)
		if project.Name == "" {
			log.Println("invalid project name", project.Name)
			return c.JSON(http.StatusBadRequest, schema.MessageResponse{Message: "project name must not be empty"})
		}
		project.Description = strings.TrimSpace(project.Description)

		if err := project.Create(db); err != nil {
			log.Println("err creating project:", err)
			return errors.New("unable to create project")
		}

		return c.JSON(http.StatusOK, project)
	})
}

func HandleGetProjectID(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		projectID := c.Param("id")
		pID, err := strconv.Atoi(projectID)
		if err != nil || pID <= 0 {
			return fmt.Errorf("invalid project ID '%s'", projectID)
		}

		project := model.Project{ID: pID}
		if err := project.ReadByID(db); err != nil {
			log.Println("err reading project by ID: ", err)
			return errors.New("unable to read project")
		}

		return c.JSON(
			http.StatusOK,
			project,
		)
	})
}

func HandlePatchProjectID(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		projectID := c.Param("id")
		pID, err := strconv.Atoi(projectID)
		if err != nil || pID <= 0 {
			return fmt.Errorf("invalid project ID '%s'", projectID)
		}

		project := model.Project{}
		if err := json.NewDecoder(c.Request().Body).Decode(&project); err != nil {
			log.Println("err decoding json: ", err)
			return c.JSON(http.StatusBadRequest, schema.MessageResponse{Message: "invalid request body"})
		}

		if project.Name == "" {
			return c.JSON(http.StatusBadRequest, schema.MessageResponse{Message: "project name must not be empty"})
		}
		project.ID = pID
		if err := project.Update(db); err != nil {
			log.Println("err updating project: ", err)
			return c.JSON(http.StatusInternalServerError, schema.MessageResponse{Message: "unable to update project"})
		}

		return c.JSON(http.StatusOK, project)
	})
}

func HandleDeleteProject(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		projectID := c.Param("id")

		pID, err := strconv.Atoi(projectID)
		if err != nil || pID <= 0 {
			return fmt.Errorf("invalid project ID '%s'", projectID)
		}

		project := model.Project{ID: pID}
		if err := project.Delete(db); err != nil {
			log.Println("err deleting project: ", err)
			return errors.New("unable to delete project")
		}
		return c.JSON(
			http.StatusNoContent,
			schema.MessageResponse{
				Message: fmt.Sprintf("project with ID '%d' deleted successfully", pID),
			},
		)
	})
}
