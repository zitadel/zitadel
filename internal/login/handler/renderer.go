package handler

import (
	"fmt"
	"html/template"

	"github.com/gorilla/csrf"

	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/renderer"
	"net/http"
	"path"

	"github.com/caos/logging"
	"golang.org/x/text/language"
)

const (
	tmplError = "error"
)

type Renderer struct {
	*renderer.Renderer
}

func CreateRenderer(staticDir http.FileSystem, cookieName string, defaultLanguage language.Tag) *Renderer {
	r := new(Renderer)
	tmplMapping := map[string]string{
		tmplError:              "error.html",
		tmplLogin:              "login.html",
		tmplUserSelection:      "select_user.html",
		tmplPassword:           "password.html",
		tmplMfaVerify:          "mfa_verify.html",
		tmplMfaPrompt:          "mfa_prompt.html",
		tmplMfaInitVerify:      "mfa_init_verify.html",
		tmplMfaInitDone:        "mfa_init_done.html",
		tmplMailVerification:   "mail_verification.html",
		tmplMailVerified:       "mail_verified.html",
		tmplInitPassword:       "init_password.html",
		tmplInitPasswordDone:   "init_password_done.html",
		tmplInitUser:           "init_user.html",
		tmplInitUserDone:       "init_user_done.html",
		tmplPasswordResetDone:  "password_reset_done.html",
		tmplChangePassword:     "change_password.html",
		tmplChangePasswordDone: "change_password_done.html",
		tmplRegister:           "register.html",
		tmplLogoutDone:         "logout_done.html",
	}
	funcs := map[string]interface{}{
		"resourceUrl": func(file string) string {
			return path.Join(EndpointResources, file)
		},
		"resourceThemeUrl": func(file, theme string) string {
			return path.Join(EndpointResources, "themes", theme, file)
		},
		"loginUrl": func() string {
			return EndpointLogin
		},
		"registerUrl": func(id string) string {
			return fmt.Sprintf("%s?%s=%s", EndpointRegister, queryAuthRequestID, id)
		},
		"usernameUrl": func() string {
			return EndpointUsername
		},
		"usernameChangeUrl": func(id string) string {
			return fmt.Sprintf("%s?%s=%s", EndpointUsername, queryAuthRequestID, id)
		},
		"userSelectionUrl": func() string {
			return EndpointUserSelection
		},
		"passwordResetUrl": func(id string) string {
			return fmt.Sprintf("%s?%s=%s", EndpointPasswordReset, queryAuthRequestID, id)
		},
		"passwordUrl": func() string {
			return EndpointPassword
		},
		"mfaVerifyUrl": func() string {
			return EndpointMfaVerify
		},
		"mfaPromptUrl": func() string {
			return EndpointMfaPrompt
		},
		"mfaInitVerifyUrl": func() string {
			return EndpointMfaInitVerify
		},
		"mailVerificationUrl": func() string {
			return EndpointMailVerification
		},
		"initPasswordUrl": func() string {
			return EndpointInitPassword
		},
		"initUserUrl": func() string {
			return EndpointInitUser
		},
		"changePasswordUrl": func() string {
			return EndpointChangePassword
		},
		"registrationUrl": func() string {
			return EndpointRegister
		},
		"selectedLanguage": func(l string) bool {
			return false
		},
		"selectedGender": func(g int32) bool {
			return false
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

func (l *Login) renderNextStep(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest) {
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), authReq.ID)
	if err != nil {
		l.renderInternalError(w, r, authReq, errors.ThrowInternal(nil, "APP-sio0W", "could not get authreq"))
	}
	if len(authReq.PossibleSteps) == 0 {
		l.renderInternalError(w, r, authReq, errors.ThrowInternal(nil, "APP-9sdp4", "no possible steps"))
		return
	}
	l.chooseNextStep(w, r, authReq, 0, nil)
}

func (l *Login) renderError(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	if authReq == nil || len(authReq.PossibleSteps) == 0 {
		l.renderInternalError(w, r, authReq, errors.ThrowInternal(err, "APP-OVOiT", "no possible steps"))
		return
	}
	l.chooseNextStep(w, r, authReq, 0, err)
}

func (l *Login) chooseNextStep(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, stepNumber int, err error) {
	switch step := authReq.PossibleSteps[stepNumber].(type) {
	case *model.LoginStep:
		if len(authReq.PossibleSteps) > 1 {
			l.chooseNextStep(w, r, authReq, 1, err)
			return
		}
		l.renderLogin(w, r, authReq, err)
	case *model.SelectUserStep:
		l.renderUserSelection(w, r, authReq, step)
	case *model.InitPasswordStep:
		l.renderInitPassword(w, r, authReq, authReq.UserID, "", err)
	case *model.PasswordStep:
		l.renderPassword(w, r, authReq, nil)
	case *model.MfaVerificationStep:
		l.renderMfaVerify(w, r, authReq, step, err)
	case *model.RedirectToCallbackStep:
		if len(authReq.PossibleSteps) > 1 {
			l.chooseNextStep(w, r, authReq, 1, err)
			return
		}
		l.redirectToCallback(w, r, authReq)
	case *model.ChangePasswordStep:
		l.renderChangePassword(w, r, authReq, err)
	case *model.VerifyEMailStep:
		l.renderMailVerification(w, r, authReq, "", err)
	case *model.MfaPromptStep:
		l.renderMfaPrompt(w, r, authReq, step, err)
	case *model.InitUserStep:
		l.renderInitUser(w, r, authReq, "", "", nil)
	default:
		l.renderInternalError(w, r, authReq, errors.ThrowInternal(nil, "APP-ds3QF", "step no possible"))
	}
}

func (l *Login) renderInternalError(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var msg string
	if err != nil {
		msg = err.Error()
	}
	data := l.getBaseData(r, authReq, "Error", "Internal", msg)
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplError], data, nil)
}

func (l *Login) getBaseData(r *http.Request, authReq *model.AuthRequest, title string, errType, errMessage string) baseData {
	return baseData{
		errorData: errorData{
			ErrType:    errType,
			ErrMessage: errMessage,
		},
		Lang:      l.renderer.Lang(r).String(),
		Title:     title,
		Theme:     l.getTheme(r),
		ThemeMode: l.getThemeMode(r),
		AuthReqID: getRequestID(authReq, r),
		CSRF:      csrf.TemplateField(r),
		Nonce:     middleware.GetNonce(r),
	}
}

func (l *Login) getTheme(r *http.Request) string {
	return "zitadel" //TODO: impl
}

func (l *Login) getThemeMode(r *http.Request) string {
	return "" //TODO: impl
}

func getRequestID(authReq *model.AuthRequest, r *http.Request) string {
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
	Lang      string
	Title     string
	Theme     string
	ThemeMode string
	AuthReqID string
	CSRF      template.HTML
	Nonce     string
}

type errorData struct {
	ErrType    string
	ErrMessage string
}

type userData struct {
	baseData
	UserName            string
	PasswordChecked     string
	MfaProviders        []model.MfaType
	SelectedMfaProvider model.MfaType
}

type userSelectionData struct {
	baseData
	Users []model.UserSelection
}

type mfaData struct {
	baseData
	UserName     string
	MfaProviders []model.MfaType
	MfaRequired  bool
}

type mfaVerifyData struct {
	baseData
	UserName string
	MfaType  model.MfaType
	otpData
}

type mfaDoneData struct {
	baseData
	UserName string
	MfaType  model.MfaType
}

type otpData struct {
	Url    string
	Secret string
	QrCode string
}
