package model

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type IdpConfigSearchRequest iam_model.IdpConfigSearchRequest
type IdpConfigSearchQuery iam_model.IdpConfigSearchQuery
type IdpConfigSearchKey iam_model.IdpConfigSearchKey

func (req IdpConfigSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req IdpConfigSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req IdpConfigSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.IdpConfigSearchKeyUnspecified {
		return nil
	}
	return IdpConfigSearchKey(req.SortingColumn)
}

func (req IdpConfigSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req IdpConfigSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = IdpConfigSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req IdpConfigSearchQuery) GetKey() repository.ColumnKey {
	return IdpConfigSearchKey(req.Key)
}

func (req IdpConfigSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req IdpConfigSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key IdpConfigSearchKey) ToColumnName() string {
	switch iam_model.IdpConfigSearchKey(key) {
	case iam_model.IdpConfigSearchKeyAggregateID:
		return IdpConfigKeyAggregateID
	case iam_model.IdpConfigSearchKeyIdpConfigID:
		return IdpConfigKeyIdpConfigID
	case iam_model.IdpConfigSearchKeyName:
		return IdpConfigKeyName
	case iam_model.IdpConfigSearchKeyIdpProviderType:
		return IdpConfigKeyProviderType
	default:
		return ""
	}
}
