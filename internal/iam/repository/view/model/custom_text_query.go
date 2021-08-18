package model

import (
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type CustomTextSearchRequest iam_model.CustomTextSearchRequest
type CustomTextSearchQuery iam_model.CustomTextSearchQuery
type CustomTextSearchKey iam_model.CustomTextSearchKey

func (req CustomTextSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req CustomTextSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req CustomTextSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.CustomTextSearchKeyUnspecified {
		return nil
	}
	return CustomTextSearchKey(req.SortingColumn)
}

func (req CustomTextSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req CustomTextSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = CustomTextSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req CustomTextSearchQuery) GetKey() repository.ColumnKey {
	return CustomTextSearchKey(req.Key)
}

func (req CustomTextSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req CustomTextSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key CustomTextSearchKey) ToColumnName() string {
	switch iam_model.CustomTextSearchKey(key) {
	case iam_model.CustomTextSearchKeyAggregateID:
		return CustomTextKeyAggregateID
	case iam_model.CustomTextSearchKeyTemplate:
		return CustomTextKeyTemplate
	case iam_model.CustomTextSearchKeyLanguage:
		return CustomTextKeyLanguage
	case iam_model.CustomTextSearchKeyKey:
		return CustomTextKeyKey
	default:
		return ""
	}
}
