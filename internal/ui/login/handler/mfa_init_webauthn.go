package handler

import (
	"encoding/base64"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"net/http"

	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	tmplMfaU2FInit             = "mfainitu2f"
	tmplMfaU2FInitVerification = "mfainitu2fverification"

	webauthnRegisterSession = "register-session"
)

type webAuthNData struct {
	userData
	CredentialCreationData string
	SessionID              string
}

type webAuthNFormData struct {
	CredentialData string `schema:"credentialData"`
	SessionID      string `schema:"sessionID"`
}

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

func (l *Login) handleRegisterU2FTest(w http.ResponseWriter, r *http.Request) {
	authReq, err := l.getAuthRequest(r)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.renderRegisterU2F(w, r, authReq, nil)
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

func (l *Login) renderLoginU2F(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	credential, sessionData, err := l.authRepo.BeginMfaU2FLogin(r.Context(), authReq.UserID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.sessionData = *sessionData
	data := &webAuthNData{
		userData:               l.getUserData(r, authReq, "Register WebAuthNToken", errType, errMessage),
		CredentialCreationData: credential,
		SessionID:              sessionData.Challenge,
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
	err = l.authRepo.VerifyMfaU2F(r.Context(), authReq.UserID, formData.SessionID, authReq.ID, userAgentID, credData, model.BrowserInfoFromRequest(r))
	if err != nil {

	}
	done := &mfaDoneData{
		//MfaType: nil,
	}
	l.renderMfaInitDone(w, r, authReq, done)
}

//
//func (l *Login) getAuthRequestAndParseWebAuthNData(r *http.Request, response interface{}) (*model.AuthRequest, error) {
//	formData := new(webAuthNFormData)
//	authReq, err := l.getAuthRequestAndParseData(r, formData)
//	if err != nil {
//		return authReq, err
//	}
//	credData, err := base64.URLEncoding.DecodeString(formData.CredentialData)
//	if err != nil {
//		return authReq, err
//	}
//	err = json.Unmarshal(credData, response)
//	return authReq, err
//}
