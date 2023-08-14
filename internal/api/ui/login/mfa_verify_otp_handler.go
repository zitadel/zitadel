package login

import (
	"net/http"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
)

const (
	tmplOTPVerification = "otpverification"
)

type mfaOTPData struct {
	userData
	MFAProviders     []domain.MFAType
	SelectedProvider domain.MFAType
}

type mfaOTPFormData struct {
	Resend           bool           `schema:"resend"`
	Code             string         `schema:"code"`
	SelectedProvider domain.MFAType `schema:"selectedProvider"`
	Provider         domain.MFAType `schema:"provider"`
}

// renderOTPVerification renders the OTP verification for SMS and Email based on the passed MFAType.
// It will send a new code to either phone or email first.
func (l *Login) renderOTPVerification(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, providers []domain.MFAType, selectedProvider domain.MFAType, err error) {
	var errID, errMessage string
	if err == nil {
		userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
		sendCode := l.authRepo.SendMFAOTPSMS
		if selectedProvider == domain.MFATypeOTPEmail {
			sendCode = l.authRepo.SendMFAOTPEmail
		}
		err = sendCode(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, authReq.ID, userAgentID)
	}
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	data := &mfaOTPData{
		userData:         l.getUserData(r, authReq, "VerifyMFAU2F.Title", "VerifyMFAU2F.Description", errID, errMessage),
		MFAProviders:     removeSelectedProviderFromList(providers, selectedProvider),
		SelectedProvider: selectedProvider,
	}
	l.renderer.RenderTemplate(w, r, l.getTranslator(r.Context(), authReq), l.renderer.Templates[tmplOTPVerification], data, nil)
}

// handleRegisterSMSCheck handles form submissions of the OTP verification.
// On successful code verification, the check will be added to the auth request.
// A user is also able to request a code resend or choose another provider.
func (l *Login) handleOTPVerification(w http.ResponseWriter, r *http.Request) {
	formData := new(mfaOTPFormData)
	authReq, err := l.getAuthRequestAndParseData(r, formData)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	step, ok := authReq.PossibleSteps[0].(*domain.MFAVerificationStep)
	if !ok {
		l.renderError(w, r, authReq, err)
		return
	}
	if formData.Resend {
		l.renderOTPVerification(w, r, authReq, step.MFAProviders, formData.SelectedProvider, nil)
		return
	}
	if formData.Code == "" {
		l.renderMFAVerifySelected(w, r, authReq, step, formData.Provider, nil)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	verifyCode := l.authRepo.VerifyMFAOTPSMS
	actionType := authMethodOTPSMS
	if formData.SelectedProvider == domain.MFATypeOTPEmail {
		verifyCode = l.authRepo.VerifyMFAOTPEmail
		actionType = authMethodOTPEmail
	}
	err = verifyCode(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, formData.Code, authReq.ID, userAgentID, domain.BrowserInfoFromRequest(r))

	metadata, actionErr := l.runPostInternalAuthenticationActions(authReq, r, actionType, err)
	if err == nil && actionErr == nil && len(metadata) > 0 {
		_, err = l.command.BulkSetUserMetadata(r.Context(), authReq.UserID, authReq.UserOrgID, metadata...)
	} else if actionErr != nil && err == nil {
		err = actionErr
	}

	if err != nil {
		l.renderOTPVerification(w, r, authReq, step.MFAProviders, formData.SelectedProvider, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
