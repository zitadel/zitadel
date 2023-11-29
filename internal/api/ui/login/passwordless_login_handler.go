package login

import (
	"encoding/base64"
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
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
		webAuthNLogin, err = l.authRepo.BeginPasswordlessLogin(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, authReq.ID, authReq.AgentID)
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
	translator := l.getTranslator(r.Context(), authReq)
	data := &passwordlessData{
		webAuthNData{
			userData:               l.getUserData(r, authReq, translator, "Passwordless.Title", "Passwordless.Description", errID, errMessage),
			CredentialCreationData: credentialData,
		},
		passwordSet,
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplPasswordlessVerification], data, nil)
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
	err = l.authRepo.VerifyPasswordless(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, authReq.ID, authReq.AgentID, credData, domain.BrowserInfoFromRequest(r))

	metadata, actionErr := l.runPostInternalAuthenticationActions(authReq, r, authMethodPasswordless, err)
	if err == nil && actionErr == nil && len(metadata) > 0 {
		_, err = l.command.BulkSetUserMetadata(r.Context(), authReq.UserID, authReq.UserOrgID, metadata...)
	} else if actionErr != nil && err == nil {
		err = actionErr
	}

	if err != nil {
		l.renderPasswordlessVerification(w, r, authReq, formData.PasswordLogin, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
