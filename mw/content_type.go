package mw

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
)

func ContentApplicationJSONOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		contentType := c.Request().Header.Get("Content-Type")
		if contentType != "application/json" {
			return c.JSON(
				http.StatusBadRequest,
				schema.MessageResponse{Message: fmt.Sprintf("invalid Content-Type: '%s'", contentType)})
		}
		return next(c)
	}
}
