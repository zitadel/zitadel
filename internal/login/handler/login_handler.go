package handler

import (
	"net/http"

	"github.com/caos/citadel/login/internal/model"
)

const (
	tmplLogin = "login"
)

type loginData struct {
	UserName string `schema:"username"`
	Register bool   `schema:"register"`
}

func (l *Login) handleLogin(w http.ResponseWriter, r *http.Request) {
	authSession, err := l.getAuthSession(r)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	if authSession == nil {
		http.Redirect(w, r, l.citadelURL, http.StatusFound)
		return
	}
	l.renderNextStep(w, r, authSession)
}

func (l *Login) handleUsername(w http.ResponseWriter, r *http.Request) {
	authSession, err := l.getAuthSession(r)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	l.renderLogin(w, r, authSession, nil)
}

func (l *Login) handleUsernameCheck(w http.ResponseWriter, r *http.Request) {
	data := new(loginData)
	authSession, err := l.getAuthSessionAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	if datl.Register {
		l.renderRegister(w, r, authSession, nil, nil)
		return
	}
	browserInfo := &model.BrowserInformation{RemoteIP: &model.IP{}} //TODO: impl
	authSession, err = l.service.Auth.VerifyUser(r.Context(), authSession, datl.UserName, browserInfo)
	if err != nil {
		l.renderLogin(w, r, authSession, err)
		return
	}
	l.renderNextStep(w, r, authSession)
}

func (l *Login) renderLogin(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, err error) {
	var errType, errMessage, username string
	if err != nil {
		errMessage = err.Error()
	}
	if authSession.PossibleSteps[0].LoginData != nil {
		errMessage = authSession.PossibleSteps[0].LoginDatl.ErrMsg
	}
	if authSession.UserSession != nil && authSession.UserSession.User != nil {
		username = authSession.UserSession.User.UserName
	}
	data := userData{
		baseData: l.getBaseData(r, authSession, "Login", errType, errMessage),
		UserName: username,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplLogin], data, nil)
}
