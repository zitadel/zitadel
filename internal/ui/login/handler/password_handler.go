package handler

import (
	"github.com/caos/zitadel/internal/domain"
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
)

const (
	tmplPassword = "password"
)

type passwordFormData struct {
	Password string `schema:"password"`
}

func (l *Login) renderPassword(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	data := l.getUserData(r, authReq, "Password", errType, errMessage)
	funcs := map[string]interface{}{
		"showPasswordReset": func() bool {
			if authReq.LoginPolicy != nil {
				return !authReq.LoginPolicy.HidePasswordReset
			}
			return true
		},
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplPassword], data, funcs)
}

func (l *Login) handlePasswordCheck(w http.ResponseWriter, r *http.Request) {
	data := new(passwordFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.authRepo.VerifyPassword(setContext(r.Context(), authReq.UserOrgID), authReq.ID, authReq.UserID, authReq.UserOrgID, data.Password, userAgentID, domain.BrowserInfoFromRequest(r))
	if err != nil {
		l.renderPassword(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
