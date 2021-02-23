package model

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/v2/domain"
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
	LoginPolicySearchKeyDefault
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

func secondFactorsToDomain(types []SecondFactorType) []domain.SecondFactorType {
	secondfactors := make([]domain.SecondFactorType, len(types))
	for i, secondfactorType := range types {
		switch secondfactorType {
		case SecondFactorTypeU2F:
			secondfactors[i] = domain.SecondFactorTypeU2F
		case SecondFactorTypeOTP:
			secondfactors[i] = domain.SecondFactorTypeOTP
		}
	}
	return secondfactors
}

func multiFactorsToDomain(types []MultiFactorType) []domain.MultiFactorType {
	multifactors := make([]domain.MultiFactorType, len(types))
	for i, multifactorType := range types {
		switch multifactorType {
		case MultiFactorTypeU2FWithPIN:
			multifactors[i] = domain.MultiFactorTypeU2FWithPIN
		}
	}
	return multifactors
}
