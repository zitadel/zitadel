package login

import (
	"net/http"

	"github.com/zitadel/logging"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	tmplLogin  = "login"
	queryOrgID = "orgID"
)

type loginData struct {
	LoginName string `schema:"loginName"`
	Register  bool   `schema:"register"`
}

func LoginLink(origin, orgID string) string {
	return externalLink(origin) + EndpointLogin + "?orgID=" + orgID
}

func externalLink(origin string) string {
	return origin + HandlerPrefix
}

func (l *Login) handleLogin(w http.ResponseWriter, r *http.Request) {
	authReq, err := l.getAuthRequest(r)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if authReq == nil {
		l.defaultRedirect(w, r)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) defaultRedirect(w http.ResponseWriter, r *http.Request) {
	orgID := r.FormValue(queryOrgID)
	policy, err := l.getLoginPolicy(r, orgID)
	logging.OnError(err).WithField("orgID", orgID).Error("error loading login policy")
	redirect := l.consolePath
	if policy != nil && policy.DefaultRedirectURI != "" {
		redirect = policy.DefaultRedirectURI
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (l *Login) handleLoginName(w http.ResponseWriter, r *http.Request) {
	authReq, err := l.getAuthRequest(r)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.renderLogin(w, r, authReq, nil)
}

func (l *Login) handleLoginNameCheck(w http.ResponseWriter, r *http.Request) {
	data := new(loginData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	if data.Register {
		if authReq != nil && authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowExternalIDP && authReq.AllowedExternalIDPs != nil && len(authReq.AllowedExternalIDPs) > 0 {
			l.handleRegisterOption(w, r)
			return
		}
		l.handleRegister(w, r)
		return
	}
	if authReq == nil {
		l.renderLogin(w, r, nil, zerrors.ThrowInvalidArgument(nil, "LOGIN-adrg3", "Errors.AuthRequest.NotFound"))
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	loginName := data.LoginName
	err = l.authRepo.CheckLoginName(r.Context(), authReq.ID, loginName, userAgentID)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) renderLogin(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	if err == nil && singleIDPAllowed(authReq) {
		l.handleIDP(w, r, authReq, authReq.AllowedExternalIDPs[0].IDPConfigID)
		return
	}
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, translator, "Login.Title", "Login.Description", errID, errMessage)
	funcs := map[string]interface{}{
		"hasUsernamePasswordLogin": func() bool {
			return authReq != nil && authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowUsernamePassword
		},
		"hasExternalLogin": func() bool {
			return authReq != nil && authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowExternalIDP && authReq.AllowedExternalIDPs != nil && len(authReq.AllowedExternalIDPs) > 0
		},
		"hasRegistration": func() bool {
			return authReq != nil && authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowRegister
		},
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplLogin], data, funcs)
}

func singleIDPAllowed(authReq *domain.AuthRequest) bool {
	return authReq != nil && authReq.LoginPolicy != nil && !authReq.LoginPolicy.AllowUsernamePassword && authReq.LoginPolicy.AllowExternalIDP && len(authReq.AllowedExternalIDPs) == 1
}
