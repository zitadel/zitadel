package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

const (
	tmplLoginSuccess = "login_success"
)

type loginSuccessData struct {
	userData
	RedirectURI string `schema:"redirect-uri"`
}

func (l *Login) redirectToLoginSuccess(w http.ResponseWriter, r *http.Request, id string) {
	http.Redirect(w, r, l.renderer.pathPrefix+EndpointLoginSuccess+"?authRequestID="+id, http.StatusFound)
}

func (l *Login) handleLoginSuccess(w http.ResponseWriter, r *http.Request) {
	authRequest, _ := l.getAuthRequest(r)
	if authRequest == nil {
		l.renderSuccessAndCallback(w, r, nil, nil)
		return
	}
	for _, step := range authRequest.PossibleSteps {
		if step.Type() != domain.NextStepLoginSucceeded && step.Type() != domain.NextStepRedirectToCallback {
			l.renderNextStep(w, r, authRequest)
			return
		}
	}
	l.renderSuccessAndCallback(w, r, authRequest, nil)
}

func (l *Login) renderSuccessAndCallback(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	data := loginSuccessData{
		userData: l.getUserData(r, authReq, "LoginSuccess.Title", "", errID, errMessage),
	}
	if authReq != nil {
		//the id will be set via the html (maybe change this with the login refactoring)
		if _, ok := authReq.Request.(*domain.AuthRequestOIDC); ok {
			data.RedirectURI = l.oidcAuthCallbackURL(r.Context(), "")
		} else if _, ok := authReq.Request.(*domain.AuthRequestSAML); ok {
			data.RedirectURI = l.samlAuthCallbackURL(r.Context(), "")
		}
	}
	l.renderer.RenderTemplate(w, r, l.getTranslator(r.Context(), authReq), l.renderer.Templates[tmplLoginSuccess], data, nil)
}

func (l *Login) redirectToCallback(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	var callback string
	switch authReq.Request.(type) {
	case *domain.AuthRequestOIDC:
		callback = l.oidcAuthCallbackURL(r.Context(), authReq.ID)
	case *domain.AuthRequestSAML:
		callback = l.samlAuthCallbackURL(r.Context(), authReq.ID)
	case *domain.AuthRequestDevice:
		callback = l.deviceAuthCallbackURL(authReq.ID)
	default:
		l.renderInternalError(w, r, authReq, caos_errs.ThrowInternal(nil, "LOGIN-rhjQF", "Errors.AuthRequest.RequestTypeNotSupported"))
		return
	}
	http.Redirect(w, r, callback, http.StatusFound)
}
