package model

import (
	"database/sql"
)

type Project struct {
	ID          int
	UserID      int
	Name        string
	Description string
}

func (p *Project) Create(db *sql.DB) error {
	stmt, err := db.Prepare(
		`
		INSERT INTO project (user_id, name, description) values ($1, $2, $3) RETURNING id
		`,
	)
	if err != nil {
		return err
	}
	return stmt.QueryRow(p.UserID, p.Name, p.Description).Scan(&p.ID)
}

func (p *Project) Delete(db *sql.DB) error {
	stmt, err := db.Prepare(
		`
		DELETE FROM project
		WHERE id = $1
		`,
	)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.ID)
	return err
}