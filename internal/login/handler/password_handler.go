package handler

import (
	"net/http"

	"github.com/caos/citadel/login/internal/model"
)

const (
	tmplPassword = "password"
)

type passwordData struct {
	Password string `schema:"password"`
}

func (l *Login) renderPassword(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, passwordData *model.PasswordData) {
	var errType, errMessage string
	if passwordData != nil {
		errMessage = passwordData.ErrMsg
	}
	data := userData{
		baseData: l.getBaseData(r, authSession, "Password", errType, errMessage),
		UserName: authSession.UserSession.User.UserName,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplPassword], data, nil)
}

func (l *Login) handlePasswordCheck(w http.ResponseWriter, r *http.Request) {
	data := new(passwordData)
	authSession, err := l.getAuthSessionAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	browserInfo := &model.BrowserInformation{RemoteIP: &model.IP{}} //TODO: impl
	authSession, err = l.service.Auth.VerifyPassword(r.Context(), authSession, data.Password, browserInfo)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	l.renderNextStep(w, r, authSession)
}
