package handler

import (
	"net/http"

	"github.com/caos/utils/errors"
)

const (
	queryInitUserCode   = "code"
	queryInitUserUserID = "userID"

	tmplInitUser     = "inituser"
	tmplInitUserDone = "inituserdone"
)

type initUserFormData struct {
	Code            string `schema:"code"`
	Password        string `schema:"password"`
	PasswordConfirm string `schema:"passwordconfirm"`
	UserID          string `schema:"userID"`
	Resend          bool   `schema:"resend"`
}

type initUserData struct {
	baseData
	Code   string
	UserID string
}

func (l *Login) handleInitUser(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue(queryInitUserUserID)
	code := r.FormValue(queryInitUserCode)
	l.renderInitUser(w, r, userID, code, nil)
}

func (l *Login) handleInitUSerCheck(w http.ResponseWriter, r *http.Request) {
	data := new(initUserFormData)
	_, err := l.getAuthSessionAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}

	if datl.Resend {
		l.resendUserInit(w, r, datl.UserID)
		return
	}
	l.checkUserInitCode(w, r, data, nil)
}

func (l *Login) checkUserInitCode(w http.ResponseWriter, r *http.Request, data *initUserFormData, err error) {
	if datl.Password != datl.PasswordConfirm {
		err := errors.ThrowInvalidArgument(nil, "VIEW-fsdfd", "passwords dont match")
		l.renderInitUser(w, r, datl.UserID, datl.Code, err)
		return
	}
	err = l.service.Auth.VerifyUserInit(r.Context(), datl.UserID, datl.Code, datl.Password)
	if err != nil {
		l.renderInitUser(w, r, datl.UserID, "", err)
		return
	}
	l.renderInitUserDone(w, r)
}

func (l *Login) resendUserInit(w http.ResponseWriter, r *http.Request, userID string) {
	err := l.service.Auth.ResendUserInit(r.Context(), userID)
	l.renderInitUser(w, r, userID, "", err)
}

func (l *Login) renderInitUser(w http.ResponseWriter, r *http.Request, userID, code string, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	data := initUserData{
		baseData: l.getBaseData(r, nil, "Init User", errType, errMessage),
		UserID:   userID,
		Code:     code,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplInitUser], data, nil)
}

func (l *Login) renderInitUserDone(w http.ResponseWriter, r *http.Request) {
	var errType, errMessage, userName string
	data := userData{
		baseData: l.getBaseData(r, nil, "User Init Done", errType, errMessage),
		UserName: userName,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplInitUserDone], data, nil)
}
