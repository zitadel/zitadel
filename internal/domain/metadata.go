package domain

import (
	"time"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Metadata struct {
	es_models.ObjectRoot

	State MetadataState
	Key   string
	Value []byte
}

type MetadataState int32

const (
	MetadataStateUnspecified MetadataState = iota
	MetadataStateActive
	MetadataStateRemoved
)

func (m *Metadata) IsValid() bool {
	return m.Key != "" && len(m.Value) > 0
}

func (s MetadataState) Exists() bool {
	return s != MetadataStateUnspecified && s != MetadataStateRemoved
}

type MetadataSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn MetadataSearchKey
	Asc           bool
	Queries       []*MetadataSearchQuery
}

type MetadataSearchKey int32

const (
	MetadataSearchKeyUnspecified MetadataSearchKey = iota
	MetadataSearchKeyAggregateID
	MetadataSearchKeyResourceOwner
	MetadataSearchKeyKey
	MetadataSearchKeyValue
)

type MetadataSearchQuery struct {
	Key    MetadataSearchKey
	Method SearchMethod
	Value  any
}

type MetadataSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*Metadata
	Sequence    uint64
	Timestamp   time.Time
}

func (r *MetadataSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return zerrors.ThrowInvalidArgument(nil, "SEARCH-0ds32", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}

func (r *MetadataSearchRequest) AppendAggregateIDQuery(aggregateID string) {
	r.Queries = append(r.Queries, &MetadataSearchQuery{Key: MetadataSearchKeyAggregateID, Method: SearchMethodEquals, Value: aggregateID})
}

func (r *MetadataSearchRequest) AppendResourceOwnerQuery(resourceOwner string) {
	r.Queries = append(r.Queries, &MetadataSearchQuery{Key: MetadataSearchKeyResourceOwner, Method: SearchMethodEquals, Value: resourceOwner})
}
