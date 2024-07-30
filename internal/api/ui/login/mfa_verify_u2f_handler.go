package login

import (
	"encoding/base64"
	"net/http"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
)

const (
	tmplU2FVerification = "u2fverification"
)

type mfaU2FData struct {
	webAuthNData
	MFAProviders     []domain.MFAType
	SelectedProvider domain.MFAType
}

type mfaU2FFormData struct {
	webAuthNFormData
	SelectedProvider domain.MFAType `schema:"provider"`
}

func (l *Login) renderU2FVerification(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, providers []domain.MFAType, err error) {
	var errID, errMessage, credentialData string
	var webAuthNLogin *domain.WebAuthNLogin
	if err == nil {
		userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
		webAuthNLogin, err = l.authRepo.BeginMFAU2FLogin(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, authReq.ID, userAgentID)
	}
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	if webAuthNLogin != nil {
		credentialData = base64.RawURLEncoding.EncodeToString(webAuthNLogin.CredentialAssertionData)
	}
	translator := l.getTranslator(r.Context(), authReq)
	data := &mfaU2FData{
		webAuthNData: webAuthNData{
			userData:               l.getUserData(r, authReq, translator, "VerifyMFAU2F.Title", "VerifyMFAU2F.Description", errID, errMessage),
			CredentialCreationData: credentialData,
		},
		MFAProviders:     providers,
		SelectedProvider: -1,
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplU2FVerification], data, nil)
}

func (l *Login) handleU2FVerification(w http.ResponseWriter, r *http.Request) {
	formData := new(mfaU2FFormData)
	authReq, err := l.ensureAuthRequestAndParseData(r, formData)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	step, ok := authReq.PossibleSteps[0].(*domain.MFAVerificationStep)
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
	err = l.authRepo.VerifyMFAU2F(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, authReq.ID, userAgentID, credData, domain.BrowserInfoFromRequest(r))

	metadata, actionErr := l.runPostInternalAuthenticationActions(authReq, r, authMethodU2F, err)
	if err == nil && actionErr == nil && len(metadata) > 0 {
		_, err = l.command.BulkSetUserMetadata(r.Context(), authReq.UserID, authReq.UserOrgID, metadata...)
	} else if actionErr != nil && err == nil {
		err = actionErr
	}

	if err != nil {
		l.renderU2FVerification(w, r, authReq, step.MFAProviders, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
