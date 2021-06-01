package handler

import (
	"net/http"

	"golang.org/x/text/language"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

const (
	tmplRegister = "register"
)

type registerFormData struct {
	Email        string `schema:"email"`
	Username     string `schema:"username"`
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
	ShowUsername              bool
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
	memberRoles := []string{domain.RoleOrgProjectCreator}

	if authRequest.RequestedOrgID != "" && authRequest.RequestedOrgID != iam.GlobalOrgID {
		memberRoles = nil
		resourceOwner = authRequest.RequestedOrgID
	}
	user, err := l.command.RegisterHuman(setContext(r.Context(), resourceOwner), resourceOwner, data.toHumanDomain(), nil, memberRoles)
	if err != nil {
		l.renderRegister(w, r, authRequest, data, err)
		return
	}
	if authRequest == nil {
		http.Redirect(w, r, l.zitadelURL, http.StatusFound)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.authRepo.SelectUser(r.Context(), authRequest.ID, user.AggregateID, userAgentID)
	if err != nil {
		l.renderRegister(w, r, authRequest, data, err)
		return
	}
	l.renderNextStep(w, r, authRequest)
}

func (l *Login) renderRegister(w http.ResponseWriter, r *http.Request, authRequest *domain.AuthRequest, formData *registerFormData, err error) {
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

	resourceOwner := authRequest.RequestedOrgID

	if resourceOwner == "" {
		iam, err := l.authRepo.GetIAM(r.Context())
		if err != nil {
			l.renderRegister(w, r, authRequest, formData, err)
			return
		}
		resourceOwner = iam.GlobalOrgID
	}

	pwPolicy, description, _ := l.getPasswordComplexityPolicy(r, resourceOwner)
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

	orgIAMPolicy, err := l.getOrgIamPolicy(r, resourceOwner)
	if err != nil {
		l.renderRegister(w, r, authRequest, formData, err)
		return
	}
	data.ShowUsername = orgIAMPolicy.UserLoginMustBeDomain

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
		Username: d.Username,
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
