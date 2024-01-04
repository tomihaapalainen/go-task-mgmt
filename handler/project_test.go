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
	testCases := []struct {
		id          int
		email       string
		password    string
		projectName string
	}{
		{testAdmin.ID, testAdminIn.Email, testAdminIn.Password, "Test project by admin"},
		{testProjectManager.ID, testProjectManagerIn.Email, testProjectManagerIn.Password, "Test project by project manager"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s %s %s", tc.email, tc.password, tc.projectName), func(t *testing.T) {
			authRes := login(t, tc.email, tc.password)

			jsonStr := fmt.Sprintf(`{"user_id": %d, "name": "%s", "description": "Test description"}`, tc.id, tc.projectName)
			rec, c := createContext("POST", "http://localhost:8080/project/create", jsonStr)
			c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", authRes.TokenType, authRes.AccessToken))
			err := mw.JwtMiddleware(mw.PermissionRequired(tDB, "create project")(HandlePostCreateProject(tDB)))(c)
			assert.AssertEq(t, err, nil)
			assert.AssertEq(t, rec.Code, http.StatusOK)
			project := model.Project{}
			err = json.NewDecoder(rec.Body).Decode(&project)
			assert.AssertEq(t, err, nil)
			assert.AssertNotEq(t, project.ID, 0)
			assert.AssertEq(t, project.UserID, tc.id)
		})
	}
}

func TestPostCreateProjectWithoutPermissionShouldFail(t *testing.T) {
	authRes := login(t, testUserIn.Email, testUserIn.Password)

	jsonStr := fmt.Sprintf(`{"user_id": %d, "name": "Test project user", "description": "Test description"}`, testUser.ID)
	rec, c := createContext("POST", "http://localhost:8080/project/create", jsonStr)
	c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", authRes.TokenType, authRes.AccessToken))
	err := mw.JwtMiddleware(mw.PermissionRequired(tDB, "create project")(HandlePostCreateProject(tDB)))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusForbidden)
}

func TestDeleteProjectWithoutPermissionShouldFail(t *testing.T) {
	authRes := login(t, testUserIn.Email, testUserIn.Password)

	rec, c := createContextWithParams(
		"DELETE",
		"http://localhost:8080/project/:id",
		"",
		[]string{"id"},
		[]string{fmt.Sprintf("%d", testProject.ID)})
	c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", authRes.TokenType, authRes.AccessToken))
	err := mw.JwtMiddleware(mw.PermissionRequired(tDB, "delete project")(HandleDeleteProject(tDB)))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusForbidden)
	r := schema.MessageResponse{}
	err = json.NewDecoder(rec.Body).Decode(&r)
	assert.AssertEq(t, err, nil)
}

func TestDeleteProjectShouldPass(t *testing.T) {
	authRes := login(t, testAdminIn.Email, testAdminIn.Password)

	rec, c := createContextWithParams(
		"DELETE",
		"http://localhost:8080/project/:id",
		"",
		[]string{"id"},
		[]string{fmt.Sprintf("%d", testProjectForDeletion.ID)})
	c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", authRes.TokenType, authRes.AccessToken))
	err := mw.JwtMiddleware(mw.PermissionRequired(tDB, "delete project")(HandleDeleteProject(tDB)))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusNoContent)
	r := schema.MessageResponse{}
	err = json.NewDecoder(rec.Body).Decode(&r)
	assert.AssertEq(t, err, nil)
}

func TestReadProjectShouldPass(t *testing.T) {
	authRes := login(t, testUserIn.Email, testUserIn.Password)

	rec, c := createContextWithParams(
		"GET",
		"http://localhost:8080/project/:id",
		"",
		[]string{"id"},
		[]string{fmt.Sprintf("%d", testProject.ID)},
	)
	c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", authRes.TokenType, authRes.AccessToken))
	err := mw.JwtMiddleware(mw.PermissionRequired(tDB, "read project")(HandleGetProjectID(tDB)))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusOK)
	p := model.Project{}
	err = json.NewDecoder(rec.Body).Decode(&p)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, p.ID, testProject.ID)
	assert.AssertEq(t, p.Name, testProject.Name)
	assert.AssertEq(t, p.Description, testProject.Description)
	assert.AssertEq(t, p.UserID, testProject.UserID)
}

func TestPatchProjectShouldPass(t *testing.T) {
	authRes := login(t, testProjectManagerIn.Email, testProjectManagerIn.Password)

	jsonStr := `{"name": "Updated project name", "description": "Updated project description"}`
	rec, c := createContextWithParams(
		"PATCH",
		"http://localhost:8080/project/:id",
		jsonStr,
		[]string{"id"},
		[]string{fmt.Sprintf("%d", testProject.ID)},
	)
	c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", authRes.TokenType, authRes.AccessToken))
	err := mw.JwtMiddleware(mw.PermissionRequired(tDB, "update project")(HandlePatchProjectID(tDB)))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusOK)
	p := model.Project{}
	err = json.NewDecoder(rec.Body).Decode(&p)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, p.ID, testProject.ID)
	assert.AssertNotEq(t, p.Name, testProject.Name)
	assert.AssertNotEq(t, p.Description, testProject.Description)
	assert.AssertEq(t, p.UserID, testProject.UserID)
}
