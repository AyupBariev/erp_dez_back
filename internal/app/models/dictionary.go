package models

// BaseDictionary - базовая модель словаря
type BaseDictionary struct {
	ID        int     `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	UpdatedAt *string `json:"updated_at" db:"updated_at"`
}

// Request/Response модели
type CreateDictionaryRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateDictionaryRequest struct {
	Name string `json:"name" binding:"required"`
}
