package handler

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/caos/logging"
	"github.com/gorilla/csrf"
	"golang.org/x/text/language"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/renderer"
	"github.com/caos/zitadel/internal/static"
)

const (
	tmplError = "error"
)

type Renderer struct {
	*renderer.Renderer
	pathPrefix    string
	staticStorage static.Storage
}

func CreateRenderer(pathPrefix string, staticDir http.FileSystem, staticStorage static.Storage, cookieName string, defaultLanguage language.Tag) *Renderer {
	r := &Renderer{
		pathPrefix:    pathPrefix,
		staticStorage: staticStorage,
	}
	tmplMapping := map[string]string{
		tmplError:                        "error.html",
		tmplLogin:                        "login.html",
		tmplUserSelection:                "select_user.html",
		tmplPassword:                     "password.html",
		tmplPasswordlessVerification:     "passwordless.html",
		tmplPasswordlessRegistration:     "passwordless_registration.html",
		tmplPasswordlessRegistrationDone: "passwordless_registration_done.html",
		tmplPasswordlessPrompt:           "passwordless_prompt.html",
		tmplMFAVerify:                    "mfa_verify_otp.html",
		tmplMFAPrompt:                    "mfa_prompt.html",
		tmplMFAInitVerify:                "mfa_init_otp.html",
		tmplMFAU2FInit:                   "mfa_init_u2f.html",
		tmplU2FVerification:              "mfa_verification_u2f.html",
		tmplMFAInitDone:                  "mfa_init_done.html",
		tmplMailVerification:             "mail_verification.html",
		tmplMailVerified:                 "mail_verified.html",
		tmplInitPassword:                 "init_password.html",
		tmplInitPasswordDone:             "init_password_done.html",
		tmplInitUser:                     "init_user.html",
		tmplInitUserDone:                 "init_user_done.html",
		tmplPasswordResetDone:            "password_reset_done.html",
		tmplChangePassword:               "change_password.html",
		tmplChangePasswordDone:           "change_password_done.html",
		tmplRegisterOption:               "register_option.html",
		tmplRegister:                     "register.html",
		tmplExternalRegisterOverview:     "external_register_overview.html",
		tmplLogoutDone:                   "logout_done.html",
		tmplRegisterOrg:                  "register_org.html",
		tmplChangeUsername:               "change_username.html",
		tmplChangeUsernameDone:           "change_username_done.html",
		tmplLinkUsersDone:                "link_users_done.html",
		tmplExternalNotFoundOption:       "external_not_found_option.html",
		tmplLoginSuccess:                 "login_success.html",
	}
	funcs := map[string]interface{}{
		"resourceUrl": func(file string) string {
			return path.Join(r.pathPrefix, EndpointResources, file)
		},
		"resourceThemeUrl": func(file, theme string) string {
			return path.Join(r.pathPrefix, EndpointResources, "themes", theme, file)
		},
		"hasCustomPolicy": func(policy *domain.LabelPolicy) bool {
			if policy != nil {
				return true
			}
			return false
		},
		"hasWatermark": func(policy *domain.LabelPolicy) bool {
			if policy != nil && policy.DisableWatermark {
				return false
			}
			return true
		},
		"variablesCssFileUrl": func(orgID string, policy *domain.LabelPolicy) string {
			cssFile := domain.CssPath + "/" + domain.CssVariablesFileName
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s&%s=%v&%s=%s", EndpointDynamicResources, "orgId", orgID, "default-policy", policy.Default, "filename", cssFile))
		},
		"customLogoResource": func(orgID string, policy *domain.LabelPolicy, darkMode bool) string {
			fileName := policy.LogoURL
			if darkMode && policy.LogoDarkURL != "" {
				fileName = policy.LogoDarkURL
			}
			if fileName == "" {
				return ""
			}
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s&%s=%v&%s=%s", EndpointDynamicResources, "orgId", orgID, "default-policy", policy.Default, "filename", fileName))
		},
		"avatarResource": func(orgID, avatar string) string {
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s&%s=%v&%s=%s", EndpointDynamicResources, "orgId", orgID, "default-policy", false, "filename", avatar))
		},
		"loginUrl": func() string {
			return path.Join(r.pathPrefix, EndpointLogin)
		},
		"externalIDPAuthURL": func(authReqID, idpConfigID string) string {
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s&%s=%s", EndpointExternalLogin, queryAuthRequestID, authReqID, queryIDPConfigID, idpConfigID))
		},
		"externalIDPRegisterURL": func(authReqID, idpConfigID string) string {
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s&%s=%s", EndpointExternalRegister, queryAuthRequestID, authReqID, queryIDPConfigID, idpConfigID))
		},
		"registerUrl": func(id string) string {
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s", EndpointRegister, queryAuthRequestID, id))
		},
		"loginNameUrl": func() string {
			return path.Join(r.pathPrefix, EndpointLoginName)
		},
		"loginNameChangeUrl": func(id string) string {
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s", EndpointLoginName, queryAuthRequestID, id))
		},
		"userSelectionUrl": func() string {
			return path.Join(r.pathPrefix, EndpointUserSelection)
		},
		"passwordLessVerificationUrl": func() string {
			return path.Join(r.pathPrefix, EndpointPasswordlessLogin)
		},
		"passwordLessRegistrationUrl": func() string {
			return path.Join(r.pathPrefix, EndpointPasswordlessRegistration)
		},
		"passwordlessPromptUrl": func() string {
			return path.Join(r.pathPrefix, EndpointPasswordlessPrompt)
		},
		"passwordResetUrl": func(id string) string {
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s", EndpointPasswordReset, queryAuthRequestID, id))
		},
		"passwordUrl": func() string {
			return path.Join(r.pathPrefix, EndpointPassword)
		},
		"mfaVerifyUrl": func() string {
			return path.Join(r.pathPrefix, EndpointMFAVerify)
		},
		"mfaPromptUrl": func() string {
			return path.Join(r.pathPrefix, EndpointMFAPrompt)
		},
		"mfaPromptChangeUrl": func(id string, provider model.MFAType) string {
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s;%s=%v", EndpointMFAPrompt, queryAuthRequestID, id, "provider", provider))
		},
		"mfaInitVerifyUrl": func() string {
			return path.Join(r.pathPrefix, EndpointMFAInitVerify)
		},
		"mfaInitU2FVerifyUrl": func() string {
			return path.Join(r.pathPrefix, EndpointMFAInitU2FVerify)
		},
		"mfaInitU2FLoginUrl": func() string {
			return path.Join(r.pathPrefix, EndpointU2FVerification)
		},
		"mailVerificationUrl": func() string {
			return path.Join(r.pathPrefix, EndpointMailVerification)
		},
		"initPasswordUrl": func() string {
			return path.Join(r.pathPrefix, EndpointInitPassword)
		},
		"initUserUrl": func() string {
			return path.Join(r.pathPrefix, EndpointInitUser)
		},
		"changePasswordUrl": func() string {
			return path.Join(r.pathPrefix, EndpointChangePassword)
		},
		"registerOptionUrl": func() string {
			return path.Join(r.pathPrefix, EndpointRegisterOption)
		},
		"registrationUrl": func() string {
			return path.Join(r.pathPrefix, EndpointRegister)
		},
		"orgRegistrationUrl": func() string {
			return path.Join(r.pathPrefix, EndpointRegisterOrg)
		},
		"externalRegistrationUrl": func() string {
			return path.Join(r.pathPrefix, EndpointExternalRegister)
		},
		"changeUsernameUrl": func() string {
			return path.Join(r.pathPrefix, EndpointChangeUsername)
		},
		"externalNotFoundOptionUrl": func() string {
			return path.Join(r.pathPrefix, EndpointExternalNotFoundOption)
		},
		"selectedLanguage": func(l string) bool {
			return false
		},
		"selectedGender": func(g int32) bool {
			return false
		},
		"hasUsernamePasswordLogin": func() bool {
			return false
		},
		"showPasswordReset": func() bool {
			return true
		},
		"hasExternalLogin": func() bool {
			return false
		},
		"idpProviderClass": func(stylingType domain.IDPConfigStylingType) string {
			return stylingType.GetCSSClass()
		},
	}
	var err error
	r.Renderer, err = renderer.NewRenderer(
		staticDir,
		tmplMapping, funcs,
		i18n.TranslatorConfig{DefaultLanguage: defaultLanguage, CookieName: cookieName},
	)
	logging.Log("APP-40tSoJ").OnError(err).WithError(err).Panic("error creating renderer")
	return r
}

