package handler

import (
	"github.com/caos/zitadel/internal/domain"
	"net/http"
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
	data := registerOptionData{
		baseData: l.getBaseData(r, authReq, "RegisterOption", errID, errMessage),
	}
	funcs := map[string]interface{}{
		"hasExternalLogin": func() bool {
			return authReq.LoginPolicy.AllowExternalIDP && authReq.AllowedExternalIDPs != nil && len(authReq.AllowedExternalIDPs) > 0
		},
	}
	translator := l.getTranslator(authReq)
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
