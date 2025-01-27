package login

import (
	"context"
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	translator := l.getTranslator(r.Context(), authReq)
	data := loginSuccessData{
		userData: l.getUserData(r, authReq, translator, "LoginSuccess.Title", "", err),
	}
	if authReq != nil {
		data.RedirectURI, err = l.authRequestCallback(r.Context(), authReq)
		if err != nil {
			l.renderInternalError(w, r, authReq, err)
			return
		}
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplLoginSuccess], data, nil)
}

func (l *Login) redirectToCallback(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	callback, err := l.authRequestCallback(r.Context(), authReq)
	if err != nil {
		l.renderInternalError(w, r, authReq, err)
		return
	}
	http.Redirect(w, r, callback, http.StatusFound)
}

func (l *Login) authRequestCallback(ctx context.Context, authReq *domain.AuthRequest) (string, error) {
	switch authReq.Request.(type) {
	case *domain.AuthRequestOIDC:
		return l.oidcAuthCallbackURL(ctx, authReq.ID), nil
	case *domain.AuthRequestSAML:
		return l.samlAuthCallbackURL(ctx, authReq.ID), nil
	case *domain.AuthRequestDevice:
		return l.deviceAuthCallbackURL(authReq.ID), nil
	default:
		return "", zerrors.ThrowInternal(nil, "LOGIN-rhjQF", "Errors.AuthRequest.RequestTypeNotSupported")
	}
}
