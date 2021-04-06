package model

import (
	"time"

	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"
)

type OrgMemberView struct {
	UserID             string
	OrgID              string
	UserName           string
	Email              string
	FirstName          string
	LastName           string
	DisplayName        string
	PreferredLoginName string
	Roles              []string
	CreationDate       time.Time
	ChangeDate         time.Time
	Sequence           uint64
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
	OrgMemberSearchKeyUnspecified OrgMemberSearchKey = iota
	OrgMemberSearchKeyUserName
	OrgMemberSearchKeyEmail
	OrgMemberSearchKeyFirstName
	OrgMemberSearchKeyLastName
	OrgMemberSearchKeyOrgID
	OrgMemberSearchKeyUserID
)

type OrgMemberSearchQuery struct {
	Key    OrgMemberSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type OrgMemberSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*OrgMemberView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *OrgMemberSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-77fu3", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}
