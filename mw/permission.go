package mw

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tomihaapalainen/go-task-mgmt/model"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
)

func PermissionRequired(db *sql.DB, permission string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			user := c.Get("user").(model.User)
			permissions := model.Permissions{}
			if err := permissions.ReadRolePermissions(db, user.RoleID); err != nil {
				return errors.New("unable to read role permissions")
			}
			hasPermission := false
			for _, p := range permissions {
				if permission == p.Name || p.Name == "all" {
					hasPermission = true
				}
			}
			if !hasPermission {
				return c.JSON(
					http.StatusForbidden,
					schema.MessageResponse{
						Message: fmt.Sprintf("user '%s' does not have permission to run this command", user.Email)})
			}
			return next(c)
		})
	}
}