func (l *Login) renderNextStep(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), authReq.ID, userAgentID)
	if err != nil {
		l.renderInternalError(w, r, authReq, caos_errs.ThrowInternal(err, "APP-sio0W", "could not get authreq"))
		return
	}
	if len(authReq.PossibleSteps) == 0 {
		l.renderInternalError(w, r, authReq, caos_errs.ThrowInternal(nil, "APP-9sdp4", "no possible steps"))
		return
	}
	l.chooseNextStep(w, r, authReq, 0, nil)
}

func (l *Login) renderError(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	if authReq == nil || len(authReq.PossibleSteps) == 0 {
		l.renderInternalError(w, r, authReq, caos_errs.ThrowInternal(err, "APP-OVOiT", "no possible steps"))
		return
	}
	l.chooseNextStep(w, r, authReq, 0, err)
}

func (l *Login) chooseNextStep(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, stepNumber int, err error) {
	switch step := authReq.PossibleSteps[stepNumber].(type) {
	case *domain.LoginStep:
		if len(authReq.PossibleSteps) > 1 {
			l.chooseNextStep(w, r, authReq, 1, err)
			return
		}
		l.renderLogin(w, r, authReq, err)
	case *domain.RegistrationStep:
		l.renderRegisterOption(w, r, authReq, nil)
	case *domain.SelectUserStep:
		l.renderUserSelection(w, r, authReq, step)
	case *domain.InitPasswordStep:
		l.renderInitPassword(w, r, authReq, authReq.UserID, "", err)
	case *domain.PasswordStep:
		l.renderPassword(w, r, authReq, nil)
	case *domain.PasswordlessStep:
		l.renderPasswordlessVerification(w, r, authReq, step.PasswordSet, nil)
	case *domain.PasswordlessRegistrationPromptStep:
		l.renderPasswordlessPrompt(w, r, authReq, nil)
	case *domain.MFAVerificationStep:
		l.renderMFAVerify(w, r, authReq, step, err)
	case *domain.RedirectToCallbackStep:
		if len(authReq.PossibleSteps) > 1 {
			l.chooseNextStep(w, r, authReq, 1, err)
			return
		}
		l.redirectToLoginSuccess(w, r, authReq.ID)
	case *domain.ChangePasswordStep:
		l.renderChangePassword(w, r, authReq, err)
	case *domain.VerifyEMailStep:
		l.renderMailVerification(w, r, authReq, "", err)
	case *domain.MFAPromptStep:
		l.renderMFAPrompt(w, r, authReq, step, err)
	case *domain.InitUserStep:
		l.renderInitUser(w, r, authReq, "", "", step.PasswordSet, nil)
	case *domain.ChangeUsernameStep:
		l.renderChangeUsername(w, r, authReq, nil)
	case *domain.LinkUsersStep:
		l.linkUsers(w, r, authReq, err)
	case *domain.ExternalNotFoundOptionStep:
		l.renderExternalNotFoundOption(w, r, authReq, err)
	case *domain.ExternalLoginStep:
		l.handleExternalLoginStep(w, r, authReq, step.SelectedIDPConfigID)
	case *domain.GrantRequiredStep:
		l.renderInternalError(w, r, authReq, caos_errs.ThrowPreconditionFailed(nil, "APP-asb43", "Errors.User.GrantRequired"))
	default:
		l.renderInternalError(w, r, authReq, caos_errs.ThrowInternal(nil, "APP-ds3QF", "step no possible"))
	}
}

