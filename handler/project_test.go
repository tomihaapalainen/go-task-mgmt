package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/tomihaapalainen/go-task-mgmt/assert"
	"github.com/tomihaapalainen/go-task-mgmt/model"
	"github.com/tomihaapalainen/go-task-mgmt/mw"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
)

func TestPostCreateProjectShouldPass(t *testing.T) {
	jsonStr := `{"email": "testproject@example.com", "password": "Testpassword1"}`

	rec, c := createContext("POST", "http://localhost:8080/auth/register", jsonStr)
	err := HandlePostRegister(tDB)(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusOK)
	user := model.User{}
	err = json.NewDecoder(rec.Body).Decode(&user)
	assert.AssertEq(t, err, nil)

	rec, c = createContext("POST", "http://localhost:8080/auth/login", jsonStr)
	err = HandlePostLogIn(tDB)(c)
	assert.AssertEq(t, err, nil)
	res := schema.AuthResponse{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	assert.AssertEq(t, err, nil)

	jsonStr = fmt.Sprintf(`{"user_id": %d, "name": "Test project", "description": "Test description"}`, user.ID)
	rec, c = createContext("POST", "http://localhost:8080/project/create", jsonStr)
	c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", res.TokenType, res.AccessToken))
	err = mw.JwtMiddleware(mw.PermissionRequired(tDB, "create project")(HandlePostCreateProject(tDB)))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusOK)
	project := model.Project{}
	err = json.NewDecoder(rec.Body).Decode(&project)
	assert.AssertEq(t, err, nil)
	assert.AssertNotEq(t, project.ID, 0)
	assert.AssertEq(t, project.UserID, user.ID)
}
