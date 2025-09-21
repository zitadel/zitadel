package repository

import (
	"encoding/json"
	"errors"
)

type JSONArray[T any] []*T

var ScanSourceErr = errors.New("unsupported scan source")

func (a *JSONArray[T]) Scan(src any) error {

	switch s := src.(type) {
	case string:
		if len(s) == 0 {
			return nil
		}
		return json.Unmarshal([]byte(s), a)
	case []byte:
		if len(s) == 0 {
			return nil
		}
		return json.Unmarshal(s, a)
	case nil:
		return nil
	default:
		return ScanSourceErr
	}
}
