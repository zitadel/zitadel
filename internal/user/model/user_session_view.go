package model

import (
	"time"

	"github.com/caos/zitadel/internal/model"
)

type UserSessionView struct {
	ID                      string
	CreationDate            time.Time
	ChangeDate              time.Time
	State                   UserState
	ResourceOwner           string
	UserAgentID             string
	UserID                  string
	UserName                string
	PasswordVerification    time.Time
	MfaSoftwareVerification time.Time
	MfaHardwareVerification time.Time
	Sequence                uint64
}

type UserSessionSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn UserSessionSearchKey
	Asc           bool
	Queries       []*UserSessionSearchQuery
}

type UserSessionSearchKey int32

const (
	USERSESSIONSEARCHKEY_UNSPECIFIED UserSessionSearchKey = iota
	USERSESSIONSEARCHKEY_SESSION_ID
	USERSESSIONSEARCHKEY_USER_AGENT_ID
	USERSESSIONSEARCHKEY_USER_ID
	USERSESSIONSEARCHKEY_STATE
	USERSESSIONSEARCHKEY_RESOURCEOWNER
)

type UserSessionSearchQuery struct {
	Key    UserSessionSearchKey
	Method model.SearchMethod
	Value  string
}

type UserSessionSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*UserSessionView
}

func (r *UserSessionSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
