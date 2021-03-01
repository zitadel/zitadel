package model

import (
	"github.com/caos/zitadel/internal/domain"
	"time"
)

type ProjectView struct {
	ProjectID            string
	Name                 string
	CreationDate         time.Time
	ChangeDate           time.Time
	State                ProjectState
	ResourceOwner        string
	ProjectRoleAssertion bool
	ProjectRoleCheck     bool
	Sequence             uint64
}

type ProjectViewSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn ProjectViewSearchKey
	Asc           bool
	Queries       []*ProjectViewSearchQuery
}

type ProjectViewSearchKey int32

const (
	ProjectViewSearchKeyUnspecified ProjectViewSearchKey = iota
	ProjectViewSearchKeyName
	ProjectViewSearchKeyProjectID
	ProjectViewSearchKeyResourceOwner
)

type ProjectViewSearchQuery struct {
	Key    ProjectViewSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type ProjectViewSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ProjectView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *ProjectViewSearchRequest) GetSearchQuery(key ProjectViewSearchKey) (int, *ProjectViewSearchQuery) {
	for i, q := range r.Queries {
		if q.Key == key {
			return i, q
		}
	}
	return -1, nil
}

func (r *ProjectViewSearchRequest) AppendMyResourceOwnerQuery(orgID string) {
	r.Queries = append(r.Queries, &ProjectViewSearchQuery{Key: ProjectViewSearchKeyResourceOwner, Method: domain.SearchMethodEquals, Value: orgID})
}

func (r *ProjectViewSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
