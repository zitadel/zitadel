package login

import (
	"net/http"

	"github.com/zitadel/logging"

	iam_model "github.com/zitadel/zitadel/internal/iam/model"
)

const (
	LowerCaseRegex = `[a-z]`
	UpperCaseRegex = `[A-Z]`
	NumberRegex    = `[0-9]`
	SymbolRegex    = `[^A-Za-z0-9]`
)

func (l *Login) getPasswordComplexityPolicy(r *http.Request, orgID string) *iam_model.PasswordComplexityPolicyView {
	policy, err := l.authRepo.GetMyPasswordComplexityPolicy(setContext(r.Context(), orgID))
	logging.WithFields("orgID", orgID).OnError(err).Error("could not load password complexity policy")
	return policy
}

func (l *Login) getPasswordComplexityPolicyByUserID(r *http.Request, userID string) *iam_model.PasswordComplexityPolicyView {
	user, err := l.query.GetUserByID(r.Context(), false, userID)
	if err != nil {
		logging.WithFields("userID", userID).OnError(err).Error("could not load user for password complexity policy")
		return nil
	}
	policy, err := l.authRepo.GetMyPasswordComplexityPolicy(setContext(r.Context(), user.ResourceOwner))
	logging.WithFields("orgID", user.ResourceOwner, "userID", userID).OnError(err).Error("could not load password complexity policy")
	return policy
}
