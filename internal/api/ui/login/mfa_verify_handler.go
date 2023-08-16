package login

import (
	"net/http"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
)

const (
	tmplMFAVerify = "mfaverify"
)

type mfaVerifyFormData struct {
	MFAType          domain.MFAType `schema:"mfaType"`
	Code             string         `schema:"code"`
	SelectedProvider domain.MFAType `schema:"provider"`
}

func (l *Login) handleMFAVerify(w http.ResponseWriter, r *http.Request) {
	data := new(mfaVerifyFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	step, ok := authReq.PossibleSteps[0].(*domain.MFAVerificationStep)
	if !ok {
		l.renderError(w, r, authReq, err)
		return
	}
	if data.Code == "" {
		l.renderMFAVerifySelected(w, r, authReq, step, data.SelectedProvider, nil)
		return
	}
	if data.MFAType == domain.MFATypeTOTP {
		userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
		err = l.authRepo.VerifyMFAOTP(setContext(r.Context(), authReq.UserOrgID), authReq.ID, authReq.UserID, authReq.UserOrgID, data.Code, userAgentID, domain.BrowserInfoFromRequest(r))

		metadata, actionErr := l.runPostInternalAuthenticationActions(authReq, r, authMethodOTP, err)
		if err == nil && actionErr == nil && len(metadata) > 0 {
			_, err = l.command.BulkSetUserMetadata(r.Context(), authReq.UserID, authReq.UserOrgID, metadata...)
		} else if actionErr != nil && err == nil {
			err = actionErr
		}

		if err != nil {
			l.renderMFAVerifySelected(w, r, authReq, step, domain.MFATypeTOTP, err)
			return
		}
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) renderMFAVerify(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, verificationStep *domain.MFAVerificationStep, err error) {
	if verificationStep == nil {
		l.renderError(w, r, authReq, err)
		return
	}
	provider := verificationStep.MFAProviders[len(verificationStep.MFAProviders)-1]
	l.renderMFAVerifySelected(w, r, authReq, verificationStep, provider, err)
}

func (l *Login) renderMFAVerifySelected(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, verificationStep *domain.MFAVerificationStep, selectedProvider domain.MFAType, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	data := l.getUserData(r, authReq, "", "", errID, errMessage)
	if verificationStep == nil {
		l.renderError(w, r, authReq, err)
		return
	}
	translator := l.getTranslator(r.Context(), authReq)

	switch selectedProvider {
	case domain.MFATypeU2F:
		data.Title = translator.LocalizeWithoutArgs("VerifyMFAU2F.Title")
		data.Description = translator.LocalizeWithoutArgs("VerifyMFAU2F.Description")
		l.renderU2FVerification(w, r, authReq, removeSelectedProviderFromList(verificationStep.MFAProviders, domain.MFATypeU2F), nil)
		return
	case domain.MFATypeTOTP:
		data.MFAProviders = removeSelectedProviderFromList(verificationStep.MFAProviders, domain.MFATypeTOTP)
		data.SelectedMFAProvider = domain.MFATypeTOTP
		data.Title = translator.LocalizeWithoutArgs("VerifyMFAOTP.Title")
		data.Description = translator.LocalizeWithoutArgs("VerifyMFAOTP.Description")
	case domain.MFATypeOTPSMS:
		l.handleOTPVerification(w, r, authReq, verificationStep.MFAProviders, domain.MFATypeOTPSMS, nil)
		return
	case domain.MFATypeOTPEmail:
		l.handleOTPVerification(w, r, authReq, verificationStep.MFAProviders, domain.MFATypeOTPEmail, nil)
		return
	default:
		l.renderError(w, r, authReq, err)
		return
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplMFAVerify], data, nil)
}

func removeSelectedProviderFromList(providers []domain.MFAType, selected domain.MFAType) []domain.MFAType {
	for i := len(providers) - 1; i >= 0; i-- {
		if providers[i] == selected {
			copy(providers[i:], providers[i+1:])
			return providers[:len(providers)-1]
		}
	}
	return providers
}
