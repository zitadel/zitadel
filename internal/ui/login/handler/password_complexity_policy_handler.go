package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/errors"
	policy_model "github.com/caos/zitadel/internal/policy/model"
	"net/http"
	"regexp"
	"strconv"
)

var (
	hasStringLowerCase = regexp.MustCompile(`[a-z]`).MatchString
	hasStringUpperCase = regexp.MustCompile(`[A-Z]`).MatchString
	hasNumber          = regexp.MustCompile(`[0-9]`).MatchString
	hasSymbol          = regexp.MustCompile(`[^A-Za-z0-9]`).MatchString
)

func (l *Login) getPasswordComplexityPolicy(r *http.Request, authReq *model.AuthRequest) (*policy_model.PasswordComplexityPolicy, string, error) {
	policy, err := l.authRepo.GetMyPasswordComplexityPolicy(setContext(r.Context(), authReq.UserOrgID))
	if err != nil {
		return nil, err.Error(), err
	}
	description := "<ul class=\"passwordcomplexity\">"
	minLength := l.renderer.Localize("Password.MinLength", nil)
	description += "<li>" + minLength + " " + strconv.Itoa(int(policy.MinLength)) + "</li>"
	if policy.HasUppercase {
		uppercase := l.renderer.Localize("Password.HasUppercase", nil)
		description += "<li>" + uppercase + "</li>"
	}
	if policy.HasLowercase {
		lowercase := l.renderer.Localize("Password.HasLowercase", nil)
		description += "<li>" + lowercase + "</li>"
	}
	if policy.HasNumber {
		hasnumber := l.renderer.Localize("Password.HasNumber", nil)
		description += "<li>" + hasnumber + "</li>"
	}
	if policy.HasSymbol {
		hassymbol := l.renderer.Localize("Password.HasSymbol", nil)
		description += "<li>" + hassymbol + "</li>"
	}

	description += "</ul>"
	return policy, description, nil
}

func (l *Login) checkPasswordComplexityPolicy(password string, r *http.Request, authReq *model.AuthRequest) error {
	policy, err := l.authRepo.GetMyPasswordComplexityPolicy(setContext(r.Context(), authReq.UserOrgID))
	if err != nil {
		return nil
	}
	if policy.MinLength != 0 && uint64(len(password)) < policy.MinLength {
		return errors.ThrowInvalidArgument(nil, "POLICY-LSo0p", "Errors.User.PasswordComplexityPolicy.MinLength")
	}

	if policy.HasLowercase && !hasStringLowerCase(password) {
		return errors.ThrowInvalidArgument(nil, "POLICY-4Sjsf", "Errors.User.PasswordComplexityPolicy.HasLower")
	}

	if policy.HasUppercase && !hasStringUpperCase(password) {
		return errors.ThrowInvalidArgument(nil, "POLICY-6Sjc9", "Errors.User.PasswordComplexityPolicy.HasUpper")
	}

	if policy.HasNumber && !hasNumber(password) {
		return errors.ThrowInvalidArgument(nil, "POLICY-2Fksi", "Errors.User.PasswordComplexityPolicy.HasNumber")
	}

	if policy.HasSymbol && !hasSymbol(password) {
		return errors.ThrowInvalidArgument(nil, "POLICY-0Js6e", "Errors.User.PasswordComplexityPolicy.HasSymbol")
	}
	return nil
}
