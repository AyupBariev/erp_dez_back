package repositories

import (
	"database/sql"
	"erp/internal/app/models"
	"fmt"
)

type EngineerRepository struct {
	db *sql.DB
}

func NewEngineerRepository(db *sql.DB) *EngineerRepository {
	return &EngineerRepository{db: db}
}

func (r *EngineerRepository) Create(engineer *models.Engineer) error {
	_, err := r.db.Exec(`
		INSERT INTO engineers (first_name, second_name, username, phone, telegram_id, is_approved)
		VALUES (?, ?, ?, ?, ?, ?)
	`, engineer.FirstName, engineer.SecondName, engineer.Username, engineer.Phone, engineer.TelegramID, engineer.IsApproved)
	return err
}

func (r *EngineerRepository) queryEngineers(where string, args ...interface{}) ([]*models.Engineer, error) {
	query := `
		SELECT e.id, e.first_name, e.second_name, e.username, e.phone, e.telegram_id, e.is_approved
# 		,		       COALESCE(s.is_working, FALSE) AS is_working
		FROM engineers e
# 		LEFT JOIN engineer_shifts s ON s.engineer_id = e.id AND s.work_date = ?
	`

	// Добавляем условие, если есть
	if where != "" {
		query += " WHERE " + where
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var engineers []*models.Engineer
	for rows.Next() {
		e := &models.Engineer{}
		if err := rows.Scan(
			&e.ID,
			&e.FirstName,
			&e.SecondName,
			&e.Username,
			&e.Phone,
			&e.TelegramID,
			&e.IsApproved,
			//TODO добавить поддержку графиков СИ &e.IsWorking,
		); err != nil {
			return nil, err
		}
		engineers = append(engineers, e)
	}

	return engineers, rows.Err()
}

func (r *EngineerRepository) ListWorking(date string) ([]*models.Engineer, error) {
	//todo вторая часть график си
	//return r.queryEngineers("e.is_approved = TRUE AND s.is_working = TRUE", date)
	return r.queryEngineers("1=1")
}

func (r *EngineerRepository) FindByTelegramID(telegramID int64) (*models.Engineer, error) {
	engineer := &models.Engineer{}
	err := r.db.QueryRow(`
	SELECT id, username, first_name, second_name, phone, telegram_id, is_approved
		FROM engineers WHERE telegram_id = ?
		`, telegramID).Scan(&engineer.ID, &engineer.Username, &engineer.FirstName, &engineer.SecondName, &engineer.Phone, &engineer.TelegramID, &engineer.IsApproved)
	return engineer, err
}

func (r *EngineerRepository) FindByID(ID int64) (*models.Engineer, error) {
	engineer := &models.Engineer{}
	err := r.db.QueryRow(`
	SELECT id, username, first_name, second_name, phone, telegram_id, is_approved
		FROM engineers WHERE id = ?
		`, ID).Scan(&engineer.ID, &engineer.Username, &engineer.FirstName, &engineer.SecondName, &engineer.Phone, &engineer.TelegramID, &engineer.IsApproved)
	return engineer, err
}

func (r *EngineerRepository) FindApprovedByID(ID int64) (*models.Engineer, error) {
	engineer := &models.Engineer{}
	err := r.db.QueryRow(`
	SELECT id, username, first_name, second_name, phone, telegram_id, is_approved
		FROM engineers WHERE id = ? AND is_approved = true
		`, ID).Scan(&engineer.ID, &engineer.Username, &engineer.FirstName, &engineer.SecondName, &engineer.Phone, &engineer.TelegramID, &engineer.IsApproved)
	return engineer, err
}

func (r *EngineerRepository) ApproveByID(engineerID int64) (*models.Engineer, error) {
	res, err := r.db.Exec(`
		UPDATE engineers 
		SET is_approved = true 
		WHERE id = ?
	`, engineerID)
	if err != nil {
		return nil, err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("engineer with id=%d not found", engineerID)
	}

	// Повторно читаем инженера, чтобы вернуть актуальные данные
	return r.FindByID(engineerID)
}