func (l *Login) renderInternalError(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var msg string
	if err != nil {
		msg = err.Error()
	}
	data := l.getBaseData(r, authReq, "Error", "Internal", msg)
	l.renderer.RenderTemplate(w, r, l.getTranslator(authReq), l.renderer.Templates[tmplError], data, nil)
}

func (l *Login) getUserData(r *http.Request, authReq *domain.AuthRequest, title string, errType, errMessage string) userData {
	userData := userData{
		baseData:    l.getBaseData(r, authReq, title, errType, errMessage),
		profileData: l.getProfileData(authReq),
	}
	if authReq != nil && authReq.LinkingUsers != nil {
		userData.Linking = len(authReq.LinkingUsers) > 0
	}
	return userData
}

func (l *Login) getBaseData(r *http.Request, authReq *domain.AuthRequest, title string, errType, errMessage string) baseData {
	baseData := baseData{
		errorData: errorData{
			ErrID:      errType,
			ErrMessage: errMessage,
		},
		Lang:                   l.renderer.ReqLang(l.getTranslator(authReq), r).String(),
		Title:                  title,
		Theme:                  l.getTheme(r),
		ThemeMode:              l.getThemeMode(r),
		DarkMode:               l.isDarkMode(r),
		OrgID:                  l.getOrgID(authReq),
		OrgName:                l.getOrgName(authReq),
		PrimaryDomain:          l.getOrgPrimaryDomain(authReq),
		DisplayLoginNameSuffix: l.isDisplayLoginNameSuffix(authReq),
		AuthReqID:              getRequestID(authReq, r),
		CSRF:                   csrf.TemplateField(r),
		Nonce:                  http_mw.GetNonce(r),
	}
	if authReq != nil {
		baseData.LoginPolicy = authReq.LoginPolicy
		baseData.LabelPolicy = authReq.LabelPolicy
		baseData.IDPProviders = authReq.AllowedExternalIDPs
		if authReq.PrivacyPolicy != nil {
			baseData.TOSLink = authReq.PrivacyPolicy.TOSLink
			baseData.PrivacyLink = authReq.PrivacyPolicy.PrivacyLink
		}
	} else {
		privacyPolicy, err := l.getDefaultPrivacyPolicy(r)
		if err != nil {
			return baseData
		}
		if privacyPolicy != nil {
			baseData.TOSLink = privacyPolicy.TOSLink
			baseData.PrivacyLink = privacyPolicy.PrivacyLink
		}
	}
	return baseData
}

