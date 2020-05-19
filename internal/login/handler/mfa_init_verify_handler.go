package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/caos/citadel/login/internal/model"
	qrcode "github.com/skip2/go-qrcode"
)

const (
	tmplMfaInitVerify = "mfainitverify"
)

type mfaInitVerifyData struct {
	MfaType model.MFAType `schema:"mfaType"`
	Code    string        `schema:"code"`
	URL     string        `schema:"url"`
	Secret  string        `schema:"secret"`
}

func (l *Login) handleMfaInitVerify(w http.ResponseWriter, r *http.Request) {
	data := new(mfaInitVerifyData)
	authSession, err :=l.getAuthSessionAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	var verifyData *mfaVerifyData
	switch data.MfaType {
	case model.MFA_OTP:
		verifyData =l.handleOtpVerify(w, r, authSession, data)
	}

	if verifyData != nil {
		l.renderMfaInitVerify(w, r, authSession, verifyData, err)
		return
	}

	done := &mfaDoneData{
		MfaType: data.MfaType,
	}
	l.renderMfaInitDone(w, r, authSession, done)
}

func (l *Login) handleOtpVerify(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, data *mfaInitVerifyData) *mfaVerifyData {
	_, err :=l.service.Auth.VerifyMfaOTP(r.Context(), data.Code, authSession.UserSession.User.UserID)
	if err == nil {
		return nil
	}
	mfadata := &mfaVerifyData{
		MfaType: data.MfaType,
		otpData: otpData{
			Secret: data.Secret,
			Url:    data.URL,
		},
	}

	return mfadata
}

func (l *Login) renderMfaInitVerify(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, data *mfaVerifyData, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	data.baseData = l.getBaseData(r, authSession, "Mfa Init Verify", errType, errMessage)
	data.UserName = authSession.UserSession.User.UserName
	if data.MfaType == model.MFA_OTP {
		qrCode, err := qrcode.Encode(data.otpData.Url, qrcode.Medium, 256)
		if err == nil {
			data.otpData.QrCode = base64.StdEncoding.EncodeToString(qrCode)
		}
	}

	l.renderer.RenderTemplate(w, r,l.renderer.Templates[tmplMfaInitVerify], data, nil)
}
