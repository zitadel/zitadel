package handler

import (
	"net/http"

	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	queryCode   = "code"
	queryUserID = "userID"

	tmplMailVerification = "mail_verification"
	tmplMailVerified     = "mail_verified"
)

type mailVerificationFormData struct {
	Code   string `schema:"code"`
	UserID string `schema:"userID"`
	Resend bool   `schema:"resend"`
}

type mailVerificationData struct {
	baseData
	profileData
	UserID string
}

func (l *Login) handleMailVerification(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue(queryUserID)
	code := r.FormValue(queryCode)
	if code != "" {
		l.checkMailCode(w, r, nil, userID, code)
		return
	}
	l.renderMailVerification(w, r, nil, userID, nil)
}

func (l *Login) handleMailVerificationCheck(w http.ResponseWriter, r *http.Request) {
	data := new(mailVerificationFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if !data.Resend {
		l.checkMailCode(w, r, authReq, data.UserID, data.Code)
		return
	}
	userOrg := login
	if authReq != nil {
		userOrg = authReq.UserOrgID
	}
	err = l.authRepo.ResendEmailVerificationMail(setContext(r.Context(), userOrg), data.UserID)
	l.renderMailVerification(w, r, authReq, data.UserID, err)
}

func (l *Login) checkMailCode(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, userID, code string) {
	userOrg := login
	if authReq != nil {
		userID = authReq.UserID
		userOrg = authReq.UserOrgID
	}
	err := l.command.VerifyHumanEmail(setContext(r.Context(), userOrg), userID, code, userOrg)
	if err != nil {
		l.renderMailVerification(w, r, authReq, userID, err)
		return
	}
	l.renderMailVerified(w, r, authReq)
}

func (l *Login) renderMailVerification(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, userID string, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	if userID == "" {
		userID = authReq.UserID
	}
	data := mailVerificationData{
		baseData:    l.getBaseData(r, authReq, "Mail Verification", errType, errMessage),
		UserID:      userID,
		profileData: l.getProfileData(authReq),
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMailVerification], data, nil)
}

func (l *Login) renderMailVerified(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest) {
	data := mailVerificationData{
		baseData:    l.getBaseData(r, authReq, "Mail Verified", "", ""),
		profileData: l.getProfileData(authReq),
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMailVerified], data, nil)
}
