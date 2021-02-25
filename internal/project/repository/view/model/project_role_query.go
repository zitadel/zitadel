package model

import (
	"github.com/caos/zitadel/internal/domain"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/view/repository"
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

func (req ProjectRoleSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == proj_model.ProjectRoleSearchKeyUnspecified {
		return nil
	}
	return ProjectRoleSearchKey(req.SortingColumn)
}

func (req ProjectRoleSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req ProjectRoleSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = ProjectRoleSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req ProjectRoleSearchQuery) GetKey() repository.ColumnKey {
	return ProjectRoleSearchKey(req.Key)
}

func (req ProjectRoleSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req ProjectRoleSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key ProjectRoleSearchKey) ToColumnName() string {
	switch proj_model.ProjectRoleSearchKey(key) {
	case proj_model.ProjectRoleSearchKeyKey:
		return ProjectRoleKeyKey
	case proj_model.ProjectRoleSearchKeyOrgID:
		return ProjectRoleKeyOrgID
	case proj_model.ProjectRoleSearchKeyProjectID:
		return ProjectRoleKeyProjectID
	case proj_model.ProjectRoleSearchKeyResourceOwner:
		return ProjectRoleKeyResourceOwner
	default:
		return ""
	}
}
