package model

import (
	"database/sql"

	"github.com/tomihaapalainen/go-task-mgmt/constants"
)

type Task struct {
	ID         int
	ProjectID  int
	AssigneeID int
	CreatorID  int
	Title      string
	Content    string
	Status     constants.TaskStatus
}

func (t *Task) Create(db *sql.DB) error {
	stmt, err := db.Prepare(
		`
		INSERT INTO task (project_id, assignee_id, creator_id, title, content, status) values(
			$1, $2, $3, $4, $5, $6
		)
		RETURNING id
		`,
	)
	if err != nil {
		return err
	}

	return stmt.QueryRow(t.ProjectID, t.AssigneeID, t.CreatorID, t.Title, t.Content, t.Status).Scan(&t.ID)
}
