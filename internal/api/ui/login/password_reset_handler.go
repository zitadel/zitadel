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
	authReq, err := l.ensureAuthRequest(r)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	// We check if the user really exists or if it is just a placeholder or an unknown user.
	// In theory, we could also check for the unknownUserID constant. However, that could disclose
	// information about the existence of a user to an attacker if they check response times,
	// since those requests would take shorter than the ones for real users.
	user, err := l.query.GetUserByID(setContext(r.Context(), authReq.UserOrgID), true, authReq.UserID)
	if err != nil {
		if authReq.LoginPolicy.IgnoreUnknownUsernames && zerrors.IsNotFound(err) {
			err = nil
		}
		l.renderPasswordResetDone(w, r, authReq, err)
		return
	}
	_, err = l.command.RequestSetPassword(setContext(r.Context(), authReq.UserOrgID), user.ID, authReq.UserOrgID, domain.NotificationTypeEmail, authReq.ID)
	l.renderPasswordResetDone(w, r, authReq, err)
}

func (l *Login) renderPasswordResetDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, translator, "PasswordResetDone.Title", "PasswordResetDone.Description", err)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplPasswordResetDone], data, nil)
}
