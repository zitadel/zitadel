package handler

import (
	"net/http"

	auth_model "github.com/caos/zitadel/internal/auth/model"
	"github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
	org_model "github.com/caos/zitadel/internal/org/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
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

	registerOrg := &auth_model.RegisterOrg{
		User: data.toUserModel(),
		Org:  data.toOrgModel(),
	}
	user, err := l.authRepo.RegisterOrg(setContext(r.Context(), ""), registerOrg)
	if err != nil {
		l.renderRegisterOrg(w, r, authRequest, data, err)
		return
	}
	if authRequest == nil {
		http.Redirect(w, r, l.zitadelURL, http.StatusFound)
		return
	}
	authRequest.LoginName = user.PreferredLoginName
	l.renderNextStep(w, r, authRequest)
}

func (l *Login) renderRegisterOrg(w http.ResponseWriter, r *http.Request, authRequest *model.AuthRequest, formData *registerOrgFormData, err error) {
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

func (d registerOrgFormData) toUserModel() *usr_model.User {
	if d.Username == "" {
		d.Username = d.Email
	}
	return &usr_model.User{
		UserName: d.Username,
		Human: &usr_model.Human{
			Profile: &usr_model.Profile{
				FirstName: d.Firstname,
				LastName:  d.Lastname,
			},
			Password: &usr_model.Password{
				SecretString: d.Password,
			},
			Email: &usr_model.Email{
				EmailAddress: d.Email,
			},
		},
	}
}

func (d registerOrgFormData) toOrgModel() *org_model.Org {
	return &org_model.Org{
		Name: d.RegisterOrgName,
	}
}
