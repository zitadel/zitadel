package handler

import (
	"encoding/base64"
	//"encoding/base64"
	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/skip2/go-qrcode"
	"net/http"
	//qrcode "github.com/skip2/go-qrcode"
)

const (
	tmplMfaInitVerify = "mfainitverify"
)

type mfaInitVerifyData struct {
	MfaType model.MfaType `schema:"mfaType"`
	Code    string        `schema:"code"`
	URL     string        `schema:"url"`
	Secret  string        `schema:"secret"`
}

func (l *Login) handleMfaInitVerify(w http.ResponseWriter, r *http.Request) {
	data := new(mfaInitVerifyData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	var verifyData *mfaVerifyData
	switch data.MfaType {
	case model.MfaTypeOTP:
		verifyData = l.handleOtpVerify(w, r, authReq, data)
	}

	if verifyData != nil {
		l.renderMfaInitVerify(w, r, authReq, verifyData, err)
		return
	}

	done := &mfaDoneData{
		MfaType: data.MfaType,
	}
	l.renderMfaInitDone(w, r, authReq, done)
}

func (l *Login) handleOtpVerify(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, data *mfaInitVerifyData) *mfaVerifyData {
	err := l.authRepo.VerifyMfaOTPSetup(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, data.Code)
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

func (l *Login) renderMfaInitVerify(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, data *mfaVerifyData, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	data.baseData = l.getBaseData(r, authReq, "Mfa Init Verify", errType, errMessage)
	data.UserName = authReq.UserName
	if data.MfaType == model.MfaTypeOTP {
		qrCode, err := qrcode.Encode(data.otpData.Url, qrcode.Medium, 256)
		if err == nil {
			data.otpData.QrCode = base64.StdEncoding.EncodeToString(qrCode)
		}
	}

	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMfaInitVerify], data, nil)
}
