package login

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	queryInitUserCode      = "code"
	queryInitUserUserID    = "userID"
	queryInitUserLoginName = "loginname"
	queryInitUserPassword  = "passwordset"

	tmplInitUser     = "inituser"
	tmplInitUserDone = "inituserdone"
)

type initUserFormData struct {
	Code            string `schema:"code"`
	LoginName       string `schema:"loginname"`
	Password        string `schema:"password"`
	PasswordConfirm string `schema:"passwordconfirm"`
	UserID          string `schema:"userID"`
	PasswordSet     bool   `schema:"passwordSet"`
	Resend          bool   `schema:"resend"`
}

type initUserData struct {
	baseData
	profileData
	Code         string
	LoginName    string
	UserID       string
	PasswordSet  bool
	MinLength    uint64
	HasUppercase string
	HasLowercase string
	HasNumber    string
	HasSymbol    string
}

func InitUserLink(origin, userID, loginName, code, orgID string, passwordSet bool, authRequestID string) string {
	v := url.Values{}
	v.Set(queryInitUserUserID, userID)
	v.Set(queryInitUserLoginName, loginName)
	v.Set(queryInitUserCode, code)
	v.Set(queryOrgID, orgID)
	v.Set(queryInitUserPassword, strconv.FormatBool(passwordSet))
	v.Set(QueryAuthRequestID, authRequestID)
	return externalLink(origin) + EndpointInitUser + "?" + v.Encode()
}

func InitUserLinkTemplate(origin, userID, orgID, authRequestID string) string {
	return fmt.Sprintf("%s%s?%s=%s&%s=%s&%s=%s&%s=%s&%s=%s&%s=%s",
		externalLink(origin), EndpointInitUser,
		queryInitUserUserID, userID,
		queryInitUserLoginName, "{{.LoginName}}",
		queryInitUserCode, "{{.Code}}",
		queryOrgID, orgID,
		queryInitUserPassword, "{{.PasswordSet}}",
		QueryAuthRequestID, authRequestID)
}

func (l *Login) handleInitUser(w http.ResponseWriter, r *http.Request) {
	authReq := l.checkOptionalAuthRequestOfEmailLinks(r)
	userID := r.FormValue(queryInitUserUserID)
	code := r.FormValue(queryInitUserCode)
	loginName := r.FormValue(queryInitUserLoginName)
	passwordSet, _ := strconv.ParseBool(r.FormValue(queryInitUserPassword))
	l.renderInitUser(w, r, authReq, userID, loginName, code, passwordSet, nil)
}

func (l *Login) handleInitUserCheck(w http.ResponseWriter, r *http.Request) {
	data := new(initUserFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}

	if data.Resend {
		l.resendUserInit(w, r, authReq, data.UserID, data.LoginName, data.PasswordSet)
		return
	}
	l.checkUserInitCode(w, r, authReq, data, nil)
}

func (l *Login) checkUserInitCode(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *initUserFormData, err error) {
	if data.Password != data.PasswordConfirm {
		err := zerrors.ThrowInvalidArgument(nil, "VIEW-fsdfd", "Errors.User.Password.ConfirmationWrong")
		l.renderInitUser(w, r, authReq, data.UserID, data.LoginName, data.Code, data.PasswordSet, err)
		return
	}
	userOrgID := ""
	if authReq != nil {
		userOrgID = authReq.UserOrgID
	}
	initCodeGenerator, err := l.query.InitEncryptionGenerator(r.Context(), domain.SecretGeneratorTypeInitCode, l.userCodeAlg)
	if err != nil {
		l.renderInitUser(w, r, authReq, data.UserID, data.LoginName, "", data.PasswordSet, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.command.HumanVerifyInitCode(setContext(r.Context(), userOrgID), data.UserID, userOrgID, data.Code, data.Password, userAgentID, initCodeGenerator)
	if err != nil {
		l.renderInitUser(w, r, authReq, data.UserID, data.LoginName, "", data.PasswordSet, err)
		return
	}
	l.renderInitUserDone(w, r, authReq, userOrgID)
}

func (l *Login) resendUserInit(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID string, loginName string, showPassword bool) {
	var userOrgID, authRequestID string
	if authReq != nil {
		userOrgID = authReq.UserOrgID
		authRequestID = authReq.ID
	}
	initCodeGenerator, err := l.query.InitEncryptionGenerator(r.Context(), domain.SecretGeneratorTypeInitCode, l.userCodeAlg)
	if err != nil {
		l.renderInitUser(w, r, authReq, userID, loginName, "", showPassword, err)
		return
	}
	_, err = l.command.ResendInitialMail(setContext(r.Context(), userOrgID), userID, "", userOrgID, initCodeGenerator, authRequestID)
	l.renderInitUser(w, r, authReq, userID, loginName, "", showPassword, err)
}

func (l *Login) renderInitUser(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID, loginName string, code string, passwordSet bool, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	if authReq != nil {
		userID = authReq.UserID
	}

	translator := l.getTranslator(r.Context(), authReq)
	data := initUserData{
		baseData:    l.getBaseData(r, authReq, translator, "InitUser.Title", "InitUser.Description", errID, errMessage),
		profileData: l.getProfileData(authReq),
		UserID:      userID,
		Code:        code,
		PasswordSet: passwordSet,
	}
	// if the user clicked on the link in the mail, we need to make sure the loginName is rendered
	if authReq == nil {
		data.LoginName = loginName
		data.UserName = loginName
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
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplInitUser], data, nil)
}

func (l *Login) renderInitUserDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, orgID string) {
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, translator, "InitUserDone.Title", "InitUserDone.Description", "", "")
	if authReq == nil {
		l.customTexts(r.Context(), translator, orgID)
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplInitUserDone], data, nil)
}
