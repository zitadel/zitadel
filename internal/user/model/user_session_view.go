package model

import (
	"time"

	req_model "github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/model"
)

type UserSessionView struct {
	CreationDate                time.Time
	ChangeDate                  time.Time
	State                       req_model.UserSessionState
	ResourceOwner               string
	UserAgentID                 string
	UserID                      string
	UserName                    string
	PasswordVerification        time.Time
	MfaSoftwareVerification     time.Time
	MfaSoftwareVerificationType req_model.MfaType
	MfaHardwareVerification     time.Time
	MfaHardwareVerificationType req_model.MfaType
	Sequence                    uint64
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
