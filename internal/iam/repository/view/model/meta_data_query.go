package model

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/view/repository"
)

type MetaDataSearchRequest domain.MetaDataSearchRequest
type MetaDataSearchQuery domain.MetaDataSearchQuery
type MetaDataSearchKey domain.MetaDataSearchKey

func (req MetaDataSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req MetaDataSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req MetaDataSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == domain.MetaDataSearchKeyUnspecified {
		return nil
	}
	return MetaDataSearchKey(req.SortingColumn)
}

func (req MetaDataSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req MetaDataSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = MetaDataSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req MetaDataSearchQuery) GetKey() repository.ColumnKey {
	return MetaDataSearchKey(req.Key)
}

func (req MetaDataSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req MetaDataSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key MetaDataSearchKey) ToColumnName() string {
	switch domain.MetaDataSearchKey(key) {
	case domain.MetaDataSearchKeyAggregateID:
		return MetaDataKeyAggregateID
	case domain.MetaDataSearchKeyResourceOwner:
		return MetaDataKeyResourceOwner
	case domain.MetaDataSearchKeyKey:
		return MetaDataKeyKey
	case domain.MetaDataSearchKeyValue:
		return MetaDataKeyValue
	default:
		return ""
	}
}
