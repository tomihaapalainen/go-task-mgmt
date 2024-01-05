package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tomihaapalainen/go-task-mgmt/model"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
)

func HandlePatchAssignRole(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		ar := schema.AssignRole{}
		if err := json.NewDecoder(c.Request().Body).Decode(&ar); err != nil {
			log.Println("err decoding body:", err)
			return errors.New("invalid request body")
		}
		u := model.User{ID: ar.UserID, RoleID: ar.RoleID}
		if err := u.UpdateRole(db); err != nil {
			log.Println("err updating role:", err)
			return errors.New("unable to update role")
		}
		if err := u.ReadByID(db); err != nil {
			log.Println("err reading user by id:", err)
			return errors.New("unable to read user by ID")
		}
		return c.JSON(http.StatusOK, u)
	})
}
