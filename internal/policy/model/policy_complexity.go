package model

import "github.com/caos/zitadel/internal/eventstore/models"

type PasswordComplexityPolicy struct {
	models.ObjectRoot

	Description  string
	State        int32
	MinLength    uint64
	HasLowercase bool
	HasUppercase bool
	HasNumber    bool
	HasSymbol    bool
}

func (p *PasswordComplexityPolicy) IsValid() bool {
	if p.Description == "" {
		return false
	}
	return true
}
