package model

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"regexp"
)

var (
	hasStringLowerCase = regexp.MustCompile(`[a-z]`).MatchString
	hasStringUpperCase = regexp.MustCompile(`[A-Z]`).MatchString
	hasNumber          = regexp.MustCompile(`[0-9]`).MatchString
	hasSymbol          = regexp.MustCompile(`[^A-Za-z0-9]`).MatchString
)

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

func (p *PasswordComplexityPolicy) Check(password string) error {
	if p.MinLength != 0 && uint64(len(password)) < p.MinLength {
		return caos_errs.ThrowInvalidArgumentf(nil, "MODEL-HuJf6", "Passwordpolicy doesn't match: Minlength %v", p.MinLength)
	}

	if p.HasLowercase && !hasStringLowerCase(password) {
		return caos_errs.ThrowInvalidArgument(nil, "MODEL-co3Xw", "Passwordpolicy doesn't match: HasLowerCase")
	}

	if p.HasUppercase && !hasStringUpperCase(password) {
		return caos_errs.ThrowInvalidArgument(nil, "MODEL-VoaRj", "Passwordpolicy doesn't match: HasUpperCase")
	}

	if p.HasNumber && !hasNumber(password) {
		return caos_errs.ThrowInvalidArgument(nil, "MODEL-ZBv4H", "Passwordpolicy doesn't match: HasNumber")
	}

	if p.HasSymbol && !hasSymbol(password) {
		return caos_errs.ThrowInvalidArgument(nil, "MODEL-ZDLwA", "Passwordpolicy doesn't match: HasSymbol")
	}
	return nil
}
