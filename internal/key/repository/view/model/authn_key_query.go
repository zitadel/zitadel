package model

import (
	"github.com/caos/zitadel/internal/domain"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type AuthNKeySearchRequest key_model.AuthNKeySearchRequest
type AuthNKeySearchQuery key_model.AuthNKeySearchQuery
type AuthNKeySearchKey key_model.AuthNKeySearchKey

func (req AuthNKeySearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req AuthNKeySearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req AuthNKeySearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == key_model.AuthNKeyKeyUnspecified {
		return nil
	}
	return AuthNKeySearchKey(req.SortingColumn)
}

func (req AuthNKeySearchRequest) GetAsc() bool {
	return req.Asc
}

func (req AuthNKeySearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = AuthNKeySearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req AuthNKeySearchQuery) GetKey() repository.ColumnKey {
	return AuthNKeySearchKey(req.Key)
}

func (req AuthNKeySearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req AuthNKeySearchQuery) GetValue() interface{} {
	return req.Value
}

func (key AuthNKeySearchKey) ToColumnName() string {
	switch key_model.AuthNKeySearchKey(key) {
	case key_model.AuthNKeyKeyID:
		return AuthNKeyKeyID
	case key_model.AuthNKeyObjectID:
		return AuthNKeyObjectID
	case key_model.AuthNKeyObjectType:
		return AuthNKeyObjectType
	default:
		return ""
	}
}
