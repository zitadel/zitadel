package model

import (
	"github.com/caos/zitadel/internal/domain"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/view/repository"
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

func (req ProjectMemberSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == proj_model.ProjectMemberSearchKeyUnspecified {
		return nil
	}
	return ProjectMemberSearchKey(req.SortingColumn)
}

func (req ProjectMemberSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req ProjectMemberSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = ProjectMemberSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req ProjectMemberSearchQuery) GetKey() repository.ColumnKey {
	return ProjectMemberSearchKey(req.Key)
}

func (req ProjectMemberSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req ProjectMemberSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key ProjectMemberSearchKey) ToColumnName() string {
	switch proj_model.ProjectMemberSearchKey(key) {
	case proj_model.ProjectMemberSearchKeyEmail:
		return ProjectMemberKeyEmail
	case proj_model.ProjectMemberSearchKeyFirstName:
		return ProjectMemberKeyFirstName
	case proj_model.ProjectMemberSearchKeyLastName:
		return ProjectMemberKeyLastName
	case proj_model.ProjectMemberSearchKeyUserName:
		return ProjectMemberKeyUserName
	case proj_model.ProjectMemberSearchKeyUserID:
		return ProjectMemberKeyUserID
	case proj_model.ProjectMemberSearchKeyProjectID:
		return ProjectMemberKeyProjectID
	default:
		return ""
	}
}
