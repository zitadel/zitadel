package model

import "github.com/caos/zitadel/internal/eventstore/models"

type PasswordComplexityPolicy struct {
	models.ObjectRoot

	Description  string
	State        PolicyState
	MinLength    uint64
	HasLowercase bool
	HasUppercase bool
	HasNumber    bool
	HasSymbol    bool
}

func (p *PasswordComplexityPolicy) IsValid() bool {
	return p.Description != ""
}
