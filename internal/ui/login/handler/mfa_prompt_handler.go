package handler

import (
	"github.com/caos/zitadel/internal/domain"
	"net/http"

	caos_errs "github.com/caos/zitadel/internal/errors"
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
	authReq, err := l.getAuthRequestAndParseData(r, data)
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
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	translator := l.getTranslator(authReq)
	data := mfaData{
		baseData:    l.getBaseData(r, authReq, translator, "MFA Prompt", errID, errMessage),
		profileData: l.getProfileData(authReq),
	}

	if mfaPromptData == nil {
		l.renderError(w, r, authReq, caos_errs.ThrowPreconditionFailed(nil, "APP-XU0tj", "Errors.User.MFA.NoProviders"))
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
	case domain.MFATypeOTP:
		l.handleOTPCreation(w, r, authReq, data)
		return
	case domain.MFATypeU2F:
		l.renderRegisterU2F(w, r, authReq, nil)
		return
	}
	l.renderError(w, r, authReq, caos_errs.ThrowPreconditionFailed(nil, "APP-Or3HO", "Errors.User.MFA.NoProviders"))
}

func (l *Login) handleOTPCreation(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *mfaVerifyData) {
	otp, err := l.command.AddHumanOTP(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	data.otpData = otpData{
		Secret: otp.SecretString,
		Url:    otp.Url,
	}
	l.renderMFAInitVerify(w, r, authReq, data, nil)
}
