package model

import (
	"github.com/caos/zitadel/internal/model"
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
)

type IDPProviderSearchQuery struct {
	Key    IDPProviderSearchKey
	Method model.SearchMethod
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

func (r *IDPProviderSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}

func (r *IDPProviderSearchRequest) AppendAggregateIDQuery(aggregateID string) {
	r.Queries = append(r.Queries, &IDPProviderSearchQuery{Key: IDPProviderSearchKeyAggregateID, Method: model.SearchMethodEquals, Value: aggregateID})
}
