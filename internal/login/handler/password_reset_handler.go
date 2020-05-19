package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	"net/http"
)

const (
	tmplPasswordResetDone = "passwordresetdone"
)

func (l *Login) handlePasswordReset(w http.ResponseWriter, r *http.Request) {
	authSession, err := l.getAuthRequest(r)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	//TODO: Change UserID to UserName
	err = l.authRepo.RequestPasswordReset(r.Context(), authSession.UserID)
	l.renderPasswordResetDone(w, r, authSession, err)
}

func (l *Login) renderPasswordResetDone(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	data := userData{
		baseData: l.getBaseData(r, authReq, "Password Reset Done", errType, errMessage),
		//TODO: Fill Username
		//UserName: authReq.UserName,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplPasswordResetDone], data, nil)
}
