package login

import (
	"net/http"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/query"
)

func (l *Login) getDefaultDomainPolicy(r *http.Request) (*query.DomainPolicy, error) {
	return l.query.DefaultDomainPolicy(r.Context())
}

func (l *Login) getOrgDomainPolicy(r *http.Request, orgID string) (*query.DomainPolicy, error) {
	if orgID == "" {
		return l.query.DefaultDomainPolicy(r.Context())
	}
	return l.query.DomainPolicyByOrg(r.Context(), orgID)
}

func (l *Login) getIDPConfigByID(r *http.Request, idpConfigID string) (*iam_model.IDPConfigView, error) {
	return l.authRepo.GetIDPConfigByID(r.Context(), idpConfigID)
}
