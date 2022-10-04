package model

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
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
)

type UserSessionSearchQuery struct {
	Key    UserSessionSearchKey
	Method domain.SearchMethod
	Value  interface{}
}
