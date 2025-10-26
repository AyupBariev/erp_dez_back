package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// MapToStruct маппит map[field]value в структуру по тегам `db:"field_name"`
func MapToStruct(data map[string]interface{}, out interface{}) error {
	v := reflect.ValueOf(out).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			continue
		}

		val, ok := data[dbTag]
		if !ok {
			continue
		}

		f := v.Field(i)

		switch f.Kind() {
		case reflect.String:
			switch v := val.(type) {
			case string:
				f.SetString(v)
			case []uint8:
				f.SetString(string(v))
			}
		case reflect.Int, reflect.Int64:
			switch v := val.(type) {
			case int64:
				f.SetInt(v)
			case int:
				f.SetInt(int64(v))
			}
		case reflect.Struct:
			switch f.Type() {
			case reflect.TypeOf(sql.NullInt64{}):
				if n, ok := val.(int64); ok {
					f.Set(reflect.ValueOf(sql.NullInt64{Int64: n, Valid: true}))
				}
			case reflect.TypeOf(time.Time{}):
				if t, ok := val.(time.Time); ok {
					f.Set(reflect.ValueOf(t))
				}
			}
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.String {
				switch v := val.(type) {
				case []byte:
					var arr []string
					_ = json.Unmarshal(v, &arr)
					f.Set(reflect.ValueOf(arr))
				case string:
					var arr []string
					_ = json.Unmarshal([]byte(v), &arr)
					f.Set(reflect.ValueOf(arr))
				}
			}
		}
	}
	return nil
}

func ToNullString(v interface{}) sql.NullString {
	if v == nil {
		return sql.NullString{Valid: false}
	}
	s := fmt.Sprintf("%v", v)
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
