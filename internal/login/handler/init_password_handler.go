package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/errors"
	"net/http"
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
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	if data.Resend {
		l.resendPasswordSet(w, r, authReq)
		return
	}
	l.checkPWCode(w, r, authReq, data, nil)
}

func (l *Login) checkPWCode(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, data *initPasswordFormData, err error) {
	if data.Password != data.PasswordConfirm {
		err := errors.ThrowInvalidArgument(nil, "VIEW-KaGue", "passwords dont match")
		l.renderInitPassword(w, r, authReq, data.UserID, data.Code, err)
		return
	}
	err = l.authRepo.SetPassword(r.Context(), data.UserID, data.Code, data.Password)
	if err != nil {
		l.renderInitPassword(w, r, authReq, data.UserID, "", err)
		return
	}
	l.renderInitPasswordDone(w, r, authReq)
}

func (l *Login) resendPasswordSet(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest) {
	err := l.authRepo.RequestPasswordReset(r.Context(), authReq.UserName)
	l.renderInitPassword(w, r, authReq, authReq.UserID, "", err)
}

func (l *Login) renderInitPassword(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, userID, code string, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	if userID == "" && authReq != nil {
		userID = authReq.UserID
	}
	data := initPasswordData{
		baseData: l.getBaseData(r, authReq, "Init Password", errType, errMessage),
		UserID:   userID,
		Code:     code,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplInitPassword], data, nil)
}

func (l *Login) renderInitPasswordDone(w http.ResponseWriter, r *http.Request, authSession *model.AuthRequest) {
	var errType, errMessage, userName string
	//TODO: fill Username
	//if authSession != nil && authSession.UserSession != nil && authSession.UserSession.User != nil {
	//	userName = authSession.UserSession.User.UserName
	//}
	data := userData{
		baseData: l.getBaseData(r, authSession, "Password Init Done", errType, errMessage),
		UserName: userName,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplInitPasswordDone], data, nil)
}
