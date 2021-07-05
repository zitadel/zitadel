package model

import (
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type PrivacyPolicySearchRequest iam_model.PrivacyPolicySearchRequest
type PrivacyPolicySearchQuery iam_model.PrivacyPolicySearchQuery
type PrivacyPolicySearchKey iam_model.PrivacyPolicySearchKey

func (req PrivacyPolicySearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req PrivacyPolicySearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req PrivacyPolicySearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.PrivacyPolicySearchKeyUnspecified {
		return nil
	}
	return PrivacyPolicySearchKey(req.SortingColumn)
}

func (req PrivacyPolicySearchRequest) GetAsc() bool {
	return req.Asc
}

func (req PrivacyPolicySearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = PrivacyPolicySearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req PrivacyPolicySearchQuery) GetKey() repository.ColumnKey {
	return PrivacyPolicySearchKey(req.Key)
}

func (req PrivacyPolicySearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req PrivacyPolicySearchQuery) GetValue() interface{} {
	return req.Value
}

func (key PrivacyPolicySearchKey) ToColumnName() string {
	switch iam_model.PrivacyPolicySearchKey(key) {
	case iam_model.PrivacyPolicySearchKeyAggregateID:
		return PrivacyKeyAggregateID
	default:
		return ""
	}
}
