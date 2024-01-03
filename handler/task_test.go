package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/tomihaapalainen/go-task-mgmt/assert"
	"github.com/tomihaapalainen/go-task-mgmt/constants"
	"github.com/tomihaapalainen/go-task-mgmt/mw"
)

func TestPostCreateTask(t *testing.T) {
	testCases := []struct {
		id       int
		email    string
		password string
		title    string
		content  string
		status   constants.TaskStatus
	}{
		{testAdmin.ID, testAdminIn.Email, testAdminIn.Password, "Admin task", "Admin task content", constants.Todo},
		{testProjectManager.ID, testProjectManagerIn.Email, testProjectManagerIn.Password, "Project manager task", "Project manager task content", constants.Todo},
		{testUser.ID, testUserIn.Email, testUserIn.Password, "User task", "User task content", constants.Todo},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s %s %s %s %s", tc.email, tc.password, tc.title, tc.content, tc.status), func(t *testing.T) {
			authRes := login(t, tc.email, tc.password)

			jsonStr := fmt.Sprintf(`{"assignee_id": %d, "title": "%s", "content": "%s", "status": "%s"}`, tc.id, tc.title, tc.content, tc.status)
			rec, c := createContextWithParams(
				"POST",
				"http://localhost:8080/project/:projectID/task/create",
				jsonStr,
				[]string{"projectID"},
				[]string{fmt.Sprintf("%d", testProject.ID)},
			)
			c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", authRes.TokenType, authRes.AccessToken))
			err := mw.JwtMiddleware(mw.PermissionRequired(tDB, "create task")(HandlePostCreateTask(tDB)))(c)
			assert.AssertEq(t, err, nil)
			assert.AssertEq(t, rec.Code, http.StatusOK)
		})
	}
}
