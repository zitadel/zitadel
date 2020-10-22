package model

import (
	"github.com/caos/zitadel/internal/model"
	"time"
)

type SoftwareMFAView struct {
	AggregateID     string
	IDPProviderType SoftwareMFAType
	Sequence        uint64
}

type HardwareMFAView struct {
	AggregateID     string
	IDPProviderType SoftwareMFAType
	Sequence        uint64
}

type SoftwareMFASearchRequest struct {
	Offset  uint64
	Limit   uint64
	Asc     bool
	Queries []*MFASearchQuery
}

type HardwareMFASearchRequest struct {
	Offset  uint64
	Limit   uint64
	Asc     bool
	Queries []*MFASearchQuery
}

type MFASearchQuery struct {
	Key    MFASearchKey
	Method model.SearchMethod
	Value  interface{}
}

type MFASearchKey int32

const (
	MFASearchKeyUnspecified MFASearchKey = iota
	MFASearchKeyAggregateID
)

type SoftwareMFASearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*SoftwareMFAView
	Sequence    uint64
	Timestamp   time.Time
}

type HardwareMFASearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*HardwareMFAView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *SoftwareMFASearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}

func (r *HardwareMFASearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}

func (r *SoftwareMFASearchRequest) AppendAggregateIDQuery(aggregateID string) {
	r.Queries = append(r.Queries, &MFASearchQuery{Key: MFASearchKeyAggregateID, Method: model.SearchMethodEquals, Value: aggregateID})
}

func (r *HardwareMFASearchRequest) AppendAggregateIDQuery(aggregateID string) {
	r.Queries = append(r.Queries, &MFASearchQuery{Key: MFASearchKeyAggregateID, Method: model.SearchMethodEquals, Value: aggregateID})
}
