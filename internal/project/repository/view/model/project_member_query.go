package model

import (
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/view"
)

type ProjectMemberSearchRequest proj_model.ProjectMemberSearchRequest
type ProjectMemberSearchQuery proj_model.ProjectMemberSearchQuery
type ProjectMemberSearchKey proj_model.ProjectMemberSearchKey

func (req ProjectMemberSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req ProjectMemberSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req ProjectMemberSearchRequest) GetSortingColumn() view.ColumnKey {
	if req.SortingColumn == proj_model.PROJECTMEMBERSEARCHKEY_UNSPECIFIED {
		return nil
	}
	return ProjectMemberSearchKey(req.SortingColumn)
}

func (req ProjectMemberSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req ProjectMemberSearchRequest) GetQueries() []view.SearchQuery {
	result := make([]view.SearchQuery, 0)
	for _, q := range req.Queries {
		result = append(result, ProjectMemberSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method})
	}
	return result
}

func (req ProjectMemberSearchQuery) GetKey() view.ColumnKey {
	return ProjectMemberSearchKey(req.Key)
}

func (req ProjectMemberSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req ProjectMemberSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key ProjectMemberSearchKey) ToColumnName() string {
	switch proj_model.ProjectMemberSearchKey(key) {
	case proj_model.PROJECTMEMBERSEARCHKEY_EMAIL:
		return ProjectMemberKeyEmail
	case proj_model.PROJECTMEMBERSEARCHKEY_FIRST_NAME:
		return ProjectMemberKeyFirstName
	case proj_model.PROJECTMEMBERSEARCHKEY_LAST_NAME:
		return ProjectMemberKeyLastName
	case proj_model.PROJECTMEMBERSEARCHKEY_USER_NAME:
		return ProjectMemberKeyUserName
	case proj_model.PROJECTMEMBERSEARCHKEY_USER_ID:
		return ProjectMemberKeyUserID
	case proj_model.PROJECTMEMBERSEARCHKEY_PROJECT_ID:
		return ProjectMemberKeyProjectID
	default:
		return ""
	}
}