func (l *Login) getTranslator(authReq *domain.AuthRequest) *i18n.Translator {
	translator, _ := l.renderer.NewTranslator()
	if authReq != nil {
		l.addLoginTranslations(translator, authReq.DefaultTranslations)
		l.addLoginTranslations(translator, authReq.OrgTranslations)
		translator.SetPreferredLanguages(authReq.UiLocales...)
	}
	return translator
}

func (l *Login) getProfileData(authReq *domain.AuthRequest) profileData {
	var userName, loginName, displayName, avatar string
	if authReq != nil {
		userName = authReq.UserName
		loginName = authReq.LoginName
		displayName = authReq.DisplayName
		avatar = authReq.AvatarKey
	}
	return profileData{
		UserName:    userName,
		LoginName:   loginName,
		DisplayName: displayName,
		AvatarKey:   avatar,
	}
}

func (l *Login) getErrorMessage(r *http.Request, err error) (errID, errMsg string) {
	caosErr := new(caos_errs.CaosError)
	if errors.As(err, &caosErr) {
		localized := l.renderer.LocalizeFromRequest(l.getTranslator(nil), r, caosErr.Message, nil)
		return caosErr.ID, localized

	}
	return "", err.Error()
}

func (l *Login) getTheme(r *http.Request) string {
	return "zitadel" //TODO: impl
}

