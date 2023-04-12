package login

import (
	errs "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/muhlemmer/gu"
	"github.com/sirupsen/logrus"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
)

const (
	tmplDeviceAuthUserCode = "device-usercode"
	tmplDeviceAuthConfirm  = "device-confirm"
)

func (l *Login) renderDeviceAuthUserCode(w io.Writer, r *http.Request, err error) {
	var errID, errMessage string
	if err != nil {
		logging.WithError(err).Error()
		errID, errMessage = l.getErrorMessage(r, err)
	}

	data := l.getBaseData(r, &domain.AuthRequest{}, "DeviceAuth.Title", "DeviceAuth.Description", errID, errMessage)
	err = l.renderer.Templates[tmplDeviceAuthUserCode].Execute(w, data)
	if err != nil {
		logging.WithError(err).Error()
	}
}

func (l *Login) renderDeviceAuthAction(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, scopes []string) {
	data := &struct {
		baseData
		AuthRequestID string
		Username      string
		ClientID      string
		Scopes        []string
	}{
		baseData:      l.getBaseData(r, authReq, "DeviceAuth.Title", "DeviceAuth.Description", "", ""),
		AuthRequestID: authReq.ID,
		Username:      authReq.UserName,
		ClientID:      authReq.ApplicationID,
		Scopes:        scopes,
	}

	err := l.renderer.Templates[tmplDeviceAuthConfirm].Execute(w, data)
	if err != nil {
		logrus.Error(err)
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
			err = errs.New(prompt)
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
		l.renderDeviceAuthUserCode(w, r, errs.New("internal error: agent ID missing"))
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
		err = errors.ThrowInvalidArgument(err, "LOGIN-OLah8", "invalid or missing auth request")
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
	case "allowed":
		_, err = l.command.ApproveDeviceAuth(r.Context(), authDev.ID, authReq.UserID)
	case "denied":
		_, err = l.command.CancelDeviceAuth(r.Context(), authDev.ID, domain.DeviceAuthCanceledDenied)
	default:
		l.renderDeviceAuthAction(w, r, authReq, authDev.Scopes)
		return
	}
	if err != nil {
		l.redirectDeviceAuthStart(w, r, err.Error())
		return
	}

	fmt.Fprintf(w, "Device authorization %s. You can now return to the device", action)
}

func (l *Login) deviceAuthCallbackURL(authRequestID string) string {
	return l.renderer.pathPrefix + EndpointDeviceAuthAction + "?authRequestID=" + authRequestID
}

// RedirectDeviceAuthToPrefix allows users to use https://domain.com/device without the /ui/login prefix
// and redirects them to the prefixed endpoint.
// https://www.rfc-editor.org/rfc/rfc8628#section-3.2 recommends the URL to be as short as possible.
func RedirectDeviceAuthToPrefix(w http.ResponseWriter, r *http.Request) {
	target := gu.PtrCopy(r.URL)
	target.Path = HandlerPrefix + EndpointDeviceAuth
	http.Redirect(w, r, target.String(), http.StatusFound)
}
