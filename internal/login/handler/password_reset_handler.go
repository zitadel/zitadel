package handler

import (
	"github.com/caos/zitadel/internal/renderer"
	"net/http"

)

const (
	tmplPasswordResetDone = "passwordresetdone"
)

func (l *Login) handlePasswordReset(w http.ResponseWriter, r *http.Request) {
	authSession, err := l.getAuthSession(r)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	err = l.service.Auth.RequestPasswordReset(r.Context(), authSession.UserSession.User.UserName)
	l.renderPasswordResetDone(w, r, authSession, err)
}

func (l *Login) renderPasswordResetDone(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	data := userData{
		baseData: l.getBaseData(r, authSession, "Password Reset Done", errType, errMessage),
		UserName: authSession.UserSession.User.UserName,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplPasswordResetDone], data, nil)
}