func (l *Login) getThemeMode(r *http.Request) string {
	if l.isDarkMode(r) {
		return "lgn-dark-theme"
	}
	return "lgn-light-theme"
}

func (l *Login) isDarkMode(r *http.Request) bool {
	cookie, err := r.Cookie("mode")
	if err != nil {
		return false
	}
	return strings.HasSuffix(cookie.Value, "dark")
}

func (l *Login) getOrgID(authReq *domain.AuthRequest) string {
	if authReq == nil {
		return ""
	}
	if authReq.RequestedOrgID != "" {
		return authReq.RequestedOrgID
	}
	return authReq.UserOrgID
}

func (l *Login) getOrgName(authReq *domain.AuthRequest) string {
	if authReq == nil {
		return ""
	}
	return authReq.RequestedOrgName
}

func (l *Login) getOrgPrimaryDomain(authReq *domain.AuthRequest) string {
	if authReq == nil {
		return ""
	}
	return authReq.RequestedPrimaryDomain
}

func (l *Login) isDisplayLoginNameSuffix(authReq *domain.AuthRequest) bool {
	if authReq == nil {
		return false
	}
	if authReq.RequestedOrgID == "" {
		return false
	}
	return authReq.LabelPolicy != nil && !authReq.LabelPolicy.HideLoginNameSuffix
}

func (l *Login) addLoginTranslations(translator *i18n.Translator, customTexts []*domain.CustomText) {
	for _, text := range customTexts {
		msg := i18n.Message{
			ID:   text.Key,
			Text: text.Text,
		}
		err := l.renderer.AddMessages(translator, text.Language, msg)
		logging.Log("HANDLE-GD3g2").OnError(err).Warn("could no add message to translator")
	}
}

func getRequestID(authReq *domain.AuthRequest, r *http.Request) string {
	if authReq != nil {
		return authReq.ID
	}
	return r.FormValue(queryAuthRequestID)
}

func (l *Login) csrfErrorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := csrf.FailureReason(r)
		l.renderInternalError(w, r, nil, err)
	})
}

func (l *Login) cspErrorHandler(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.renderInternalError(w, r, nil, err)
	})
}

type baseData struct {
	errorData
	Lang                   string
	Title                  string
	Theme                  string
	ThemeMode              string
	DarkMode               bool
	OrgID                  string
	OrgName                string
	PrimaryDomain          string
	DisplayLoginNameSuffix bool
	TOSLink                string
	PrivacyLink            string
	AuthReqID              string
	CSRF                   template.HTML
	Nonce                  string
	LoginPolicy            *domain.LoginPolicy
	IDPProviders           []*domain.IDPProvider
	LabelPolicy            *domain.LabelPolicy
	LoginTexts             []*domain.CustomLoginText
}

type errorData struct {
	ErrID      string
	ErrMessage string
}

type userData struct {
	baseData
	profileData
	PasswordChecked     string
	MFAProviders        []domain.MFAType
	SelectedMFAProvider domain.MFAType
	Linking             bool
}

type profileData struct {
	LoginName   string
	UserName    string
	DisplayName string
	AvatarKey   string
}

type passwordData struct {
	baseData
	profileData
	PasswordPolicyDescription string
	MinLength                 uint64
	HasUppercase              string
	HasLowercase              string
	HasNumber                 string
	HasSymbol                 string
}

type userSelectionData struct {
	baseData
	Users   []domain.UserSelection
	Linking bool
}

type mfaData struct {
	baseData
	profileData
	MFAProviders []domain.MFAType
	MFARequired  bool
}

type mfaVerifyData struct {
	baseData
	profileData
	MFAType domain.MFAType
	otpData
}

type mfaDoneData struct {
	baseData
	profileData
	MFAType domain.MFAType
}

type otpData struct {
	Url    string
	Secret string
	QrCode string
}
