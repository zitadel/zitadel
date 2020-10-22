package model

import (
	"github.com/caos/zitadel/internal/model"
)

type SoftwareMFASearchRequest struct {
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
	TotalResult uint64
	Result      []SoftwareMFAType
}

type HardwareMFASearchResponse struct {
	TotalResult uint64
	Result      []HardwareMFAType
}

func (r *SoftwareMFASearchRequest) AppendAggregateIDQuery(aggregateID string) {
	r.Queries = append(r.Queries, &MFASearchQuery{Key: MFASearchKeyAggregateID, Method: model.SearchMethodEquals, Value: aggregateID})
}

func (r *HardwareMFASearchRequest) AppendAggregateIDQuery(aggregateID string) {
	r.Queries = append(r.Queries, &MFASearchQuery{Key: MFASearchKeyAggregateID, Method: model.SearchMethodEquals, Value: aggregateID})
}
