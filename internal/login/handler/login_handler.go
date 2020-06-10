package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"net/http"
)

const (
	tmplLogin = "login"
)

type loginData struct {
	UserName string `schema:"username"`
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

func (l *Login) handleUsername(w http.ResponseWriter, r *http.Request) {
	authSession, err := l.getAuthRequest(r)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	l.renderLogin(w, r, authSession, nil)
}

func (l *Login) handleUsernameCheck(w http.ResponseWriter, r *http.Request) {
	data := new(loginData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	err = l.authRepo.CheckUsername(r.Context(), authReq.ID, data.UserName)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) renderLogin(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		if caos_errs.IsNotFound(err) {
			errMessage = l.renderer.LocalizeFromRequest(r, "Errors.UserNotFound", nil)
		} else {
			errMessage = l.getErrorMessage(err)
		}
	}
	data := userData{
		baseData: l.getBaseData(r, authReq, "Login", errType, errMessage),
		UserName: authReq.UserName,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplLogin], data, nil)
}
