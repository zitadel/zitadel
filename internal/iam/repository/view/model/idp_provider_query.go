package model

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type IdpProviderSearchRequest iam_model.IdpProviderSearchRequest
type IdpProviderSearchQuery iam_model.IdpProviderSearchQuery
type IdpProviderSearchKey iam_model.IdpProviderSearchKey

func (req IdpProviderSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req IdpProviderSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req IdpProviderSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.IdpProviderSearchKeyUnspecified {
		return nil
	}
	return IdpProviderSearchKey(req.SortingColumn)
}

func (req IdpProviderSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req IdpProviderSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = IdpProviderSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req IdpProviderSearchQuery) GetKey() repository.ColumnKey {
	return IdpProviderSearchKey(req.Key)
}

func (req IdpProviderSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req IdpProviderSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key IdpProviderSearchKey) ToColumnName() string {
	switch iam_model.IdpProviderSearchKey(key) {
	case iam_model.IdpProviderSearchKeyAggregateID:
		return IdpProviderKeyAggregateID
	case iam_model.IdpProviderSearchKeyIdpConfigID:
		return IdpProviderKeyIdpConfigID
	default:
		return ""
	}
}
