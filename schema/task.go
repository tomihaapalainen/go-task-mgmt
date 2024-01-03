package schema

import "github.com/tomihaapalainen/go-task-mgmt/constants"

type TaskIn struct {
	AssigneeID int                  `json:"assignee_id"`
	Title      string               `json:"title"`
	Content    string               `json:"content"`
	Status     constants.TaskStatus `json:"status"`
}
