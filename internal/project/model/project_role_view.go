package model

import (
	"time"

	"github.com/caos/zitadel/internal/model"
)

type ProjectRoleView struct {
	ResourceOwner string
	OrgID         string
	ProjectID     string
	Key           string
	DisplayName   string
	Group         string
	CreationDate  time.Time
	ChangeDate    time.Time
	Sequence      uint64
}

type ProjectRoleSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn ProjectRoleSearchKey
	Asc           bool
	Queries       []*ProjectRoleSearchQuery
}

type ProjectRoleSearchKey int32

const (
	ProjectRoleSearchKeyUnspecified ProjectRoleSearchKey = iota
	ProjectRoleSearchKeyKey
	ProjectRoleSearchKeyProjectID
	ProjectRoleSearchKeyOrgID
	ProjectRoleSearchKeyResourceOwner
	ProjectRoleSearchKeyDisplayName
)

type ProjectRoleSearchQuery struct {
	Key    ProjectRoleSearchKey
	Method model.SearchMethod
	Value  interface{}
}

type ProjectRoleSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ProjectRoleView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *ProjectRoleSearchRequest) AppendMyOrgQuery(orgID string) {
	r.Queries = append(r.Queries, &ProjectRoleSearchQuery{Key: ProjectRoleSearchKeyOrgID, Method: model.SearchMethodEquals, Value: orgID})
}
func (r *ProjectRoleSearchRequest) AppendProjectQuery(projectID string) {
	r.Queries = append(r.Queries, &ProjectRoleSearchQuery{Key: ProjectRoleSearchKeyProjectID, Method: model.SearchMethodEquals, Value: projectID})
}

func (r *ProjectRoleSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
