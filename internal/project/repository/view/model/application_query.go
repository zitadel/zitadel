package model

import (
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/view"
)

type ApplicationSearchRequest proj_model.ApplicationSearchRequest
type ApplicationSearchQuery proj_model.ApplicationSearchQuery
type ApplicationSearchKey proj_model.ApplicationSearchKey

func (req ApplicationSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req ApplicationSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req ApplicationSearchRequest) GetSortingColumn() view.ColumnKey {
	if req.SortingColumn == proj_model.APPLICATIONSEARCHKEY_UNSPECIFIED {
		return nil
	}
	return ApplicationSearchKey(req.SortingColumn)
}

func (req ApplicationSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req ApplicationSearchRequest) GetQueries() []view.SearchQuery {
	result := make([]view.SearchQuery, 0)
	for _, q := range req.Queries {
		result = append(result, ApplicationSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method})
	}
	return result
}

func (req ApplicationSearchQuery) GetKey() view.ColumnKey {
	return ApplicationSearchKey(req.Key)
}

func (req ApplicationSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req ApplicationSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key ApplicationSearchKey) ToColumnName() string {
	switch proj_model.ApplicationSearchKey(key) {
	case proj_model.APPLICATIONSEARCHKEY_APP_ID:
		return ApplicationKeyID
	case proj_model.APPLICATIONSEARCHKEY_NAME:
		return ApplicationKeyName
	case proj_model.APPLICATIONSEARCHKEY_PROJECT_ID:
		return ApplicationKeyProjectID
	case proj_model.APPLICATIONSEARCHKEY_OIDC_CLIENT_ID:
		return ApplicationKeyOIDCClientID
	default:
		return ""
	}
}
