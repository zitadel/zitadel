package model

import (
	"github.com/caos/zitadel/internal/model"
	"time"
)

type LoginPolicyView struct {
	AggregateID           string
	AllowUsernamePassword bool
	AllowRegister         bool
	AllowExternalIDP      bool
	ForceMFA              bool
	SoftwareMFAs          []SoftwareMFAType
	HardwareMFAs          []HardwareMFAType
	Default               bool

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type LoginPolicySearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn LoginPolicySearchKey
	Asc           bool
	Queries       []*LoginPolicySearchQuery
}

type LoginPolicySearchKey int32

const (
	LoginPolicySearchKeyUnspecified LoginPolicySearchKey = iota
	LoginPolicySearchKeyAggregateID
)

type LoginPolicySearchQuery struct {
	Key    LoginPolicySearchKey
	Method model.SearchMethod
	Value  interface{}
}

type LoginPolicySearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*LoginPolicyView
	Sequence    uint64
	Timestamp   time.Time
}
