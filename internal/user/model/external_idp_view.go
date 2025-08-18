package model

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ExternalIDPView struct {
	UserID          string
	IDPConfigID     string
	ExternalUserID  string
	IDPName         string
	UserDisplayName string
	CreationDate    time.Time
	ChangeDate      time.Time
	ResourceOwner   string
	Sequence        uint64
}

type ExternalIDPSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn ExternalIDPSearchKey
	Asc           bool
	Queries       []*ExternalIDPSearchQuery
}

type ExternalIDPSearchKey int32

const (
	ExternalIDPSearchKeyUnspecified ExternalIDPSearchKey = iota
	ExternalIDPSearchKeyExternalUserID
	ExternalIDPSearchKeyUserID
	ExternalIDPSearchKeyIdpConfigID
	ExternalIDPSearchKeyResourceOwner
	ExternalIDPSearchKeyInstanceID
	ExternalIDPSearchKeyOwnerRemoved
)

type ExternalIDPSearchQuery struct {
	Key    ExternalIDPSearchKey
	Method domain.SearchMethod
	Value  any
}

type ExternalIDPSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ExternalIDPView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *ExternalIDPSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return zerrors.ThrowInvalidArgument(nil, "SEARCH-3n8fM", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}

func (r *ExternalIDPSearchRequest) AppendUserQuery(userID string) {
	r.Queries = append(r.Queries, &ExternalIDPSearchQuery{Key: ExternalIDPSearchKeyUserID, Method: domain.SearchMethodEquals, Value: userID})
}
