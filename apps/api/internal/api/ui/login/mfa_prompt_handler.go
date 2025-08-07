package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	tmplMFAPrompt = "mfaprompt"
)

type mfaPromptData struct {
	MFAProvider domain.MFAType `schema:"provider"`
	Skip        bool           `schema:"skip"`
}

func (l *Login) handleMFAPrompt(w http.ResponseWriter, r *http.Request) {
	data := new(mfaPromptData)
	authReq, err := l.ensureAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if !data.Skip {
		mfaVerifyData := new(mfaVerifyData)
		mfaVerifyData.MFAType = data.MFAProvider
		l.handleMFACreation(w, r, authReq, mfaVerifyData)
		return
	}
	err = l.command.HumanSkipMFAInit(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.handleLogin(w, r)
}

func (l *Login) handleMFAPromptSelection(w http.ResponseWriter, r *http.Request) {
	data := new(mfaPromptData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	l.renderNextStep(w, r, authReq)
}

func (l *Login) renderMFAPrompt(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, mfaPromptData *domain.MFAPromptStep, err error) {
	translator := l.getTranslator(r.Context(), authReq)
	data := mfaData{
		baseData:    l.getBaseData(r, authReq, translator, "InitMFAPrompt.Title", "InitMFAPrompt.Description", err),
		profileData: l.getProfileData(authReq),
	}

	if mfaPromptData == nil {
		l.renderError(w, r, authReq, zerrors.ThrowPreconditionFailed(nil, "APP-XU0tj", "Errors.User.MFA.NoProviders"))
		return
	}

	data.MFAProviders = mfaPromptData.MFAProviders
	data.MFARequired = mfaPromptData.Required

	if len(mfaPromptData.MFAProviders) == 1 && mfaPromptData.Required {
		data := &mfaVerifyData{
			MFAType: mfaPromptData.MFAProviders[0],
		}
		l.handleMFACreation(w, r, authReq, data)
		return
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplMFAPrompt], data, nil)
}

func (l *Login) handleMFACreation(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *mfaVerifyData) {
	switch data.MFAType {
	case domain.MFATypeTOTP:
		l.handleTOTPCreation(w, r, authReq, data)
		return
	case domain.MFATypeOTPSMS:
		l.handleRegisterOTPSMS(w, r, authReq)
		return
	case domain.MFATypeOTPEmail:
		l.handleRegisterOTPEmail(w, r, authReq)
		return
	case domain.MFATypeU2F:
		l.renderRegisterU2F(w, r, authReq, nil)
		return
	}
	l.renderError(w, r, authReq, zerrors.ThrowPreconditionFailed(nil, "APP-Or3HO", "Errors.User.MFA.NoProviders"))
}

func (l *Login) handleTOTPCreation(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *mfaVerifyData) {
	otp, err := l.command.AddHumanTOTP(setUserContext(r.Context(), authReq.UserID, authReq.UserOrgID), authReq.UserID, authReq.UserOrgID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	data.totpData = totpData{
		Secret: otp.Secret,
		Url:    otp.URI,
	}
	l.renderMFAInitVerify(w, r, authReq, data, nil)
}

// handleRegisterOTPEmail will directly add OTP Email as 2FA.
// It will also add a successful OTP Email check to the auth request.
func (l *Login) handleRegisterOTPEmail(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	_, err := l.command.AddHumanOTPEmailWithCheckSucceeded(setUserContext(r.Context(), authReq.UserID, authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, authReq)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	done := &mfaDoneData{
		MFAType: domain.MFATypeOTPEmail,
	}
	l.renderMFAInitDone(w, r, authReq, done)
}
