package login

import (
	"bytes"
	"html/template"
	"net/http"

	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/qrcode"
)

const (
	tmplMFAInitVerify = "mfainitverify"
)

type mfaInitVerifyData struct {
	MFAType domain.MFAType `schema:"mfaType"`
	Code    string         `schema:"code"`
	URL     string         `schema:"url"`
	Secret  string         `schema:"secret"`
}

func (l *Login) handleMFAInitVerify(w http.ResponseWriter, r *http.Request) {
	data := new(mfaInitVerifyData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	var verifyData *mfaVerifyData
	switch data.MFAType {
	case domain.MFATypeTOTP:
		verifyData = l.handleOTPVerify(w, r, authReq, data)
	}

	if verifyData != nil {
		l.renderMFAInitVerify(w, r, authReq, verifyData, err)
		return
	}

	done := &mfaDoneData{
		MFAType: data.MFAType,
	}
	l.renderMFAInitDone(w, r, authReq, done)
}

func (l *Login) handleOTPVerify(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *mfaInitVerifyData) *mfaVerifyData {
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	_, err := l.command.HumanCheckMFATOTPSetup(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, data.Code, userAgentID, authReq.UserOrgID)
	if err == nil {
		return nil
	}
	mfadata := &mfaVerifyData{
		MFAType: data.MFAType,
		totpData: totpData{
			Secret: data.Secret,
			Url:    data.URL,
		},
	}

	return mfadata
}

func (l *Login) renderMFAInitVerify(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *mfaVerifyData, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	translator := l.getTranslator(r.Context(), authReq)
	data.baseData = l.getBaseData(r, authReq, translator, "InitMFAOTP.Title", "InitMFAOTP.Description", errID, errMessage)
	data.profileData = l.getProfileData(authReq)
	if data.MFAType == domain.MFATypeTOTP {
		code, err := generateQrCode(data.totpData.Url)
		if err == nil {
			data.totpData.QrCode = template.HTML(code)
		}
	}

	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplMFAInitVerify], data, nil)
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
