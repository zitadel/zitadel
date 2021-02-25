package model

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"time"
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

func (p *PasswordComplexityPolicyView) Check(password string) error {
	if p.MinLength != 0 && uint64(len(password)) < p.MinLength {
		return caos_errs.ThrowInvalidArgument(nil, "MODEL-HuJf6", "Errors.User.PasswordComplexityPolicy.MinLength")
	}

	if p.HasLowercase && !hasStringLowerCase(password) {
		return caos_errs.ThrowInvalidArgument(nil, "MODEL-co3Xw", "Errors.User.PasswordComplexityPolicy.HasLower")
	}

	if p.HasUppercase && !hasStringUpperCase(password) {
		return caos_errs.ThrowInvalidArgument(nil, "MODEL-VoaRj", "Errors.User.PasswordComplexityPolicy.HasUpper")
	}

	if p.HasNumber && !hasNumber(password) {
		return caos_errs.ThrowInvalidArgument(nil, "MODEL-ZBv4H", "Errors.User.PasswordComplexityPolicy.HasNumber")
	}

	if p.HasSymbol && !hasSymbol(password) {
		return caos_errs.ThrowInvalidArgument(nil, "MODEL-ZDLwA", "Errors.User.PasswordComplexityPolicy.HasSymbol")
	}
	return nil
}
