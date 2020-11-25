package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	EndpointRoot                     = "/"
	EndpointHealthz                  = "/healthz"
	EndpointReadiness                = "/ready"
	EndpointLogin                    = "/login"
	EndpointExternalLogin            = "/login/externalidp"
	EndpointExternalLoginCallback    = "/login/externalidp/callback"
	EndpointPasswordlessLogin        = "/login/passwordless"
	EndpointLoginName                = "/loginname"
	EndpointUserSelection            = "/userselection"
	EndpointChangeUsername           = "/username/change"
	EndpointPassword                 = "/password"
	EndpointInitPassword             = "/password/init"
	EndpointChangePassword           = "/password/change"
	EndpointPasswordReset            = "/password/reset"
	EndpointInitUser                 = "/user/init"
	EndpointMfaVerify                = "/mfa/verify"
	EndpointMfaPrompt                = "/mfa/prompt"
	EndpointMfaInitVerify            = "/mfa/init/verify"
	EndpointMfaInitU2FVerify         = "/mfa/init/u2f/verify"
	EndpointU2FVerification          = "/mfa/u2f/verify"
	EndpointMailVerification         = "/mail/verification"
	EndpointMailVerified             = "/mail/verified"
	EndpointRegisterOption           = "/register/option"
	EndpointRegister                 = "/register"
	EndpointExternalRegister         = "/register/externalidp"
	EndpointExternalRegisterCallback = "/register/externalidp/callback"
	EndpointRegisterOrg              = "/register/org"
	EndpointLogoutDone               = "/logout/done"
	EndpointExternalNotFoundOption   = "/externaluser/option"

	EndpointResources = "/resources"
)

func CreateRouter(login *Login, staticDir http.FileSystem, interceptors ...mux.MiddlewareFunc) *mux.Router {
	router := mux.NewRouter()
	router.Use(interceptors...)
	router.HandleFunc(EndpointRoot, login.handleLogin).Methods(http.MethodGet)
	router.HandleFunc(EndpointHealthz, login.handleHealthz).Methods(http.MethodGet)
	router.HandleFunc(EndpointReadiness, login.handleReadiness).Methods(http.MethodGet)
	router.HandleFunc(EndpointLogin, login.handleLogin).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc(EndpointExternalLogin, login.handleExternalLogin).Methods(http.MethodGet)
	router.HandleFunc(EndpointExternalLoginCallback, login.handleExternalLoginCallback).Methods(http.MethodGet)
	router.HandleFunc(EndpointPasswordlessLogin, login.handlePasswordlessVerification).Methods(http.MethodPost)
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
	router.HandleFunc(EndpointMfaVerify, login.handleMfaVerify).Methods(http.MethodPost)
	router.HandleFunc(EndpointMfaPrompt, login.handleMfaPromptSelection).Methods(http.MethodGet)
	router.HandleFunc(EndpointMfaPrompt, login.handleMfaPrompt).Methods(http.MethodPost)
	router.HandleFunc(EndpointMfaInitVerify, login.handleMfaInitVerify).Methods(http.MethodPost)
	router.HandleFunc(EndpointMfaInitU2FVerify, login.handleRegisterU2F).Methods(http.MethodPost)
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
	router.HandleFunc(EndpointExternalRegisterCallback, login.handleExternalRegisterCallback).Methods(http.MethodGet)
	router.HandleFunc(EndpointLogoutDone, login.handleLogoutDone).Methods(http.MethodGet)
	router.PathPrefix(EndpointResources).Handler(login.handleResources(staticDir)).Methods(http.MethodGet)
	router.HandleFunc(EndpointRegisterOrg, login.handleRegisterOrg).Methods(http.MethodGet)
	router.HandleFunc(EndpointRegisterOrg, login.handleRegisterOrgCheck).Methods(http.MethodPost)
	return router
}
