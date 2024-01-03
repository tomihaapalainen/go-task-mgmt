package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/tomihaapalainen/go-task-mgmt/assert"
	"github.com/tomihaapalainen/go-task-mgmt/constants"
	"github.com/tomihaapalainen/go-task-mgmt/model"
	"github.com/tomihaapalainen/go-task-mgmt/mw"
	"github.com/tomihaapalainen/go-task-mgmt/schema"
)

func TestAssignRoleShouldPass(t *testing.T) {
	authRes := login(t, testAdminIn.Email, testAdminIn.Password)

	jsonStr := fmt.Sprintf(`{"role_id": %d, "user_id": %d}`, constants.ProjectManagerRoleID, testUserForRole.ID)
	rec, c := createContext("PATCH", "http://localhost:8080/role/assign", jsonStr)
	c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", authRes.TokenType, authRes.AccessToken))
	err := mw.JwtMiddleware(mw.PermissionRequired(tDB, "manage roles")(HandlePatchAssignRole(tDB)))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusOK)
	u := model.User{ID: testUserForRole.ID}
	err = json.NewDecoder(rec.Body).Decode(&u)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, u.RoleID, constants.ProjectManagerRoleID)
	err = u.ReadByID(tDB)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, u.RoleID, constants.ProjectManagerRoleID)
}

func TestAssignRoleWithoutPermissionShouldFail(t *testing.T) {
	authRes := login(t, testUserIn.Email, testUserIn.Password)

	jsonStr := fmt.Sprintf(`{"role_id": %d, "user_id": %d}`, constants.ProjectManagerRoleID, testUserForRole.ID)
	rec, c := createContext("PATCH", "http://localhost:8080/role/assign", jsonStr)
	c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", authRes.TokenType, authRes.AccessToken))
	err := mw.JwtMiddleware(mw.PermissionRequired(tDB, "manage roles")(HandlePatchAssignRole(tDB)))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusForbidden)
	m := schema.MessageResponse{}
	err = json.NewDecoder(rec.Body).Decode(&m)
	assert.AssertEq(t, err, nil)
	if !strings.Contains(m.Message, "does not have permission") {
		t.Fatalf("'%s' did not contain 'does not have permission'", m.Message)
	}
}
