package handler

import (
	org_model "github.com/caos/zitadel/internal/org/model"
	"net/http"
)

func (l *Login) getOrgIamPolicy(r *http.Request, orgID string) (*org_model.OrgIAMPolicy, error) {
	return l.authRepo.GetOrgIamPolicy(r.Context(), orgID)
}
