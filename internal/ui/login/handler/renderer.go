package handler

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/static"

	"github.com/caos/logging"
	"github.com/gorilla/csrf"
	"golang.org/x/text/language"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/renderer"
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
		tmplError:                    "error.html",
		tmplLogin:                    "login.html",
		tmplUserSelection:            "select_user.html",
		tmplPassword:                 "password.html",
		tmplPasswordlessVerification: "passwordless.html",
		tmplMFAVerify:                "mfa_verify.html",
		tmplMFAPrompt:                "mfa_prompt.html",
		tmplMFAInitVerify:            "mfa_init_verify.html",
		tmplMFAU2FInit:               "mfa_init_u2f.html",
		tmplU2FVerification:          "mfa_verification_u2f.html",
		tmplMFAInitDone:              "mfa_init_done.html",
		tmplMailVerification:         "mail_verification.html",
		tmplMailVerified:             "mail_verified.html",
		tmplInitPassword:             "init_password.html",
		tmplInitPasswordDone:         "init_password_done.html",
		tmplInitUser:                 "init_user.html",
		tmplInitUserDone:             "init_user_done.html",
		tmplPasswordResetDone:        "password_reset_done.html",
		tmplChangePassword:           "change_password.html",
		tmplChangePasswordDone:       "change_password_done.html",
		tmplRegisterOption:           "register_option.html",
		tmplRegister:                 "register.html",
		tmplLogoutDone:               "logout_done.html",
		tmplRegisterOrg:              "register_org.html",
		tmplChangeUsername:           "change_username.html",
		tmplChangeUsernameDone:       "change_username_done.html",
		tmplLinkUsersDone:            "link_users_done.html",
		tmplExternalNotFoundOption:   "external_not_found_option.html",
	}
	funcs := map[string]interface{}{
		"resourceUrl": func(file string) string {
			return path.Join(r.pathPrefix, EndpointResources, file)
		},
		"resourceThemeUrl": func(file, theme string) string {
			return path.Join(r.pathPrefix, EndpointResources, "themes", theme, file)
		},
		"variablesCssFileUrl": func(orgID, theme string) string {
			defaultURL := path.Join(r.pathPrefix, EndpointResources, "themes", theme, "variables.css")

			if r.staticStorage == nil {
				return defaultURL
			}

			if orgID != "" {
				presignedURL, err := r.staticStorage.GetObjectPresignedURL(context.TODO(), orgID, domain.OrgCssPath+"/"+domain.CssVariablesFileName, time.Hour*1)
				if err == nil {
					return presignedURL.String()
				}
			}

			presignedURL, err := r.staticStorage.GetObjectPresignedURL(context.TODO(), domain.IAMID, domain.IAMCssPath+"/"+domain.CssVariablesFileName, time.Hour*1)
			if err != nil {
				return defaultURL
			}
			return presignedURL.String()
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
	case *domain.SelectUserStep:
		l.renderUserSelection(w, r, authReq, step)
	case *domain.InitPasswordStep:
		l.renderInitPassword(w, r, authReq, authReq.UserID, "", err)
	case *domain.PasswordStep:
		l.renderPassword(w, r, authReq, nil)
	case *domain.PasswordlessStep:
		l.renderPasswordlessVerification(w, r, authReq, nil)
	case *domain.MFAVerificationStep:
		l.renderMFAVerify(w, r, authReq, step, err)
	case *domain.RedirectToCallbackStep:
		if len(authReq.PossibleSteps) > 1 {
			l.chooseNextStep(w, r, authReq, 1, err)
			return
		}
		l.redirectToCallback(w, r, authReq)
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
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplError], data, nil)
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
			ErrType:    errType,
			ErrMessage: errMessage,
		},
		Lang:                   l.renderer.Lang(r).String(),
		Title:                  title,
		Theme:                  l.getTheme(r),
		ThemeMode:              l.getThemeMode(r),
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
		baseData.IDPProviders = authReq.AllowedExternalIDPs
	}
	return baseData
}

func (l *Login) getProfileData(authReq *domain.AuthRequest) profileData {
	var userName, loginName, displayName string
	if authReq != nil {
		userName = authReq.UserName
		loginName = authReq.LoginName
		displayName = authReq.DisplayName
	}
	return profileData{
		UserName:    userName,
		LoginName:   loginName,
		DisplayName: displayName,
	}
}

func (l *Login) getErrorMessage(r *http.Request, err error) (errMsg string) {
	caosErr := new(caos_errs.CaosError)
	if errors.As(err, &caosErr) {
		localized := l.renderer.LocalizeFromRequest(r, caosErr.Message, nil)
		return localized + " (" + caosErr.ID + ")"

	}
	return err.Error()
}

func (l *Login) getTheme(r *http.Request) string {
	return "zitadel" //TODO: impl
}

func (l *Login) getThemeMode(r *http.Request) string {
	return "lgn-dark-theme" //TODO: impl
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
	OrgID                  string
	OrgName                string
	PrimaryDomain          string
	DisplayLoginNameSuffix bool
	AuthReqID              string
	CSRF                   template.HTML
	Nonce                  string
	LoginPolicy            *domain.LoginPolicy
	IDPProviders           []*domain.IDPProvider
}

type errorData struct {
	ErrType    string
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
