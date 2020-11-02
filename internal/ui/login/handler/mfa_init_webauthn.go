package handler

import (
	"bytes"
	"encoding/base64"
	"net/http"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"

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
}

type webAuthNFormData struct {
	CredentialData string `schema:"credentialData"`
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
	if err = l.webAuthnCookieHandler.SetEncryptedCookie(w, webauthnRegisterSession, u2f.SessionData); err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	data := &webAuthNData{
		userData:               l.getUserData(r, authReq, "Register U2F", errType, errMessage),
		CredentialCreationData: u2f.CredentialCreationDataString,
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
	authReq, data, err := l.getAuthRequestAndParseWebAuthNData(r)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	sessionData := new(webauthn.SessionData)
	if err = l.webAuthnCookieHandler.GetEncryptedCookieValue(r, webauthnRegisterSession, sessionData); err != nil {

	}
	if err = l.authRepo.VerifyMfaU2FSetup(r.Context(), authReq.UserID, *sessionData, data); err != nil {
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
	err = l.authRepo.VerifyMfaU2F(r.Context(), authReq.UserID, l.sessionData, data)
	if err != nil {

	}
	done := &mfaDoneData{
		//MfaType: nil,
	}
	l.renderMfaInitDone(w, r, authReq, done)
}

func (l *Login) getAuthRequestAndParseWebAuthNData(r *http.Request) (*model.AuthRequest, *protocol.ParsedCredentialCreationData, error) {
	formData := new(webAuthNFormData)
	authReq, err := l.getAuthRequestAndParseData(r, formData)
	if err != nil {
		return authReq, nil, err
	}
	credData, err := base64.URLEncoding.DecodeString(formData.CredentialData)
	if err != nil {
		return authReq, nil, err
	}
	data, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(credData))
	return authReq, data, err
}
