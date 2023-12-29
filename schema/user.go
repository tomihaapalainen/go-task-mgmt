package schema

type UserIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserOut struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	RoleID int    `json:"role_id"`
}
