package login

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
)

func (l *Login) renderDeviceAuthUserCode(w io.Writer, err error) {
	data := struct {
		Error string
	}{}
	if err != nil {
		data.Error = err.Error()
	}
	err = l.renderer.Templates["device-usercode"].Execute(w, data)
	if err != nil {
		logrus.Error(err)
	}
}

func (l *Login) renderDeviceAuthConfirm(w http.ResponseWriter, username, clientID string, scopes []string) {
	data := &struct {
		Username string
		ClientID string
		Scopes   []string
	}{
		Username: username,
		ClientID: clientID,
		Scopes:   scopes,
	}

	err := l.renderer.Templates["device-confirm"].Execute(w, data)
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
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		l.renderDeviceAuthUserCode(w, err)
		return
	}
	userCode := r.Form.Get("user_code")
	if userCode == "" {
		if prompt, _ := url.QueryUnescape(r.Form.Get("prompt")); prompt != "" {
			err = errors.New(prompt)
		}
		l.renderDeviceAuthUserCode(w, err)
		return
	}
	deviceAuth, err := l.query.DeviceAuthByUserCode(r.Context(), userCode)
	if err != nil {
		l.renderDeviceAuthUserCode(w, err)
		return
	}
	agentID, ok := middleware.UserAgentIDFromCtx(r.Context())
	if !ok {
		l.renderDeviceAuthUserCode(w, errors.New("internal error: agent ID missing"))
		return
	}
	authRequest, err := l.authRepo.CreateAuthRequest(r.Context(), &domain.AuthRequest{
		AgentID:       agentID,
		ApplicationID: deviceAuth.ClientID,
		Request: &domain.AuthRequestDevice{
			DeviceCode: deviceAuth.DeviceCode,
			UserCode:   deviceAuth.UserCode,
			Scopes:     deviceAuth.Scopes,
		},
	})
	if err != nil {
		l.renderDeviceAuthUserCode(w, err)
		return
	}

	http.Redirect(w, r, l.renderer.pathPrefix+EndpointLogin+"?authRequestID="+authRequest.ID, http.StatusFound)
}

// redirectDeviceAuthStart redirects the user to the start point of
// the device authorization flow. A prompt can be set to inform the user
// of the reason why they are redirected back.
func redirectDeviceAuthStart(w http.ResponseWriter, r *http.Request, prompt string) {
	values := make(url.Values)
	values.Set("prompt", url.QueryEscape(prompt))

	url := url.URL{
		Path:     "/device",
		RawQuery: values.Encode(),
	}
	http.Redirect(w, r, url.String(), http.StatusSeeOther)
}

type deviceConfirmRequest struct {
	AuthRequestID string `schema:"authRequestID"`
	Action        string `schema:"action"`
}

// handleDeviceAuthConfirm is the handler where the user is redirected after login.
// The authRequest is checked if the login was indeed completed.
// When the action of "allowed" or "denied", the device authorization is updated accordingly.
// Else the user is presented with a page where they can choose / submit either action.
func (l *Login) handleDeviceAuthConfirm(w http.ResponseWriter, r *http.Request) {
	req := new(deviceConfirmRequest)
	if err := l.getParseData(r, req); err != nil {
		redirectDeviceAuthStart(w, r, err.Error())
		return

	}
	agentID, ok := middleware.UserAgentIDFromCtx(r.Context())
	if !ok {
		redirectDeviceAuthStart(w, r, "internal error: agent ID missing")
		return
	}
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), req.AuthRequestID, agentID)
	if err != nil {
		redirectDeviceAuthStart(w, r, err.Error())
		return
	}
	if !authReq.Done() {
		redirectDeviceAuthStart(w, r, "authentication not completed")
		return
	}
	authDev, ok := authReq.Request.(*domain.AuthRequestDevice)
	if !ok {
		redirectDeviceAuthStart(w, r, fmt.Sprintf("wrong auth request type: %T", authReq.Request))
		return
	}

	switch req.Action {
	case "allowed":
		_, err = l.command.ApproveDeviceAuth(r.Context(), authDev.ID, authReq.UserID)
	case "denied":
		_, err = l.command.DenyDeviceAuth(r.Context(), authDev.ID)
	default:
		l.renderDeviceAuthConfirm(w, authReq.UserName, authReq.ApplicationID, authDev.Scopes)
		return
	}
	if err != nil {
		redirectDeviceAuthStart(w, r, err.Error())
		return
	}

	fmt.Fprintf(w, "Device authorization %s. You can now return to the device", req.Action)
}

func (l *Login) deviceAuthCallbackURL(authRequestID string) string {
	return l.renderer.pathPrefix + EndpointDeviceAuthConfirm + "?authRequestID=" + authRequestID
}
