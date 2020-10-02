package model

import (
	"github.com/caos/zitadel/internal/model"
	"time"
)

type OrgIAMPolicyView struct {
	AggregateID           string
	UserLoginMustBeDomain bool
	Default               bool

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type OrgIAMPolicySearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn OrgIAMPolicySearchKey
	Asc           bool
	Queries       []*OrgIAMPolicySearchQuery
}

type OrgIAMPolicySearchKey int32

const (
	OrgIAMPolicySearchKeyUnspecified OrgIAMPolicySearchKey = iota
	OrgIAMPolicySearchKeyAggregateID
)

type OrgIAMPolicySearchQuery struct {
	Key    OrgIAMPolicySearchKey
	Method model.SearchMethod
	Value  interface{}
}

type OrgIAMPolicySearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*OrgIAMPolicyView
	Sequence    uint64
	Timestamp   time.Time
}
