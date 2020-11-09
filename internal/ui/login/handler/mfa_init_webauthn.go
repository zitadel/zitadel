package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	tmplMfaU2FInit             = "mfainitu2f"
	tmplMfaU2FInitVerification = "mfainitu2fverification"
)

func (l *Login) renderRegisterU2F(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	u2f, err := l.authRepo.AddMfaU2F(setContext(r.Context(), authReq.UserOrgID), authReq.UserID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	data := &webAuthNData{
		userData:               l.getUserData(r, authReq, "Register WebAuthNToken", errType, errMessage),
		CredentialCreationData: base64.RawURLEncoding.EncodeToString(u2f.CredentialCreationData),
		SessionID:              u2f.WebAuthNTokenID,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMfaU2FInit], data, nil)
}

func (l *Login) handleRegisterU2F(w http.ResponseWriter, r *http.Request) {
	data := new(webAuthNFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if data.Resend {
		l.renderRegisterU2F(w, r, authReq, nil)
		return
	}
	credData, err := base64.URLEncoding.DecodeString(data.CredentialData)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	if err = l.authRepo.VerifyMfaU2FSetup(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, credData); err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	//done := &mfaDoneData{
	//	//MfaType: nil,
	//}
	l.renderLoginU2F(w, r, authReq, nil)
}

//TODO: remove
func (l *Login) handleRegisterU2FTest(w http.ResponseWriter, r *http.Request) {
	authReq, err := l.getAuthRequest(r)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.renderRegisterU2F(w, r, authReq, nil)
}
