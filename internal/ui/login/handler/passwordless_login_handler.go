package handler

import (
	"encoding/base64"
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
	user_model "github.com/caos/zitadel/internal/user/model"
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

func (l *Login) renderPasswordlessVerification(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage, credentialData string
	var webAuthNLogin *user_model.WebAuthNLogin
	if err == nil {
		userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
		webAuthNLogin, err = l.authRepo.BeginPasswordlessLogin(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.ID, userAgentID)
	}
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	if webAuthNLogin != nil {
		credentialData = base64.RawURLEncoding.EncodeToString(webAuthNLogin.CredentialAssertionData)
	}
	var passwordLogin bool
	if authReq.LoginPolicy != nil {
		passwordLogin = authReq.LoginPolicy.AllowUsernamePassword
	}
	data := &passwordlessData{
		webAuthNData{
			userData:               l.getUserData(r, authReq, "Login Passwordless", errType, errMessage),
			CredentialCreationData: credentialData,
		},
		passwordLogin,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplPasswordlessVerification], data, nil)
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
		l.renderPasswordlessVerification(w, r, authReq, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.authRepo.VerifyPasswordless(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.ID, userAgentID, credData, model.BrowserInfoFromRequest(r))
	if err != nil {
		l.renderPasswordlessVerification(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
