package model

import (
	"database/sql"

	"github.com/tomihaapalainen/go-task-mgmt/constants"
)

type User struct {
	ID           int              `json:"id"`
	Email        string           `json:"email"`
	PasswordHash string           `json:"-"`
	RoleID       constants.RoleID `json:"role_id"`
}

func (u *User) Create(db *sql.DB) error {
	stmt, err := db.Prepare(
		`
		INSERT INTO user (email, password_hash, role_id) values ($1, $2, $3) RETURNING id
		`,
	)
	if err != nil {
		return err
	}
	return stmt.QueryRow(u.Email, u.PasswordHash, u.RoleID).Scan(&u.ID)
}

func (u *User) ReadByID(db *sql.DB) error {
	stmt, err := db.Prepare(
		`
		SELECT email, password_hash, role_id
		FROM user
		WHERE id = $1
		`,
	)
	if err != nil {
		return err
	}
	return stmt.QueryRow(u.ID).Scan(&u.Email, &u.PasswordHash, &u.RoleID)
}

func (u *User) ReadByEmail(db *sql.DB) error {
	stmt, err := db.Prepare(
		`
		SELECT id, password_hash, role_id
		FROM user
		WHERE email = $1
		`,
	)
	if err != nil {
		return err
	}
	return stmt.QueryRow(u.Email).Scan(&u.ID, &u.PasswordHash, &u.RoleID)
}

func (u *User) UpdateRole(db *sql.DB) error {
	stmt, err := db.Prepare(
		`
		UPDATE user
		SET role_id = $1
		WHERE id = $2
		`,
	)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.RoleID, u.ID)
	return err
}
