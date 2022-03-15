package command

import (
	"regexp"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

var (
	hasStringLowerCase = regexp.MustCompile(`[a-z]`).MatchString
	hasStringUpperCase = regexp.MustCompile(`[A-Z]`).MatchString
	hasNumber          = regexp.MustCompile(`[0-9]`).MatchString
	hasSymbol          = regexp.MustCompile(`[^A-Za-z0-9]`).MatchString
)

type PasswordComplexityPolicyWriteModel struct {
	eventstore.WriteModel

	MinLength    uint64
	HasLowercase bool
	HasUppercase bool
	HasNumber    bool
	HasSymbol    bool
	State        domain.PolicyState
}

func (wm *PasswordComplexityPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.PasswordComplexityPolicyAddedEvent:
			wm.MinLength = e.MinLength
			wm.HasLowercase = e.HasLowercase
			wm.HasUppercase = e.HasUppercase
			wm.HasNumber = e.HasNumber
			wm.HasSymbol = e.HasSymbol
			wm.State = domain.PolicyStateActive
		case *policy.PasswordComplexityPolicyChangedEvent:
			if e.MinLength != nil {
				wm.MinLength = *e.MinLength
			}
			if e.HasLowercase != nil {
				wm.HasLowercase = *e.HasLowercase
			}
			if e.HasUppercase != nil {
				wm.HasUppercase = *e.HasUppercase
			}
			if e.HasNumber != nil {
				wm.HasNumber = *e.HasNumber
			}
			if e.HasSymbol != nil {
				wm.HasSymbol = *e.HasSymbol
			}
		case *policy.PasswordComplexityPolicyRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *PasswordComplexityPolicyWriteModel) Validate(password string) error {
	if wm.MinLength != 0 && uint64(len(password)) < wm.MinLength {
		return errors.ThrowInvalidArgument(nil, "COMMA-HuJf6", "Errors.User.PasswordComplexityPolicy.MinLength")
	}

	if wm.HasLowercase && !hasStringLowerCase(password) {
		return errors.ThrowInvalidArgument(nil, "COMMA-co3Xw", "Errors.User.PasswordComplexityPolicy.HasLower")
	}

	if wm.HasUppercase && !hasStringUpperCase(password) {
		return errors.ThrowInvalidArgument(nil, "COMMA-VoaRj", "Errors.User.PasswordComplexityPolicy.HasUpper")
	}

	if wm.HasNumber && !hasNumber(password) {
		return errors.ThrowInvalidArgument(nil, "COMMA-ZBv4H", "Errors.User.PasswordComplexityPolicy.HasNumber")
	}

	if wm.HasSymbol && !hasSymbol(password) {
		return errors.ThrowInvalidArgument(nil, "COMMA-ZDLwA", "Errors.User.PasswordComplexityPolicy.HasSymbol")
	}
	return nil
}
