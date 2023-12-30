package utils

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func ReadAuthorizationToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	split := strings.Split(authHeader, "Bearer ")
	if len(split) != 2 {
		return ""
	}
	return split[1]
}

func ParseToken(s string) (*jwt.Token, error) {
	token, err := jwt.Parse(s, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("GO_TASK_MGMT_SIGNING_SECRET")), nil
	})
	return token, err
}

func ParseClaims(s string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	var err error
	_, err = jwt.ParseWithClaims(s, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("GO_TASK_MGMT_SIGNING_SECRET")), nil
	})
	return claims, err
}
