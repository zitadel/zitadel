package model

import (
	"time"

	"github.com/caos/zitadel/internal/domain"
)

type FeaturesView struct {
	AggregateID  string
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
	Default      bool

	TierName                 string
	TierDescription          string
	TierState                domain.TierState
	TierStateDescription     string
	LoginPolicyFactors       bool
	LoginPolicyIDP           bool
	LoginPolicyPasswordless  bool
	LoginPolicyRegistration  bool
	LoginPolicyUsernameLogin bool
}

type FeaturesSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn FeaturesSearchKey
	Asc           bool
	Queries       []*FeaturesSearchQuery
}

type FeaturesSearchKey int32

const (
	FeaturesSearchKeyUnspecified FeaturesSearchKey = iota
	FeaturesSearchKeyAggregateID
	FeaturesSearchKeyDefault
	FeaturesSearchKeyResourceOwner
)

type FeaturesSearchQuery struct {
	Key    FeaturesSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type FeaturesSearchResult struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*FeaturesView
	Sequence    uint64
	Timestamp   time.Time
}
