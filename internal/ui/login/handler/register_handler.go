package handler

import (
	"github.com/caos/zitadel/internal/v2/domain"
	"net/http"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
)

const (
	tmplRegister          = "register"
	orgProjectCreatorRole = "ORG_PROJECT_CREATOR"
)

type registerFormData struct {
	Email        string `schema:"email"`
	Firstname    string `schema:"firstname"`
	Lastname     string `schema:"lastname"`
	Language     string `schema:"language"`
	Gender       int32  `schema:"gender"`
	Password     string `schema:"register-password"`
	Password2    string `schema:"register-password-confirmation"`
	TermsConfirm bool   `schema:"terms-confirm"`
}

type registerData struct {
	baseData
	registerFormData
	PasswordPolicyDescription string
	MinLength                 uint64
	HasUppercase              string
	HasLowercase              string
	HasNumber                 string
	HasSymbol                 string
}

func (l *Login) handleRegister(w http.ResponseWriter, r *http.Request) {
	data := new(registerFormData)
	authRequest, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authRequest, err)
		return
	}
	l.renderRegister(w, r, authRequest, data, nil)
}

func (l *Login) handleRegisterCheck(w http.ResponseWriter, r *http.Request) {
	data := new(registerFormData)
	authRequest, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authRequest, err)
		return
	}
	if data.Password != data.Password2 {
		err := caos_errs.ThrowInvalidArgument(nil, "VIEW-KaGue", "Errors.User.Password.ConfirmationWrong")
		l.renderRegister(w, r, authRequest, data, err)
		return
	}
	iam, err := l.authRepo.GetIAM(r.Context())
	if err != nil {
		l.renderRegister(w, r, authRequest, data, err)
		return
	}

	resourceOwner := iam.GlobalOrgID
	member := &domain.Member{
		ObjectRoot: models.ObjectRoot{AggregateID: iam.GlobalOrgID},
		Roles:      []string{orgProjectCreatorRole},
	}
	if authRequest.RequestedOrgID != "" && authRequest.RequestedOrgID != iam.GlobalOrgID {
		member = nil
		resourceOwner = authRequest.RequestedOrgID
	}
	user, err := l.command.RegisterHuman(setContext(r.Context(), resourceOwner), resourceOwner, data.toHumanDomain(), nil, member)
	if err != nil {
		l.renderRegister(w, r, authRequest, data, err)
		return
	}
	if authRequest == nil {
		http.Redirect(w, r, l.zitadelURL, http.StatusFound)
		return
	}
	authRequest.LoginName = user.PreferredLoginName
	l.renderNextStep(w, r, authRequest)
}

func (l *Login) renderRegister(w http.ResponseWriter, r *http.Request, authRequest *model.AuthRequest, formData *registerFormData, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	if formData == nil {
		formData = new(registerFormData)
	}
	if formData.Language == "" {
		formData.Language = l.renderer.Lang(r).String()
	}
	data := registerData{
		baseData:         l.getBaseData(r, authRequest, "Register", errType, errMessage),
		registerFormData: *formData,
	}
	iam, _ := l.authRepo.GetIAM(r.Context())
	if iam != nil {
		policy, description, _ := l.getPasswordComplexityPolicy(r, iam.GlobalOrgID)
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
	}

	funcs := map[string]interface{}{
		"selectedLanguage": func(l string) bool {
			if formData == nil {
				return false
			}
			return formData.Language == l
		},
		"selectedGender": func(g int32) bool {
			if formData == nil {
				return false
			}
			return formData.Gender == g
		},
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplRegister], data, funcs)
}

func (d registerFormData) toHumanDomain() *domain.Human {
	return &domain.Human{
		Profile: &domain.Profile{
			FirstName:         d.Firstname,
			LastName:          d.Lastname,
			PreferredLanguage: language.Make(d.Language),
			Gender:            domain.Gender(d.Gender),
		},
		Password: &domain.Password{
			SecretString: d.Password,
		},
		Email: &domain.Email{
			EmailAddress: d.Email,
		},
	}
}
