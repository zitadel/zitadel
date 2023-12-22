package model

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UserSessionView struct {
	CreationDate                 time.Time
	ChangeDate                   time.Time
	State                        domain.UserSessionState
	ResourceOwner                string
	UserAgentID                  string
	UserID                       string
	UserName                     string
	LoginName                    string
	DisplayName                  string
	AvatarKey                    string
	SelectedIDPConfigID          string
	PasswordVerification         time.Time
	PasswordlessVerification     time.Time
	ExternalLoginVerification    time.Time
	SecondFactorVerification     time.Time
	SecondFactorVerificationType domain.MFAType
	MultiFactorVerification      time.Time
	MultiFactorVerificationType  domain.MFAType
	Sequence                     uint64
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
	UserSessionSearchKeyUnspecified UserSessionSearchKey = iota
	UserSessionSearchKeyUserAgentID
	UserSessionSearchKeyUserID
	UserSessionSearchKeyState
	UserSessionSearchKeyResourceOwner
	UserSessionSearchKeyInstanceID
	UserSessionSearchKeyOwnerRemoved
)

type UserSessionSearchQuery struct {
	Key    UserSessionSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type UserSessionSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*UserSessionView
}

func (r *UserSessionSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return zerrors.ThrowInvalidArgument(nil, "SEARCH-27ifs", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}
