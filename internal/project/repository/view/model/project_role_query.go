package model

import (
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/view"
)

type ProjectRoleSearchRequest proj_model.ProjectRoleSearchRequest
type ProjectRoleSearchQuery proj_model.ProjectRoleSearchQuery
type ProjectRoleSearchKey proj_model.ProjectRoleSearchKey

func (req ProjectRoleSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req ProjectRoleSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req ProjectRoleSearchRequest) GetSortingColumn() view.ColumnKey {
	if req.SortingColumn == proj_model.PROJECTROLESEARCHKEY_UNSPECIFIED {
		return nil
	}
	return ProjectRoleSearchKey(req.SortingColumn)
}

func (req ProjectRoleSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req ProjectRoleSearchRequest) GetQueries() []view.SearchQuery {
	result := make([]view.SearchQuery, 0)
	for _, q := range req.Queries {
		result = append(result, ProjectRoleSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method})
	}
	return result
}

func (req ProjectRoleSearchQuery) GetKey() view.ColumnKey {
	return ProjectRoleSearchKey(req.Key)
}

func (req ProjectRoleSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req ProjectRoleSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key ProjectRoleSearchKey) ToColumnName() string {
	switch proj_model.ProjectRoleSearchKey(key) {
	case proj_model.PROJECTROLESEARCHKEY_KEY:
		return ProjectRoleKeyKey
	case proj_model.PROJECTROLESEARCHKEY_ORGID:
		return ProjectRoleKeyOrgID
	case proj_model.PROJECTROLESEARCHKEY_PROJECTID:
		return ProjectRoleKeyProjectID
	case proj_model.PROJECTROLESEARCHKEY_RESOURCEOWNER:
		return ProjectRoleKeyResourceOwner
	default:
		return ""
	}
}
