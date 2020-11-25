package handler

import (
	"net/http"

	"github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

const (
	tmplMfaPrompt = "mfaprompt"
)

type mfaPromptData struct {
	MfaProvider model.MFAType `schema:"provider"`
	Skip        bool          `schema:"skip"`
}

func (l *Login) handleMfaPrompt(w http.ResponseWriter, r *http.Request) {
	data := new(mfaPromptData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if !data.Skip {
		mfaVerifyData := new(mfaVerifyData)
		mfaVerifyData.MfaType = data.MfaProvider
		l.handleMfaCreation(w, r, authReq, mfaVerifyData)
		return
	}
	err = l.authRepo.SkipMfaInit(setContext(r.Context(), authReq.UserOrgID), authReq.UserID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.handleLogin(w, r)
}

func (l *Login) handleMfaPromptSelection(w http.ResponseWriter, r *http.Request) {
	data := new(mfaPromptData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	l.renderNextStep(w, r, authReq)
}

func (l *Login) renderMfaPrompt(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, mfaPromptData *model.MfaPromptStep, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	data := mfaData{
		baseData:    l.getBaseData(r, authReq, "Mfa Prompt", errType, errMessage),
		profileData: l.getProfileData(authReq),
	}

	if mfaPromptData == nil {
		l.renderError(w, r, authReq, caos_errs.ThrowPreconditionFailed(nil, "APP-XU0tj", "Errors.User.Mfa.NoProviders"))
		return
	}

	data.MfaProviders = mfaPromptData.MfaProviders
	data.MfaRequired = mfaPromptData.Required

	if len(mfaPromptData.MfaProviders) == 1 && mfaPromptData.Required {
		data := &mfaVerifyData{
			MfaType: mfaPromptData.MfaProviders[0],
		}
		l.handleMfaCreation(w, r, authReq, data)
		return
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMfaPrompt], data, nil)
}

func (l *Login) handleMfaCreation(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, data *mfaVerifyData) {
	switch data.MfaType {
	case model.MFATypeOTP:
		l.handleOtpCreation(w, r, authReq, data)
		return
	case model.MFATypeU2F:
		l.renderRegisterU2F(w, r, authReq, nil)
		return
	}
	l.renderError(w, r, authReq, caos_errs.ThrowPreconditionFailed(nil, "APP-Or3HO", "Errors.User.Mfa.NoProviders"))
}

func (l *Login) handleOtpCreation(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, data *mfaVerifyData) {
	otp, err := l.authRepo.AddMFAOTP(setContext(r.Context(), authReq.UserOrgID), authReq.UserID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	data.otpData = otpData{
		Secret: otp.SecretString,
		Url:    otp.Url,
	}
	l.renderMfaInitVerify(w, r, authReq, data, nil)
}
