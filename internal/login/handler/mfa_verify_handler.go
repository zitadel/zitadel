package handler

import (
	"net/http"

	"github.com/caos/citadel/login/internal/model"
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
	authSession, err := l.getAuthSessionAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	authSession.UserSession.MfaType = data.MfaType
	browserInfo := &model.BrowserInformation{RemoteIP: &model.IP{}} //TODO: impl
	authSession, err = l.service.Auth.VerifyMfa(r.Context(), authSession, data.Code, browserInfo)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	l.renderNextStep(w, r, authSession)
}

func (l *Login) renderMfaVerify(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, verifyData *model.MfaVerifyData, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	if verifyData != nil {
		errMessage = verifyData.ErrMsg
	}
	data := userData{
		baseData: l.getBaseData(r, authSession, "Mfa Verify", errType, errMessage),
		UserName: authSession.UserSession.User.UserName,
	}
	if verifyData != nil {
		data.MfaProviders = verifyData.MfaProviders
		data.SelectedMfaProvider = verifyData.MfaProviders[0]
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMfaVerify], data, nil)
}
