package model

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
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
	Value  interface{}
}

type ExternalIDPSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ExternalIDPView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *ExternalIDPSearchRequest) AppendUserQuery(userID string) {
	r.Queries = append(r.Queries, &ExternalIDPSearchQuery{Key: ExternalIDPSearchKeyUserID, Method: domain.SearchMethodEquals, Value: userID})
}
