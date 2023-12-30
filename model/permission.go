package model

import "database/sql"

type Permission struct {
	ID   int
	Name string
}

type Permissions []Permission

func (ps *Permissions) ReadRolePermissions(db *sql.DB, roleID int) error {
	stmt, err := db.Prepare(
		`
		SELECT p.id, p.name
		FROM role_permission rp
		INNER JOIN permission p
		ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		`,
	)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(roleID)
	if err != nil {
		return err
	}

	for rows.Next() {
		p := Permission{}
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return err
		}
		*ps = append(*ps, p)
	}
	return nil
}
