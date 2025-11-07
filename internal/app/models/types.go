package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// StringArray — тип для хранения []string в JSON-поле
type StringArray []string

// Value — преобразует Go → SQL
func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan — преобразует SQL → Go
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringArray: %v", value)
	}

	if err := json.Unmarshal(bytes, a); err != nil {
		return err
	}

	return nil
}
