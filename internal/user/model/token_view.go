package model

import (
	"github.com/zitadel/zitadel/internal/domain"

	"time"
)

type TokenView struct {
	ID                string
	CreationDate      time.Time
	ChangeDate        time.Time
	ResourceOwner     string
	UserID            string
	ApplicationID     string
	UserAgentID       string
	Audience          []string
	Expiration        time.Time
	Scopes            []string
	PreferredLanguage string
	RefreshTokenID    string
	IsPAT             bool
}

type TokenSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn TokenSearchKey
	Asc           bool
	Queries       []*TokenSearchQuery
}

type TokenSearchKey int32

const (
	TokenSearchKeyUnspecified TokenSearchKey = iota
	TokenSearchKeyTokenID
	TokenSearchKeyUserID
	TokenSearchKeyRefreshTokenID
	TokenSearchKeyApplicationID
	TokenSearchKeyUserAgentID
	TokenSearchKeyExpiration
	TokenSearchKeyResourceOwner
	TokenSearchKeyInstanceID
)

type TokenSearchQuery struct {
	Key    TokenSearchKey
	Method domain.SearchMethod
	Value  interface{}
}
