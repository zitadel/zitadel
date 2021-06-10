package model

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"

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
	GrantedProjectSearchKeyUnspecified ProjectGrantViewSearchKey = iota
	GrantedProjectSearchKeyName
	GrantedProjectSearchKeyProjectID
	GrantedProjectSearchKeyGrantID
	GrantedProjectSearchKeyOrgID
	GrantedProjectSearchKeyResourceOwner
	GrantedProjectSearchKeyRoleKeys
)

type ProjectGrantViewSearchQuery struct {
	Key    ProjectGrantViewSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type ProjectGrantViewSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ProjectGrantView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *ProjectGrantViewSearchRequest) GetSearchQuery(key ProjectGrantViewSearchKey) (int, *ProjectGrantViewSearchQuery) {
	for i, q := range r.Queries {
		if q.Key == key {
			return i, q
		}
	}
	return -1, nil
}

func (r *ProjectGrantViewSearchRequest) AppendMyOrgQuery(orgID string) {
	r.Queries = append(r.Queries, &ProjectGrantViewSearchQuery{Key: GrantedProjectSearchKeyOrgID, Method: domain.SearchMethodEquals, Value: orgID})
}

func (r *ProjectGrantViewSearchRequest) AppendNotMyOrgQuery(orgID string) {
	r.Queries = append(r.Queries, &ProjectGrantViewSearchQuery{Key: GrantedProjectSearchKeyOrgID, Method: domain.SearchMethodNotEquals, Value: orgID})
}

func (r *ProjectGrantViewSearchRequest) AppendMyResourceOwnerQuery(orgID string) {
	r.Queries = append(r.Queries, &ProjectGrantViewSearchQuery{Key: GrantedProjectSearchKeyResourceOwner, Method: domain.SearchMethodEquals, Value: orgID})
}

func (r *ProjectGrantViewSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-0fj3s", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}
