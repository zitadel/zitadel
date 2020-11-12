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
	PasswordlessType      PasswordlessType
	SecondFactors         []SecondFactorType
	MultiFactors          []MultiFactorType
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

func (p *LoginPolicyView) HasSecondFactors() bool {
	if p.SecondFactors == nil || len(p.SecondFactors) == 0 {
		return false
	}
	return true
}

func (p *LoginPolicyView) HasMultiFactors() bool {
	if p.MultiFactors == nil || len(p.MultiFactors) == 0 {
		return false
	}
	return true
}
