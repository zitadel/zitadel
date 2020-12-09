package handler

import (
	"encoding/base64"
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
	user_model "github.com/caos/zitadel/internal/user/model"
)

const (
	tmplU2FVerification = "u2fverification"
)

type mfaU2FData struct {
	webAuthNData
	MFAProviders     []model.MFAType
	SelectedProvider model.MFAType
}

type mfaU2FFormData struct {
	webAuthNFormData
	SelectedProvider model.MFAType `schema:"provider"`
}

func (l *Login) renderU2FVerification(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, providers []model.MFAType, err error) {
	var errType, errMessage, credentialData string
	var webAuthNLogin *user_model.WebAuthNLogin
	if err == nil {
		userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
		webAuthNLogin, err = l.authRepo.BeginMFAU2FLogin(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.ID, userAgentID)
	}
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	if webAuthNLogin != nil {
		credentialData = base64.RawURLEncoding.EncodeToString(webAuthNLogin.CredentialAssertionData)
	}
	data := &mfaU2FData{
		webAuthNData: webAuthNData{
			userData:               l.getUserData(r, authReq, "Login WebAuthNToken", errType, errMessage),
			CredentialCreationData: credentialData,
		},
		MFAProviders:     providers,
		SelectedProvider: -1,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplU2FVerification], data, nil)
}

func (l *Login) handleU2FVerification(w http.ResponseWriter, r *http.Request) {
	formData := new(mfaU2FFormData)
	authReq, err := l.getAuthRequestAndParseData(r, formData)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	step, ok := authReq.PossibleSteps[0].(*model.MFAVerificationStep)
	if !ok {
		l.renderError(w, r, authReq, err)
		return
	}
	if formData.CredentialData == "" {
		l.renderMFAVerifySelected(w, r, authReq, step, formData.SelectedProvider, nil)
		return
	}
	credData, err := base64.URLEncoding.DecodeString(formData.CredentialData)
	if err != nil {
		l.renderU2FVerification(w, r, authReq, step.MFAProviders, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.authRepo.VerifyMFAU2F(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.ID, userAgentID, credData, model.BrowserInfoFromRequest(r))
	if err != nil {
		l.renderU2FVerification(w, r, authReq, step.MFAProviders, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
