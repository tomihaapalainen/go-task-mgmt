package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tomihaapalainen/go-task-mgmt/model"
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
