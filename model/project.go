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

func (p *Project) ReadByID(db *sql.DB) error {
	stmt, err := db.Prepare(
		`
		SELECT user_id, name, description
		FROM project
		WHERE id = $1
		`,
	)
	if err != nil {
		return err
	}
	return stmt.QueryRow(p.ID).Scan(&p.UserID, &p.Name, &p.Description)
}

func (p *Project) Update(db *sql.DB) error {
	stmt, err := db.Prepare(
		`
		UPDATE project
		SET name = $1,
			description = $2
		WHERE id = $3
		RETURNING user_id, name, description
		`,
	)
	if err != nil {
		return err
	}
	return stmt.QueryRow(p.Name, p.Description, p.ID).Scan(&p.UserID, &p.Name, &p.Description)
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
