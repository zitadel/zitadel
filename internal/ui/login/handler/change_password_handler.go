package handler

import (
	"github.com/caos/zitadel/internal/domain"
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
)

const (
	tmplChangePassword     = "changepassword"
	tmplChangePasswordDone = "changepassworddone"
)

type changePasswordData struct {
	OldPassword             string `schema:"change-old-password"`
	NewPassword             string `schema:"change-new-password"`
	NewPasswordConfirmation string `schema:"change-password-confirmation"`
}

func (l *Login) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	data := new(changePasswordData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	_, err = l.command.ChangePassword(setContext(r.Context(), authReq.UserOrgID), authReq.UserOrgID, authReq.UserID, data.OldPassword, data.NewPassword, userAgentID)
	if err != nil {
		l.renderChangePassword(w, r, authReq, err)
		return
	}
	l.renderChangePasswordDone(w, r, authReq)
}

func (l *Login) renderChangePassword(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	data := passwordData{
		baseData:    l.getBaseData(r, authReq, "Change Password", errType, errMessage),
		profileData: l.getProfileData(authReq),
	}
	policy, description, _ := l.getPasswordComplexityPolicy(r, authReq.UserOrgID)
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
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplChangePassword], data, nil)
}

func (l *Login) renderChangePasswordDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	var errType, errMessage string
	data := l.getUserData(r, authReq, "Password Change Done", errType, errMessage)
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplChangePasswordDone], data, nil)
}
