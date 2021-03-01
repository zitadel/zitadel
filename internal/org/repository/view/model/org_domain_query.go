package model

import (
	"github.com/caos/zitadel/internal/domain"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type OrgDomainSearchRequest org_model.OrgDomainSearchRequest
type OrgDomainSearchQuery org_model.OrgDomainSearchQuery
type OrgDomainSearchKey org_model.OrgDomainSearchKey

func (req OrgDomainSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req OrgDomainSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req OrgDomainSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == org_model.OrgDomainSearchKeyUnspecified {
		return nil
	}
	return OrgDomainSearchKey(req.SortingColumn)
}

func (req OrgDomainSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req OrgDomainSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = OrgDomainSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req OrgDomainSearchQuery) GetKey() repository.ColumnKey {
	return OrgDomainSearchKey(req.Key)
}

func (req OrgDomainSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req OrgDomainSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key OrgDomainSearchKey) ToColumnName() string {
	switch org_model.OrgDomainSearchKey(key) {
	case org_model.OrgDomainSearchKeyDomain:
		return OrgDomainKeyDomain
	case org_model.OrgDomainSearchKeyOrgID:
		return OrgDomainKeyOrgID
	case org_model.OrgDomainSearchKeyVerified:
		return OrgDomainKeyVerified
	case org_model.OrgDomainSearchKeyPrimary:
		return OrgDomainKeyPrimary
	default:
		return ""
	}
}
