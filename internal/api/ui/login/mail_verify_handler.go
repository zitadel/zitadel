package login

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/zitadel/logging"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	queryCode   = "code"
	queryUserID = "userID"

	tmplMailVerification = "mail_verification"
	tmplMailVerified     = "mail_verified"
)

type mailVerificationFormData struct {
	Code            string `schema:"code"`
	UserID          string `schema:"userID"`
	Resend          bool   `schema:"resend"`
	PasswordInit    bool   `schema:"passwordInit"`
	Password        string `schema:"password"`
	PasswordConfirm string `schema:"passwordconfirm"`
}

type mailVerificationData struct {
	baseData
	profileData
	UserID       string
	Code         string
	PasswordInit bool
	MinLength    uint64
	HasUppercase string
	HasLowercase string
	HasNumber    string
	HasSymbol    string
}

func MailVerificationLink(origin, userID, code, orgID, authRequestID string) string {
	v := url.Values{}
	v.Set(queryUserID, userID)
	v.Set(queryCode, code)
	v.Set(queryOrgID, orgID)
	v.Set(QueryAuthRequestID, authRequestID)
	return externalLink(origin) + EndpointMailVerification + "?" + v.Encode()
}

func MailVerificationLinkTemplate(origin, userID, orgID, authRequestID string) string {
	return fmt.Sprintf("%s%s?%s=%s&%s=%s&%s=%s&%s=%s",
		externalLink(origin), EndpointMailVerification,
		queryUserID, userID,
		queryCode, "{{.Code}}",
		queryOrgID, orgID,
		QueryAuthRequestID, authRequestID)
}

func (l *Login) handleMailVerification(w http.ResponseWriter, r *http.Request) {
	authReq := l.checkOptionalAuthRequestOfEmailLinks(r)
	userID := r.FormValue(queryUserID)
	code := r.FormValue(queryCode)
	if userID == "" && authReq == nil {
		l.renderError(w, r, authReq, nil)
		return
	}
	if userID == "" {
		userID = authReq.UserID
	}
	passwordInit := l.checkUserNoFirstFactor(r.Context(), userID)
	if code != "" && !passwordInit {
		l.checkMailCode(w, r, authReq, userID, code, "")
		return
	}
	l.renderMailVerification(w, r, authReq, userID, code, passwordInit, nil)
}

func (l *Login) checkUserNoFirstFactor(ctx context.Context, userID string) bool {
	authMethods, err := l.query.ListUserAuthMethodTypes(setUserContext(ctx, userID, ""), userID, false, false, "")
	if err != nil {
		logging.WithFields("userID", userID).OnError(err).Warn("unable to load user's auth methods for mail verification")
		return false
	}
	return !slices.ContainsFunc(authMethods.AuthMethodTypes, func(m domain.UserAuthMethodType) bool {
		return m == domain.UserAuthMethodTypeIDP ||
			m == domain.UserAuthMethodTypePassword ||
			m == domain.UserAuthMethodTypePasswordless
	})
}

func (l *Login) handleMailVerificationCheck(w http.ResponseWriter, r *http.Request) {
	data := new(mailVerificationFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if !data.Resend {
		if data.PasswordInit && data.Password != data.PasswordConfirm {
			err := zerrors.ThrowInvalidArgument(nil, "VIEW-fsdfd", "Errors.User.Password.ConfirmationWrong")
			l.renderMailVerification(w, r, authReq, data.UserID, data.Code, data.PasswordInit, err)
			return
		}
		l.checkMailCode(w, r, authReq, data.UserID, data.Code, data.Password)
		return
	}
	var userOrg, authReqID string
	if authReq != nil {
		userOrg = authReq.UserOrgID
		authReqID = authReq.ID
	}
	emailCodeGenerator, err := l.query.InitEncryptionGenerator(r.Context(), domain.SecretGeneratorTypeVerifyEmailCode, l.userCodeAlg)
	if err != nil {
		l.renderMailVerification(w, r, authReq, data.UserID, "", data.PasswordInit, err)
		return
	}
	_, err = l.command.CreateHumanEmailVerificationCode(setContext(r.Context(), userOrg), data.UserID, userOrg, emailCodeGenerator, authReqID)
	l.renderMailVerification(w, r, authReq, data.UserID, "", data.PasswordInit, err)
}

func (l *Login) checkMailCode(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID, code, password string) {
	userOrg := ""
	if authReq != nil {
		userID = authReq.UserID
		userOrg = authReq.UserOrgID
	}
	emailCodeGenerator, err := l.query.InitEncryptionGenerator(r.Context(), domain.SecretGeneratorTypeVerifyEmailCode, l.userCodeAlg)
	if err != nil {
		l.renderMailVerification(w, r, authReq, userID, "", password != "", err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	_, err = l.command.VerifyHumanEmail(setContext(r.Context(), userOrg), userID, code, userOrg, password, userAgentID, emailCodeGenerator)
	if err != nil {
		l.renderMailVerification(w, r, authReq, userID, "", password != "", err)
		return
	}
	l.renderMailVerified(w, r, authReq, userOrg)
}

func (l *Login) renderMailVerification(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID, code string, passwordInit bool, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	if userID == "" && authReq != nil {
		userID = authReq.UserID
	}

	translator := l.getTranslator(r.Context(), authReq)
	data := mailVerificationData{
		baseData:     l.getBaseData(r, authReq, translator, "EmailVerification.Title", "EmailVerification.Description", errID, errMessage),
		UserID:       userID,
		profileData:  l.getProfileData(authReq),
		Code:         code,
		PasswordInit: passwordInit,
	}
	if passwordInit {
		policy := l.getPasswordComplexityPolicyByUserID(r, userID)
		if policy != nil {
			data.MinLength = policy.MinLength
			if policy.HasUppercase {
				data.HasUppercase = UpperCaseRegex
			}
			if policy.HasLowercase {
				data.HasLowercase = LowerCaseRegex
			}
			if policy.HasSymbol {
				data.HasSymbol = SymbolRegex
			}
			if policy.HasNumber {
				data.HasNumber = NumberRegex
			}
		}
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
