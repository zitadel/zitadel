package model

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/features/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type FeaturesSearchRequest model.FeaturesSearchRequest
type FeaturesSearchQuery model.FeaturesSearchQuery
type FeaturesSearchKey model.FeaturesSearchKey

func (req FeaturesSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req FeaturesSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req FeaturesSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == model.FeaturesSearchKeyUnspecified {
		return nil
	}
	return FeaturesSearchKey(req.SortingColumn)
}

func (req FeaturesSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req FeaturesSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = FeaturesSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req FeaturesSearchQuery) GetKey() repository.ColumnKey {
	return FeaturesSearchKey(req.Key)
}

func (req FeaturesSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req FeaturesSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key FeaturesSearchKey) ToColumnName() string {
	switch model.FeaturesSearchKey(key) {
	case model.FeaturesSearchKeyAggregateID:
		return FeaturesKeyAggregateID
	case model.FeaturesSearchKeyDefault:
		return FeaturesKeyDefault
	default:
		return ""
	}
}
