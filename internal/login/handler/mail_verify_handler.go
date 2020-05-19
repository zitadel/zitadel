package handler

import (
	"context"
	"net/http"

	"github.com/caos/citadel/login/internal/model"
	"github.com/caos/citadel/utils/auth"
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
	authSession, err := l.getAuthSessionAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	if !datl.Resend {
		l.checkMailCode(w, r, authSession, datl.UserID, datl.Code)
		return
	}
	if authSession == nil || authSession.UserSession != nil && authSession.UserSession.User == nil {
		err = l.service.Auth.ResendEmailVerificationMail(r.Context(), datl.UserID)
	} else {
		ctx := context.WithValue(r.Context(), auth.CtxKeyData{}, &auth.CtxData{UserID: authSession.UserSession.User.UserID, OrgID: "LOGIN"})
		err = l.service.Auth.ResendMyEmailVerificationMail(ctx)
	}
	l.renderMailVerification(w, r, authSession, datl.UserID, err)
}

func (l *Login) checkMailCode(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, userID, code string) {
	var err error
	if authSession != nil && authSession.UserSession != nil && authSession.UserSession.User != nil {
		ctx := context.WithValue(r.Context(), auth.CtxKeyData{}, &auth.CtxData{UserID: authSession.UserSession.User.UserID, OrgID: "LOGIN"})
		err = l.service.Auth.VerifyMyEmail(ctx, code)
	} else {
		err = l.service.Auth.VerifyEmail(r.Context(), userID, code)
	}
	if err != nil {
		l.renderMailVerification(w, r, authSession, userID, err)
		return
	}
	l.renderMailVerified(w, r, authSession)
}

func (l *Login) renderMailVerification(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, userID string, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	if userID == "" && authSession != nil && authSession.UserSession != nil && authSession.UserSession.User != nil {
		userID = authSession.UserSession.User.UserID
	}
	data := mailVerificationData{
		baseData: l.getBaseData(r, authSession, "Mail Verification", errType, errMessage),
		UserID:   userID,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMailVerification], data, nil)
}

func (l *Login) renderMailVerified(w http.ResponseWriter, r *http.Request, authReq *model.AuthSession) {
	data := mailVerificationData{
		baseData: l.getBaseData(r, authReq, "Mail Verified", "", ""),
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMailVerified], data, nil)
}
