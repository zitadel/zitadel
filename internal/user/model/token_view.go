package model

import (
	"github.com/caos/zitadel/internal/domain"
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
	TokenSearchKeyApplicationID
	TokenSearchKeyUserAgentID
	TokenSearchKeyExpiration
	TokenSearchKeyResourceOwner
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

func (r *TokenSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
