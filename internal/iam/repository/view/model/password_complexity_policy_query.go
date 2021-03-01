package model

import (
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type PasswordComplexityPolicySearchRequest iam_model.PasswordComplexityPolicySearchRequest
type PasswordComplexityPolicySearchQuery iam_model.PasswordComplexityPolicySearchQuery
type PasswordComplexityPolicySearchKey iam_model.PasswordComplexityPolicySearchKey

func (req PasswordComplexityPolicySearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req PasswordComplexityPolicySearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req PasswordComplexityPolicySearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.PasswordComplexityPolicySearchKeyUnspecified {
		return nil
	}
	return PasswordComplexityPolicySearchKey(req.SortingColumn)
}

func (req PasswordComplexityPolicySearchRequest) GetAsc() bool {
	return req.Asc
}

func (req PasswordComplexityPolicySearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = PasswordComplexityPolicySearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req PasswordComplexityPolicySearchQuery) GetKey() repository.ColumnKey {
	return PasswordComplexityPolicySearchKey(req.Key)
}

func (req PasswordComplexityPolicySearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req PasswordComplexityPolicySearchQuery) GetValue() interface{} {
	return req.Value
}

func (key PasswordComplexityPolicySearchKey) ToColumnName() string {
	switch iam_model.PasswordComplexityPolicySearchKey(key) {
	case iam_model.PasswordComplexityPolicySearchKeyAggregateID:
		return PasswordComplexityKeyAggregateID
	default:
		return ""
	}
}
