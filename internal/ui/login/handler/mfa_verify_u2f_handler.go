package handler

import (
	"encoding/base64"
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	tmplMfaU2FInitVerification = "mfainitu2fverification"
)

func (l *Login) renderLoginU2F(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	webAuthNLogin, err := l.authRepo.BeginMfaU2FLogin(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.ID, userAgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	//l.sessionData = *sessionData
	data := &webAuthNData{
		userData:               l.getUserData(r, authReq, "Register WebAuthNToken", errType, errMessage),
		CredentialCreationData: base64.URLEncoding.EncodeToString(webAuthNLogin.CredentialAssertionData),
		SessionID:              webAuthNLogin.Challenge,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMfaU2FInitVerification], data, nil)
}

func (l *Login) handleLoginU2F(w http.ResponseWriter, r *http.Request) {
	formData := new(webAuthNFormData)
	authReq, err := l.getAuthRequestAndParseData(r, formData)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	credData, err := base64.URLEncoding.DecodeString(formData.CredentialData)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.authRepo.VerifyMfaU2F(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.ID, userAgentID, credData, model.BrowserInfoFromRequest(r))
	if err != nil {

	}
	done := &mfaDoneData{
		//MfaType: nil,
	}
	l.renderMfaInitDone(w, r, authReq, done)
}

//TODO: remove
func (l *Login) handleLoginU2FTest(w http.ResponseWriter, r *http.Request) {
	authReq, err := l.getAuthRequest(r)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.renderLoginU2F(w, r, authReq, nil)
}
