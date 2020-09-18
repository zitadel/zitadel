package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
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

func (l *Login) renderRegisterOption(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	data := registerOptionData{
		baseData: l.getBaseData(r, authReq, "RegisterOption", errType, errMessage),
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplRegisterOption], data, nil)
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
