package handler

import (
	"net/http"

	"github.com/caos/citadel/login/internal/model"
	"github.com/caos/utils/errors"
)

const (
	queryInitPWCode   = "code"
	queryInitPWUserID = "userID"

	tmplInitPassword     = "initpassword"
	tmplInitPasswordDone = "initpassworddone"
)

type initPasswordFormData struct {
	Code            string `schema:"code"`
	Password        string `schema:"password"`
	PasswordConfirm string `schema:"passwordconfirm"`
	UserID          string `schema:"userID"`
	Resend          bool   `schema:"resend"`
}

type initPasswordData struct {
	baseData
	Code   string
	UserID string
}

func (l *Login) handleInitPassword(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue(queryInitPWUserID)
	code := r.FormValue(queryInitPWCode)
	l.renderInitPassword(w, r, nil, userID, code, nil)
}

func (l *Login) handleInitPasswordCheck(w http.ResponseWriter, r *http.Request) {
	data := new(initPasswordFormData)
	authReq, err := l.getAuthSessionAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	if datl.Resend {
		l.resendPasswordSet(w, r, authReq)
		return
	}
	l.checkPWCode(w, r, authReq, data, nil)
}

func (l *Login) checkPWCode(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, data *initPasswordFormData, err error) {
	if datl.Password != datl.PasswordConfirm {
		err := errors.ThrowInvalidArgument(nil, "VIEW-KaGue", "passwords dont match")
		l.renderInitPassword(w, r, authSession, datl.UserID, datl.Code, err)
		return
	}
	err = l.service.Auth.PasswordReset(r.Context(), datl.UserID, datl.Code, datl.Password)
	if err != nil {
		l.renderInitPassword(w, r, authSession, datl.UserID, "", err)
		return
	}
	l.renderInitPasswordDone(w, r, authSession)
}

func (l *Login) resendPasswordSet(w http.ResponseWriter, r *http.Request, authReq *model.AuthSession) {
	err := l.service.Auth.RequestPasswordReset(r.Context(), authReq.UserSession.User.UserName)
	l.renderInitPassword(w, r, authReq, authReq.UserSession.User.UserID, "", err)
}

func (l *Login) renderInitPassword(w http.ResponseWriter, r *http.Request, authReq *model.AuthSession, userID, code string, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	if userID == "" && authReq != nil && authReq.UserSession != nil && authReq.UserSession.User != nil {
		userID = authReq.UserSession.User.UserID
	}
	data := initPasswordData{
		baseData: l.getBaseData(r, authReq, "Init Password", errType, errMessage),
		UserID:   userID,
		Code:     code,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplInitPassword], data, nil)
}

func (l *Login) renderInitPasswordDone(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession) {
	var errType, errMessage, userName string
	if authSession != nil && authSession.UserSession != nil && authSession.UserSession.User != nil {
		userName = authSession.UserSession.User.UserName
	}
	data := userData{
		baseData: l.getBaseData(r, authSession, "Password Init Done", errType, errMessage),
		UserName: userName,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplInitPasswordDone], data, nil)
}
