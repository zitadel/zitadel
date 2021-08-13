package model

import (
	"github.com/caos/zitadel/internal/domain"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type OrgProjectMappingSearchRequest proj_model.OrgProjectMappingViewSearchRequest
type OrgProjectMappingSearchQuery proj_model.OrgProjectMappingViewSearchQuery
type OrgProjectMappingSearchKey proj_model.OrgProjectMappingViewSearchKey

func (req OrgProjectMappingSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req OrgProjectMappingSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req OrgProjectMappingSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == proj_model.OrgProjectMappingSearchKeyUnspecified {
		return nil
	}
	return OrgProjectMappingSearchKey(req.SortingColumn)
}

func (req OrgProjectMappingSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req OrgProjectMappingSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = OrgProjectMappingSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req OrgProjectMappingSearchQuery) GetKey() repository.ColumnKey {
	return OrgProjectMappingSearchKey(req.Key)
}

func (req OrgProjectMappingSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req OrgProjectMappingSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key OrgProjectMappingSearchKey) ToColumnName() string {
	switch proj_model.OrgProjectMappingViewSearchKey(key) {
	case proj_model.OrgProjectMappingSearchKeyOrgID:
		return OrgProjectMappingKeyOrgID
	case proj_model.OrgProjectMappingSearchKeyProjectID:
		return OrgProjectMappingKeyProjectID
	default:
		return ""
	}
}
