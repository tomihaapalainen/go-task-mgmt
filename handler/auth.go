package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"unicode"

	"github.com/labstack/echo/v4"
	"github.com/mattn/go-sqlite3"
	"github.com/tomihaapalainen/go-task-mgmt/model"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
	"golang.org/x/crypto/bcrypt"
)

func emailIsValid(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func passwordIsValid(s string) bool {
	count := 0
	nums := 0
	upper := 0
	lower := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			nums++
		}
		if unicode.IsLower(r) {
			lower++
		}
		if unicode.IsUpper(r) {
			upper++
		}
		count++
	}
	return count >= 8 && nums >= 1 && upper >= 1 && lower >= 1
}

func HandlePostRegister(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userIn := schema.UserIn{}
		if err := json.NewDecoder(c.Request().Body).Decode(&userIn); err != nil {
			log.Println("err decoding request body: ", err)
			return c.JSON(
				http.StatusBadRequest,
				schema.ErrorResponse{
					Message: "Invalid request data",
				},
			)
		}

		userIn.Email = strings.TrimSpace(userIn.Email)
		if !emailIsValid(userIn.Email) {
			return c.JSON(
				http.StatusBadRequest,
				schema.ErrorResponse{
					Message: "'email' must be a valid email address",
				},
			)
		}
		if !passwordIsValid(userIn.Password) {
			return c.JSON(
				http.StatusBadRequest,
				schema.ErrorResponse{
					Message: `'password' must be at least 8 characters long,
						contain 1 upper case character, 1 lower case character and 1 digit`,
				},
			)
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(userIn.Password), 4)
		if err != nil {
			log.Println("err generating password hash:", err)
			return c.JSON(
				http.StatusInternalServerError,
				schema.ErrorResponse{Message: "Error generating password hash"},
			)
		}

		user := model.User{Email: userIn.Email, PasswordHash: string(passwordHash), RoleID: 1}
		if err := user.Create(db); err != nil {
			log.Printf("err creating user %+v: %+v\n", user, err)
			if errors.Is(err, sqlite3.ErrConstraintUnique) {
				return c.JSON(
					http.StatusBadRequest,
					schema.ErrorResponse{
						Message: fmt.Sprintf("User with email '%s' already exists", user.Email),
					},
				)
			}
			return c.JSON(
				http.StatusInternalServerError,
				schema.ErrorResponse{Message: "Unable to create user"},
			)
		}

		return c.JSON(
			http.StatusOK,
			schema.UserOut{ID: user.ID, Email: user.Email, RoleID: user.RoleID},
		)
	})
}
