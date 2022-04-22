package saml

import (
	"fmt"
	"github.com/caos/logging"
	httpapi "github.com/caos/zitadel/internal/api/http"
	"net"
	"net/http"
)

func (p *IdentityProvider) callbackHandleFunc(w http.ResponseWriter, r *http.Request) {
	response := &Response{
		PostTemplate: p.postTemplate,
		ErrorFunc: func(err error) {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		},
		Issuer: p.EntityID,
	}

	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		logging.Log("SAML-91j1kk").Error(err)
		http.Error(w, fmt.Errorf("failed to parse form: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	requestID := r.Form.Get("id")
	if requestID == "" {
		err := fmt.Errorf("no requestID provided")
		logging.Log("SAML-91j1dk").Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authRequest, err := p.storage.AuthRequestByID(r.Context(), requestID)
	response.RequestID = authRequest.GetAuthRequestID()
	response.RelayState = authRequest.GetRelayState()
	response.ProtocolBinding = authRequest.GetBindingType()
	response.AcsUrl = authRequest.GetAccessConsumerServiceURL()
	if err != nil {
		logging.Log("SAML-91jp3k").Error(err)
		response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("failed to get request: %w", err).Error()))
		return
	}

	if !authRequest.Done() {
		logging.Log("SAML-91jp2k").Error(err)
		http.Error(w, fmt.Errorf("failed to get entityID: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	entityID, err := p.storage.GetEntityIDByAppID(r.Context(), authRequest.GetApplicationID())
	if err != nil {
		logging.Log("SAML-91jpdk").Error(err)
		http.Error(w, fmt.Errorf("failed to get entityID: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	response.Audience = entityID

	attrs := &Attributes{}
	if err := p.storage.SetUserinfoWithUserID(ctx, attrs, authRequest.GetUserID(), []int{}); err != nil {
		logging.Log("SAML-91jplp").Error(err)
		http.Error(w, fmt.Errorf("failed to get userinfo: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	response.SendIP = getIP(r).String()
	samlResponse := response.makeSuccessfulResponse(attrs)

	switch response.ProtocolBinding {
	case PostBinding:
		if err := createPostSignature(samlResponse, p); err != nil {
			logging.Log("SAML-120dk2").Error(err)
			response.sendBackResponse(r, w, response.makeResponderFailResponse(fmt.Errorf("failed to sign response: %w", err).Error()))
			return
		}
	case RedirectBinding:
		if err := createRedirectSignature(samlResponse, p, response); err != nil {
			logging.Log("SAML-jwnu2i").Error(err)
			response.sendBackResponse(r, w, response.makeResponderFailResponse(fmt.Errorf("failed to sign response: %w", err).Error()))
			return
		}
	}

	response.sendBackResponse(r, w, samlResponse)
	return
}

func getIP(request *http.Request) net.IP {
	return httpapi.RemoteIPFromRequest(request)
}
