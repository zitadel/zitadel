package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
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
	err = l.authRepo.RequestPasswordReset(setContext(r.Context(), authReq.UserOrgID), authReq.LoginName)
	l.renderPasswordResetDone(w, r, authReq, err)
}

func (l *Login) renderPasswordResetDone(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	data := userData{
		baseData:  l.getBaseData(r, authReq, "Password Reset Done", errType, errMessage),
		LoginName: authReq.LoginName,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplPasswordResetDone], data, nil)
}
