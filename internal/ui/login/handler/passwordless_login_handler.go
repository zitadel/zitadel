package handler

import (
	"encoding/base64"
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	tmplPasswordLessVerification = "passwordlessverification"
)

func (l *Login) renderPasswordLessVerification(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest) {
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	webAuthNLogin, err := l.authRepo.BeginPasswordlessLogin(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.ID, userAgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	data := &webAuthNData{
		userData:               l.getUserData(r, authReq, "Login PasswordLess", "", ""),
		CredentialCreationData: base64.URLEncoding.EncodeToString(webAuthNLogin.CredentialAssertionData),
		SessionID:              webAuthNLogin.Challenge,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplPasswordLessVerification], data, nil)
}

func (l *Login) handlePasswordLessVerification(w http.ResponseWriter, r *http.Request) {
	formData := new(webAuthNFormData)
	authReq, err := l.getAuthRequestAndParseData(r, formData)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if formData.Resend {
		l.renderPasswordLessVerification(w, r, authReq)
		return
	}
	credData, err := base64.URLEncoding.DecodeString(formData.CredentialData)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.authRepo.VerifyPasswordless(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.ID, userAgentID, credData, model.BrowserInfoFromRequest(r))
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
