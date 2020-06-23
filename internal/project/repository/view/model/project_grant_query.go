package model

import (
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/view"
)

type ProjectGrantSearchRequest proj_model.ProjectGrantViewSearchRequest
type ProjectGrantSearchQuery proj_model.ProjectGrantViewSearchQuery
type ProjectGrantSearchKey proj_model.ProjectGrantViewSearchKey

func (req ProjectGrantSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req ProjectGrantSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req ProjectGrantSearchRequest) GetSortingColumn() view.ColumnKey {
	if req.SortingColumn == proj_model.GRANTEDPROJECTSEARCHKEY_UNSPECIFIED {
		return nil
	}
	return ProjectGrantSearchKey(req.SortingColumn)
}

func (req ProjectGrantSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req ProjectGrantSearchRequest) GetQueries() []view.SearchQuery {
	result := make([]view.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = ProjectGrantSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req ProjectGrantSearchQuery) GetKey() view.ColumnKey {
	return ProjectGrantSearchKey(req.Key)
}

func (req ProjectGrantSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req ProjectGrantSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key ProjectGrantSearchKey) ToColumnName() string {
	switch proj_model.ProjectGrantViewSearchKey(key) {
	case proj_model.GRANTEDPROJECTSEARCHKEY_NAME:
		return ProjectGrantKeyName
	case proj_model.GRANTEDPROJECTSEARCHKEY_GRANTID:
		return ProjectGrantKeyGrantID
	case proj_model.GRANTEDPROJECTSEARCHKEY_ORGID:
		return ProjectGrantKeyOrgID
	case proj_model.GRANTEDPROJECTSEARCHKEY_PROJECTID:
		return ProjectGrantKeyProjectID
	case proj_model.GRANTEDPROJECTSEARCHKEY_RESOURCE_OWNER:
		return ProjectGrantKeyResourceOwner
	case proj_model.GRANTEDPROJECTSEARCHKEY_ROLE_KEYS:
		return ProjectGrantKeyRoleKeys
	default:
		return ""
	}
}
