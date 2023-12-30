package mw

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tomihaapalainen/go-task-mgmt/model"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
	"github.com/tomihaapalainen/go-task-mgmt/utils"
)

func JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr := utils.ReadAuthorizationToken(c)
		if tokenStr == "" {
			return c.JSON(
				http.StatusUnauthorized,
				schema.ErrorResponse{Message: "Missing authorization header"},
			)
		}

		claims, err := utils.ParseClaims(tokenStr)
		if err != nil {
			log.Println("err parsing auth token: ", err)
			return c.JSON(
				http.StatusUnauthorized,
				schema.ErrorResponse{Message: "Error parsing authorization token"},
			)
		}
		err = claims.Valid()
		if err != nil {
			log.Println("err invalid claims: ", err)
			return c.JSON(
				http.StatusUnauthorized,
				schema.ErrorResponse{Message: "Invalid claims"},
			)
		}
		userData := claims["data"].(string)
		user := model.User{}
		if err := json.Unmarshal([]byte(userData), &user); err != nil {
			log.Println("err unmarshaling: ", err)
			return errors.New("internal server error")
		}
		c.Set("user", user)
		return next(c)
	}
}
