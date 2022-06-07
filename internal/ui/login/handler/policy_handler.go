package handler

import (
	"net/http"

	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/query"
)

func (l *Login) getDefaultOrgIamPolicy(r *http.Request) (*query.OrgIAMPolicy, error) {
	return l.query.DefaultOrgIAMPolicy(r.Context())
}

func (l *Login) getOrgIamPolicy(r *http.Request, orgID string) (*query.OrgIAMPolicy, error) {
	if orgID == "" {
		return l.query.DefaultOrgIAMPolicy(r.Context())
	}
	return l.query.OrgIAMPolicyByOrg(r.Context(), false, orgID)
}

func (l *Login) getIDPConfigByID(r *http.Request, idpConfigID string) (*iam_model.IDPConfigView, error) {
	return l.authRepo.GetIDPConfigByID(r.Context(), idpConfigID)
}

func (l *Login) getLoginPolicy(r *http.Request, orgID string) (*query.LoginPolicy, error) {
	if orgID == "" {
		return l.query.DefaultLoginPolicy(r.Context())
	}
	return l.query.LoginPolicyByID(r.Context(), false, orgID)
}
