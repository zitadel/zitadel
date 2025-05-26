package login

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	EndpointRoot                          = "/"
	EndpointHealthz                       = "/healthz"
	EndpointReadiness                     = "/ready"
	EndpointLogin                         = "/login"
	EndpointExternalLogin                 = "/login/externalidp"
	EndpointExternalLoginCallback         = "/login/externalidp/callback"
	EndpointExternalLoginCallbackFormPost = "/login/externalidp/callback/form"
	EndpointExternalLogout                = "/logout/externalidp"
	EndpointSAMLACS                       = "/login/externalidp/saml/acs"
	EndpointJWTAuthorize                  = "/login/jwt/authorize"
	EndpointJWTCallback                   = "/login/jwt/callback"
	EndpointLDAPLogin                     = "/login/ldap"
	EndpointLDAPCallback                  = "/login/ldap/callback"
	EndpointPasswordlessLogin             = "/login/passwordless"
	EndpointPasswordlessRegistration      = "/login/passwordless/init"
	EndpointPasswordlessPrompt            = "/login/passwordless/prompt"
	EndpointLoginName                     = "/loginname"
	EndpointUserSelection                 = "/userselection"
	EndpointChangeUsername                = "/username/change"
	EndpointPassword                      = "/password"
	EndpointInitPassword                  = "/password/init"
	EndpointChangePassword                = "/password/change"
	EndpointPasswordReset                 = "/password/reset"
	EndpointInitUser                      = "/user/init"
	EndpointInviteUser                    = "/user/invite"
	EndpointMFAVerify                     = "/mfa/verify"
	EndpointMFAPrompt                     = "/mfa/prompt"
	EndpointMFAInitVerify                 = "/mfa/init/verify"
	EndpointMFASMSInitVerify              = "/mfa/init/sms/verify"
	EndpointMFAOTPVerify                  = "/mfa/otp/verify"
	EndpointMFAInitU2FVerify              = "/mfa/init/u2f/verify"
	EndpointU2FVerification               = "/mfa/u2f/verify"
	EndpointMailVerification              = "/mail/verification"
	EndpointMailVerified                  = "/mail/verified"
	EndpointRegisterOption                = "/register/option"
	EndpointRegister                      = "/register"
	EndpointExternalRegister              = "/register/externalidp"
	EndpointExternalRegisterCallback      = "/register/externalidp/callback"
	EndpointRegisterOrg                   = "/register/org"
	EndpointLogoutDone                    = "/logout/done"
	EndpointLoginSuccess                  = "/login/success"
	EndpointExternalNotFoundOption        = "/externaluser/option"

	EndpointResources        = "/resources"
	EndpointDynamicResources = "/resources/dynamic"

	EndpointDeviceAuth       = "/device"
	EndpointDeviceAuthAction = "/device/{action}"

	EndpointLinkingUserPrompt = "/link/user"
)

var (
	IgnoreInstanceEndpoints = []string{
		EndpointResources + "/fonts",
		EndpointResources + "/images",
		EndpointResources + "/scripts",
		EndpointResources + "/themes",
	}
)

