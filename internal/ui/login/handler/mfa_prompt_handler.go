package handler

import (
	"net/http"

	"github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

const (
	tmplMFAPrompt = "mfaprompt"
)

type mfaPromptData struct {
	MFAProvider model.MFAType `schema:"provider"`
	Skip        bool          `schema:"skip"`
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
	err = l.authRepo.SkipMFAInit(setContext(r.Context(), authReq.UserOrgID), authReq.UserID)
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

func (l *Login) renderMFAPrompt(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, mfaPromptData *model.MFAPromptStep, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	data := mfaData{
		baseData:    l.getBaseData(r, authReq, "MFA Prompt", errType, errMessage),
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
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMFAPrompt], data, nil)
}

func (l *Login) handleMFACreation(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, data *mfaVerifyData) {
	switch data.MFAType {
	case model.MFATypeOTP:
		l.handleOTPCreation(w, r, authReq, data)
		return
	case model.MFATypeU2F:
		l.renderRegisterU2F(w, r, authReq, nil)
		return
	}
	l.renderError(w, r, authReq, caos_errs.ThrowPreconditionFailed(nil, "APP-Or3HO", "Errors.User.MFA.NoProviders"))
}

func (l *Login) handleOTPCreation(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, data *mfaVerifyData) {
	otp, err := l.authRepo.AddMFAOTP(setContext(r.Context(), authReq.UserOrgID), authReq.UserID)
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
