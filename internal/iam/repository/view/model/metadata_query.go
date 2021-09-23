package model

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/view/repository"
)

type MetadataSearchRequest domain.MetadataSearchRequest
type MetadataSearchQuery domain.MetadataSearchQuery
type MetadataSearchKey domain.MetadataSearchKey

func (req MetadataSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req MetadataSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req MetadataSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == domain.MetadataSearchKeyUnspecified {
		return nil
	}
	return MetadataSearchKey(req.SortingColumn)
}

func (req MetadataSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req MetadataSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = MetadataSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req MetadataSearchQuery) GetKey() repository.ColumnKey {
	return MetadataSearchKey(req.Key)
}

func (req MetadataSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req MetadataSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key MetadataSearchKey) ToColumnName() string {
	switch domain.MetadataSearchKey(key) {
	case domain.MetadataSearchKeyAggregateID:
		return MetadataKeyAggregateID
	case domain.MetadataSearchKeyResourceOwner:
		return MetadataKeyResourceOwner
	case domain.MetadataSearchKeyKey:
		return MetadataKeyKey
	case domain.MetadataSearchKeyValue:
		return MetadataKeyValue
	default:
		return ""
	}
}
