package handler

import (
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	tmplMfaVerify = "mfaverify"
)

type mfaVerifyFormData struct {
	MfaType model.MFAType `schema:"mfaType"`
	Code    string        `schema:"code"`
}

func (l *Login) handleMfaVerify(w http.ResponseWriter, r *http.Request) {
	data := new(mfaVerifyFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if data.MfaType == model.MFATypeOTP {
		userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
		err = l.authRepo.VerifyMFAOTP(setContext(r.Context(), authReq.UserOrgID), authReq.ID, authReq.UserID, data.Code, userAgentID, model.BrowserInfoFromRequest(r))
	}
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) renderMfaVerify(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, verificationStep *model.MfaVerificationStep, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	data := l.getUserData(r, authReq, "Mfa Verify", errType, errMessage)
	if verificationStep == nil {
		l.renderError(w, r, authReq, err)
		return
	}
	switch verificationStep.MfaProviders[len(verificationStep.MfaProviders)-1] {
	case model.MFATypeU2F:
		l.renderU2FVerification(w, r, authReq)
		return
	case model.MFATypeOTP:
		data.MfaProviders = verificationStep.MfaProviders
		data.SelectedMfaProvider = model.MFATypeOTP
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMfaVerify], data, nil)
}
