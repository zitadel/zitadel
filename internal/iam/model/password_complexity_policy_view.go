package model

import (
	"time"

	"github.com/zitadel/zitadel/v2/internal/domain"
)

type PasswordComplexityPolicyView struct {
	AggregateID  string
	MinLength    uint64
	HasLowercase bool
	HasUppercase bool
	HasNumber    bool
	HasSymbol    bool
	Default      bool

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type PasswordComplexityPolicySearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn PasswordComplexityPolicySearchKey
	Asc           bool
	Queries       []*PasswordComplexityPolicySearchQuery
}

type PasswordComplexityPolicySearchKey int32

const (
	PasswordComplexityPolicySearchKeyUnspecified PasswordComplexityPolicySearchKey = iota
	PasswordComplexityPolicySearchKeyAggregateID
)

type PasswordComplexityPolicySearchQuery struct {
	Key    PasswordComplexityPolicySearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type PasswordComplexityPolicySearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*PasswordComplexityPolicyView
	Sequence    uint64
	Timestamp   time.Time
}
