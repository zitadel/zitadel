package model

import (
	"github.com/caos/zitadel/internal/model"
	"time"
)

type IdpProviderView struct {
	AggregateID     string
	IdpConfigID     string
	IdpProviderType IdpProviderType
	Name            string
	IdpConfigType   IdpConfigType

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type IdpProviderSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn IdpProviderSearchKey
	Asc           bool
	Queries       []*IdpProviderSearchQuery
}

type IdpProviderSearchKey int32

const (
	IdpProviderSearchKeyUnspecified IdpProviderSearchKey = iota
	IdpProviderSearchKeyAggregateID
	IdpProviderSearchKeyIdpConfigID
)

type IdpProviderSearchQuery struct {
	Key    IdpProviderSearchKey
	Method model.SearchMethod
	Value  interface{}
}

type IdpProviderSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*IdpProviderView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *IdpProviderSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}

func (r *IdpProviderSearchRequest) AppendAggregateIDQuery(aggregateID string) {
	r.Queries = append(r.Queries, &IdpProviderSearchQuery{Key: IdpProviderSearchKeyAggregateID, Method: model.SearchMethodEquals, Value: aggregateID})
}
