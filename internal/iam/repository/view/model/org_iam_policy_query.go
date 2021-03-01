package model

import (
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type OrgIAMPolicySearchRequest iam_model.OrgIAMPolicySearchRequest
type OrgIAMPolicySearchQuery iam_model.OrgIAMPolicySearchQuery
type OrgIAMPolicySearchKey iam_model.OrgIAMPolicySearchKey

func (req OrgIAMPolicySearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req OrgIAMPolicySearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req OrgIAMPolicySearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.OrgIAMPolicySearchKeyUnspecified {
		return nil
	}
	return OrgIAMPolicySearchKey(req.SortingColumn)
}

func (req OrgIAMPolicySearchRequest) GetAsc() bool {
	return req.Asc
}

func (req OrgIAMPolicySearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = OrgIAMPolicySearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req OrgIAMPolicySearchQuery) GetKey() repository.ColumnKey {
	return OrgIAMPolicySearchKey(req.Key)
}

func (req OrgIAMPolicySearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req OrgIAMPolicySearchQuery) GetValue() interface{} {
	return req.Value
}

func (key OrgIAMPolicySearchKey) ToColumnName() string {
	switch iam_model.OrgIAMPolicySearchKey(key) {
	case iam_model.OrgIAMPolicySearchKeyAggregateID:
		return OrgIAMPolicyKeyAggregateID
	default:
		return ""
	}
}
