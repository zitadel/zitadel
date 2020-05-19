package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	EndpointRoot             = "/"
	EndpointHealthz          = "/healthz"
	EndpointReadiness        = "/ready"
	EndpointLogin            = "/login"
	EndpointUsername         = "/username"
	EndpointUserSelection    = "/userselection"
	EndpointPassword         = "/password"
	EndpointInitPassword     = "/password/init"
	EndpointChangePassword   = "/password/change"
	EndpointPasswordReset    = "/password/reset"
	EndpointInitUser         = "/user/init"
	EndpointMfaVerify        = "/mfa/verify"
	EndpointMfaPrompt        = "/mfa/prompt"
	EndpointMfaInitVerify    = "/mfa/init/verify"
	EndpointMailVerification = "/mail/verification"
	EndpointMailVerified     = "/mail/verified"
	EndpointRegister         = "/register"
	EndpointLogoutDone       = "/logout/done"

	EndpointResources = "/resources"
)

func CreateRouter(login *Login, staticDir string) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(EndpointRoot, login.handleLogin).Methods(http.MethodGet)
	router.HandleFunc(EndpointHealthz, login.handleHealthz).Methods(http.MethodGet)
	router.HandleFunc(EndpointReadiness, login.handleReadiness).Methods(http.MethodGet)
	router.HandleFunc(EndpointLogin, login.handleLogin).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc(EndpointUsername, login.handleUsername).Methods(http.MethodGet)
	router.HandleFunc(EndpointUsername, login.handleUsernameCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointUserSelection, login.handleSelectUser).Methods(http.MethodPost)
	router.HandleFunc(EndpointPassword, login.handlePasswordCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointInitPassword, login.handleInitPassword).Methods(http.MethodGet)
	router.HandleFunc(EndpointInitPassword, login.handleInitPasswordCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointPasswordReset, login.handlePasswordReset).Methods(http.MethodGet)
	router.HandleFunc(EndpointInitUser, login.handleInitUser).Methods(http.MethodGet)
	router.HandleFunc(EndpointInitUser, login.handleInitUserCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointMfaVerify, login.handleMfaVerify).Methods(http.MethodPost)
	router.HandleFunc(EndpointMfaPrompt, login.handleMfaPrompt).Methods(http.MethodPost)
	router.HandleFunc(EndpointMfaInitVerify, login.handleMfaInitVerify).Methods(http.MethodPost)
	router.HandleFunc(EndpointMailVerification, login.handleMailVerification).Methods(http.MethodGet)
	router.HandleFunc(EndpointMailVerification, login.handleMailVerificationCheck).Methods(http.MethodPost)
	router.HandleFunc(EndpointChangePassword, login.handleChangePassword).Methods(http.MethodPost)
	router.HandleFunc(EndpointRegister, login.handleRegister).Methods(http.MethodPost)
	router.HandleFunc(EndpointLogoutDone, login.handleLogoutDone).Methods(http.MethodGet)
	router.PathPrefix(EndpointResources).Handler(login.handleResources(staticDir)).Methods(http.MethodGet)
	return router
}
