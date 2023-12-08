package model

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"

	"time"
)

type RefreshTokenView struct {
	ID                    string
	CreationDate          time.Time
	ChangeDate            time.Time
	ResourceOwner         string
	UserID                string
	ClientID              string
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
	RefreshTokenSearchKeyInstanceID
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
	Sequence    uint64
	Timestamp   time.Time
	Result      []*RefreshTokenView
}

func (r *RefreshTokenSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return zerrors.ThrowInvalidArgument(nil, "SEARCH-M0fse", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}
