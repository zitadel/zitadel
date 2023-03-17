package login

import (
	"net/http"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

const (
	tmplRegister = "register"
)

type registerFormData struct {
	Email        domain.EmailAddress `schema:"email"`
	Username     string              `schema:"username"`
	Firstname    string              `schema:"firstname"`
	Lastname     string              `schema:"lastname"`
	Language     string              `schema:"language"`
	Password     string              `schema:"register-password"`
	Password2    string              `schema:"register-password-confirmation"`
	TermsConfirm bool                `schema:"terms-confirm"`
}

type registerData struct {
	baseData
	registerFormData
	MinLength          uint64
	HasUppercase       string
	HasLowercase       string
	HasNumber          string
	HasSymbol          string
	ShowUsername       bool
	ShowUsernameSuffix bool
	OrgRegister        bool
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

	resourceOwner := authz.GetInstance(r.Context()).DefaultOrganisationID()

	if authRequest != nil && authRequest.RequestedOrgID != "" && authRequest.RequestedOrgID != resourceOwner {
		resourceOwner = authRequest.RequestedOrgID
	}
	initCodeGenerator, err := l.query.InitEncryptionGenerator(r.Context(), domain.SecretGeneratorTypeInitCode, l.userCodeAlg)
	if err != nil {
		l.renderRegister(w, r, authRequest, data, err)
		return
	}
	emailCodeGenerator, err := l.query.InitEncryptionGenerator(r.Context(), domain.SecretGeneratorTypeVerifyEmailCode, l.userCodeAlg)
	if err != nil {
		l.renderRegister(w, r, authRequest, data, err)
		return
	}
	phoneCodeGenerator, err := l.query.InitEncryptionGenerator(r.Context(), domain.SecretGeneratorTypeVerifyPhoneCode, l.userCodeAlg)
	if err != nil {
		l.renderRegister(w, r, authRequest, data, err)
		return
	}

	// For consistency with the external authentication flow,
	// the setMetadata() function is provided on the pre creation hook, for now,
	// like for the ExternalAuthentication flow.
	// If there is a need for additional context after registration,
	// we could provide that method in the PostCreation trigger too,
	// without breaking existing actions.
	// Also, if that field is needed, we probably also should provide it
	// for ExternalAuthentication.
	user, metadatas, err := l.runPreCreationActions(authRequest, r, data.toHumanDomain(), make([]*domain.Metadata, 0), resourceOwner, domain.FlowTypeInternalAuthentication)
	if err != nil {
		l.renderRegister(w, r, authRequest, data, err)
		return
	}

	user, err = l.command.RegisterHuman(setContext(r.Context(), resourceOwner), resourceOwner, user, nil, nil, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator)
	if err != nil {
		l.renderRegister(w, r, authRequest, data, err)
		return
	}

	if len(metadatas) > 0 {
		_, err = l.command.BulkSetUserMetadata(r.Context(), user.AggregateID, resourceOwner, metadatas...)
		if err != nil {
			// TODO: What if action is configured to be allowed to fail? Same question for external registration.
			l.renderRegister(w, r, authRequest, data, err)
			return
		}
	}

	userGrants, err := l.runPostCreationActions(user.AggregateID, authRequest, r, resourceOwner, domain.FlowTypeInternalAuthentication)
	if err != nil {
		l.renderError(w, r, authRequest, err)
		return
	}

	err = l.appendUserGrants(r.Context(), userGrants, resourceOwner)
	if err != nil {
		l.renderError(w, r, authRequest, err)
		return
	}

	if authRequest == nil {
		l.defaultRedirect(w, r)
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
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	translator := l.getTranslator(r.Context(), authRequest)
	if formData == nil {
		formData = new(registerFormData)
	}
	if formData.Language == "" {
		formData.Language = l.renderer.ReqLang(translator, r).String()
	}

	var resourceOwner string
	if authRequest != nil {
		resourceOwner = authRequest.RequestedOrgID
	}

	if resourceOwner == "" {
		resourceOwner = authz.GetInstance(r.Context()).DefaultOrganisationID()
	}

	data := registerData{
		baseData:         l.getBaseData(r, authRequest, "RegistrationUser.Title", "RegistrationUser.Description", errID, errMessage),
		registerFormData: *formData,
	}

	pwPolicy := l.getPasswordComplexityPolicy(r, resourceOwner)
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

	orgIAMPolicy, err := l.getOrgDomainPolicy(r, resourceOwner)
	if err != nil {
		l.renderRegister(w, r, authRequest, formData, err)
		return
	}
	data.ShowUsername = orgIAMPolicy.UserLoginMustBeDomain
	data.OrgRegister = orgIAMPolicy.UserLoginMustBeDomain

	labelPolicy, err := l.getLabelPolicy(r, resourceOwner)
	if err != nil {
		l.renderRegister(w, r, authRequest, formData, err)
		return
	}
	data.ShowUsernameSuffix = !labelPolicy.HideLoginNameSuffix

	funcs := map[string]interface{}{
		"selectedLanguage": func(l string) bool {
			if formData == nil {
				return false
			}
			return formData.Language == l
		},
	}
	if authRequest == nil {
		l.customTexts(r.Context(), translator, resourceOwner)
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplRegister], data, funcs)
}

func (d registerFormData) toHumanDomain() *domain.Human {
	return &domain.Human{
		Username: d.Username,
		Profile: &domain.Profile{
			FirstName:         d.Firstname,
			LastName:          d.Lastname,
			PreferredLanguage: language.Make(d.Language),
		},
		Password: &domain.Password{
			SecretString: d.Password,
		},
		Email: &domain.Email{
			EmailAddress: d.Email,
		},
	}
}
