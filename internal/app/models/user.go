package models

type Role struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

type Permission struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Type        string `db:"type" json:"type"` // backend | frontend
	Description string `db:"description" json:"description"`
}

type RolePermission struct {
	RoleID       int64 `db:"role_id"`
	PermissionID int64 `db:"permission_id"`
}

type User struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Login      string `json:"login"`
	Password   string `json:"-"`
	RoleID     int64  `db:"role_id" json:"role_id"`
	Role       *Role  `json:"role,omitempty"`
}
