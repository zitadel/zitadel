package model

import (
	"time"

	"github.com/caos/zitadel/internal/domain"
)

type PrivacyPolicyView struct {
	AggregateID string
	TOSLink     string
	PrivacyLink string
	Default     bool

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type PrivacyPolicySearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn PrivacyPolicySearchKey
	Asc           bool
	Queries       []*PrivacyPolicySearchQuery
}

type PrivacyPolicySearchKey int32

const (
	PrivacyPolicySearchKeyUnspecified PrivacyPolicySearchKey = iota
	PrivacyPolicySearchKeyAggregateID
)

type PrivacyPolicySearchQuery struct {
	Key    PrivacyPolicySearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type PrivacyPolicySearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*PrivacyPolicyView
	Sequence    uint64
	Timestamp   time.Time
}
