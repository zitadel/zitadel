package handler

import (
	"encoding/base64"
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
	user_model "github.com/caos/zitadel/internal/user/model"
)

const (
	tmplMFAU2FInit = "mfainitu2f"
)

func (l *Login) renderRegisterU2F(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage, credentialData string
	var u2f *user_model.WebAuthNToken
	if err == nil {
		u2f, err = l.authRepo.AddMFAU2F(setContext(r.Context(), authReq.UserOrgID), authReq.UserID)
	}
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	if u2f != nil {
		credentialData = base64.RawURLEncoding.EncodeToString(u2f.CredentialCreationData)
	}
	data := &webAuthNData{
		userData:               l.getUserData(r, authReq, "Register WebAuthNToken", errType, errMessage),
		CredentialCreationData: credentialData,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMFAU2FInit], data, nil)
}

func (l *Login) handleRegisterU2F(w http.ResponseWriter, r *http.Request) {
	data := new(webAuthNFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	credData, err := base64.URLEncoding.DecodeString(data.CredentialData)
	if err != nil {
		l.renderRegisterU2F(w, r, authReq, err)
		return
	}

	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	if err = l.authRepo.VerifyMFAU2FSetup(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, data.Name, userAgentID, credData); err != nil {
		l.renderRegisterU2F(w, r, authReq, err)
		return
	}
	done := &mfaDoneData{
		MFAType: model.MFATypeU2F,
	}
	l.renderMFAInitDone(w, r, authReq, done)
}
