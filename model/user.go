package model

import "database/sql"

type User struct {
	ID           int
	Email        string
	PasswordHash string
	RoleID       int
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
