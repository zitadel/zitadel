package model

import (
	"github.com/caos/zitadel/internal/model"
	"time"
)

type UserGrantView struct {
	ID            string
	ResourceOwner string
	UserID        string
	ProjectID     string
	UserName      string
	FirstName     string
	LastName      string
	Email         string
	ProjectName   string
	OrgName       string
	OrgDomain     string
	RoleKeys      []string

	CreationDate time.Time
	ChangeDate   time.Time
	State        UserGrantState

	Sequence uint64
}

type UserGrantSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn UserGrantSearchKey
	Asc           bool
	Queries       []*UserGrantSearchQuery
}

type UserGrantSearchKey int32

const (
	USERGRANTSEARCHKEY_UNSPECIFIED UserGrantSearchKey = iota
	USERGRANTSEARCHKEY_USER_ID
	USERGRANTSEARCHKEY_PROJECT_ID
	USERGRANTSEARCHKEY_RESOURCEOWNER
	USERGRANTSEARCHKEY_STATE
	USERGRANTSEARCHKEY_GRANT_ID
	USERGRANTSEARCHKEY_ORG_NAME
	USERGRANTSEARCHKEY_ROLE_KEY
)

type UserGrantSearchQuery struct {
	Key    UserGrantSearchKey
	Method model.SearchMethod
	Value  interface{}
}

type UserGrantSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*UserGrantView
}

func (r *UserGrantSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}

func (r *UserGrantSearchRequest) AppendMyOrgQuery(orgID string) {
	r.Queries = append(r.Queries, &UserGrantSearchQuery{Key: USERGRANTSEARCHKEY_RESOURCEOWNER, Method: model.SEARCHMETHOD_EQUALS, Value: orgID})
}
