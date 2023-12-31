package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
		if err := project.Create(db); err != nil {
			log.Println("err creating project: ", err)
			return errors.New("unable to create project")
		}

		return c.JSON(http.StatusOK, project)
	})
}

func HandleDeleteProject(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		projectID := c.Param("id")

		pID, err := strconv.Atoi(projectID)
		if err != nil {
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
