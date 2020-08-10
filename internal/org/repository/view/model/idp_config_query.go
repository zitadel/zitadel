package model

import (
	global_model "github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type IdpConfigSearchRequest org_model.IdpConfigSearchRequest
type IdpConfigSearchQuery org_model.IdpConfigSearchQuery
type IdpConfigSearchKey org_model.IdpConfigSearchKey

func (req IdpConfigSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req IdpConfigSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req IdpConfigSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == org_model.IdpConfigSearchKeyUnspecified {
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
	switch org_model.IdpConfigSearchKey(key) {
	case org_model.IdpConfigSearchKeyIamID:
		return IdpConfigKeyIamID
	case org_model.IdpConfigSearchKeyIdpConfigID:
		return IdpConfigKeyIdpConfigID
	case org_model.IdpConfigSearchKeyName:
		return IdpConfigKeyName
	default:
		return ""
	}
}
