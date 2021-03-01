package model

import (
	"github.com/caos/zitadel/internal/domain"
	"time"
)

type PasswordLockoutPolicyView struct {
	AggregateID         string
	MaxAttempts         uint64
	ShowLockOutFailures bool
	Default             bool

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type PasswordLockoutPolicySearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn PasswordLockoutPolicySearchKey
	Asc           bool
	Queries       []*PasswordLockoutPolicySearchQuery
}

type PasswordLockoutPolicySearchKey int32

const (
	PasswordLockoutPolicySearchKeyUnspecified PasswordLockoutPolicySearchKey = iota
	PasswordLockoutPolicySearchKeyAggregateID
)

type PasswordLockoutPolicySearchQuery struct {
	Key    PasswordLockoutPolicySearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type PasswordLockoutPolicySearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*PasswordLockoutPolicyView
	Sequence    uint64
	Timestamp   time.Time
}
