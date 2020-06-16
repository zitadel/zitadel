package model

import (
	global_model "github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/view"
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

func (req OrgDomainSearchRequest) GetSortingColumn() view.ColumnKey {
	if req.SortingColumn == org_model.ORGDOMAINSEARCHKEY_UNSPECIFIED {
		return nil
	}
	return OrgDomainSearchKey(req.SortingColumn)
}

func (req OrgDomainSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req OrgDomainSearchRequest) GetQueries() []view.SearchQuery {
	result := make([]view.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = OrgDomainSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req OrgDomainSearchQuery) GetKey() view.ColumnKey {
	return OrgDomainSearchKey(req.Key)
}

func (req OrgDomainSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req OrgDomainSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key OrgDomainSearchKey) ToColumnName() string {
	switch org_model.OrgDomainSearchKey(key) {
	case org_model.ORGDOMAINSEARCHKEY_DOMAIN:
		return OrgDomainKeyDomain
	case org_model.ORGDOMAINSEARCHKEY_ORG_ID:
		return OrgDomainKeyOrgID
	case org_model.ORGDOMAINSEARCHKEY_VERIFIED:
		return OrgDomainKeyVerified
	case org_model.ORGDOMAINSEARCHKEY_PRIMARY:
		return OrgDomainKeyPrimary
	default:
		return ""
	}
}
