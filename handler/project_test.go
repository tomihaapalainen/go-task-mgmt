package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/tomihaapalainen/go-task-mgmt/assert"
	"github.com/tomihaapalainen/go-task-mgmt/model"
	"github.com/tomihaapalainen/go-task-mgmt/mw"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
)

func TestPostCreateProjectShouldPass(t *testing.T) {
	jsonStr := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, testUserIn.Email, testUserIn.Password)

	rec, c := createContext("POST", "http://localhost:8080/auth/login", jsonStr)
	err := HandlePostLogIn(tDB)(c)
	assert.AssertEq(t, err, nil)
	res := schema.AuthResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)

	jsonStr = fmt.Sprintf(`{"user_id": %d, "name": "Test project", "description": "Test description"}`, testUser.ID)
	rec, c = createContext("POST", "http://localhost:8080/project/create", jsonStr)
	c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", res.TokenType, res.AccessToken))
	err = mw.JwtMiddleware(mw.PermissionRequired(tDB, "create project")(HandlePostCreateProject(tDB)))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusOK)
	project := model.Project{}
	err = json.NewDecoder(rec.Body).Decode(&project)
	assert.AssertEq(t, err, nil)
	assert.AssertNotEq(t, project.ID, 0)
	assert.AssertEq(t, project.UserID, testUser.ID)
}

func TestDeleteProjectShouldPass(t *testing.T) {
	jsonStr := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, testUserIn.Email, testUserIn.Password)

	rec, c := createContext("POST", "http://localhost:8080/auth/login", jsonStr)
	err := HandlePostLogIn(tDB)(c)
	assert.AssertEq(t, err, nil)
	res := schema.AuthResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)

	e := echo.New()
	req := httptest.NewRequest("DELETE", "http://localhost:8080/project/:id", nil)
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", testProject.ID))
	c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", res.TokenType, res.AccessToken))
	err = mw.JwtMiddleware(mw.PermissionRequired(tDB, "delete project")(HandleDeleteProject(tDB)))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusOK)
	r := schema.MessageResponse{}
	err = json.NewDecoder(rec.Body).Decode(&r)
	assert.AssertEq(t, err, nil)
}
