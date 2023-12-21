package model

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"

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
	Sequence          uint64
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

type TokenSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*Token
}

func (r *TokenSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return zerrors.ThrowInvalidArgument(nil, "SEARCH-M0fse", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}
