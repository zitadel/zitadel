package repository

import (
	"encoding/json"
	"errors"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type JSONArray[T any] []*T

var ErrScanSource = errors.New("unsupported scan source")

func (a *JSONArray[T]) Scan(src any) (err error) {
	var rawJSON []byte
	switch s := src.(type) {
	case string:
		if len(s) == 0 {
			return nil
		}
		rawJSON = []byte(s)
	case []byte:
		if len(s) == 0 {
			return nil
		}
		rawJSON = s
	case nil:
		return nil
	default:
		return ErrScanSource
	}
	err = json.Unmarshal(rawJSON, a)
	if err != nil {
		return database.NewScanError(err)
	}
	return nil
}
