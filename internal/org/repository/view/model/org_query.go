package model

import (
	global_model "github.com/caos/zitadel/internal/model"
	usr_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/view"
)

type OrgSearchRequest usr_model.OrgSearchRequest
type OrgSearchQuery usr_model.OrgSearchQuery
type OrgSearchKey usr_model.OrgSearchKey

func (req OrgSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req OrgSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req OrgSearchRequest) GetSortingColumn() view.ColumnKey {
	if req.SortingColumn == usr_model.ORGSEARCHKEY_UNSPECIFIED {
		return nil
	}
	return OrgSearchKey(req.SortingColumn)
}

func (req OrgSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req OrgSearchRequest) GetQueries() []view.SearchQuery {
	result := make([]view.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = OrgSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req OrgSearchQuery) GetKey() view.ColumnKey {
	return OrgSearchKey(req.Key)
}

func (req OrgSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req OrgSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key OrgSearchKey) ToColumnName() string {
	switch usr_model.OrgSearchKey(key) {
	case usr_model.ORGSEARCHKEY_ORG_DOMAIN:
		return OrgKeyOrgDomain
	case usr_model.ORGSEARCHKEY_ORG_ID:
		return OrgKeyOrgID
	case usr_model.ORGSEARCHKEY_ORG_NAME:
		return OrgKeyOrgName
	case usr_model.ORGSEARCHKEY_RESOURCEOWNER:
		return OrgKeyResourceOwner
	case usr_model.ORGSEARCHKEY_STATE:
		return OrgKeyState
	default:
		return ""
	}
}
