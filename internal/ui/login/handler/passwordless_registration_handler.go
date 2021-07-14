package handler

import (
	"encoding/base64"
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/domain"
)

const (
	tmplPasswordlessRegistration        = "passwordlessregistration"
	queryPasswordlessRegistrationCode   = "code"
	queryPasswordlessRegistrationUserID = "userID"
)

type passwordlessRegistrationData struct {
	webAuthNData
	Code   string
	UserID string
}

type passwordlessRegistrationFormData struct {
	webAuthNFormData
	Code      string `schema:"code"`
	UserID    string `schema:"userID"`
	TokenName string `schema:"name"`
	Resend    bool   `schema:"resend"`
}

func (l *Login) handlePasswordlessRegistration(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue(queryPasswordlessRegistrationUserID)
	code := r.FormValue(queryPasswordlessRegistrationCode)
	l.renderPasswordlessRegistration(w, r, nil, userID, code, nil)
}

func (l *Login) renderPasswordlessRegistration(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID, code string, err error) {
	var errID, errMessage, credentialData, userOrgID string
	if authReq != nil {
		userID = authReq.UserID
		userOrgID = authReq.UserOrgID
	}
	var webAuthNToken *domain.WebAuthNToken
	if err == nil {
		webAuthNToken, err = l.authRepo.BeginPasswordlessSetup(setContext(r.Context(), userOrgID), userID, userOrgID)
	}
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	if webAuthNToken != nil {
		credentialData = base64.RawURLEncoding.EncodeToString(webAuthNToken.CredentialCreationData)
	}
	data := &passwordlessRegistrationData{
		webAuthNData{
			userData:               l.getUserData(r, authReq, "Login Passwordless", errID, errMessage),
			CredentialCreationData: credentialData,
		},
		code,
		userID,
	}
	l.renderer.RenderTemplate(w, r, l.getTranslator(authReq), l.renderer.Templates[tmplPasswordlessRegistration], data, nil)
}

func (l *Login) handlePasswordlessRegistrationCheck(w http.ResponseWriter, r *http.Request) {
	formData := new(passwordlessRegistrationFormData)
	authReq, err := l.getAuthRequestAndParseData(r, formData)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if formData.Resend {
		l.resendPasswordlessRegistration(w, r, authReq, formData.UserID)
		return
	}
	l.checkPasswordlessRegistration(w, r, authReq, formData, nil)
}

func (l *Login) checkPasswordlessRegistration(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, formData *passwordlessRegistrationFormData, err error) {
	credData, err := base64.URLEncoding.DecodeString(formData.CredentialData)
	if err != nil {
		l.renderPasswordlessRegistration(w, r, authReq, formData.UserID, formData.Code, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	userOrgID := ""
	if authReq != nil {
		userOrgID = authReq.UserOrgID
	}
	err = l.authRepo.VerifyPasswordlessSetup(setContext(r.Context(), userOrgID), formData.UserID, userOrgID, userAgentID, formData.TokenName, credData)
	if err != nil {
		l.renderPasswordlessRegistration(w, r, authReq, formData.UserID, formData.Code, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) resendPasswordlessRegistration(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID string) {
	userOrgID := ""
	if authReq != nil {
		userOrgID = authReq.UserOrgID
	}
	_, err := l.command.ResendInitialMail(setContext(r.Context(), userOrgID), userID, "", userOrgID) //TODO: resend pw less
	l.renderPasswordlessRegistration(w, r, authReq, userID, "", err)
}
