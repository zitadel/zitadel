package handler

import (
	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/auth_request/model"
	"net/http"
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
	authSession, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	if !data.Resend {
		l.checkMailCode(w, r, authSession, data.UserID, data.Code)
		return
	}
	//TODO: Check UserSession?
	if authSession == nil /*|| authSession.UserSession != nil && authSession.UserSession.User == nil*/ {
		err = l.authRepo.ResendEmailVerificationMail(r.Context(), data.UserID)
	} else {
		ctx := auth.SetCtxData(r.Context(), auth.CtxData{UserID: authSession.UserID, OrgID: "LOGIN"})
		err = l.authRepo.ResendMyEmailVerificationMail(ctx)
	}
	l.renderMailVerification(w, r, authSession, data.UserID, err)
}

func (l *Login) checkMailCode(w http.ResponseWriter, r *http.Request, authSession *model.AuthRequest, userID, code string) {
	var err error
	//TODO: Check UserSession
	if authSession != nil /* && authSession.UserSession != nil && authSession.UserSession.User != nil */ {
		ctx := auth.SetCtxData(r.Context(), auth.CtxData{UserID: authSession.UserID, OrgID: "LOGIN"})
		err = l.authRepo.VerifyMyEmail(ctx, code)
	} else {
		err = l.authRepo.VerifyEmail(r.Context(), userID, code)
	}
	if err != nil {
		l.renderMailVerification(w, r, authSession, userID, err)
		return
	}
	l.renderMailVerified(w, r, authSession)
}

func (l *Login) renderMailVerification(w http.ResponseWriter, r *http.Request, authSession *model.AuthRequest, userID string, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	//TODO: Check UserSession?
	if userID == "" /* && authSession != nil && authSession.UserSession != nil && authSession.UserSession.User != nil */ {
		userID = authSession.UserID
	}
	data := mailVerificationData{
		baseData: l.getBaseData(r, authSession, "Mail Verification", errType, errMessage),
		UserID:   userID,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMailVerification], data, nil)
}

func (l *Login) renderMailVerified(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest) {
	data := mailVerificationData{
		baseData: l.getBaseData(r, authReq, "Mail Verified", "", ""),
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMailVerified], data, nil)
}
