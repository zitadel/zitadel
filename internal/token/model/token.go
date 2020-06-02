package model

import (
	"time"

	"github.com/caos/zitadel/internal/model"
)

type Token struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	UserID        string
	ApplicationID string
	UserAgentID   string
	Audience      []string
	Expiration    time.Time
	Scopes        []string
	Sequence      uint64
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
	TOKENSEARCHKEY_UNSPECIFIED TokenSearchKey = iota
	TOKENSEARCHKEY_TOKEN_ID
	TOKENSEARCHKEY_USER_ID
	TOKENSEARCHKEY_APPLICATION_ID
	TOKENSEARCHKEY_USER_AGENT_ID
	TOKENSEARCHKEY_EXPIRATION
	TOKENSEARCHKEY_RESOURCEOWNER
)

type TokenSearchQuery struct {
	Key    TokenSearchKey
	Method model.SearchMethod
	Value  string
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
