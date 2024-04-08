package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
)

const (
	tmplMFASMSInit = "mfainitsms"
)

type smsInitData struct {
	userData
	Edit    bool
	MFAType domain.MFAType
	Phone   string
}

type smsInitFormData struct {
	Edit     bool   `schema:"edit"`
	Resend   bool   `schema:"resend"`
	Phone    string `schema:"phone"`
	NewPhone string `schema:"newPhone"`
	Code     string `schema:"code"`
}

// handleRegisterOTPSMS checks if the user has a verified phone number and will directly add OTP SMS as 2FA.
// It will also add a successful OTP SMS check to the auth request.
// If there's no verified phone number, the potential last phone number will be used to render the registration page
func (l *Login) handleRegisterOTPSMS(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	user, err := l.query.GetNotifyUserByID(r.Context(), true, authReq.UserID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if user.VerifiedPhone == "" {
		data := new(smsInitData)
		data.Phone = user.LastPhone
		data.Edit = user.LastPhone == ""
		l.renderRegisterSMS(w, r, authReq, data, nil)
		return
	}
	_, err = l.command.AddHumanOTPSMSWithCheckSucceeded(setUserContext(r.Context(), authReq.UserID, authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, authReq)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	done := &mfaDoneData{
		MFAType: domain.MFATypeOTPSMS,
	}
	l.renderMFAInitDone(w, r, authReq, done)
}

func (l *Login) renderRegisterSMS(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *smsInitData, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	translator := l.getTranslator(r.Context(), authReq)
	data.baseData = l.getBaseData(r, authReq, translator, "InitMFAOTP.Title", "InitMFAOTP.Description", errID, errMessage)
	data.profileData = l.getProfileData(authReq)
	data.MFAType = domain.MFATypeOTPSMS
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplMFASMSInit], data, nil)
}

// handleRegisterSMSCheck handles form submissions of the SMS registration.
// The user can be either in edit mode, where a phone number can be entered / changed.
// If a phone was set, the user can either switch to edit mode, have a resend of the code or verify the code by entering it.
// On successful code verification, the phone will be added to the user as well as his MFA
// and a successful OTP SMS check will be added to the auth request.
func (l *Login) handleRegisterSMSCheck(w http.ResponseWriter, r *http.Request) {
	formData := new(smsInitFormData)
	authReq, err := l.getAuthRequestAndParseData(r, formData)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	ctx := setUserContext(r.Context(), authReq.UserID, authReq.UserOrgID)
	// save the current state
	data := &smsInitData{Phone: formData.Phone}

	if formData.Edit {
		data.Edit = true
		l.renderRegisterSMS(w, r, authReq, data, err)
		return
	}

	if formData.Resend {
		_, err = l.command.CreateHumanPhoneVerificationCode(ctx, authReq.UserID, authReq.UserOrgID)
		l.renderRegisterSMS(w, r, authReq, data, err)
		return
	}

	// if the user is currently in edit mode,
	// he can either change the phone number
	// or just return to the code verification again
	if formData.Code == "" {
		data.Phone = formData.NewPhone
		if formData.NewPhone != formData.Phone {
			_, err = l.command.ChangeUserPhone(ctx, authReq.UserID, formData.NewPhone, l.userCodeAlg)
			if err != nil {
				// stay in edit more
				data.Edit = true
			}
		}
		l.renderRegisterSMS(w, r, authReq, data, err)
		return
	}

	_, err = l.command.VerifyUserPhone(ctx, authReq.UserID, formData.Code, l.userCodeAlg)
	if err != nil {
		l.renderRegisterSMS(w, r, authReq, data, err)
		return
	}
	_, err = l.command.AddHumanOTPSMSWithCheckSucceeded(ctx, authReq.UserID, authReq.UserOrgID, authReq)
	if err != nil {
		l.renderRegisterSMS(w, r, authReq, data, err)
		return
	}
	done := &mfaDoneData{
		MFAType: domain.MFATypeOTPSMS,
	}
	l.renderMFAInitDone(w, r, authReq, done)
}
