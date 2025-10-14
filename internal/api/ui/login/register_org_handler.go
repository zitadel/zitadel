package login

import (
	"net/http"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	tmplRegisterOrg = "registerorg"
)

type registerOrgFormData struct {
	RegisterOrgName string              `schema:"orgname"`
	Email           domain.EmailAddress `schema:"email"`
	Username        string              `schema:"username"`
	Firstname       string              `schema:"firstname"`
	Lastname        string              `schema:"lastname"`
	Password        string              `schema:"register-password"`
	Password2       string              `schema:"register-password-confirmation"`
	TermsConfirm    bool                `schema:"terms-confirm"`
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
	restrictions, err := l.query.GetInstanceRestrictions(r.Context())
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	if restrictions.DisallowPublicOrgRegistration {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data := new(registerOrgFormData)
	authRequest, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authRequest, err)
		return
	}
	l.renderRegisterOrg(w, r, authRequest, data, nil)
}

func (l *Login) handleRegisterOrgCheck(w http.ResponseWriter, r *http.Request) {
	restrictions, err := l.query.GetInstanceRestrictions(r.Context())
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	if restrictions.DisallowPublicOrgRegistration {
		w.WriteHeader(http.StatusConflict)
		return
	}
	data := new(registerOrgFormData)
	authRequest, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authRequest, err)
		return
	}
	if data.Password != data.Password2 {
		err := zerrors.ThrowInvalidArgument(nil, "VIEW-KaGue", "Errors.User.Password.ConfirmationWrong")
		l.renderRegisterOrg(w, r, authRequest, data, err)
		return
	}

	ctx := setContext(r.Context(), "")
	userIDs, err := l.getClaimedUserIDsOfOrgDomain(ctx, data.RegisterOrgName)
	if err != nil {
		l.renderRegisterOrg(w, r, authRequest, data, err)
		return
	}
	_, err = l.command.SetUpOrg(ctx, data.toCommandOrg(), true, userIDs...)
	if err != nil {
		l.renderRegisterOrg(w, r, authRequest, data, err)
		return
	}
	if authRequest == nil {
		l.defaultRedirect(w, r)
		return
	}
	l.renderNextStep(w, r, authRequest)
}

func (l *Login) renderRegisterOrg(w http.ResponseWriter, r *http.Request, authRequest *domain.AuthRequest, formData *registerOrgFormData, err error) {
	if formData == nil {
		formData = new(registerOrgFormData)
	}
	translator := l.getTranslator(r.Context(), authRequest)
	data := registerOrgData{
		baseData:            l.getBaseData(r, authRequest, translator, "RegistrationOrg.Title", "RegistrationOrg.Description", err),
		registerOrgFormData: *formData,
	}
	pwPolicy := l.getPasswordComplexityPolicy(r, "0")
	if pwPolicy != nil {
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
	orgPolicy, _ := l.getDefaultDomainPolicy(r)
	if orgPolicy != nil {
		data.UserLoginMustBeDomain = orgPolicy.UserLoginMustBeDomain
		data.IamDomain = http_util.DomainContext(r.Context()).RequestedDomain()
	}

	if authRequest == nil {
		l.customTexts(r.Context(), translator, "")
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplRegisterOrg], data, nil)
}

func (d registerOrgFormData) toUserDomain() *domain.Human {
	if d.Username == "" {
		d.Username = string(d.Email)
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

func (d registerOrgFormData) toCommandOrg() *command.OrgSetup {
	if d.Username == "" {
		d.Username = string(d.Email)
	}
	return &command.OrgSetup{
		Name: d.RegisterOrgName,
		Admins: []*command.OrgSetupAdmin{
			{
				Human: &command.AddHuman{
					Username:  d.Username,
					FirstName: d.Firstname,
					LastName:  d.Lastname,
					Email: command.Email{
						Address: d.Email,
					},
					Password: d.Password,
					Register: true,
				},
			},
		},
	}
}
