package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/query"
)

func (l *Login) getDefaultDomainPolicy(r *http.Request) (*query.DomainPolicy, error) {
	return l.query.DefaultDomainPolicy(r.Context())
}

func (l *Login) getOrgDomainPolicy(r *http.Request, orgID string) (*query.DomainPolicy, error) {
	if orgID == "" {
		return l.query.DefaultDomainPolicy(r.Context())
	}
	return l.query.DomainPolicyByOrg(r.Context(), false, orgID, false)
}

func (l *Login) getIDPByID(r *http.Request, id string) (*query.IDPTemplate, error) {
	return l.query.IDPTemplateByID(r.Context(), false, id, false)
}

func (l *Login) getLoginPolicy(r *http.Request, orgID string) (*query.LoginPolicy, error) {
	if orgID == "" {
		return l.query.DefaultLoginPolicy(r.Context())
	}
	return l.query.LoginPolicyByID(r.Context(), false, orgID, false)
}

func (l *Login) getLabelPolicy(r *http.Request, orgID string) (*query.LabelPolicy, error) {
	if orgID == "" {
		return l.query.DefaultActiveLabelPolicy(r.Context())
	}
	return l.query.ActiveLabelPolicyByOrg(r.Context(), orgID, false)
}
