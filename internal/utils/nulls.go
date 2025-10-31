package utils

import "database/sql"

// Int64ToNullInt64 — конвертирует int64 → sql.NullInt64
func Int64ToNullInt64(value int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: value,
		Valid: true,
	}
}
