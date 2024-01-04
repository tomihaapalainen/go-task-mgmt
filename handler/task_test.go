package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/tomihaapalainen/go-task-mgmt/assert"
	"github.com/tomihaapalainen/go-task-mgmt/constants"
	"github.com/tomihaapalainen/go-task-mgmt/model"
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
			task := model.Task{}
			err = json.NewDecoder(rec.Body).Decode(&task)
			assert.AssertEq(t, err, nil)
			assert.AssertNotEq(t, task.ID, 0)
			assert.AssertEq(t, task.AssigneeID, tc.id)
			assert.AssertEq(t, task.CreatorID, tc.id)
			assert.AssertEq(t, task.Title, tc.title)
			assert.AssertEq(t, task.Content, tc.content)
			assert.AssertEq(t, task.Status, tc.status)
		})
	}
}

func TestDeleteTaskShouldPass(t *testing.T) {
	authRes := login(t, testUserIn.Email, testUserIn.Password)

	rec, c := createContextWithParams(
		"DELETE",
		"http://localhost:8080/project/:projectID/task/:id",
		"",
		[]string{"projectID", "id"},
		[]string{fmt.Sprintf("%d", testProject.ID), fmt.Sprintf("%d", testTaskForDeletion.ID)},
	)

	c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", authRes.TokenType, authRes.AccessToken))
	err := mw.JwtMiddleware(mw.PermissionRequired(tDB, "delete task")(HandleDeleteTask(tDB)))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusNoContent)
}

func TestReadTaskShouldPass(t *testing.T) {
	authRes := login(t, testUserIn.Email, testUserIn.Password)

	rec, c := createContextWithParams(
		"GET",
		"http://localhost:8080/project/:projectID/task/:id",
		"",
		[]string{"projectID", "id"},
		[]string{fmt.Sprintf("%d", testProject.ID), fmt.Sprintf("%d", testTask.ID)},
	)
	c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", authRes.TokenType, authRes.AccessToken))
	err := mw.JwtMiddleware(mw.PermissionRequired(tDB, "read task")(HandleGetTaskID(tDB)))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusOK)
	task := model.Task{}
	err = json.NewDecoder(rec.Body).Decode(&task)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, task.ID, testTask.ID)
	assert.AssertEq(t, task.AssigneeID, testTask.AssigneeID)
	assert.AssertEq(t, task.CreatorID, testTask.CreatorID)
	assert.AssertEq(t, task.Title, testTask.Title)
	assert.AssertEq(t, task.Content, testTask.Content)
	assert.AssertEq(t, task.Status, testTask.Status)
}

func TestPatchTaskShouldPass(t *testing.T) {
	authRes := login(t, testUserIn.Email, testUserIn.Password)

	jsonStr := fmt.Sprintf(
		`{"assignee_id": %d, "title": "Updated task title", "content": "Updated task content", "status": "doing"}`,
		testTask.AssigneeID,
	)
	rec, c := createContextWithParams(
		"PATCH",
		"http://localhost:8080/project/:projectID/task/:id",
		jsonStr,
		[]string{"projectID", "id"},
		[]string{fmt.Sprintf("%d", testProject.ID), fmt.Sprintf("%d", testTask.ID)},
	)
	c.Request().Header.Set("Authorization", fmt.Sprintf("%s %s", authRes.TokenType, authRes.AccessToken))
	err := mw.JwtMiddleware(mw.PermissionRequired(tDB, "update task")(HandlePatchTaskID(tDB)))(c)
	assert.AssertEq(t, err, nil)
	assert.AssertEq(t, rec.Code, http.StatusOK)
}
