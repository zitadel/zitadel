package handler

import (
	org_model "github.com/caos/zitadel/internal/org/model"
	"net/http"
)

func (l *Login) getOrgIamPolicy(r *http.Request, orgID string) (*org_model.OrgIAMPolicy, error) {
	return l.authRepo.GetOrgIamPolicy(r.Context(), orgID)
}

//
//func (l *Login) getLoginPolicy(r *http.Request, authReq *model.AuthRequest) (*iam_model.LoginPolicyView, []*iam_model.IDPConfigView) {
//	orgID := l.getOrgID(authReq)
//	loginPolicy, err := l.authRepo.GetLoginPolicy(r.Context(), orgID)
//	if err != nil {
//		return nil, nil
//	}
//	if !loginPolicy.AllowExternalIDP {
//		return loginPolicy, nil
//	}
//	idpConfigs, err := l.authRepo.GetLoginPolicyIDPConfigs(r.Context(), orgID)
//	if err != nil {
//		return loginPolicy, nil
//	}
//	return loginPolicy, idpConfigs
//}
