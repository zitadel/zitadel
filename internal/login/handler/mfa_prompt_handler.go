package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"net/http"
)

const (
	tmplMfaPrompt = "mfaprompt"
)

type mfaPromptData struct {
	MfaProvider model.MfaType `schema:"provider"`
	Skip        bool          `schema:"skip"`
}

func (l *Login) handleMfaPrompt(w http.ResponseWriter, r *http.Request) {
	data := new(mfaPromptData)
	authSession, err := l.getAuthRequestAndParseData(r, data)
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
	err = l.authRepo.SkipMfaInit(setContext(r.Context(), authSession.UserOrgID), authSession.UserID)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	l.handleLogin(w, r)
}

func (l *Login) renderMfaPrompt(w http.ResponseWriter, r *http.Request, authSession *model.AuthRequest, mfaPromptData *model.MfaPromptStep, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	data := mfaData{
		baseData: l.getBaseData(r, authSession, "Mfa Prompt", errType, errMessage),
		UserName: authSession.UserName,
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
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMfaPrompt], data, nil)
}

func (l *Login) handleMfaCreation(w http.ResponseWriter, r *http.Request, authSession *model.AuthRequest, data *mfaVerifyData) {
	switch data.MfaType {
	case model.MfaTypeOTP:
		l.handleOtpCreation(w, r, authSession, data)
		return
	}
	l.renderError(w, r, authSession, caos_errs.ThrowPreconditionFailed(nil, "APP-Or3HO", "No available mfa providers"))
}

func (l *Login) handleOtpCreation(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, data *mfaVerifyData) {
	otp, err := l.authRepo.AddMfaOTP(setContext(r.Context(), authReq.UserOrgID), authReq.UserID)
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
