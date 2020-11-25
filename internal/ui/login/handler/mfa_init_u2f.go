package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	tmplMfaU2FInit = "mfainitu2f"
)

func (l *Login) renderRegisterU2F(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	u2f, err := l.authRepo.AddMFAU2F(setContext(r.Context(), authReq.UserOrgID), authReq.UserID)
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
	if data.Recreate {
		l.renderRegisterU2F(w, r, authReq, nil)
		return
	}
	credData, err := base64.URLEncoding.DecodeString(data.CredentialData)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	if err = l.authRepo.VerifyMFAU2FSetup(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, data.Name, credData); err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	done := &mfaDoneData{
		MfaType: model.MFATypeU2F,
	}
	l.renderMfaInitDone(w, r, authReq, done)
}
