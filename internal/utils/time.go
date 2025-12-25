package utils

import (
	"fmt"
	"time"
)

func ParseScheduledAt(input string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		time.RFC3339,
	}

	var t time.Time
	var err error

	for _, f := range formats {
		t, err = time.ParseInLocation(f, input, time.Local)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format: %s", input)
}
