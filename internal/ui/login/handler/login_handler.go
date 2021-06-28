package handler

import (
	"github.com/caos/zitadel/internal/domain"
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
)

const (
	tmplLogin = "login"
)

type loginData struct {
	LoginName string `schema:"loginName"`
	Register  bool   `schema:"register"`
}

func (l *Login) handleLogin(w http.ResponseWriter, r *http.Request) {
	authReq, err := l.getAuthRequest(r)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if authReq == nil {
		http.Redirect(w, r, l.zitadelURL, http.StatusFound)
		return
	}
	l.renderNextStep(w, r, authReq)
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
		if authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowExternalIDP && authReq.AllowedExternalIDPs != nil && len(authReq.AllowedExternalIDPs) > 0 {
			l.handleRegisterOption(w, r)
			return
		}
		l.handleRegister(w, r)
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
	data := l.getUserData(r, authReq, "Login", errID, errMessage)
	funcs := map[string]interface{}{
		"hasUsernamePasswordLogin": func() bool {
			return authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowUsernamePassword
		},
		"hasExternalLogin": func() bool {
			return authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowExternalIDP && authReq.AllowedExternalIDPs != nil && len(authReq.AllowedExternalIDPs) > 0
		},
	}
	l.renderer.RenderTemplate(w, r, l.getTranslator(authReq), l.renderer.Templates[tmplLogin], data, funcs)
}
