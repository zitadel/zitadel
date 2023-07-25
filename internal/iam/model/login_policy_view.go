package model

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type LoginPolicyView struct {
	AggregateID           string
	AllowUsernamePassword bool
	AllowRegister         bool
	AllowExternalIDP      bool
	ForceMFA              bool
	HidePasswordReset     bool
	PasswordlessType      PasswordlessType
	SecondFactors         []domain.SecondFactorType
	MultiFactors          []domain.MultiFactorType
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
	LoginPolicySearchKeyDefault
)

type LoginPolicySearchQuery struct {
	Key    LoginPolicySearchKey
	Method domain.SearchMethod
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

func (p *LoginPolicyView) ToLoginPolicyDomain() *domain.LoginPolicy {
	return &domain.LoginPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID:  p.AggregateID,
			CreationDate: p.CreationDate,
			ChangeDate:   p.ChangeDate,
			Sequence:     p.Sequence,
		},
		Default:               p.Default,
		AllowUsernamePassword: p.AllowUsernamePassword,
		AllowRegister:         p.AllowRegister,
		AllowExternalIDP:      p.AllowExternalIDP,
		ForceMFA:              p.ForceMFA,
		HidePasswordReset:     p.HidePasswordReset,
		PasswordlessType:      passwordLessTypeToDomain(p.PasswordlessType),
		SecondFactors:         secondFactorsToDomain(p.SecondFactors),
		MultiFactors:          multiFactorsToDomain(p.MultiFactors),
	}
}

func passwordLessTypeToDomain(passwordless PasswordlessType) domain.PasswordlessType {
	switch passwordless {
	case PasswordlessTypeNotAllowed:
		return domain.PasswordlessTypeNotAllowed
	case PasswordlessTypeAllowed:
		return domain.PasswordlessTypeAllowed
	default:
		return domain.PasswordlessTypeNotAllowed
	}
}

func secondFactorsToDomain(types []domain.SecondFactorType) []domain.SecondFactorType {
	secondfactors := make([]domain.SecondFactorType, len(types))
	for i, secondfactorType := range types {
		switch secondfactorType {
		case domain.SecondFactorTypeU2F:
			secondfactors[i] = domain.SecondFactorTypeU2F
		case domain.SecondFactorTypeTOTP:
			secondfactors[i] = domain.SecondFactorTypeTOTP
		case domain.SecondFactorTypeOTPEmail:
			secondfactors[i] = domain.SecondFactorTypeOTPEmail
		case domain.SecondFactorTypeOTPSMS:
			secondfactors[i] = domain.SecondFactorTypeOTPSMS
		}
	}
	return secondfactors
}

func multiFactorsToDomain(types []domain.MultiFactorType) []domain.MultiFactorType {
	multifactors := make([]domain.MultiFactorType, len(types))
	for i, multifactorType := range types {
		switch multifactorType {
		case domain.MultiFactorTypeU2FWithPIN:
			multifactors[i] = domain.MultiFactorTypeU2FWithPIN
		}
	}
	return multifactors
}
