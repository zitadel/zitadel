package domain

import (
	"regexp"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	hasStringLowerCase = regexp.MustCompile(`[a-z]`).MatchString
	hasStringUpperCase = regexp.MustCompile(`[A-Z]`).MatchString
	hasNumber          = regexp.MustCompile(`[0-9]`).MatchString
	hasSymbol          = regexp.MustCompile(`[^A-Za-z0-9]`).MatchString
)

type PasswordComplexityPolicy struct {
	models.ObjectRoot

	MinLength    uint64
	HasLowercase bool
	HasUppercase bool
	HasNumber    bool
	HasSymbol    bool

	Default bool
}

func (p *PasswordComplexityPolicy) IsValid() error {
	if p.MinLength == 0 || p.MinLength > 72 {
		return zerrors.ThrowInvalidArgument(nil, "MODEL-Lsp0e", "Errors.User.PasswordComplexityPolicy.MinLengthNotAllowed")
	}
	return nil
}

func (p *PasswordComplexityPolicy) Check(password string) error {
	if p.MinLength != 0 && uint64(len(password)) < p.MinLength {
		return zerrors.ThrowInvalidArgument(nil, "DOMAIN-HuJf6", "Errors.User.PasswordComplexityPolicy.MinLength")
	}

	if p.HasLowercase && !hasStringLowerCase(password) {
		return zerrors.ThrowInvalidArgument(nil, "DOMAIN-co3Xw", "Errors.User.PasswordComplexityPolicy.HasLower")
	}

	if p.HasUppercase && !hasStringUpperCase(password) {
		return zerrors.ThrowInvalidArgument(nil, "DOMAIN-VoaRj", "Errors.User.PasswordComplexityPolicy.HasUpper")
	}

	if p.HasNumber && !hasNumber(password) {
		return zerrors.ThrowInvalidArgument(nil, "DOMAIN-ZBv4H", "Errors.User.PasswordComplexityPolicy.HasNumber")
	}

	if p.HasSymbol && !hasSymbol(password) {
		return zerrors.ThrowInvalidArgument(nil, "DOMAIN-ZDLwA", "Errors.User.PasswordComplexityPolicy.HasSymbol")
	}
	return nil
}
