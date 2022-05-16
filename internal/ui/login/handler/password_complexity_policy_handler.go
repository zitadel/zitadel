package handler

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/zitadel/zitadel/internal/domain"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
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

func (l *Login) getPasswordComplexityPolicy(r *http.Request, authReq *domain.AuthRequest, orgID string) (*iam_model.PasswordComplexityPolicyView, string, error) {
	policy, err := l.authRepo.GetMyPasswordComplexityPolicy(setContext(r.Context(), orgID))
	if err != nil {
		return nil, err.Error(), err
	}
	description, err := l.generatePolicyDescription(r, authReq, policy)
	return policy, description, nil
}

func (l *Login) getPasswordComplexityPolicyByUserID(r *http.Request, authReq *domain.AuthRequest, userID string) (*iam_model.PasswordComplexityPolicyView, string, error) {
	user, err := l.query.GetUserByID(r.Context(), userID, false)
	if err != nil {
		return nil, "", nil
	}
	policy, err := l.authRepo.GetMyPasswordComplexityPolicy(setContext(r.Context(), user.ResourceOwner))
	if err != nil {
		return nil, err.Error(), err
	}
	description, err := l.generatePolicyDescription(r, authReq, policy)
	return policy, description, nil
}

func (l *Login) generatePolicyDescription(r *http.Request, authReq *domain.AuthRequest, policy *iam_model.PasswordComplexityPolicyView) (string, error) {
	description := "<ul class=\"lgn-no-dots lgn-policy\" id=\"passwordcomplexity\">"
	translator := l.getTranslator(authReq)
	minLength := l.renderer.LocalizeFromRequest(translator, r, "Password.MinLength", nil)
	description += "<li id=\"minlength\" class=\"invalid\"><i class=\"lgn-icon-times-solid lgn-warn\"></i><span>" + minLength + " " + strconv.Itoa(int(policy.MinLength)) + "</span></li>"
	if policy.HasUppercase {
		uppercase := l.renderer.LocalizeFromRequest(translator, r, "Password.HasUppercase", nil)
		description += "<li id=\"uppercase\" class=\"invalid\"><i class=\"lgn-icon-times-solid lgn-warn\"></i><span>" + uppercase + "</span></li>"
	}
	if policy.HasLowercase {
		lowercase := l.renderer.LocalizeFromRequest(translator, r, "Password.HasLowercase", nil)
		description += "<li id=\"lowercase\" class=\"invalid\"><i class=\"lgn-icon-times-solid lgn-warn\"></i><span>" + lowercase + "</span></li>"
	}
	if policy.HasNumber {
		hasnumber := l.renderer.LocalizeFromRequest(translator, r, "Password.HasNumber", nil)
		description += "<li id=\"number\" class=\"invalid\"><i class=\"lgn-icon-times-solid lgn-warn\"></i><span>" + hasnumber + "</span></li>"
	}
	if policy.HasSymbol {
		hassymbol := l.renderer.LocalizeFromRequest(translator, r, "Password.HasSymbol", nil)
		description += "<li id=\"symbol\" class=\"invalid\"><i class=\"lgn-icon-times-solid lgn-warn\"></i><span>" + hassymbol + "</span></li>"
	}
	confirmation := l.renderer.LocalizeFromRequest(translator, r, "Password.Confirmation", nil)
	description += "<li id=\"confirmation\" class=\"invalid\"><i class=\"lgn-icon-times-solid lgn-warn\"></i><span>" + confirmation + "</span></li>"

	description += "</ul>"
	return description, nil
}
