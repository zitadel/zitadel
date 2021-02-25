package model

import (
	"github.com/caos/zitadel/internal/domain"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type ApplicationSearchRequest proj_model.ApplicationSearchRequest
type ApplicationSearchQuery proj_model.ApplicationSearchQuery
type ApplicationSearchKey proj_model.AppSearchKey

func (req ApplicationSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req ApplicationSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req ApplicationSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == proj_model.AppSearchKeyUnspecified {
		return nil
	}
	return ApplicationSearchKey(req.SortingColumn)
}

func (req ApplicationSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req ApplicationSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = ApplicationSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req ApplicationSearchQuery) GetKey() repository.ColumnKey {
	return ApplicationSearchKey(req.Key)
}

func (req ApplicationSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req ApplicationSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key ApplicationSearchKey) ToColumnName() string {
	switch proj_model.AppSearchKey(key) {
	case proj_model.AppSearchKeyAppID:
		return ApplicationKeyID
	case proj_model.AppSearchKeyName:
		return ApplicationKeyName
	case proj_model.AppSearchKeyProjectID:
		return ApplicationKeyProjectID
	case proj_model.AppSearchKeyOIDCClientID:
		return ApplicationKeyOIDCClientID
	default:
		return ""
	}
}
