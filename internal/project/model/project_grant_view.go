package model

import (
	"github.com/caos/zitadel/internal/model"
	"time"
)

type ProjectGrantView struct {
	ProjectID         string
	Name              string
	CreationDate      time.Time
	ChangeDate        time.Time
	State             ProjectState
	ResourceOwner     string
	ResourceOwnerName string
	OrgID             string
	OrgName           string
	OrgDomain         string
	Sequence          uint64
	GrantID           string
	GrantedRoleKeys   []string
}

type ProjectGrantViewSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn ProjectGrantViewSearchKey
	Asc           bool
	Queries       []*ProjectGrantViewSearchQuery
}

type ProjectGrantViewSearchKey int32

const (
	GRANTEDPROJECTSEARCHKEY_UNSPECIFIED ProjectGrantViewSearchKey = iota
	GRANTEDPROJECTSEARCHKEY_NAME
	GRANTEDPROJECTSEARCHKEY_PROJECTID
	GRANTEDPROJECTSEARCHKEY_GRANTID
	GRANTEDPROJECTSEARCHKEY_ORGID
	GRANTEDPROJECTSEARCHKEY_RESOURCE_OWNER
	GRANTEDPROJECTSEARCHKEY_ROLE_KEYS
)

type ProjectGrantViewSearchQuery struct {
	Key    ProjectGrantViewSearchKey
	Method model.SearchMethod
	Value  interface{}
}

type ProjectGrantViewSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ProjectGrantView
}

func (r *ProjectGrantViewSearchRequest) AppendMyOrgQuery(orgID string) {
	r.Queries = append(r.Queries, &ProjectGrantViewSearchQuery{Key: GRANTEDPROJECTSEARCHKEY_ORGID, Method: model.SEARCHMETHOD_EQUALS, Value: orgID})
}

func (r *ProjectGrantViewSearchRequest) AppendNotMyOrgQuery(orgID string) {
	r.Queries = append(r.Queries, &ProjectGrantViewSearchQuery{Key: GRANTEDPROJECTSEARCHKEY_ORGID, Method: model.SEARCHMETHOD_NOT_EQUALS, Value: orgID})
}

func (r *ProjectGrantViewSearchRequest) AppendMyResourceOwnerQuery(orgID string) {
	r.Queries = append(r.Queries, &ProjectGrantViewSearchQuery{Key: GRANTEDPROJECTSEARCHKEY_RESOURCE_OWNER, Method: model.SEARCHMETHOD_EQUALS, Value: orgID})
}

func (r *ProjectGrantViewSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
