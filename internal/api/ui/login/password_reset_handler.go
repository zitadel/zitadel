package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	user, err := l.query.GetUserByLoginName(setContext(r.Context(), authReq.UserOrgID), true, authReq.LoginName)
	if err != nil {
		if authReq.LoginPolicy.IgnoreUnknownUsernames && zerrors.IsNotFound(err) {
			err = nil
		}
		l.renderPasswordResetDone(w, r, authReq, err)
		return
	}
	passwordCodeGenerator, err := l.query.InitEncryptionGenerator(r.Context(), domain.SecretGeneratorTypePasswordResetCode, l.userCodeAlg)
	if err != nil {
		if authReq.LoginPolicy.IgnoreUnknownUsernames && zerrors.IsNotFound(err) {
			err = nil
		}
		l.renderPasswordResetDone(w, r, authReq, err)
		return
	}
	_, err = l.command.RequestSetPassword(setContext(r.Context(), authReq.UserOrgID), user.ID, authReq.UserOrgID, domain.NotificationTypeEmail, passwordCodeGenerator, authReq.ID)
	l.renderPasswordResetDone(w, r, authReq, err)
}

func (l *Login) renderPasswordResetDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, translator, "PasswordResetDone.Title", "PasswordResetDone.Description", errID, errMessage)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplPasswordResetDone], data, nil)
}
