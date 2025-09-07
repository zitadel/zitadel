package repository

import (
	"encoding/json"
	"errors"
)

type JSONArray[T any] []*T

func (a JSONArray[T]) Scan(src any) error {
	switch s := src.(type) {
	case string:
		return json.Unmarshal([]byte(s), &a)
	case []byte:
		return json.Unmarshal(s, &a)
	case nil:
		return nil
	default:
		return errors.New("unsupported scan source")
	}
}
