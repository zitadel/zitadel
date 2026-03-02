package repository

import (
	"encoding/json"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type JSON[T any] struct {
	Value T
}

func (j *JSON[T]) Scan(src any) (err error) {
	var rawJSON []byte
	switch s := src.(type) {
	case string:
		rawJSON = []byte(s)
	case []byte:
		rawJSON = s
	case nil:
		return nil
	default:
		return ErrScanSource
	}
	if err = j.UnmarshalJSON(rawJSON); err != nil {
		return database.NewScanError(err)
	}
	return nil
}

// UnmarshalJSON allows encoding/json to decode nested JSON into JSON[T] fields.
// This makes [json.Unmarshal] work for fields inside aggregated JSON (e.g. JSON arrays).
func (j *JSON[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	j.Value = v
	return nil
}

// MarshalJSON encodes the wrapped value or null if nil.
func (j JSON[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Value)
}
