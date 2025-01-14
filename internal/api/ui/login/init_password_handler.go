package login

import (
	"fmt"
	"net/http"
	"net/url"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	queryInitPWCode   = "code"
	queryInitPWUserID = "userID"

	tmplInitPassword     = "initpassword"
	tmplInitPasswordDone = "initpassworddone"
)

type initPasswordFormData struct {
	Code            string `schema:"code"`
	Password        string `schema:"password"`
	PasswordConfirm string `schema:"passwordconfirm"`
	UserID          string `schema:"userID"`
	OrgID           string `schema:"orgID"`
	Resend          bool   `schema:"resend"`
}

type initPasswordData struct {
	baseData
	profileData
	Code         string
	UserID       string
	MinLength    uint64
	HasUppercase string
	HasLowercase string
	HasNumber    string
	HasSymbol    string
}

func InitPasswordLink(origin, userID, code, orgID, authRequestID string) string {
	v := url.Values{}
	v.Set(queryInitPWUserID, userID)
	v.Set(queryInitPWCode, code)
	v.Set(queryOrgID, orgID)
	v.Set(QueryAuthRequestID, authRequestID)
	return externalLink(origin) + EndpointInitPassword + "?" + v.Encode()
}

func InitPasswordLinkTemplate(origin, userID, orgID, authRequestID string) string {
	return fmt.Sprintf("%s%s?%s=%s&%s=%s&%s=%s&%s=%s",
		externalLink(origin), EndpointInitPassword,
		queryInitPWUserID, userID,
		queryInitPWCode, "{{.Code}}",
		queryOrgID, orgID,
		QueryAuthRequestID, authRequestID)
}

func (l *Login) handleInitPassword(w http.ResponseWriter, r *http.Request) {
	authReq := l.checkOptionalAuthRequestOfEmailLinks(r)
	userID := r.FormValue(queryInitPWUserID)
	code := r.FormValue(queryInitPWCode)
	l.renderInitPassword(w, r, authReq, userID, code, nil)
}

func (l *Login) handleInitPasswordCheck(w http.ResponseWriter, r *http.Request) {
	data := new(initPasswordFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	if data.Resend {
		l.resendPasswordSet(w, r, authReq, data)
		return
	}
	l.checkPWCode(w, r, authReq, data)
}

func (l *Login) checkPWCode(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *initPasswordFormData) {
	if data.Password != data.PasswordConfirm {
		err := zerrors.ThrowInvalidArgument(nil, "VIEW-KaGue", "Errors.User.Password.ConfirmationWrong")
		l.renderInitPassword(w, r, authReq, data.UserID, data.Code, err)
		return
	}
	userOrg := ""
	if authReq != nil {
		userOrg = authReq.UserOrgID
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	_, err := l.command.SetPasswordWithVerifyCode(setContext(r.Context(), userOrg), userOrg, data.UserID, data.Code, data.Password, userAgentID, false)
	if err != nil {
		l.renderInitPassword(w, r, authReq, data.UserID, "", err)
		return
	}
	l.renderInitPasswordDone(w, r, authReq, userOrg)
}

func (l *Login) resendPasswordSet(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *initPasswordFormData) {
	userOrg := data.OrgID
	userID := data.UserID
	var authReqID string
	if authReq != nil {
		userOrg = authReq.UserOrgID
		userID = authReq.UserID
		authReqID = authReq.ID
	}
	_, err := l.command.RequestSetPassword(setContext(r.Context(), userOrg), userID, userOrg, domain.NotificationTypeEmail, authReqID)
	l.renderInitPassword(w, r, authReq, userID, "", err)
}

func (l *Login) renderInitPassword(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID, code string, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	if userID == "" && authReq != nil {
		userID = authReq.UserID
	}

	translator := l.getTranslator(r.Context(), authReq)

	data := initPasswordData{
		baseData:    l.getBaseData(r, authReq, translator, "InitPassword.Title", "InitPassword.Description", errID, errMessage),
		profileData: l.getProfileData(authReq),
		UserID:      userID,
		Code:        code,
	}
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
	if authReq == nil {
		user, err := l.query.GetUserByID(r.Context(), false, userID)
		if err == nil {
			l.customTexts(r.Context(), translator, user.ResourceOwner)
		}
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplInitPassword], data, nil)
}

func (l *Login) renderInitPasswordDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, orgID string) {
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, translator, "InitPasswordDone.Title", "InitPasswordDone.Description", "", "")
	if authReq == nil {
		l.customTexts(r.Context(), translator, orgID)
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplInitPasswordDone], data, nil)
}
