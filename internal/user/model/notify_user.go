package model

import (
	"github.com/caos/zitadel/internal/model"
	"time"
)

type NotifyUser struct {
	ID                string
	CreationDate      time.Time
	ChangeDate        time.Time
	ResourceOwner     string
	UserName          string
	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	PreferredLanguage string
	Gender            Gender
	LastEmail         string
	VerifiedEmail     string
	LastPhone         string
	VerifiedPhone     string
	PasswordSet       bool
	Sequence          uint64
}

type NotifyUserSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn NotifyUserSearchKey
	Asc           bool
	Queries       []*NotifyUserSearchQuery
}

type NotifyUserSearchKey int32

const (
	NOTIFYUSERSEARCHKEY_UNSPECIFIED UserSearchKey = iota
	NOTIFYUSERSEARCHKEY_USER_ID
)

type NotifyUserSearchQuery struct {
	Key    NotifyUserSearchKey
	Method model.SearchMethod
	Value  string
}

type NotifyUserSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*UserView
}
