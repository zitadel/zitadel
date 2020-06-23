package model

import (
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/view"
)

type ProjectSearchRequest proj_model.ProjectViewSearchRequest
type ProjectSearchQuery proj_model.ProjectViewSearchQuery
type ProjectSearchKey proj_model.ProjectViewSearchKey

func (req ProjectSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req ProjectSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req ProjectSearchRequest) GetSortingColumn() view.ColumnKey {
	if req.SortingColumn == proj_model.ProjectViewSearchKeyUnspecified {
		return nil
	}
	return ProjectSearchKey(req.SortingColumn)
}

func (req ProjectSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req ProjectSearchRequest) GetQueries() []view.SearchQuery {
	result := make([]view.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = ProjectSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req ProjectSearchQuery) GetKey() view.ColumnKey {
	return ProjectSearchKey(req.Key)
}

func (req ProjectSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req ProjectSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key ProjectSearchKey) ToColumnName() string {
	switch proj_model.ProjectViewSearchKey(key) {
	case proj_model.ProjectViewSearchKeyName:
		return ProjectKeyName
	case proj_model.ProjectViewSearchKeyProjectID:
		return ProjectKeyProjectID
	case proj_model.ProjectViewSearchKeyResourceOwner:
		return ProjectKeyResourceOwner
	default:
		return ""
	}
}
