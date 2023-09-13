package login

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
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

func InitUserLink(origin, userID, loginName, code, orgID string, passwordSet bool) string {
	return fmt.Sprintf("%s%s?userID=%s&loginname=%s&code=%s&orgID=%s&passwordset=%t", externalLink(origin), EndpointInitUser, userID, loginName, code, orgID, passwordSet)
}

func (l *Login) handleInitUser(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue(queryInitUserUserID)
	code := r.FormValue(queryInitUserCode)
	loginName := r.FormValue(queryInitUserLoginName)
	passwordSet, _ := strconv.ParseBool(r.FormValue(queryInitUserPassword))
	l.renderInitUser(w, r, nil, userID, loginName, code, passwordSet, nil)
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
		err := caos_errs.ThrowInvalidArgument(nil, "VIEW-fsdfd", "Errors.User.Password.ConfirmationWrong")
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
	err = l.command.HumanVerifyInitCode(setContext(r.Context(), userOrgID), data.UserID, userOrgID, data.Code, data.Password, initCodeGenerator)
	if err != nil {
		l.renderInitUser(w, r, authReq, data.UserID, data.LoginName, "", data.PasswordSet, err)
		return
	}
	l.renderInitUserDone(w, r, authReq, userOrgID)
}

func (l *Login) resendUserInit(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID string, loginName string, showPassword bool) {
	userOrgID := ""
	if authReq != nil {
		userOrgID = authReq.UserOrgID
	}
	initCodeGenerator, err := l.query.InitEncryptionGenerator(r.Context(), domain.SecretGeneratorTypeInitCode, l.userCodeAlg)
	if err != nil {
		l.renderInitUser(w, r, authReq, userID, loginName, "", showPassword, err)
		return
	}
	_, err = l.command.ResendInitialMail(setContext(r.Context(), userOrgID), userID, "", userOrgID, initCodeGenerator)
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
		baseData:    l.getBaseData(r, authReq, "InitUser.Title", "InitUser.Description", errID, errMessage),
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
		user, err := l.query.GetUserByID(r.Context(), false, userID, false)
		if err == nil {
			l.customTexts(r.Context(), translator, user.ResourceOwner)
		}
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplInitUser], data, nil)
}

func (l *Login) renderInitUserDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, orgID string) {
	data := l.getUserData(r, authReq, "InitUserDone.Title", "InitUserDone.Description", "", "")
	translator := l.getTranslator(r.Context(), authReq)
	if authReq == nil {
		l.customTexts(r.Context(), translator, orgID)
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplInitUserDone], data, nil)
}
