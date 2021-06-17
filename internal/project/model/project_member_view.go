package model

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"

	"time"
)

type ProjectMemberView struct {
	UserID             string
	ProjectID          string
	UserName           string
	Email              string
	FirstName          string
	LastName           string
	DisplayName        string
	PreferredLoginName string
	AvatarURL          string
	UserResourceOwner  string
	Roles              []string
	CreationDate       time.Time
	ChangeDate         time.Time
	Sequence           uint64
}

type ProjectMemberSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn ProjectMemberSearchKey
	Asc           bool
	Queries       []*ProjectMemberSearchQuery
}

type ProjectMemberSearchKey int32

const (
	ProjectMemberSearchKeyUnspecified ProjectMemberSearchKey = iota
	ProjectMemberSearchKeyUserName
	ProjectMemberSearchKeyEmail
	ProjectMemberSearchKeyFirstName
	ProjectMemberSearchKeyLastName
	ProjectMemberSearchKeyProjectID
	ProjectMemberSearchKeyUserID
)

type ProjectMemberSearchQuery struct {
	Key    ProjectMemberSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type ProjectMemberSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ProjectMemberView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *ProjectMemberSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-389Nd", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}
func (r *ProjectMemberSearchRequest) AppendProjectQuery(projectID string) {
	r.Queries = append(r.Queries, &ProjectMemberSearchQuery{Key: ProjectMemberSearchKeyProjectID, Method: domain.SearchMethodEquals, Value: projectID})
}
