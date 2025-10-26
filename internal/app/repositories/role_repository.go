package repositories

import (
	"database/sql"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetRoleIDByName(name string) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT id FROM roles WHERE name = ?", name).Scan(&id)
	return id, err
}
