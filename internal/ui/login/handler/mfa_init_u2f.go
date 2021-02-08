package handler

import (
	"encoding/base64"
	"github.com/caos/zitadel/internal/v2/domain"
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	tmplMFAU2FInit = "mfainitu2f"
)

type u2fInitData struct {
	webAuthNData
	MFAType model.MFAType
}

func (l *Login) renderRegisterU2F(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errType, errMessage, credentialData string
	var u2f *domain.WebAuthNToken
	if err == nil {
		u2f, err = l.command.HumanAddU2FSetup(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, true)
	}
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	if u2f != nil {
		credentialData = base64.RawURLEncoding.EncodeToString(u2f.CredentialCreationData)
	}
	data := &u2fInitData{
		webAuthNData: webAuthNData{
			userData:               l.getUserData(r, authReq, "Register WebAuthNToken", errType, errMessage),
			CredentialCreationData: credentialData,
		},
		MFAType: model.MFATypeU2F,
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
	if err = l.command.HumanVerifyU2FSetup(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, data.Name, userAgentID, credData); err != nil {
		l.renderRegisterU2F(w, r, authReq, err)
		return
	}
	done := &mfaDoneData{
		MFAType: domain.MFATypeU2F,
	}
	l.renderMFAInitDone(w, r, authReq, done)
}
