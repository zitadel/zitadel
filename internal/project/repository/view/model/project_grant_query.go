package model

import (
	"github.com/caos/zitadel/internal/domain"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/view/repository"
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

func (req ProjectGrantSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == proj_model.GrantedProjectSearchKeyUnspecified {
		return nil
	}
	return ProjectGrantSearchKey(req.SortingColumn)
}

func (req ProjectGrantSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req ProjectGrantSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = ProjectGrantSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req ProjectGrantSearchQuery) GetKey() repository.ColumnKey {
	return ProjectGrantSearchKey(req.Key)
}

func (req ProjectGrantSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req ProjectGrantSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key ProjectGrantSearchKey) ToColumnName() string {
	switch proj_model.ProjectGrantViewSearchKey(key) {
	case proj_model.GrantedProjectSearchKeyName:
		return ProjectGrantKeyName
	case proj_model.GrantedProjectSearchKeyGrantID:
		return ProjectGrantKeyGrantID
	case proj_model.GrantedProjectSearchKeyOrgID:
		return ProjectGrantKeyOrgID
	case proj_model.GrantedProjectSearchKeyProjectID:
		return ProjectGrantKeyProjectID
	case proj_model.GrantedProjectSearchKeyResourceOwner:
		return ProjectGrantKeyResourceOwner
	case proj_model.GrantedProjectSearchKeyRoleKeys:
		return ProjectGrantKeyRoleKeys
	default:
		return ""
	}
}
