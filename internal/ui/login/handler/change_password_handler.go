package handler

import (
	"github.com/caos/zitadel/internal/errors"
	"net/http"

	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	tmplChangePassword     = "changepassword"
	tmplChangePasswordDone = "changepassworddone"
)

type changePasswordData struct {
	OldPassword             string `schema:"old_password"`
	NewPassword             string `schema:"new_password"`
	NewPasswordConfirmation string `schema:"new_password_confirmation"`
}

func (l *Login) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	data := new(changePasswordData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if data.NewPassword != data.NewPasswordConfirmation {
		err := errors.ThrowInvalidArgument(nil, "ERR-sj2Sq", "Errors.User.Password.ConfirmationWrong")
		l.renderChangePassword(w, r, authReq, err)
		return
	}
	err = l.checkPasswordComplexityPolicy(data.NewPassword, r, authReq)
	if err != nil {
		l.renderChangePassword(w, r, authReq, err)
	}

	err = l.authRepo.ChangePassword(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, data.OldPassword, data.NewPassword)
	if err != nil {
		l.renderChangePassword(w, r, authReq, err)
		return
	}
	l.renderChangePasswordDone(w, r, authReq)
}

func (l *Login) renderChangePassword(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	policy, description, _ := l.getPasswordComplexityPolicy(r, authReq)
	data := passwordData{
		baseData:                  l.getBaseData(r, authReq, "Change Password", errType, errMessage),
		LoginName:                 authReq.LoginName,
		PasswordPolicyDescription: description,
		MinLength:                 policy.MinLength,
	}
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
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplChangePassword], data, nil)
}

func (l *Login) renderChangePasswordDone(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest) {
	var errType, errMessage string
	data := userData{
		baseData:  l.getBaseData(r, authReq, "Password Change Done", errType, errMessage),
		LoginName: authReq.LoginName,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplChangePasswordDone], data, nil)
}
