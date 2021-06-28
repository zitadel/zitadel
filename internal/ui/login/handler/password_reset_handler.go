package handler

import (
	"github.com/caos/zitadel/internal/domain"
	"net/http"
)

const (
	tmplPasswordResetDone = "passwordresetdone"
)

func (l *Login) handlePasswordReset(w http.ResponseWriter, r *http.Request) {
	authReq, err := l.getAuthRequest(r)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	user, err := l.authRepo.UserByLoginName(setContext(r.Context(), authReq.UserOrgID), authReq.LoginName)
	if err != nil {
		l.renderPasswordResetDone(w, r, authReq, err)
		return
	}
	_, err = l.command.RequestSetPassword(setContext(r.Context(), authReq.UserOrgID), user.ID, authReq.UserOrgID, domain.NotificationTypeEmail)
	l.renderPasswordResetDone(w, r, authReq, err)
}

func (l *Login) renderPasswordResetDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	data := l.getUserData(r, authReq, "Password Reset Done", errID, errMessage)
	l.renderer.RenderTemplate(w, r, l.getTranslator(authReq), l.renderer.Templates[tmplPasswordResetDone], data, nil)
}
