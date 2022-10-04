package model

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
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
