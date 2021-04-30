package model

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"

	"time"
)

type RefreshTokenView struct {
	ID                    string
	CreationDate          time.Time
	ChangeDate            time.Time
	ResourceOwner         string
	UserID                string
	ApplicationID         string
	UserAgentID           string
	AuthMethodsReferences []string
	Audience              []string
	AuthTime              time.Time
	IdleExpiration        time.Time
	Expiration            time.Time
	Scopes                []string
	Sequence              uint64
	Token                 string
}

type RefreshTokenSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn RefreshTokenSearchKey
	Asc           bool
	Queries       []*RefreshTokenSearchQuery
}

type RefreshTokenSearchKey int32

const (
	RefreshTokenSearchKeyUnspecified RefreshTokenSearchKey = iota
	RefreshTokenSearchKeyRefreshTokenID
	RefreshTokenSearchKeyUserID
	RefreshTokenSearchKeyApplicationID
	RefreshTokenSearchKeyUserAgentID
	RefreshTokenSearchKeyExpiration
	RefreshTokenSearchKeyResourceOwner
)

type RefreshTokenSearchQuery struct {
	Key    RefreshTokenSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type RefreshTokenSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*RefreshToken
}

func (r *RefreshTokenSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-M0fse", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}
