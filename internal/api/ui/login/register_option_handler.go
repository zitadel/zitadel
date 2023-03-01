package login

import (
	"net/http"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

const (
	tmplRegisterOption = "registeroption"
)

type registerOptionFormData struct {
	UsernamePassword bool `schema:"usernamepassword"`
}

type registerOptionData struct {
	baseData
}

func (l *Login) handleRegisterOption(w http.ResponseWriter, r *http.Request) {
	data := new(registerOptionFormData)
	authRequest, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authRequest, err)
		return
	}
	l.renderRegisterOption(w, r, authRequest, nil)
}

func (l *Login) renderRegisterOption(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	allowed := registrationAllowed(authReq)
	externalAllowed := externalRegistrationAllowed(authReq)
	if err == nil {
		// if only external allowed with a single idp then use that
		if !allowed && externalAllowed && len(authReq.AllowedExternalIDPs) == 1 {
			l.handleIDP(w, r, authReq, authReq.AllowedExternalIDPs[0].IDPConfigID)
			return
		}
		// if only direct registration is allowed, show the form
		if allowed && !externalAllowed {
			data := l.passLoginHintToRegistration(r, authReq)
			l.renderRegister(w, r, authReq, data, nil)
			return
		}
	}
	translator := l.getTranslator(r.Context(), authReq)
	data := registerOptionData{
		baseData: l.getBaseData(r, authReq, "RegisterOption.Title", "RegisterOption.Description", errID, errMessage),
	}
	funcs := map[string]interface{}{
		"hasRegistration": func() bool {
			return allowed
		},
		"hasExternalLogin": func() bool {
			return externalAllowed
		},
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplRegisterOption], data, funcs)
}

func (l *Login) handleRegisterOptionCheck(w http.ResponseWriter, r *http.Request) {
	data := new(registerOptionFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if data.UsernamePassword {
		l.handleRegister(w, r)
		return
	}
	l.handleRegisterOption(w, r)
}

func registrationAllowed(authReq *domain.AuthRequest) bool {
	return authReq != nil && authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowRegister && authReq.LoginPolicy.AllowUsernamePassword
}

func externalRegistrationAllowed(authReq *domain.AuthRequest) bool {
	return authReq != nil && authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowExternalIDP && authReq.AllowedExternalIDPs != nil && len(authReq.AllowedExternalIDPs) > 0
}

func (l *Login) passLoginHintToRegistration(r *http.Request, authReq *domain.AuthRequest) *registerFormData {
	data := &registerFormData{}
	if authReq == nil {
		return data
	}
	data.Email = authReq.LoginHint
	domainPolicy, err := l.getOrgDomainPolicy(r, authReq.RequestedOrgID)
	if err != nil {
		logging.WithFields("authRequest", authReq.ID, "org", authReq.RequestedOrgID).Error("unable to load domain policy for registration loginHint")
		return data
	}
	data.Username = authReq.LoginHint
	if !domainPolicy.UserLoginMustBeDomain {
		return data
	}
	searchQuery, err := query.NewOrgDomainOrgIDSearchQuery(authReq.RequestedOrgID)
	if err != nil {
		logging.WithFields("authRequest", authReq.ID, "org", authReq.RequestedOrgID).Error("unable to search query for registration loginHint")
		return data
	}
	domains, err := l.query.SearchOrgDomains(r.Context(), &query.OrgDomainSearchQueries{Queries: []query.SearchQuery{searchQuery}}, false)
	if err != nil {
		logging.WithFields("authRequest", authReq.ID, "org", authReq.RequestedOrgID).Error("unable to load domains for registration loginHint")
		return data
	}
	for _, orgDomain := range domains.Domains {
		if orgDomain.IsVerified && strings.HasSuffix(authReq.LoginHint, "@"+orgDomain.Domain) {
			data.Username = strings.TrimSuffix(authReq.LoginHint, "@"+orgDomain.Domain)
			return data
		}
	}
	return data
}
