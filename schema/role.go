package schema

import "github.com/tomihaapalainen/go-task-mgmt/constants"

type AssignRole struct {
	RoleID constants.RoleID `json:"role_id"`
	UserID int              `json:"user_id"`
}
