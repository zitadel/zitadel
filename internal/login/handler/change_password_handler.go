package handler

import (
	"net/http"

	"github.com/caos/citadel/login/internal/model"
)

const (
	tmplChangePassword     = "changepassword"
	tmplChangePasswordDone = "changepassworddone"
)

type changePasswordData struct {
	OldPassword string `schema:"old_password"`
	NewPassword string `schema:"new_password"`
}

func (l *Login) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	data := new(changePasswordData)
	authSession, err := l.getAuthSessionAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	err = l.service.Auth.ChangePassword(r.Context(), authSession.UserSession.User.UserID, datl.OldPassword, datl.NewPassword)
	if err != nil {
		l.renderChangePassword(w, r, authSession, err)
		return
	}
	l.renderChangePasswordDone(w, r, authSession)
}

func (l *Login) renderChangePassword(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	data := userData{
		baseData: l.getBaseData(r, authSession, "Change Password", errType, errMessage),
		UserName: authSession.UserSession.User.UserName,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplChangePassword], data, nil)
}

func (l *Login) renderChangePasswordDone(w http.ResponseWriter, r *http.Request, authReq *model.AuthSession) {
	var errType, errMessage string
	data := userData{
		baseData: l.getBaseData(r, authReq, "Password Change Done", errType, errMessage),
		UserName: authReq.UserSession.User.UserName,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplChangePasswordDone], data, nil)
}
