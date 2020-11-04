package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/duo-labs/webauthn/protocol"

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
	u2f, err := l.authRepo.AddMfaU2F(r.Context(), authReq.UserID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	data := &webAuthNData{
		userData:               l.getUserData(r, authReq, "Register U2F", errType, errMessage),
		CredentialCreationData: u2f.CredentialCreationDataString,
		SessionID:              u2f.SessionID,
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
	credentialData := new(protocol.CredentialCreationResponse)
	err = json.Unmarshal(credData, credentialData)
	if err = l.authRepo.VerifyMfaU2FSetup(r.Context(), authReq.UserID, data.SessionID, credentialData); err != nil {
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
		userData:               l.getUserData(r, authReq, "Register U2F", errType, errMessage),
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
	data, err := protocol.ParseCredentialRequestResponseBody(bytes.NewReader(credData))
	err = l.authRepo.VerifyMfaU2F(r.Context(), authReq.UserID, formData.SessionID, data)
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
