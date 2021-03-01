package repository

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/view/model"
)

type GeneralSearchRequest model.GeneralSearchRequest
type GeneralSearchQuery model.GeneralSearchQuery
type GeneralSearchKey model.GeneralSearchKey

func (req GeneralSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req GeneralSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req GeneralSearchRequest) GetSortingColumn() ColumnKey {
	if req.SortingColumn == model.GeneralSearchKeyUnspecified {
		return nil
	}
	return GeneralSearchKey(req.SortingColumn)
}

func (req GeneralSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req GeneralSearchRequest) GetQueries() []SearchQuery {
	result := make([]SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = GeneralSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req GeneralSearchQuery) GetKey() ColumnKey {
	return GeneralSearchKey(req.Key)
}

func (req GeneralSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req GeneralSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key GeneralSearchKey) ToColumnName() string {
	return ""
}
