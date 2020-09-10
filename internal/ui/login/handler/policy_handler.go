package handler

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"net/http"
)

func (l *Login) getOrgIamPolicy(r *http.Request, orgID string) (*org_model.OrgIAMPolicy, error) {
	if orgID == "" {
		return l.authRepo.GetDefaultOrgIamPolicy(r.Context())
	}
	return l.authRepo.GetOrgIamPolicy(r.Context(), orgID)
}

func (l *Login) getIDPConfigByID(r *http.Request, idpConfigID string) (*iam_model.IDPConfigView, error) {
	return l.authRepo.GetIDPConfigByID(r.Context(), idpConfigID)
}