func CreateRouter(login *Login, interceptors ...mux.MiddlewareFunc) *mux.Router {
	router := mux.NewRouter()
	router.Use(interceptors...)
	router.HandleFunc(EndpointRoot, login.handleLogin).Methods(http.MethodGet)
	router.HandleFunc(EndpointHealthz, login.handleHealthz).Methods(http.MethodGet)
	router.HandleFunc(EndpointReadiness, login.handleReadiness).Methods(http.MethodGet)
	router.HandleFunc(EndpointLogin, login.handleLogin).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc(EndpointExternalLogin, login.handleExternalLogin).Methods(http.MethodGet)
	router.HandleFunc(EndpointExternalLoginCallback, login.handleExternalLoginCallback).Methods(http.MethodGet)
	router.HandleFunc(EndpointExternalLoginCallbackFormPost, login.handleExternalLoginCallbackForm).Methods(http.MethodPost)
	router.HandleFunc(EndpointExternalLogout, login.handleExternalLogout).Methods(http.MethodGet)
	router.HandleFunc(EndpointSAMLACS, login.handleExternalLoginCallback).Methods(http.MethodGet)
	router.HandleFunc(EndpointSAMLACS, login.handleExternalLoginCallbackForm).Methods(http.MethodPost)
	router.HandleFunc(EndpointJWTAuthorize, login.handleJWTRequest).Methods(http.MethodGet)
	router.HandleFunc(EndpointJWTCallback, login.handleJWTCallback).Methods(http.MethodGet)
	router.HandleFunc(EndpointPasswordlessLogin, login.handlePasswordlessVerification).Methods(http.MethodPost)
	router.HandleFunc(EndpointPasswordlessRegistration, login.handlePasswordlessRegistration).Methods(http.MethodGet)
	router.HandleFunc(EndpointPasswordlessRegistration, login.handlePasswordlessRegistrationCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointPasswordlessPrompt, login.handlePasswordlessPrompt).Methods(http.MethodPost)
	router.HandleFunc(EndpointLoginName, login.handleLoginName).Methods(http.MethodGet)
	router.HandleFunc(EndpointLoginName, login.handleLoginNameCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointUserSelection, login.handleSelectUser).Methods(http.MethodPost)
	router.HandleFunc(EndpointChangeUsername, login.handleChangeUsername).Methods(http.MethodPost)
	router.HandleFunc(EndpointPassword, login.handlePasswordCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointInitPassword, login.handleInitPassword).Methods(http.MethodGet)
	router.HandleFunc(EndpointInitPassword, login.handleInitPasswordCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointPasswordReset, login.handlePasswordReset).Methods(http.MethodGet)
	router.HandleFunc(EndpointInitUser, login.handleInitUser).Methods(http.MethodGet)
	router.HandleFunc(EndpointInitUser, login.handleInitUserCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointInviteUser, login.handleInviteUser).Methods(http.MethodGet)
	router.HandleFunc(EndpointInviteUser, login.handleInviteUserCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointMFAVerify, login.handleMFAVerify).Methods(http.MethodPost)
	router.HandleFunc(EndpointMFAPrompt, login.handleMFAPromptSelection).Methods(http.MethodGet)
	router.HandleFunc(EndpointMFAPrompt, login.handleMFAPrompt).Methods(http.MethodPost)
	router.HandleFunc(EndpointMFAInitVerify, login.handleMFAInitVerify).Methods(http.MethodPost)
	router.HandleFunc(EndpointMFASMSInitVerify, login.handleRegisterSMSCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointMFAOTPVerify, login.handleOTPVerificationCheck).Methods(http.MethodGet)
	router.HandleFunc(EndpointMFAOTPVerify, login.handleOTPVerificationCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointMFAInitU2FVerify, login.handleRegisterU2F).Methods(http.MethodPost)
	router.HandleFunc(EndpointU2FVerification, login.handleU2FVerification).Methods(http.MethodPost)
	router.HandleFunc(EndpointMailVerification, login.handleMailVerification).Methods(http.MethodGet)
	router.HandleFunc(EndpointMailVerification, login.handleMailVerificationCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointChangePassword, login.handleChangePassword).Methods(http.MethodPost)
	router.HandleFunc(EndpointRegisterOption, login.handleRegisterOption).Methods(http.MethodGet)
	router.HandleFunc(EndpointRegisterOption, login.handleRegisterOptionCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointExternalNotFoundOption, login.handleExternalNotFoundOptionCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointRegister, login.handleRegister).Methods(http.MethodGet)
	router.HandleFunc(EndpointRegister, login.handleRegisterCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointExternalRegister, login.handleExternalRegister).Methods(http.MethodGet)
	router.HandleFunc(EndpointExternalRegisterCallback, login.handleExternalLoginCallback).Methods(http.MethodGet)
	router.HandleFunc(EndpointLogoutDone, login.handleLogoutDone).Methods(http.MethodGet)
	router.HandleFunc(EndpointDynamicResources, login.handleDynamicResources).Methods(http.MethodGet)
	router.PathPrefix(EndpointResources).Handler(login.handleResources()).Methods(http.MethodGet)
	router.HandleFunc(EndpointRegisterOrg, login.handleRegisterOrg).Methods(http.MethodGet)
	router.HandleFunc(EndpointRegisterOrg, login.handleRegisterOrgCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointLoginSuccess, login.handleLoginSuccess).Methods(http.MethodGet)
	router.HandleFunc(EndpointLDAPLogin, login.handleLDAP).Methods(http.MethodGet)
	router.HandleFunc(EndpointLDAPCallback, login.handleLDAPCallback).Methods(http.MethodPost)
	router.SkipClean(true).Handle("", http.RedirectHandler(HandlerPrefix+"/", http.StatusMovedPermanently))
	router.HandleFunc(EndpointDeviceAuth, login.handleDeviceAuthUserCode).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc(EndpointDeviceAuthAction, login.handleDeviceAuthAction).Methods(http.MethodGet, http.MethodPost)
	return router
}
