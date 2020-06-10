package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"net/http"
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
	l.renderInitUser(w, r, nil, userID, code, nil)
}

func (l *Login) handleInitUserCheck(w http.ResponseWriter, r *http.Request) {
	data := new(initUserFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}

	if data.Resend {
		l.resendUserInit(w, r, authReq, data.UserID)
		return
	}
	l.checkUserInitCode(w, r, authReq, data, nil)
}

func (l *Login) checkUserInitCode(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, data *initUserFormData, err error) {
	if data.Password != data.PasswordConfirm {
		err := caos_errs.ThrowInvalidArgument(nil, "VIEW-fsdfd", "Errors.User.Password.ConfirmationWrong")
		l.renderInitUser(w, r, nil, data.UserID, data.Code, err)
		return
	}
	userOrgID := login
	if authReq != nil {
		userOrgID = authReq.UserOrgID
	}
	err = l.authRepo.VerifyInitCode(setContext(r.Context(), userOrgID), data.UserID, data.Code, data.Password)
	if err != nil {
		l.renderInitUser(w, r, nil, data.UserID, "", err)
		return
	}
	l.renderInitUserDone(w, r, nil)
}

func (l *Login) resendUserInit(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, userID string) {
	userOrgID := login
	if authReq != nil {
		userOrgID = authReq.UserOrgID
	}
	err := l.authRepo.ResendInitVerificationMail(setContext(r.Context(), userOrgID), userID)
	l.renderInitUser(w, r, authReq, userID, "", err)
}

func (l *Login) renderInitUser(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, userID, code string, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	if authReq != nil {
		userID = authReq.UserID
	}
	data := initUserData{
		baseData: l.getBaseData(r, nil, "Init User", errType, errMessage),
		UserID:   userID,
		Code:     code,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplInitUser], data, nil)
}

func (l *Login) renderInitUserDone(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest) {
	var errType, errMessage, userName string
	if authReq != nil {
		userName = authReq.UserName
	}
	data := userData{
		baseData: l.getBaseData(r, authReq, "User Init Done", errType, errMessage),
		UserName: userName,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplInitUserDone], data, nil)
}
