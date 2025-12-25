package repositories

import (
	"database/sql"
	"erp/internal/app/dto"
	"errors"
	"fmt"
	"time"
)

type DictionaryRepository struct {
	db *sql.DB
}

func NewDictionaryRepository(db *sql.DB) *DictionaryRepository {
	return &DictionaryRepository{
		db: db,
	}
}

var allowedTables = map[string]bool{
	"aggregators": true,
	"problems":    true,
}

func validateTable(table string) error {
	if !allowedTables[table] {
		return errors.New("invalid dictionary type")
	}
	return nil
}

// Универсальные методы для работы с любой таблицей словарей
func (r *DictionaryRepository) GetAll(tableName string) ([]dto.BaseDictionary, error) {
	var items []dto.BaseDictionary // убрать *
	if err := validateTable(tableName); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT id, name, created_at, updated_at FROM %s ORDER BY name`, tableName)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item dto.BaseDictionary // убрать *
		var dbCreatedAt time.Time
		var dbUpdatedAt sql.NullTime

		err := rows.Scan(&item.ID, &item.Name, &dbCreatedAt, &dbUpdatedAt)
		if err != nil {
			return nil, err
		}

		item.CreatedAt = dbCreatedAt.Format("2006-01-02 15:04:05")
		if dbUpdatedAt.Valid {
			updatedAtStr := dbUpdatedAt.Time.Format("2006-01-02 15:04:05")
			item.UpdatedAt = &updatedAtStr
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *DictionaryRepository) GetByID(tableName string, id int) (*dto.BaseDictionary, error) {
	var item dto.BaseDictionary // убрать *
	var dbCreatedAt time.Time
	var dbUpdatedAt sql.NullTime
	if err := validateTable(tableName); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT id, name, created_at, updated_at FROM %s WHERE id = ?`, tableName)
	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.Name, &dbCreatedAt, &dbUpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	item.CreatedAt = dbCreatedAt.Format("2006-01-02 15:04:05")
	if dbUpdatedAt.Valid {
		updatedAtStr := dbUpdatedAt.Time.Format("2006-01-02 15:04:05")
		item.UpdatedAt = &updatedAtStr
	}

	return &item, nil // возвращаем указатель
}

func (r *DictionaryRepository) Create(tableName string, item *dto.BaseDictionary) error {
	if err := validateTable(tableName); err != nil {
		return err
	}
	query := fmt.Sprintf(`INSERT INTO %s (name) VALUES (?)`, tableName)
	result, err := r.db.Exec(query, item.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	item.ID = int(id)

	// Получаем created_at из базы
	var createdAt time.Time
	err = r.db.QueryRow("SELECT created_at FROM "+tableName+" WHERE id = ?", item.ID).Scan(&createdAt)
	if err != nil {
		return err
	}
	item.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

	return nil
}

func (r *DictionaryRepository) Update(tableName string, item *dto.BaseDictionary) error {
	if err := validateTable(tableName); err != nil {
		return err
	}
	query := fmt.Sprintf(`UPDATE %s SET name = ?, updated_at = NOW() WHERE id = ?`, tableName)
	_, err := r.db.Exec(query, item.Name, item.ID)
	return err
}

func (r *DictionaryRepository) Delete(tableName string, id int) error {
	if err := validateTable(tableName); err != nil {
		return err
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, tableName)
	_, err := r.db.Exec(query, id)
	return err
}
