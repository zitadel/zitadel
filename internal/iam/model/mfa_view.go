package model

import (
	"github.com/caos/zitadel/internal/model"
)

type SecondFactorsSearchRequest struct {
	Queries []*MFASearchQuery
}

type MultiFactorsSearchRequest struct {
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

type SecondFactorsSearchResponse struct {
	TotalResult uint64
	Result      []SecondFactorType
}

type MultiFactorsSearchResponse struct {
	TotalResult uint64
	Result      []MultiFactorType
}

func (r *SecondFactorsSearchRequest) AppendAggregateIDQuery(aggregateID string) {
	r.Queries = append(r.Queries, &MFASearchQuery{Key: MFASearchKeyAggregateID, Method: model.SearchMethodEquals, Value: aggregateID})
}

func (r *MultiFactorsSearchRequest) AppendAggregateIDQuery(aggregateID string) {
	r.Queries = append(r.Queries, &MFASearchQuery{Key: MFASearchKeyAggregateID, Method: model.SearchMethodEquals, Value: aggregateID})
}
