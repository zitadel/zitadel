package model

import (
	"github.com/zitadel/zitadel/internal/domain"
	caos_errors "github.com/zitadel/zitadel/internal/errors"

	"time"
)

type IDPProviderView struct {
	AggregateID     string
	IDPConfigID     string
	IDPProviderType IDPProviderType
	Name            string
	StylingType     IDPStylingType
	IDPConfigType   IdpConfigType
	IDPState        IDPConfigState

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type IDPProviderSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn IDPProviderSearchKey
	Asc           bool
	Queries       []*IDPProviderSearchQuery
}

type IDPProviderSearchKey int32

const (
	IDPProviderSearchKeyUnspecified IDPProviderSearchKey = iota
	IDPProviderSearchKeyAggregateID
	IDPProviderSearchKeyIdpConfigID
	IDPProviderSearchKeyState
	IDPProviderSearchKeyInstanceID
	IDPProviderSearchKeyOwnerRemoved
)

type IDPProviderSearchQuery struct {
	Key    IDPProviderSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type IDPProviderSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*IDPProviderView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *IDPProviderSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-3n8fs", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}

func (r *IDPProviderSearchRequest) AppendAggregateIDQuery(aggregateID string) {
	r.Queries = append(r.Queries, &IDPProviderSearchQuery{Key: IDPProviderSearchKeyAggregateID, Method: domain.SearchMethodEquals, Value: aggregateID})
}
