package handler

import (
	"net/http"

	"github.com/caos/citadel/login/internal/model"

	caos_errs "github.com/caos/utils/errors"
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
	authSession, err := l.getAuthSessionAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	if !data.Skip {
		mfaVerifyData := new(mfaVerifyData)
		mfaVerifyData.MfaType = data.MfaProvider
		l.handleMfaCreation(w, r, authSession, mfaVerifyData)
		return
	}
	err = l.service.Auth.SkipMfaInit(r.Context(), authSession.UserSession.User.UserID)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	l.handleLogin(w, r)
}

func (l *Login) renderMfaPrompt(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, mfaPromptData *model.MfaPromptData, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	data := mfaData{
		baseData: l.getBaseData(r, authSession, "Mfa Prompt", errType, errMessage),
		UserName: authSession.UserSession.User.UserName,
	}

	if mfaPromptData == nil {
		l.renderError(w, r, authSession, caos_errs.ThrowPreconditionFailed(nil, "APP-XU0tj", "No available mfa providers"))
		return
	}

	data.MfaProviders = mfaPromptData.MfaProviders
	data.MfaRequired = mfaPromptData.Required

	if len(mfaPromptData.MfaProviders) == 1 && mfaPromptData.Required {
		data := &mfaVerifyData{
			MfaType: mfaPromptData.MfaProviders[0],
		}
		l.handleMfaCreation(w, r, authSession, data)
		return
	}
	l.renderer.RenderTemplate(w, r, a.renderer.Templates[tmplMfaPrompt], data, nil)
}

func (l *Login) handleMfaCreation(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, data *mfaVerifyData) {
	switch data.MfaType {
	case model.MFA_OTP:
		l.handleOtpCreation(w, r, authSession, data)
		return
	}
	l.renderError(w, r, authSession, caos_errs.ThrowPreconditionFailed(nil, "APP-Or3HO", "No available mfa providers"))
}

func (l *Login) handleOtpCreation(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, data *mfaVerifyData) {
	otp, err := l.service.Auth.AddMfaOTP(r.Context(), authSession.UserSession.User.UserID)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}

	data.otpData = otpData{
		Secret: otp.Secret,
		Url:    otp.Url,
	}

	l.renderMfaInitVerify(w, r, authSession, data, nil)
}
