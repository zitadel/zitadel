package model

import (
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type LockoutPolicySearchRequest iam_model.LockoutPolicySearchRequest
type LockoutPolicySearchQuery iam_model.LockoutPolicySearchQuery
type LockoutPolicySearchKey iam_model.LockoutPolicySearchKey

func (req LockoutPolicySearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req LockoutPolicySearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req LockoutPolicySearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.LockoutPolicySearchKeyUnspecified {
		return nil
	}
	return LockoutPolicySearchKey(req.SortingColumn)
}

func (req LockoutPolicySearchRequest) GetAsc() bool {
	return req.Asc
}

func (req LockoutPolicySearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = LockoutPolicySearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req LockoutPolicySearchQuery) GetKey() repository.ColumnKey {
	return LockoutPolicySearchKey(req.Key)
}

func (req LockoutPolicySearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req LockoutPolicySearchQuery) GetValue() interface{} {
	return req.Value
}

func (key LockoutPolicySearchKey) ToColumnName() string {
	switch iam_model.LockoutPolicySearchKey(key) {
	case iam_model.LockoutPolicySearchKeyAggregateID:
		return LockoutKeyAggregateID
	default:
		return ""
	}
}
