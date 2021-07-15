package handler

import (
	"encoding/base64"
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/domain"
)

const (
	tmplPasswordlessRegistration        = "passwordlessregistration"
	tmplPasswordlessRegistrationDone    = "passwordlessregistrationdone"
	queryPasswordlessRegistrationCode   = "code"
	queryPasswordlessRegistrationCodeID = "codeID"
	queryPasswordlessRegistrationUserID = "userID"
)

type passwordlessRegistrationData struct {
	webAuthNData
	Code   string
	CodeID string
	UserID string
}

type passwordlessRegistrationFormData struct {
	webAuthNFormData
	Code      string `schema:"code"`
	CodeID    string `schema:"codeID"`
	UserID    string `schema:"userID"`
	TokenName string `schema:"name"`
	Resend    bool   `schema:"resend"`
}

func (l *Login) handlePasswordlessRegistration(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue(queryPasswordlessRegistrationUserID)
	codeID := r.FormValue(queryPasswordlessRegistrationCodeID)
	code := r.FormValue(queryPasswordlessRegistrationCode)
	l.renderPasswordlessRegistration(w, r, nil, userID, codeID, code, nil)
}

func (l *Login) renderPasswordlessRegistration(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID, codeID, code string, err error) {
	var errID, errMessage, credentialData string
	if authReq != nil {
		userID = authReq.UserID
	}
	var webAuthNToken *domain.WebAuthNToken
	if err == nil {
		if authReq != nil {
			webAuthNToken, err = l.authRepo.BeginPasswordlessSetup(setContext(r.Context(), authReq.UserOrgID), userID, authReq.UserOrgID)
		} else {
			webAuthNToken, err = l.authRepo.BeginPasswordlessInitCodeSetup(setContext(r.Context(), ""), userID, "", codeID, code)
		}
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
		codeID,
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
		l.renderPasswordlessRegistration(w, r, authReq, formData.UserID, formData.CodeID, formData.Code, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	if authReq != nil {
		err = l.authRepo.VerifyPasswordlessSetup(setContext(r.Context(), authReq.UserOrgID), formData.UserID, authReq.UserOrgID, userAgentID, formData.TokenName, credData)
	} else {
		err = l.authRepo.VerifyPasswordlessInitCodeSetup(setContext(r.Context(), ""), formData.UserID, "", userAgentID, formData.TokenName, formData.CodeID, formData.Code, credData)
	}
	if err != nil {
		l.renderPasswordlessRegistration(w, r, authReq, formData.UserID, formData.CodeID, formData.Code, err)
		return
	}
	l.renderPasswordResetDone(w, r, authReq, nil)
}

func (l *Login) renderPasswordlessRegistrationDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	data := l.getUserData(r, authReq, "Passwordless Registration Done", errID, errMessage)
	l.renderer.RenderTemplate(w, r, l.getTranslator(authReq), l.renderer.Templates[tmplPasswordlessRegistrationDone], data, nil)
}

func (l *Login) resendPasswordlessRegistration(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID string) {
	userOrgID := ""
	if authReq != nil {
		userOrgID = authReq.UserOrgID
	}
	_, err := l.command.ResendInitialMail(setContext(r.Context(), userOrgID), userID, "", userOrgID) //TODO: resend pw less
	l.renderPasswordlessRegistration(w, r, authReq, userID, "", "", err)
}
