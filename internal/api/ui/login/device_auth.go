package login

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/muhlemmer/gu"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	tmplDeviceAuthUserCode = "device-usercode"
	tmplDeviceAuthAction   = "device-action"
)

func (l *Login) renderDeviceAuthUserCode(w http.ResponseWriter, r *http.Request, err error) {
	var errID, errMessage string
	if err != nil {
		logging.WithError(err).Error()
		errID, errMessage = l.getErrorMessage(r, err)
	}
	translator := l.getTranslator(r.Context(), nil)
	data := l.getBaseData(r, nil, translator, "DeviceAuth.Title", "DeviceAuth.UserCode.Description", errID, errMessage)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplDeviceAuthUserCode], data, nil)
}

func (l *Login) renderDeviceAuthAction(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, scopes []string) {
	translator := l.getTranslator(r.Context(), authReq)
	data := &struct {
		baseData
		AuthRequestID string
		Username      string
		ClientID      string
		Scopes        []string
	}{
		baseData:      l.getBaseData(r, authReq, translator, "DeviceAuth.Title", "DeviceAuth.Action.Description", "", ""),
		AuthRequestID: authReq.ID,
		Username:      authReq.UserName,
		ClientID:      authReq.ApplicationID,
		Scopes:        scopes,
	}

	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplDeviceAuthAction], data, nil)
}

const (
	deviceAuthAllowed = "allowed"
	deviceAuthDenied  = "denied"
)

// renderDeviceAuthDone renders success.html when the action was allowed and error.html when it was denied.
func (l *Login) renderDeviceAuthDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, action string) {
	translator := l.getTranslator(r.Context(), authReq)
	data := &struct {
		baseData
		Message string
	}{
		baseData: l.getBaseData(r, authReq, translator, "DeviceAuth.Title", "DeviceAuth.Done.Description", "", ""),
	}
	switch action {
	case deviceAuthAllowed:
		data.Message = translator.LocalizeFromRequest(r, "DeviceAuth.Done.Approved", nil)
		l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplSuccess], data, nil)
	case deviceAuthDenied:
		data.ErrMessage = translator.LocalizeFromRequest(r, "DeviceAuth.Done.Denied", nil)
		l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplError], data, nil)
	}
}

// handleDeviceUserCode serves the Device Authorization user code submission form.
// The "user_code" may be submitted by URL (GET) or form (POST).
// When a "user_code" is received and found through query,
// handleDeviceAuthUserCode will create a new AuthRequest in the repository.
// The user is then redirected to the /login endpoint to complete authentication.
//
// The agent ID from the context is set to the authentication request
// to ensure the complete login flow is completed from the same browser.
func (l *Login) handleDeviceAuthUserCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		l.renderDeviceAuthUserCode(w, r, err)
		return
	}
	userCode := r.Form.Get("user_code")
	if userCode == "" {
		if prompt, _ := url.QueryUnescape(r.Form.Get("prompt")); prompt != "" {
			err = errors.New(prompt)
		}
		l.renderDeviceAuthUserCode(w, r, err)
		return
	}
	deviceAuth, err := l.query.DeviceAuthByUserCode(ctx, userCode)
	if err != nil {
		l.renderDeviceAuthUserCode(w, r, err)
		return
	}
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		l.renderDeviceAuthUserCode(w, r, errors.New("internal error: agent ID missing"))
		return
	}
	authRequest, err := l.authRepo.CreateAuthRequest(ctx, &domain.AuthRequest{
		CreationDate:  time.Now(),
		AgentID:       userAgentID,
		ApplicationID: deviceAuth.ClientID,
		InstanceID:    authz.GetInstance(ctx).InstanceID(),
		Request: &domain.AuthRequestDevice{
			ID:         deviceAuth.AggregateID,
			DeviceCode: deviceAuth.DeviceCode,
			UserCode:   deviceAuth.UserCode,
			Scopes:     deviceAuth.Scopes,
		},
	})
	if err != nil {
		l.renderDeviceAuthUserCode(w, r, err)
		return
	}

	http.Redirect(w, r, l.renderer.pathPrefix+EndpointLogin+"?authRequestID="+authRequest.ID, http.StatusFound)
}

// redirectDeviceAuthStart redirects the user to the start point of
// the device authorization flow. A prompt can be set to inform the user
// of the reason why they are redirected back.
func (l *Login) redirectDeviceAuthStart(w http.ResponseWriter, r *http.Request, prompt string) {
	values := make(url.Values)
	values.Set("prompt", url.QueryEscape(prompt))

	url := url.URL{
		Path:     l.renderer.pathPrefix + EndpointDeviceAuth,
		RawQuery: values.Encode(),
	}
	http.Redirect(w, r, url.String(), http.StatusSeeOther)
}

// handleDeviceAuthAction is the handler where the user is redirected after login.
// The authRequest is checked if the login was indeed completed.
// When the action of "allowed" or "denied", the device authorization is updated accordingly.
// Else the user is presented with a page where they can choose / submit either action.
func (l *Login) handleDeviceAuthAction(w http.ResponseWriter, r *http.Request) {
	authReq, err := l.getAuthRequest(r)
	if authReq == nil {
		err = zerrors.ThrowInvalidArgument(err, "LOGIN-OLah8", "invalid or missing auth request")
		l.redirectDeviceAuthStart(w, r, err.Error())
		return
	}
	if !authReq.Done() {
		l.redirectDeviceAuthStart(w, r, "authentication not completed")
		return
	}
	authDev, ok := authReq.Request.(*domain.AuthRequestDevice)
	if !ok {
		l.redirectDeviceAuthStart(w, r, fmt.Sprintf("wrong auth request type: %T", authReq.Request))
		return
	}

	action := mux.Vars(r)["action"]
	switch action {
	case deviceAuthAllowed:
		_, err = l.command.ApproveDeviceAuth(r.Context(), authDev.ID, authReq.UserID)
	case deviceAuthDenied:
		_, err = l.command.CancelDeviceAuth(r.Context(), authDev.ID, domain.DeviceAuthCanceledDenied)
	default:
		l.renderDeviceAuthAction(w, r, authReq, authDev.Scopes)
		return
	}
	if err != nil {
		l.redirectDeviceAuthStart(w, r, err.Error())
		return
	}

	l.renderDeviceAuthDone(w, r, authReq, action)
}

// deviceAuthCallbackURL creates the callback URL with which the user
// is redirected back to the device authorization flow.
func (l *Login) deviceAuthCallbackURL(authRequestID string) string {
	return l.renderer.pathPrefix + EndpointDeviceAuthAction + "?authRequestID=" + authRequestID
}

// RedirectDeviceAuthToPrefix allows users to use https://domain.com/device without the /ui/login prefix
// and redirects them to the prefixed endpoint.
// [rfc 8628](https://www.rfc-editor.org/rfc/rfc8628#section-3.2) recommends the URL to be as short as possible.
func RedirectDeviceAuthToPrefix(w http.ResponseWriter, r *http.Request) {
	target := gu.PtrCopy(r.URL)
	target.Path = HandlerPrefix + EndpointDeviceAuth
	http.Redirect(w, r, target.String(), http.StatusFound)
}
