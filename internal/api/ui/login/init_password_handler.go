package login

import (
	"fmt"
	"net/http"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
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

func InitPasswordLink(origin, userID, code, orgID string) string {
	return fmt.Sprintf("%s%s?userID=%s&code=%s&orgID=%s", externalLink(origin), EndpointInitPassword, userID, code, orgID)
}

func (l *Login) handleInitPassword(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue(queryInitPWUserID)
	code := r.FormValue(queryInitPWCode)
	l.renderInitPassword(w, r, nil, userID, code, nil)
}

func (l *Login) handleInitPasswordCheck(w http.ResponseWriter, r *http.Request) {
	data := new(initPasswordFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	if data.Resend {
		l.resendPasswordSet(w, r, authReq)
		return
	}
	l.checkPWCode(w, r, authReq, data)
}

func (l *Login) checkPWCode(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *initPasswordFormData) {
	if data.Password != data.PasswordConfirm {
		err := errors.ThrowInvalidArgument(nil, "VIEW-KaGue", "Errors.User.Password.ConfirmationWrong")
		l.renderInitPassword(w, r, authReq, data.UserID, data.Code, err)
		return
	}
	userOrg := ""
	if authReq != nil {
		userOrg = authReq.UserOrgID
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	_, err := l.command.SetPasswordWithVerifyCode(setContext(r.Context(), userOrg), userOrg, data.UserID, data.Code, data.Password, userAgentID)
	if err != nil {
		l.renderInitPassword(w, r, authReq, data.UserID, "", err)
		return
	}
	l.renderInitPasswordDone(w, r, authReq, userOrg)
}

func (l *Login) resendPasswordSet(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	if authReq == nil {
		l.renderError(w, r, nil, errors.ThrowInternal(nil, "LOGIN-8sn7s", "Errors.AuthRequest.NotFound"))
		return
	}
	userOrg := login
	if authReq != nil {
		userOrg = authReq.UserOrgID
	}
	loginName, err := query.NewUserLoginNamesSearchQuery(authReq.LoginName)
	if err != nil {
		l.renderInitPassword(w, r, authReq, authReq.UserID, "", err)
		return
	}
	passwordCodeGenerator, err := l.query.InitEncryptionGenerator(r.Context(), domain.SecretGeneratorTypePasswordResetCode, l.userCodeAlg)
	if err != nil {
		l.renderInitPassword(w, r, authReq, authReq.UserID, "", err)
		return
	}
	user, err := l.query.GetUser(setContext(r.Context(), userOrg), false, loginName)
	if err != nil {
		l.renderInitPassword(w, r, authReq, authReq.UserID, "", err)
		return
	}
	_, err = l.command.RequestSetPassword(setContext(r.Context(), userOrg), user.ID, user.ResourceOwner, domain.NotificationTypeEmail, passwordCodeGenerator)
	l.renderInitPassword(w, r, authReq, authReq.UserID, "", err)
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
		baseData:    l.getBaseData(r, authReq, "InitPassword.Title", "InitPassword.Description", errID, errMessage),
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
	data := l.getUserData(r, authReq, "InitPasswordDone.Title", "InitPasswordDone.Description", "", "")
	translator := l.getTranslator(r.Context(), authReq)
	if authReq == nil {
		l.customTexts(r.Context(), translator, orgID)
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplInitPasswordDone], data, nil)
}
