package model

import (
	"github.com/caos/zitadel/internal/domain"
	"time"
)

type LockoutPolicyView struct {
	AggregateID         string
	MaxPasswordAttempts uint64
	ShowLockOutFailures bool
	Default             bool

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type LockoutPolicySearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn LockoutPolicySearchKey
	Asc           bool
	Queries       []*LockoutPolicySearchQuery
}

type LockoutPolicySearchKey int32

const (
	LockoutPolicySearchKeyUnspecified LockoutPolicySearchKey = iota
	LockoutPolicySearchKeyAggregateID
)

type LockoutPolicySearchQuery struct {
	Key    LockoutPolicySearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type LockoutPolicySearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*LockoutPolicyView
	Sequence    uint64
	Timestamp   time.Time
}
