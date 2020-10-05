package handler

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"net/http"
)

func (l *Login) getDefaultOrgIamPolicy(r *http.Request) (*iam_model.OrgIAMPolicyView, error) {
	return l.authRepo.GetDefaultOrgIamPolicy(r.Context())
}

func (l *Login) getOrgIamPolicy(r *http.Request, orgID string) (*iam_model.OrgIAMPolicyView, error) {
	if orgID == "" {
		return l.authRepo.GetDefaultOrgIamPolicy(r.Context())
	}
	return l.authRepo.GetOrgIamPolicy(r.Context(), orgID)
}

func (l *Login) getIDPConfigByID(r *http.Request, idpConfigID string) (*iam_model.IDPConfigView, error) {
	return l.authRepo.GetIDPConfigByID(r.Context(), idpConfigID)
}
