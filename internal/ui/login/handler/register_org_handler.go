package handler

import (
	"github.com/caos/zitadel/internal/domain"
	"net/http"

	caos_errs "github.com/caos/zitadel/internal/errors"
)

const (
	tmplRegisterOrg = "registerorg"
)

type registerOrgFormData struct {
	RegisterOrgName string `schema:"orgname"`
	Email           string `schema:"email"`
	Username        string `schema:"username"`
	Firstname       string `schema:"firstname"`
	Lastname        string `schema:"lastname"`
	Password        string `schema:"register-password"`
	Password2       string `schema:"register-password-confirmation"`
	TermsConfirm    bool   `schema:"terms-confirm"`
}

type registerOrgData struct {
	baseData
	registerOrgFormData
	PasswordPolicyDescription string
	MinLength                 uint64
	HasUppercase              string
	HasLowercase              string
	HasNumber                 string
	HasSymbol                 string
	UserLoginMustBeDomain     bool
	IamDomain                 string
}

func (l *Login) handleRegisterOrg(w http.ResponseWriter, r *http.Request) {
	data := new(registerOrgFormData)
	authRequest, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authRequest, err)
		return
	}
	l.renderRegisterOrg(w, r, authRequest, data, nil)
}

func (l *Login) handleRegisterOrgCheck(w http.ResponseWriter, r *http.Request) {
	data := new(registerOrgFormData)
	authRequest, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authRequest, err)
		return
	}
	if data.Password != data.Password2 {
		err := caos_errs.ThrowInvalidArgument(nil, "VIEW-KaGue", "Errors.User.Password.ConfirmationWrong")
		l.renderRegisterOrg(w, r, authRequest, data, err)
		return
	}

	_, err = l.command.SetUpOrg(setContext(r.Context(), ""), data.toOrgDomain(), data.toUserDomain())
	if err != nil {
		l.renderRegisterOrg(w, r, authRequest, data, err)
		return
	}
	if authRequest == nil {
		http.Redirect(w, r, l.zitadelURL, http.StatusFound)
		return
	}
	l.renderNextStep(w, r, authRequest)
}

func (l *Login) renderRegisterOrg(w http.ResponseWriter, r *http.Request, authRequest *domain.AuthRequest, formData *registerOrgFormData, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	if formData == nil {
		formData = new(registerOrgFormData)
	}

	data := registerOrgData{
		baseData:            l.getBaseData(r, authRequest, "Register", errType, errMessage),
		registerOrgFormData: *formData,
	}
	pwPolicy, description, _ := l.getPasswordComplexityPolicy(r, "0")
	if pwPolicy != nil {
		data.PasswordPolicyDescription = description
		data.MinLength = pwPolicy.MinLength
		if pwPolicy.HasUppercase {
			data.HasUppercase = UpperCaseRegex
		}
		if pwPolicy.HasLowercase {
			data.HasLowercase = LowerCaseRegex
		}
		if pwPolicy.HasSymbol {
			data.HasSymbol = SymbolRegex
		}
		if pwPolicy.HasNumber {
			data.HasNumber = NumberRegex
		}
	}
	orgPolicy, err := l.getDefaultOrgIamPolicy(r)
	if orgPolicy != nil {
		data.UserLoginMustBeDomain = orgPolicy.UserLoginMustBeDomain
		data.IamDomain = orgPolicy.IAMDomain
	}

	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplRegisterOrg], data, nil)
}

func (d registerOrgFormData) toUserDomain() *domain.Human {
	if d.Username == "" {
		d.Username = d.Email
	}
	return &domain.Human{
		Username: d.Username,
		Profile: &domain.Profile{
			FirstName: d.Firstname,
			LastName:  d.Lastname,
		},
		Password: &domain.Password{
			SecretString: d.Password,
		},
		Email: &domain.Email{
			EmailAddress: d.Email,
		},
	}
}

func (d registerOrgFormData) toOrgDomain() *domain.Org {
	return &domain.Org{
		Name: d.RegisterOrgName,
	}
}
