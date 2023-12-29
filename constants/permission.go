package constants

type PermissionID int

const (
	All           PermissionID = 1
	CreateProject PermissionID = 2
	ReadProject   PermissionID = 3
	UpdateProject PermissionID = 4
	DeleteProject PermissionID = 5
	CreateTask    PermissionID = 6
	ReadTask      PermissionID = 7
	UpdateTask    PermissionID = 8
	DeleteTask    PermissionID = 9
	ManageUsers   PermissionID = 10
	ManageRoles   PermissionID = 11
)
