package handler

import (
	"github.com/caos/zitadel/internal/domain"
	"net/http"
	"strconv"

	caos_errs "github.com/caos/zitadel/internal/errors"
)

const (
	queryInitUserCode     = "code"
	queryInitUserUserID   = "userID"
	queryInitUserPassword = "passwordset"

	tmplInitUser     = "inituser"
	tmplInitUserDone = "inituserdone"
)

type initUserFormData struct {
	Code            string `schema:"code"`
	Password        string `schema:"password"`
	PasswordConfirm string `schema:"passwordconfirm"`
	UserID          string `schema:"userID"`
	PasswordSet     bool   `schema:"passwordSet"`
	Resend          bool   `schema:"resend"`
}

type initUserData struct {
	baseData
	profileData
	Code                      string
	UserID                    string
	PasswordSet               bool
	PasswordPolicyDescription string
	MinLength                 uint64
	HasUppercase              string
	HasLowercase              string
	HasNumber                 string
	HasSymbol                 string
}

func (l *Login) handleInitUser(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue(queryInitUserUserID)
	code := r.FormValue(queryInitUserCode)
	passwordSet, _ := strconv.ParseBool(r.FormValue(queryInitUserPassword))
	l.renderInitUser(w, r, nil, userID, code, passwordSet, nil)
}

func (l *Login) handleInitUserCheck(w http.ResponseWriter, r *http.Request) {
	data := new(initUserFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}

	if data.Resend {
		l.resendUserInit(w, r, authReq, data.UserID, data.PasswordSet)
		return
	}
	l.checkUserInitCode(w, r, authReq, data, nil)
}

func (l *Login) checkUserInitCode(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *initUserFormData, err error) {
	if data.Password != data.PasswordConfirm {
		err := caos_errs.ThrowInvalidArgument(nil, "VIEW-fsdfd", "Errors.User.Password.ConfirmationWrong")
		l.renderInitUser(w, r, authReq, data.UserID, data.Code, data.PasswordSet, err)
		return
	}
	userOrgID := ""
	if authReq != nil {
		userOrgID = authReq.UserOrgID
	}
	err = l.command.HumanVerifyInitCode(setContext(r.Context(), userOrgID), data.UserID, userOrgID, data.Code, data.Password)
	if err != nil {
		l.renderInitUser(w, r, authReq, data.UserID, "", data.PasswordSet, err)
		return
	}
	l.renderInitUserDone(w, r, authReq)
}

func (l *Login) resendUserInit(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID string, showPassword bool) {
	userOrgID := ""
	if authReq != nil {
		userOrgID = authReq.UserOrgID
	}
	_, err := l.command.ResendInitialMail(setContext(r.Context(), userOrgID), userID, "", userOrgID)
	l.renderInitUser(w, r, authReq, userID, "", showPassword, err)
}

func (l *Login) renderInitUser(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID, code string, passwordSet bool, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	if authReq != nil {
		userID = authReq.UserID
	}
	data := initUserData{
		baseData:    l.getBaseData(r, authReq, "Init User", errID, errMessage),
		profileData: l.getProfileData(authReq),
		UserID:      userID,
		Code:        code,
		PasswordSet: passwordSet,
	}
	policy, description, _ := l.getPasswordComplexityPolicyByUserID(r, nil, userID)
	if policy != nil {
		data.PasswordPolicyDescription = description
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
	translator := l.getTranslator(authReq)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplInitUser], data, nil)
}

func (l *Login) renderInitUserDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	data := l.getUserData(r, authReq, "User Init Done", "", "")
	l.renderer.RenderTemplate(w, r, l.getTranslator(authReq), l.renderer.Templates[tmplInitUserDone], data, nil)
}
