package login

import (
	"context"
	"fmt"
	"net/http"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
)

const (
	tmplOTPVerification   = "otpverification"
	querySelectedProvider = "selectedProvider"
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

func OTPLink(origin, authRequestID, code string, provider domain.MFAType) string {
	return fmt.Sprintf("%s%s?%s=%s&%s=%s&%s=%d", externalLink(origin), EndpointMFAOTPVerify, QueryAuthRequestID, authRequestID, queryCode, code, querySelectedProvider, provider)
}

func OTPLinkTemplate(origin, authRequestID string, provider domain.MFAType) string {
	return fmt.Sprintf("%s%s?%s=%s&%s=%s&%s=%d", externalLink(origin), EndpointMFAOTPVerify, QueryAuthRequestID, authRequestID, queryCode, "{{.Code}}", querySelectedProvider, provider)
}

// renderOTPVerification renders the OTP verification for SMS and Email based on the passed MFAType.
// It will send a new code to either phone or email first.
func (l *Login) handleOTPVerification(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, providers []domain.MFAType, selectedProvider domain.MFAType, err error) {
	if err != nil {
		l.renderOTPVerification(w, r, authReq, providers, selectedProvider, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	var sendCode func(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID string) error
	switch selectedProvider {
	case domain.MFATypeOTPSMS:
		sendCode = l.authRepo.SendMFAOTPSMS
	case domain.MFATypeOTPEmail:
		sendCode = l.authRepo.SendMFAOTPEmail
		// another type should never be passed, but just making sure
	case domain.MFATypeU2F,
		domain.MFATypeTOTP,
		domain.MFATypeU2FUserVerification:
		l.renderError(w, r, authReq, err)
		return
	}
	err = sendCode(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, authReq.ID, userAgentID)
	l.renderOTPVerification(w, r, authReq, providers, selectedProvider, err)
}

func (l *Login) renderOTPVerification(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, providers []domain.MFAType, selectedProvider domain.MFAType, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	translator := l.getTranslator(r.Context(), authReq)
	data := &mfaOTPData{
		userData:         l.getUserData(r, authReq, translator, "VerifyMFAU2F.Title", "VerifyMFAU2F.Description", errID, errMessage),
		MFAProviders:     removeSelectedProviderFromList(providers, selectedProvider),
		SelectedProvider: selectedProvider,
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplOTPVerification], data, nil)
}

// handleOTPVerificationCheck handles form submissions of the OTP verification.
// On successful code verification, the check will be added to the auth request.
// A user is also able to request a code resend or choose another provider.
func (l *Login) handleOTPVerificationCheck(w http.ResponseWriter, r *http.Request) {
	formData := new(mfaOTPFormData)
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
	if formData.Resend {
		l.handleOTPVerification(w, r, authReq, step.MFAProviders, formData.SelectedProvider, nil)
		return
	}
	if formData.Code == "" {
		l.renderMFAVerifySelected(w, r, authReq, step, formData.Provider, nil)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	var actionType authMethod
	var verifyCode func(ctx context.Context, userID, resourceOwner, code, authRequestID, userAgentID string, info *domain.BrowserInfo) error
	switch formData.SelectedProvider {
	case domain.MFATypeOTPSMS:
		actionType = authMethodOTPSMS
		verifyCode = l.authRepo.VerifyMFAOTPSMS
	case domain.MFATypeOTPEmail:
		actionType = authMethodOTPEmail
		verifyCode = l.authRepo.VerifyMFAOTPEmail
		// another type should never be passed, but just making sure
	case domain.MFATypeU2F,
		domain.MFATypeTOTP,
		domain.MFATypeU2FUserVerification:
		l.renderOTPVerification(w, r, authReq, step.MFAProviders, formData.SelectedProvider, err)
		return
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
