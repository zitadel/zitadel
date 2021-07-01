package handler

import (
	"net/http"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

func (l *Login) getDefaultPrivacyPolicy(r *http.Request) (*iam_model.PrivacyPolicyView, error) {
	policy, err := l.authRepo.GetDefaultPrivacyPolicy(r.Context())
	if err != nil {
		return nil, err
	}
	return policy, nil
}
