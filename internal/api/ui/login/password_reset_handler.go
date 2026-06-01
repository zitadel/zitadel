package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
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
	userID := authReq.UserID
	if userID == "" {
		var user *query.User
		user, err = l.query.GetUserByLoginName(setContext(r.Context(), authReq.UserOrgID), true, authReq.LoginName)
		if err != nil {
			if authReq.LoginPolicy.IgnoreUnknownUsernames && zerrors.IsNotFound(err) {
				l.renderPasswordResetDone(w, r, authReq, nil)
				return
			}
			l.renderError(w, r, authReq, err)
			return
		}
		userID = user.ID
	}
	_, err = l.command.RequestSetPassword(setContext(r.Context(), authReq.UserOrgID), userID, authReq.UserOrgID, domain.NotificationTypeEmail, authReq.ID)
	l.renderPasswordResetDone(w, r, authReq, err)
}

func (l *Login) renderPasswordResetDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, translator, "PasswordResetDone.Title", "PasswordResetDone.Description", err)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplPasswordResetDone], data, nil)
}
