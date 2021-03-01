package model

import (
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type PasswordLockoutPolicySearchRequest iam_model.PasswordLockoutPolicySearchRequest
type PasswordLockoutPolicySearchQuery iam_model.PasswordLockoutPolicySearchQuery
type PasswordLockoutPolicySearchKey iam_model.PasswordLockoutPolicySearchKey

func (req PasswordLockoutPolicySearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req PasswordLockoutPolicySearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req PasswordLockoutPolicySearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.PasswordLockoutPolicySearchKeyUnspecified {
		return nil
	}
	return PasswordLockoutPolicySearchKey(req.SortingColumn)
}

func (req PasswordLockoutPolicySearchRequest) GetAsc() bool {
	return req.Asc
}

func (req PasswordLockoutPolicySearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = PasswordLockoutPolicySearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req PasswordLockoutPolicySearchQuery) GetKey() repository.ColumnKey {
	return PasswordLockoutPolicySearchKey(req.Key)
}

func (req PasswordLockoutPolicySearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req PasswordLockoutPolicySearchQuery) GetValue() interface{} {
	return req.Value
}

func (key PasswordLockoutPolicySearchKey) ToColumnName() string {
	switch iam_model.PasswordLockoutPolicySearchKey(key) {
	case iam_model.PasswordLockoutPolicySearchKeyAggregateID:
		return PasswordLockoutKeyAggregateID
	default:
		return ""
	}
}
