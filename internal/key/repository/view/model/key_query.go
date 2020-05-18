package model

import (
	key_model "github.com/caos/zitadel/internal/key/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view"
)

type KeySearchRequest key_model.KeySearchRequest
type KeySearchQuery key_model.KeySearchQuery
type KeySearchKey key_model.KeySearchKey

func (req KeySearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req KeySearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req KeySearchRequest) GetSortingColumn() view.ColumnKey {
	if req.SortingColumn == key_model.KEYSEARCHKEY_UNSPECIFIED {
		return nil
	}
	return KeySearchKey(req.SortingColumn)
}

func (req KeySearchRequest) GetAsc() bool {
	return req.Asc
}

func (req KeySearchRequest) GetQueries() []view.SearchQuery {
	result := make([]view.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = KeySearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req KeySearchQuery) GetKey() view.ColumnKey {
	return KeySearchKey(req.Key)
}

func (req KeySearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req KeySearchQuery) GetValue() interface{} {
	return req.Value
}

func (key KeySearchKey) ToColumnName() string {
	switch key_model.KeySearchKey(key) {
	case key_model.KEYSEARCHKEY_APP_ID:
		return KeyKeyID
	case key_model.KEYSEARCHKEY_PRIVATE:
		return KeyPrivate
	case key_model.KEYSEARCHKEY_USAGE:
		return KeyUsage
	case key_model.KEYSEARCHKEY_EXPIRY:
		return KeyExpiry
	default:
		return ""
	}
}
