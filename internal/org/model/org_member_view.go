package model

import (
	"time"

	"github.com/caos/zitadel/internal/model"
)

type OrgMemberView struct {
	UserID       string
	OrgID        string
	UserName     string
	Email        string
	FirstName    string
	LastName     string
	Roles        []string
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type OrgMemberSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn OrgMemberSearchKey
	Asc           bool
	Queries       []*OrgMemberSearchQuery
}

type OrgMemberSearchKey int32

const (
	ORGMEMBERSEARCHKEY_UNSPECIFIED OrgMemberSearchKey = iota
	ORGMEMBERSEARCHKEY_USER_NAME
	ORGMEMBERSEARCHKEY_EMAIL
	ORGMEMBERSEARCHKEY_FIRST_NAME
	ORGMEMBERSEARCHKEY_LAST_NAME
	ORGMEMBERSEARCHKEY_ORG_ID
	ORGMEMBERSEARCHKEY_USER_ID
)

type OrgMemberSearchQuery struct {
	Key    OrgMemberSearchKey
	Method model.SearchMethod
	Value  interface{}
}

type OrgMemberSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*OrgMemberView
}

func (r *OrgMemberSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
