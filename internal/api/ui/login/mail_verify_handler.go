package login

import (
	"net/http"
	"net/url"

	"github.com/zitadel/zitadel/internal/domain"
)

const (
	queryCode   = "code"
	queryUserID = "userID"

	tmplMailVerification = "mail_verification"
	tmplMailVerified     = "mail_verified"
)

type mailVerificationFormData struct {
	Code   string `schema:"code"`
	UserID string `schema:"userID"`
	Resend bool   `schema:"resend"`
}

type mailVerificationData struct {
	baseData
	profileData
	UserID string
}

func MailVerificationLink(origin, userID, code, orgID, authRequestID string) string {
	v := url.Values{}
	v.Set(queryUserID, userID)
	v.Set(queryCode, code)
	v.Set(queryOrgID, orgID)
	v.Set(QueryAuthRequestID, authRequestID)
	return externalLink(origin) + EndpointMailVerification + "?" + v.Encode()
}

func (l *Login) handleMailVerification(w http.ResponseWriter, r *http.Request) {
	authReq := l.checkOptionalAuthRequestOfEmailLinks(r)
	userID := r.FormValue(queryUserID)
	code := r.FormValue(queryCode)
	if code != "" {
		l.checkMailCode(w, r, authReq, userID, code)
		return
	}
	l.renderMailVerification(w, r, authReq, userID, nil)
}

func (l *Login) handleMailVerificationCheck(w http.ResponseWriter, r *http.Request) {
	data := new(mailVerificationFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if !data.Resend {
		l.checkMailCode(w, r, authReq, data.UserID, data.Code)
		return
	}
	var userOrg, authReqID string
	if authReq != nil {
		userOrg = authReq.UserOrgID
		authReqID = authReq.ID
	}
	emailCodeGenerator, err := l.query.InitEncryptionGenerator(r.Context(), domain.SecretGeneratorTypeVerifyEmailCode, l.userCodeAlg)
	if err != nil {
		l.checkMailCode(w, r, authReq, data.UserID, data.Code)
		return
	}
	_, err = l.command.CreateHumanEmailVerificationCode(setContext(r.Context(), userOrg), data.UserID, userOrg, emailCodeGenerator, authReqID)
	l.renderMailVerification(w, r, authReq, data.UserID, err)
}

func (l *Login) checkMailCode(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID, code string) {
	userOrg := ""
	if authReq != nil {
		userID = authReq.UserID
		userOrg = authReq.UserOrgID
	}
	emailCodeGenerator, err := l.query.InitEncryptionGenerator(r.Context(), domain.SecretGeneratorTypeVerifyEmailCode, l.userCodeAlg)
	if err != nil {
		l.renderMailVerification(w, r, authReq, userID, err)
		return
	}
	_, err = l.command.VerifyHumanEmail(setContext(r.Context(), userOrg), userID, code, userOrg, emailCodeGenerator)
	if err != nil {
		l.renderMailVerification(w, r, authReq, userID, err)
		return
	}
	l.renderMailVerified(w, r, authReq, userOrg)
}

func (l *Login) renderMailVerification(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID string, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	if userID == "" {
		userID = authReq.UserID
	}

	translator := l.getTranslator(r.Context(), authReq)
	data := mailVerificationData{
		baseData:    l.getBaseData(r, authReq, translator, "EmailVerification.Title", "EmailVerification.Description", errID, errMessage),
		UserID:      userID,
		profileData: l.getProfileData(authReq),
	}
	if authReq == nil {
		user, err := l.query.GetUserByID(r.Context(), false, userID)
		if err == nil {
			l.customTexts(r.Context(), translator, user.ResourceOwner)
		}
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplMailVerification], data, nil)
}

func (l *Login) renderMailVerified(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, orgID string) {
	translator := l.getTranslator(r.Context(), authReq)
	data := mailVerificationData{
		baseData:    l.getBaseData(r, authReq, translator, "EmailVerificationDone.Title", "EmailVerificationDone.Description", "", ""),
		profileData: l.getProfileData(authReq),
	}
	if authReq == nil {
		l.customTexts(r.Context(), translator, orgID)
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplMailVerified], data, nil)
}
