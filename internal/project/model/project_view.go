package model

import (
	"github.com/caos/zitadel/internal/model"
	"time"
)

type ProjectView struct {
	ProjectID     string
	Name          string
	CreationDate  time.Time
	ChangeDate    time.Time
	State         ProjectState
	ResourceOwner string
	Sequence      uint64
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
	PROJECTSEARCHKEY_UNSPECIFIED ProjectViewSearchKey = iota
	PROJECTSEARCHKEY_NAME
	PROJECTSEARCHKEY_PROJECTID
	PROJECTSEARCHKEY_RESOURCE_OWNER
)

type ProjectViewSearchQuery struct {
	Key    ProjectViewSearchKey
	Method model.SearchMethod
	Value  interface{}
}

type ProjectViewSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ProjectView
}

func (r *ProjectViewSearchRequest) AppendMyResourceOwnerQuery(orgID string) {
	r.Queries = append(r.Queries, &ProjectViewSearchQuery{Key: PROJECTSEARCHKEY_RESOURCE_OWNER, Method: model.SEARCHMETHOD_EQUALS, Value: orgID})
}

func (r *ProjectViewSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
