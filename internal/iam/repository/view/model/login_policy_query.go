package model

import (
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type LoginPolicySearchRequest iam_model.LoginPolicySearchRequest
type LoginPolicySearchQuery iam_model.LoginPolicySearchQuery
type LoginPolicySearchKey iam_model.LoginPolicySearchKey

func (req LoginPolicySearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req LoginPolicySearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req LoginPolicySearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.LoginPolicySearchKeyUnspecified {
		return nil
	}
	return LoginPolicySearchKey(req.SortingColumn)
}

func (req LoginPolicySearchRequest) GetAsc() bool {
	return req.Asc
}

func (req LoginPolicySearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = LoginPolicySearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req LoginPolicySearchQuery) GetKey() repository.ColumnKey {
	return LoginPolicySearchKey(req.Key)
}

func (req LoginPolicySearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req LoginPolicySearchQuery) GetValue() interface{} {
	return req.Value
}

func (key LoginPolicySearchKey) ToColumnName() string {
	switch iam_model.LoginPolicySearchKey(key) {
	case iam_model.LoginPolicySearchKeyAggregateID:
		return LoginPolicyKeyAggregateID
	case iam_model.LoginPolicySearchKeyDefault:
		return LoginPolicyKeyDefault
	default:
		return ""
	}
}
