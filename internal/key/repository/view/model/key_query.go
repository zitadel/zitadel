package model

import (
	"github.com/caos/zitadel/internal/domain"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/view/repository"
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

func (req KeySearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == key_model.KeySearchKeyUnspecified {
		return nil
	}
	return KeySearchKey(req.SortingColumn)
}

func (req KeySearchRequest) GetAsc() bool {
	return req.Asc
}

func (req KeySearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = KeySearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req KeySearchQuery) GetKey() repository.ColumnKey {
	return KeySearchKey(req.Key)
}

func (req KeySearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req KeySearchQuery) GetValue() interface{} {
	return req.Value
}

func (key KeySearchKey) ToColumnName() string {
	switch key_model.KeySearchKey(key) {
	case key_model.KeySearchKeyID:
		return KeyKeyID
	case key_model.KeySearchKeyPrivate:
		return KeyPrivate
	case key_model.KeySearchKeyUsage:
		return KeyUsage
	case key_model.KeySearchKeyExpiry:
		return KeyExpiry
	default:
		return ""
	}
}
