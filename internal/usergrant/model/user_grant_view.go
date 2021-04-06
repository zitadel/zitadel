package model

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"

	"time"
)

type UserGrantView struct {
	ID               string
	ResourceOwner    string
	UserID           string
	ProjectID        string
	GrantID          string
	UserName         string
	FirstName        string
	LastName         string
	DisplayName      string
	Email            string
	ProjectName      string
	OrgName          string
	OrgPrimaryDomain string
	RoleKeys         []string

	CreationDate time.Time
	ChangeDate   time.Time
	State        UserGrantState

	Sequence uint64
}

type UserGrantState int32

const (
	UserGrantStateActive UserGrantState = iota
	UserGrantStateInactive
	UserGrantStateRemoved
)

type UserGrantSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn UserGrantSearchKey
	Asc           bool
	Queries       []*UserGrantSearchQuery
}

type UserGrantSearchKey int32

const (
	UserGrantSearchKeyUnspecified UserGrantSearchKey = iota
	UserGrantSearchKeyUserID
	UserGrantSearchKeyProjectID
	UserGrantSearchKeyResourceOwner
	UserGrantSearchKeyState
	UserGrantSearchKeyGrantID
	UserGrantSearchKeyOrgName
	UserGrantSearchKeyRoleKey
	UserGrantSearchKeyID
	UserGrantSearchKeyUserName
	UserGrantSearchKeyFirstName
	UserGrantSearchKeyLastName
	UserGrantSearchKeyEmail
	UserGrantSearchKeyOrgDomain
	UserGrantSearchKeyProjectName
	UserGrantSearchKeyDisplayName
	UserGrantSearchKeyWithGranted
)

type UserGrantSearchQuery struct {
	Key    UserGrantSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type UserGrantSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*UserGrantView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *UserGrantSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-1N9ds", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}

func (r *UserGrantSearchRequest) GetSearchQuery(key UserGrantSearchKey) (int, *UserGrantSearchQuery) {
	for i, q := range r.Queries {
		if q.Key == key {
			return i, q
		}
	}
	return -1, nil
}

func (r *UserGrantSearchRequest) AppendMyOrgQuery(orgID string) {
	r.Queries = append(r.Queries, &UserGrantSearchQuery{Key: UserGrantSearchKeyResourceOwner, Method: domain.SearchMethodEquals, Value: orgID})
}

func (r *UserGrantSearchRequest) AppendProjectIDQuery(projectID string) {
	r.Queries = append(r.Queries, &UserGrantSearchQuery{Key: UserGrantSearchKeyProjectID, Method: domain.SearchMethodEquals, Value: projectID})
}
