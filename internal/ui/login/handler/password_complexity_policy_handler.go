package handler

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

const (
	LowerCaseRegex = `[a-z]`
	UpperCaseRegex = `[A-Z]`
	NumberRegex    = `[0-9]`
	SymbolRegex    = `[^A-Za-z0-9]`
)

var (
	hasStringLowerCase = regexp.MustCompile(LowerCaseRegex).MatchString
	hasStringUpperCase = regexp.MustCompile(UpperCaseRegex).MatchString
	hasNumber          = regexp.MustCompile(NumberRegex).MatchString
	hasSymbol          = regexp.MustCompile(SymbolRegex).MatchString
)

func (l *Login) getPasswordComplexityPolicy(r *http.Request, orgID string) (*iam_model.PasswordComplexityPolicyView, string, error) {
	policy, err := l.authRepo.GetMyPasswordComplexityPolicy(setContext(r.Context(), orgID))
	if err != nil {
		return nil, err.Error(), err
	}
	description, err := l.generatePolicyDescription(r, policy)
	return policy, description, nil
}

func (l *Login) getPasswordComplexityPolicyByUserID(r *http.Request, userID string) (*iam_model.PasswordComplexityPolicyView, string, error) {
	user, err := l.authRepo.UserByID(r.Context(), userID)
	if err != nil {
		return nil, "", nil
	}
	policy, err := l.authRepo.GetMyPasswordComplexityPolicy(setContext(r.Context(), user.ResourceOwner))
	if err != nil {
		return nil, err.Error(), err
	}
	description, err := l.generatePolicyDescription(r, policy)
	return policy, description, nil
}

func (l *Login) generatePolicyDescription(r *http.Request, policy *iam_model.PasswordComplexityPolicyView) (string, error) {
	description := "<ul class=\"lgn-no-dots lgn-policy\" id=\"passwordcomplexity\">"
	minLength := l.renderer.LocalizeFromRequest(r, "Password.MinLength", nil)
	description += "<li id=\"minlength\" class=\"invalid\"><i class=\"lgn-icon-times-solid lgn-warn\"></i><span>" + minLength + " " + strconv.Itoa(int(policy.MinLength)) + "</span></li>"
	if policy.HasUppercase {
		uppercase := l.renderer.LocalizeFromRequest(r, "Password.HasUppercase", nil)
		description += "<li id=\"uppercase\" class=\"invalid\"><i class=\"lgn-icon-times-solid lgn-warn\"></i><span>" + uppercase + "</span></li>"
	}
	if policy.HasLowercase {
		lowercase := l.renderer.LocalizeFromRequest(r, "Password.HasLowercase", nil)
		description += "<li id=\"lowercase\" class=\"invalid\"><i class=\"lgn-icon-times-solid lgn-warn\"></i><span>" + lowercase + "</span></li>"
	}
	if policy.HasNumber {
		hasnumber := l.renderer.LocalizeFromRequest(r, "Password.HasNumber", nil)
		description += "<li id=\"number\" class=\"invalid\"><i class=\"lgn-icon-times-solid lgn-warn\"></i><span>" + hasnumber + "</span></li>"
	}
	if policy.HasSymbol {
		hassymbol := l.renderer.LocalizeFromRequest(r, "Password.HasSymbol", nil)
		description += "<li id=\"symbol\" class=\"invalid\"><i class=\"lgn-icon-times-solid lgn-warn\"></i><span>" + hassymbol + "</span></li>"
	}
	confirmation := l.renderer.LocalizeFromRequest(r, "Password.Confirmation", nil)
	description += "<li id=\"confirmation\" class=\"invalid\"><i class=\"lgn-icon-times-solid lgn-warn\"></i><span>" + confirmation + "</span></li>"

	description += "</ul>"
	return description, nil
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
