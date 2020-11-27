package handler

import (
	"bytes"
	"net/http"

	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/qrcode"
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
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	var verifyData *mfaVerifyData
	switch data.MfaType {
	case model.MFATypeOTP:
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
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err := l.authRepo.VerifyMfaOTPSetup(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, data.Code, userAgentID)
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
		errMessage = l.getErrorMessage(r, err)
	}
	data.baseData = l.getBaseData(r, authReq, "Mfa Init Verify", errType, errMessage)
	data.profileData = l.getProfileData(authReq)
	if data.MfaType == model.MFATypeOTP {
		code, err := generateQrCode(data.otpData.Url)
		if err == nil {
			data.otpData.QrCode = code
		}
	}

	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMfaInitVerify], data, nil)
}

func generateQrCode(url string) (string, error) {
	var b bytes.Buffer
	s := svg.New(&b)

	qrCode, err := qr.Encode(url, qr.M, qr.Auto)
	if err != nil {
		return "", err
	}
	qs := qrcode.NewQrSVG(qrCode, 5)
	qs.StartQrSVG(s)
	qs.WriteQrSVG(s)

	s.End()
	return string(b.Bytes()), nil
}
