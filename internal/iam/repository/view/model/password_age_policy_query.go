package model

import (
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type PasswordAgePolicySearchRequest iam_model.PasswordAgePolicySearchRequest
type PasswordAgePolicySearchQuery iam_model.PasswordAgePolicySearchQuery
type PasswordAgePolicySearchKey iam_model.PasswordAgePolicySearchKey

func (req PasswordAgePolicySearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req PasswordAgePolicySearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req PasswordAgePolicySearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.PasswordAgePolicySearchKeyUnspecified {
		return nil
	}
	return PasswordAgePolicySearchKey(req.SortingColumn)
}

func (req PasswordAgePolicySearchRequest) GetAsc() bool {
	return req.Asc
}

func (req PasswordAgePolicySearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = PasswordAgePolicySearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req PasswordAgePolicySearchQuery) GetKey() repository.ColumnKey {
	return PasswordAgePolicySearchKey(req.Key)
}

func (req PasswordAgePolicySearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req PasswordAgePolicySearchQuery) GetValue() interface{} {
	return req.Value
}

func (key PasswordAgePolicySearchKey) ToColumnName() string {
	switch iam_model.PasswordAgePolicySearchKey(key) {
	case iam_model.PasswordAgePolicySearchKeyAggregateID:
		return PasswordAgeKeyAggregateID
	default:
		return ""
	}
}
