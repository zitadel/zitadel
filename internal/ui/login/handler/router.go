package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	EndpointRoot                  = "/"
	EndpointHealthz               = "/healthz"
	EndpointReadiness             = "/ready"
	EndpointLogin                 = "/login"
	EndpointExternalLogin         = "/login/externalidp"
	EndpointExternalLoginCallback = "/login/externalidp/callback"
	EndpointLoginName             = "/loginname"
	EndpointUserSelection         = "/userselection"
	EndpointChangeUsername        = "/username/change"
	EndpointPassword              = "/password"
	EndpointInitPassword          = "/password/init"
	EndpointChangePassword        = "/password/change"
	EndpointPasswordReset         = "/password/reset"
	EndpointInitUser              = "/user/init"
	EndpointMfaVerify             = "/mfa/verify"
	EndpointMfaPrompt             = "/mfa/prompt"
	EndpointMfaInitVerify         = "/mfa/init/verify"
	EndpointMailVerification      = "/mail/verification"
	EndpointMailVerified          = "/mail/verified"
	EndpointRegister              = "/register"
	EndpointRegisterOrg           = "/register/org"
	EndpointLogoutDone            = "/logout/done"

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
	router.HandleFunc(EndpointMailVerification, login.handleMailVerification).Methods(http.MethodGet)
	router.HandleFunc(EndpointMailVerification, login.handleMailVerificationCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointChangePassword, login.handleChangePassword).Methods(http.MethodPost)
	router.HandleFunc(EndpointRegister, login.handleRegister).Methods(http.MethodGet)
	router.HandleFunc(EndpointRegister, login.handleRegisterCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointLogoutDone, login.handleLogoutDone).Methods(http.MethodGet)
	router.PathPrefix(EndpointResources).Handler(login.handleResources(staticDir)).Methods(http.MethodGet)
	router.HandleFunc(EndpointRegisterOrg, login.handleRegisterOrg).Methods(http.MethodGet)
	router.HandleFunc(EndpointRegisterOrg, login.handleRegisterOrgCheck).Methods(http.MethodPost)
	return router
}
