package model

import (
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/view"
)

type ProjectGrantMemberSearchRequest proj_model.ProjectGrantMemberSearchRequest
type ProjectGrantMemberSearchQuery proj_model.ProjectGrantMemberSearchQuery
type ProjectGrantMemberSearchKey proj_model.ProjectGrantMemberSearchKey

func (req ProjectGrantMemberSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req ProjectGrantMemberSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req ProjectGrantMemberSearchRequest) GetSortingColumn() view.ColumnKey {
	if req.SortingColumn == proj_model.ProjectGrantMemberSearchKeyUnspecified {
		return nil
	}
	return ProjectGrantMemberSearchKey(req.SortingColumn)
}

func (req ProjectGrantMemberSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req ProjectGrantMemberSearchRequest) GetQueries() []view.SearchQuery {
	result := make([]view.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = ProjectGrantMemberSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req ProjectGrantMemberSearchQuery) GetKey() view.ColumnKey {
	return ProjectGrantMemberSearchKey(req.Key)
}

func (req ProjectGrantMemberSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req ProjectGrantMemberSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key ProjectGrantMemberSearchKey) ToColumnName() string {
	switch proj_model.ProjectGrantMemberSearchKey(key) {
	case proj_model.ProjectGrantMemberSearchKeyEmail:
		return ProjectGrantMemberKeyEmail
	case proj_model.ProjectGrantMemberSearchKeyFirstName:
		return ProjectGrantMemberKeyFirstName
	case proj_model.ProjectGrantMemberSearchKeyLastName:
		return ProjectGrantMemberKeyLastName
	case proj_model.ProjectGrantMemberSearchKeyUserName:
		return ProjectGrantMemberKeyUserName
	case proj_model.ProjectGrantMemberSearchKeyUserID:
		return ProjectGrantMemberKeyUserID
	case proj_model.ProjectGrantMemberSearchKeyGrantID:
		return ProjectGrantMemberKeyGrantID
	default:
		return ""
	}
}
