package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
)

const (
	tmplPassword = "password"
)

type passwordFormData struct {
	Password string `schema:"password"`
}

func (l *Login) renderPassword(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, translator, "Password.Title", "Password.Description", err)
	funcs := map[string]interface{}{
		"showPasswordReset": func() bool {
			if authReq.LoginPolicy != nil {
				return !authReq.LoginPolicy.HidePasswordReset
			}
			return true
		},
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplPassword], data, funcs)
}

func (l *Login) handlePasswordCheck(w http.ResponseWriter, r *http.Request) {
	data := new(passwordFormData)
	authReq, err := l.ensureAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	err = l.authRepo.VerifyPassword(setContext(r.Context(), authReq.UserOrgID), authReq.ID, authReq.UserID, authReq.UserOrgID, data.Password, authReq.AgentID, domain.BrowserInfoFromRequest(r))

	metadata, actionErr := l.runPostInternalAuthenticationActions(authReq, r, authMethodPassword, err)
	if err == nil && actionErr == nil && len(metadata) > 0 {
		err = l.bulkSetUserMetadata(r.Context(), authReq.UserID, authReq.UserOrgID, metadata)
	} else if actionErr != nil && err == nil {
		err = actionErr
	}

	if err != nil {
		if authReq.LoginPolicy.IgnoreUnknownUsernames {
			l.renderLogin(w, r, authReq, err)
			return
		}
		l.renderPassword(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
