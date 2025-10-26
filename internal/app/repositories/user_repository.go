package repositories

import (
	"database/sql"
	"erp/internal/app/models"
	"errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(login string) (*models.User, error) {
	row := r.db.QueryRow(`
		SELECT id, login, password, first_name, second_name, role_id
		FROM users
		WHERE login = ?
	`, login)

	var user models.User
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.FirstName, &user.SecondName, &user.RoleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	result, err := r.db.Exec(`
		INSERT INTO users (login, password, first_name, second_name, role_id)
		VALUES (?, ?, ?, ?, ?)
	`, user.Login, user.Password, user.FirstName, user.SecondName, user.RoleID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}
