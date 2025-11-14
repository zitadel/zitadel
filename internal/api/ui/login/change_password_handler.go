package login

import (
	"net/http"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
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
	authReq, err := l.ensureAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	_, err = l.command.ChangePassword(setContext(r.Context(), authReq.UserOrgID), authReq.UserOrgID, authReq.UserID, data.OldPassword, data.NewPassword, userAgentID, false)
	if err != nil {
		l.renderChangePassword(w, r, authReq, err)
		return
	}
	l.renderChangePasswordDone(w, r, authReq)
}

func (l *Login) renderChangePassword(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	translator := l.getTranslator(r.Context(), authReq)
	if authReq == nil || len(authReq.PossibleSteps) < 1 {
		l.renderError(w, r, authReq, err)
		return
	}
	step, ok := authReq.PossibleSteps[0].(*domain.ChangePasswordStep)
	if !ok {
		l.renderError(w, r, authReq, err)
		return
	}
	data := passwordData{
		baseData:    l.getBaseData(r, authReq, translator, "PasswordChange.Title", "PasswordChange.Description", err),
		profileData: l.getProfileData(authReq),
		Expired:     step.Expired,
	}
	policy := l.getPasswordComplexityPolicy(r, authReq.UserOrgID)
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
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplChangePassword], data, nil)
}

func (l *Login) renderChangePasswordDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, translator, "PasswordChange.Title", "PasswordChange.Description", nil)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplChangePasswordDone], data, nil)
}
