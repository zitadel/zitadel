package model

import (
	"github.com/caos/zitadel/internal/domain"
	"time"
)

type PasswordAgePolicyView struct {
	AggregateID    string
	MaxAgeDays     uint64
	ExpireWarnDays uint64
	Default        bool

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type PasswordAgePolicySearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn PasswordAgePolicySearchKey
	Asc           bool
	Queries       []*PasswordAgePolicySearchQuery
}

type PasswordAgePolicySearchKey int32

const (
	PasswordAgePolicySearchKeyUnspecified PasswordAgePolicySearchKey = iota
	PasswordAgePolicySearchKeyAggregateID
)

type PasswordAgePolicySearchQuery struct {
	Key    PasswordAgePolicySearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type PasswordAgePolicySearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*PasswordAgePolicyView
	Sequence    uint64
	Timestamp   time.Time
}
