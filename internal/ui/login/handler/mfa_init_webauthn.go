package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"

	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	tmplMfaU2FInit             = "mfainitu2f"
	tmplMfaU2FInitVerification = "mfainitu2fverification"
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
	user, err := l.authRepo.UserByID(r.Context(), authReq.UserID)
	if err != nil {

	}
	credential, sessionData, _ := l.webAuthN.BeginRegistration(user, protocol.Platform, protocol.VerificationDiscouraged, l.creds...)
	l.sessionData = *sessionData
	credentialData, _ := json.Marshal(credential)
	data := &webAuthNData{
		userData:               l.getUserData(r, authReq, "Register U2F", errType, errMessage),
		CredentialCreationData: base64.URLEncoding.EncodeToString(credentialData),
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
	user, err := l.authRepo.UserByID(r.Context(), authReq.UserID)
	if err != nil {

	}
	credential, _ := l.webAuthN.FinishRegistration(user, l.sessionData, data)
	if l.creds == nil {
		l.creds = make([]webauthn.Credential, 0)
	}
	l.creds = append(l.creds, *credential)
	if err = l.authRepo.AddU2FKey(r.Context(), authReq.UserID, data); err != nil {
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
	user, err := l.authRepo.UserByID(r.Context(), authReq.UserID)
	if err != nil {

	}
	credential, sessionData, _ := l.webAuthN.BeginLogin(user, protocol.VerificationDiscouraged, l.creds...)
	l.sessionData = *sessionData
	credentialData, _ := json.Marshal(credential)
	data := &webAuthNData{
		userData:               l.getUserData(r, authReq, "Register U2F", errType, errMessage),
		CredentialCreationData: base64.URLEncoding.EncodeToString(credentialData),
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
	user, err := l.authRepo.UserByID(r.Context(), authReq.UserID)
	if err != nil {

	}
	err = l.webAuthN.FinishLogin(user, l.sessionData, data, l.creds...)
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
