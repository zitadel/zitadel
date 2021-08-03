package handler

import (
	"encoding/base64"
	"github.com/caos/zitadel/internal/domain"
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
)

const (
	tmplPasswordlessVerification = "passwordlessverification"
)

type passwordlessData struct {
	webAuthNData
	PasswordLogin bool
}

type passwordlessFormData struct {
	webAuthNFormData
	PasswordLogin bool `schema:"passwordlogin"`
}

func (l *Login) renderPasswordlessVerification(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, passwordSet bool, err error) {
	var errID, errMessage, credentialData string
	var webAuthNLogin *domain.WebAuthNLogin
	if err == nil {
		userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
		webAuthNLogin, err = l.authRepo.BeginPasswordlessLogin(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, authReq.ID, userAgentID)
	}
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	if webAuthNLogin != nil {
		credentialData = base64.RawURLEncoding.EncodeToString(webAuthNLogin.CredentialAssertionData)
	}
	if passwordSet && authReq.LoginPolicy != nil {
		passwordSet = authReq.LoginPolicy.AllowUsernamePassword
	}
	data := &passwordlessData{
		webAuthNData{
			userData:               l.getUserData(r, authReq, "Login Passwordless", errID, errMessage),
			CredentialCreationData: credentialData,
		},
		passwordSet,
	}
	l.renderer.RenderTemplate(w, r, l.getTranslator(authReq), l.renderer.Templates[tmplPasswordlessVerification], data, nil)
}

func (l *Login) handlePasswordlessVerification(w http.ResponseWriter, r *http.Request) {
	formData := new(passwordlessFormData)
	authReq, err := l.getAuthRequestAndParseData(r, formData)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if formData.PasswordLogin {
		l.renderPassword(w, r, authReq, nil)
		return
	}
	credData, err := base64.URLEncoding.DecodeString(formData.CredentialData)
	if err != nil {
		l.renderPasswordlessVerification(w, r, authReq, formData.PasswordLogin, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.authRepo.VerifyPasswordless(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, authReq.ID, userAgentID, credData, domain.BrowserInfoFromRequest(r))
	if err != nil {
		l.renderPasswordlessVerification(w, r, authReq, formData.PasswordLogin, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
