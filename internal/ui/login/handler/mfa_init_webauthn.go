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
	tmplMfaU2FInitVerify = "mfainitu2f"
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
	credential, _, _ := l.webAuthN.BeginRegistration(user, protocol.Platform, protocol.VerificationDiscouraged)
	credentialData, _ := json.Marshal(credential)
	data := &webAuthNData{
		userData:               l.getUserData(r, authReq, "Register U2F", errType, errMessage),
		CredentialCreationData: base64.URLEncoding.EncodeToString(credentialData),
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMfaU2FInitVerify], data, nil)
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
	data := new(protocol.ParsedCredentialCreationData)
	authReq, err := l.getAuthRequestAndParseWebAuthNData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if err = l.authRepo.AddU2FKey(r.Context(), authReq.UserID, data); err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	done := &mfaDoneData{
		//MfaType: nil,
	}
	l.renderMfaInitDone(w, r, authReq, done)
}

func (l *Login) getAuthRequestAndParseWebAuthNData(r *http.Request, data interface{}) (*model.AuthRequest, error) {
	formData := new(webAuthNFormData)
	authReq, err := l.getAuthRequestAndParseData(r, formData)
	if err != nil {
		return authReq, err
	}
	credData, err := base64.URLEncoding.DecodeString(formData.CredentialData)
	if err != nil {
		return authReq, err
	}
	data, err = protocol.ParseCredentialCreationResponseBody(bytes.NewReader(credData))
	return authReq, err
}
