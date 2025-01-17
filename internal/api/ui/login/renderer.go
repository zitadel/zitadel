package login

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/renderer"
	"github.com/zitadel/zitadel/internal/static"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	tmplError   = "error"
	tmplSuccess = "success"
)

type Renderer struct {
	*renderer.Renderer
	pathPrefix    string
	staticStorage static.Storage
}

type LanguageData struct {
	Lang string
}

func CreateRenderer(pathPrefix string, staticStorage static.Storage, cookieName string) *Renderer {
	r := &Renderer{
		pathPrefix:    pathPrefix,
		staticStorage: staticStorage,
	}
	tmplMapping := map[string]string{
		tmplError:                        "error.html",
		tmplSuccess:                      "success.html",
		tmplLogin:                        "login.html",
		tmplUserSelection:                "select_user.html",
		tmplPassword:                     "password.html",
		tmplPasswordlessVerification:     "passwordless.html",
		tmplPasswordlessRegistration:     "passwordless_registration.html",
		tmplPasswordlessRegistrationDone: "passwordless_registration_done.html",
		tmplPasswordlessPrompt:           "passwordless_prompt.html",
		tmplMFAVerify:                    "mfa_verify_totp.html",
		tmplMFAPrompt:                    "mfa_prompt.html",
		tmplMFAInitVerify:                "mfa_init_otp.html",
		tmplMFASMSInit:                   "mfa_init_otp_sms.html",
		tmplOTPVerification:              "mfa_verify_otp.html",
		tmplMFAU2FInit:                   "mfa_init_u2f.html",
		tmplU2FVerification:              "mfa_verification_u2f.html",
		tmplMFAInitDone:                  "mfa_init_done.html",
		tmplMailVerification:             "mail_verification.html",
		tmplMailVerified:                 "mail_verified.html",
		tmplInitPassword:                 "init_password.html",
		tmplInitPasswordDone:             "init_password_done.html",
		tmplInitUser:                     "init_user.html",
		tmplInitUserDone:                 "init_user_done.html",
		tmplInviteUser:                   "invite_user.html",
		tmplPasswordResetDone:            "password_reset_done.html",
		tmplChangePassword:               "change_password.html",
		tmplChangePasswordDone:           "change_password_done.html",
		tmplRegisterOption:               "register_option.html",
		tmplRegister:                     "register.html",
		tmplLogoutDone:                   "logout_done.html",
		tmplRegisterOrg:                  "register_org.html",
		tmplChangeUsername:               "change_username.html",
		tmplChangeUsernameDone:           "change_username_done.html",
		tmplLinkUsersDone:                "link_users_done.html",
		tmplExternalNotFoundOption:       "external_not_found_option.html",
		tmplLoginSuccess:                 "login_success.html",
		tmplLDAPLogin:                    "ldap_login.html",
		tmplDeviceAuthUserCode:           "device_usercode.html",
		tmplDeviceAuthAction:             "device_action.html",
	}
	funcs := map[string]interface{}{
		"resourceUrl": func(file string) string {
			return path.Join(r.pathPrefix, EndpointResources, file)
		},
		"resourceThemeUrl": func(file, theme string) string {
			return path.Join(r.pathPrefix, EndpointResources, "themes", theme, file)
		},
		"hasCustomPolicy": func(policy *domain.LabelPolicy) bool {
			return policy != nil
		},
		"hasWatermark": func(policy *domain.LabelPolicy) bool {
			return policy == nil || !policy.DisableWatermark
		},
		"variablesCssFileUrl": func(orgID string, policy *domain.LabelPolicy) string {
			cssFile := domain.CssPath + "/" + domain.CssVariablesFileName + "?v=" + policy.ChangeDate.Format(time.RFC3339)
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
		"customIconResource": func(orgID string, policy *domain.LabelPolicy, darkMode bool) string {
			fileName := policy.IconURL
			if darkMode && policy.IconDarkURL != "" {
				fileName = policy.IconDarkURL
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
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s&%s=%s", EndpointExternalLogin, QueryAuthRequestID, authReqID, queryIDPConfigID, idpConfigID))
		},
		"externalIDPRegisterURL": func(authReqID, idpConfigID string) string {
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s&%s=%s", EndpointExternalRegister, QueryAuthRequestID, authReqID, queryIDPConfigID, idpConfigID))
		},
		"registerUrl": func(id string) string {
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s", EndpointRegister, QueryAuthRequestID, id))
		},
		"loginNameUrl": func() string {
			return path.Join(r.pathPrefix, EndpointLoginName)
		},
		"loginNameChangeUrl": func(id string) string {
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s", EndpointLoginName, QueryAuthRequestID, id))
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
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s", EndpointPasswordReset, QueryAuthRequestID, id))
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
		"mfaPromptChangeUrl": func(id string, provider domain.MFAType) string {
			return path.Join(r.pathPrefix, fmt.Sprintf("%s?%s=%s&%s=%v", EndpointMFAPrompt, QueryAuthRequestID, id, "provider", provider))
		},
		"mfaInitVerifyUrl": func() string {
			return path.Join(r.pathPrefix, EndpointMFAInitVerify)
		},
		"mfaInitSMSVerifyUrl": func() string {
			return path.Join(r.pathPrefix, EndpointMFASMSInitVerify)
		},
		"mfaOTPVerifyUrl": func() string {
			return path.Join(r.pathPrefix, EndpointMFAOTPVerify)
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
		"inviteUserUrl": func() string {
			return path.Join(r.pathPrefix, EndpointInviteUser)
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
		"changeUsernameUrl": func() string {
			return path.Join(r.pathPrefix, EndpointChangeUsername)
		},
		"externalNotFoundOptionUrl": func(action string) string {
			return path.Join(r.pathPrefix, EndpointExternalNotFoundOption+"?"+action+"=true")
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
		"hasRegistration": func() bool {
			return true
		},
		"idpProviderClass": func(idpType domain.IDPType) string {
			return idpType.GetCSSClass()
		},
		"ldapUrl": func() string {
			return path.Join(r.pathPrefix, EndpointLDAPCallback)
		},
		"linkingUserPromptUrl": func() string {
			return path.Join(r.pathPrefix, EndpointLinkingUserPrompt)
		},
	}
	var err error
	r.Renderer, err = renderer.NewRenderer(
		tmplMapping, funcs,
		cookieName,
	)
	logging.New().OnError(err).WithError(err).Panic("error creating renderer")
	return r
}

func (l *Login) renderNextStep(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	if authReq == nil {
		l.renderInternalError(w, r, nil, zerrors.ThrowInvalidArgument(nil, "LOGIN-Df3f2", "Errors.AuthRequest.NotFound"))
		return
	}
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
	if err != nil {
		l.renderInternalError(w, r, authReq, err)
		return
	}
	if len(authReq.PossibleSteps) == 0 {
		l.renderInternalError(w, r, authReq, zerrors.ThrowInternal(nil, "APP-9sdp4", "no possible steps"))
		return
	}
	l.chooseNextStep(w, r, authReq, 0, nil)
}

func (l *Login) renderError(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	if err != nil {
		l.renderInternalError(w, r, authReq, err)
		return
	}
	if authReq == nil || len(authReq.PossibleSteps) == 0 {
		l.renderInternalError(w, r, authReq, zerrors.ThrowInternal(err, "APP-OVOiT", "no possible steps"))
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
	case *domain.RedirectToExternalIDPStep:
		l.handleIDP(w, r, authReq, authReq.SelectedIDPConfigID)
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
		l.redirectToCallback(w, r, authReq)
	case *domain.LoginSucceededStep:
		l.redirectToLoginSuccess(w, r, authReq.ID)
	case *domain.ChangePasswordStep:
		l.renderChangePassword(w, r, authReq, err)
	case *domain.VerifyEMailStep:
		l.renderMailVerification(w, r, authReq, authReq.UserID, "", step.InitPassword, err)
	case *domain.MFAPromptStep:
		l.renderMFAPrompt(w, r, authReq, step, err)
	case *domain.InitUserStep:
		l.renderInitUser(w, r, authReq, "", "", "", step.PasswordSet, nil)
	case *domain.ChangeUsernameStep:
		l.renderChangeUsername(w, r, authReq, nil)
	case *domain.LinkUsersStep:
		l.linkUsers(w, r, authReq, err)
	case *domain.ExternalNotFoundOptionStep:
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
	case *domain.ExternalLoginStep:
		l.handleExternalLoginStep(w, r, authReq, step.SelectedIDPConfigID)
	case *domain.GrantRequiredStep:
		l.renderInternalError(w, r, authReq, zerrors.ThrowPreconditionFailed(nil, "APP-asb43", "Errors.User.GrantRequired"))
	case *domain.ProjectRequiredStep:
		l.renderInternalError(w, r, authReq, zerrors.ThrowPreconditionFailed(nil, "APP-m92d", "Errors.User.ProjectRequired"))
	case *domain.VerifyInviteStep:
		l.renderInviteUser(w, r, authReq, "", "", "", "", nil)
	default:
		l.renderInternalError(w, r, authReq, zerrors.ThrowInternal(nil, "APP-ds3QF", "step no possible"))
	}
}

func (l *Login) renderInternalError(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	if err != nil {
		log := logging.WithError(err)
		if authReq != nil {
			log = log.WithField("auth_req_id", authReq.ID)
		}
		if zerrors.IsInternal(err) {
			log.Error()
		} else {
			log.Info()
		}
	}
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getBaseData(r, authReq, translator, "Errors.Internal", "", err)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplError], data, nil)
}

func (l *Login) getUserData(r *http.Request, authReq *domain.AuthRequest, translator *i18n.Translator, titleI18nKey string, descriptionI18nKey string, err error) userData {
	userData := userData{
		baseData:    l.getBaseData(r, authReq, translator, titleI18nKey, descriptionI18nKey, err),
		profileData: l.getProfileData(authReq),
	}
	if authReq != nil && authReq.LinkingUsers != nil {
		userData.Linking = len(authReq.LinkingUsers) > 0
	}
	return userData
}

func (l *Login) getBaseData(r *http.Request, authReq *domain.AuthRequest, translator *i18n.Translator, titleI18nKey string, descriptionI18nKey string, err error) baseData {
	title := ""
	if titleI18nKey != "" {
		title = translator.LocalizeWithoutArgs(titleI18nKey)
	}

	description := ""
	if descriptionI18nKey != "" {
		description = translator.LocalizeWithoutArgs(descriptionI18nKey)
	}

	lang, _ := l.renderer.ReqLang(translator, r).Base()
	var errID, errMessage string
	var errPopup bool
	if err != nil {
		errID, errMessage, errPopup = l.getErrorMessage(r, err)
	}
	baseData := baseData{
		errorData: errorData{
			ErrID:      errID,
			ErrMessage: errMessage,
			ErrPopup:   errPopup,
		},
		Lang:                   lang.String(),
		Title:                  title,
		Description:            description,
		Theme:                  l.getTheme(r),
		PrivateLabelingOrgID:   l.getPrivateLabelingID(r, authReq),
		OrgID:                  l.getOrgID(r, authReq),
		OrgName:                l.getOrgName(authReq),
		PrimaryDomain:          l.getOrgPrimaryDomain(r, authReq),
		DisplayLoginNameSuffix: l.isDisplayLoginNameSuffix(authReq),
		AuthReqID:              getRequestID(authReq, r),
		CSRF:                   csrf.TemplateField(r),
		Nonce:                  http_mw.GetNonce(r),
	}
	var privacyPolicy *domain.PrivacyPolicy
	if authReq != nil {
		baseData.LoginPolicy = authReq.LoginPolicy
		baseData.LabelPolicy = authReq.LabelPolicy
		baseData.IDPProviders = authReq.AllowedExternalIDPs
		if authReq.PrivacyPolicy == nil {
			return baseData
		}
		privacyPolicy = authReq.PrivacyPolicy
	} else {
		labelPolicy, _ := l.query.ActiveLabelPolicyByOrg(r.Context(), baseData.PrivateLabelingOrgID, false)
		if labelPolicy != nil {
			baseData.LabelPolicy = labelPolicy.ToDomain()
		}
		policy, err := l.query.DefaultPrivacyPolicy(r.Context(), false)
		if err != nil {
			return baseData
		}
		privacyPolicy = policy.ToDomain()
	}
	baseData.ThemeMode = l.getThemeMode(baseData.LabelPolicy)
	baseData.ThemeClass = l.getThemeClass(r, baseData.LabelPolicy)
	baseData.DarkMode = l.isDarkMode(r, baseData.LabelPolicy)
	baseData = l.setLinksOnBaseData(baseData, privacyPolicy)
	return baseData
}

func (l *Login) getTranslator(ctx context.Context, authReq *domain.AuthRequest) *i18n.Translator {
	restrictions, err := l.query.GetInstanceRestrictions(ctx)
	if err != nil {
		logging.OnError(err).Warn("cannot load instance restrictions to retrieve allowed languages for creating the translator")
	}
	translator, err := l.renderer.NewTranslator(ctx, restrictions.AllowedLanguages)
	logging.OnError(err).Warn("cannot load translator")
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

func (l *Login) setLinksOnBaseData(baseData baseData, privacyPolicy *domain.PrivacyPolicy) baseData {
	lang := LanguageData{
		Lang: baseData.Lang,
	}
	baseData.TOSLink = privacyPolicy.TOSLink
	baseData.PrivacyLink = privacyPolicy.PrivacyLink
	baseData.HelpLink = privacyPolicy.HelpLink

	if link, err := templates.ParseTemplateText(privacyPolicy.TOSLink, lang); err == nil {
		baseData.TOSLink = link
	}
	if link, err := templates.ParseTemplateText(privacyPolicy.PrivacyLink, lang); err == nil {
		baseData.PrivacyLink = link
	}
	if link, err := templates.ParseTemplateText(privacyPolicy.HelpLink, lang); err == nil {
		baseData.HelpLink = link
	}
	if link, err := templates.ParseTemplateText(string(privacyPolicy.SupportEmail), lang); err == nil {
		baseData.SupportEmail = link
	}
	return baseData
}

func (l *Login) getErrorMessage(r *http.Request, err error) (errID, errMsg string, popup bool) {
	idpErr := new(IdPError)
	if errors.Is(err, idpErr) {
		popup = true
	}
	caosErr := new(zerrors.ZitadelError)
	if errors.As(err, &caosErr) {
		localized := l.renderer.LocalizeFromRequest(l.getTranslator(r.Context(), nil), r, caosErr.Message, map[string]interface{}{"Details": caosErr.Parent})
		return caosErr.ID, localized, popup
	}
	return "", err.Error(), popup
}

func (l *Login) getTheme(r *http.Request) string {
	return "zitadel"
}

// getThemeClass returns the css class for the login html.
// Possible values are `lgn-light-theme` and `lgn-dark-theme` and are based on the policy first
// and if it's set to auto the cookie is checked.
func (l *Login) getThemeClass(r *http.Request, policy *domain.LabelPolicy) string {
	if l.isDarkMode(r, policy) {
		return "lgn-dark-theme"
	}
	return "lgn-light-theme"
}

// isDarkMode checks policy first and if not set to specifically use dark or light only,
// it will also check the cookie.
func (l *Login) isDarkMode(r *http.Request, policy *domain.LabelPolicy) bool {
	if mode := l.getThemeMode(policy); mode != domain.LabelPolicyThemeAuto {
		return mode == domain.LabelPolicyThemeDark
	}
	cookie, err := r.Cookie("mode")
	if err != nil {
		return false
	}
	return strings.HasSuffix(cookie.Value, "dark")
}

func (l *Login) getThemeMode(policy *domain.LabelPolicy) domain.LabelPolicyThemeMode {
	if policy != nil {
		return policy.ThemeMode
	}
	return domain.LabelPolicyThemeAuto
}

func (l *Login) getOrgID(r *http.Request, authReq *domain.AuthRequest) string {
	if authReq == nil {
		return r.FormValue(queryOrgID)
	}
	if authReq.RequestedOrgID != "" {
		return authReq.RequestedOrgID
	}
	return authReq.UserOrgID
}

func (l *Login) getPrivateLabelingID(r *http.Request, authReq *domain.AuthRequest) string {
	instance := authz.GetInstance(r.Context())
	defaultID := instance.DefaultOrganisationID()
	if !instance.Features().LoginDefaultOrg {
		defaultID = instance.InstanceID()
	}

	if authReq != nil {
		return authReq.PrivateLabelingOrgID(defaultID)
	}
	if id := r.FormValue(queryOrgID); id != "" {
		return id
	}
	return defaultID
}

func (l *Login) getOrgName(authReq *domain.AuthRequest) string {
	if authReq == nil {
		return ""
	}
	return authReq.RequestedOrgName
}

func (l *Login) getOrgPrimaryDomain(r *http.Request, authReq *domain.AuthRequest) string {
	orgID := authz.GetInstance(r.Context()).DefaultOrganisationID()
	if authReq != nil && authReq.RequestedPrimaryDomain != "" {
		return authReq.RequestedPrimaryDomain
	}
	org, err := l.query.OrgByID(r.Context(), false, orgID)
	if err != nil {
		logging.New().WithError(err).Error("cannot get default org")
		return ""
	}
	return org.Domain
}

func (l *Login) isDisplayLoginNameSuffix(authReq *domain.AuthRequest) bool {
	if authReq == nil {
		return false
	}
	if authReq.RequestedOrgID == "" || !authReq.RequestedOrgDomain {
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
		logging.OnError(err).Warn("could no add message to translator")
	}
}

func (l *Login) customTexts(ctx context.Context, translator *i18n.Translator, orgID string) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceTexts, err := l.query.CustomTextListByTemplate(ctx, instanceID, domain.LoginCustomText, false)
	if err != nil {
		logging.WithFields("instanceID", instanceID).Warn("unable to load custom texts for instance")
		return
	}
	l.addLoginTranslations(translator, query.CustomTextsToDomain(instanceTexts))
	if orgID == "" {
		return
	}
	orgTexts, err := l.query.CustomTextListByTemplate(ctx, orgID, domain.LoginCustomText, false)
	if err != nil {
		logging.WithFields("instanceID", instanceID, "org", orgID).Warn("unable to load custom texts for org")
		return
	}
	l.addLoginTranslations(translator, query.CustomTextsToDomain(orgTexts))
}

func getRequestID(authReq *domain.AuthRequest, r *http.Request) string {
	if authReq != nil {
		return authReq.ID
	}
	return r.FormValue(QueryAuthRequestID)
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
	Description            string
	Theme                  string
	ThemeMode              domain.LabelPolicyThemeMode
	ThemeClass             string
	DarkMode               bool
	PrivateLabelingOrgID   string
	OrgID                  string
	OrgName                string
	PrimaryDomain          string
	DisplayLoginNameSuffix bool
	TOSLink                string
	PrivacyLink            string
	HelpLink               string
	SupportEmail           string
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
	ErrPopup   bool
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
	MinLength    uint64
	HasUppercase string
	HasLowercase string
	HasNumber    string
	HasSymbol    string
	Expired      bool
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
	totpData
}

type mfaDoneData struct {
	baseData
	profileData
	MFAType domain.MFAType
}

type totpData struct {
	Url    string
	Secret string
	QrCode template.HTML
}
