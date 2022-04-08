package saml

import (
	"encoding/base64"
	"fmt"
	"github.com/caos/logging"
	httpapi "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/api/saml/signature"
	"github.com/caos/zitadel/internal/api/saml/xml"
	"net"
	"net/http"
	"net/url"
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
	case "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST":
		sig, err := signature.Create(p.signingContext, samlResponse)
		if err != nil {
			logging.Log("SAML-120dk2").Error(err)
			response.sendBackResponse(r, w, response.makeResponderFailResponse(fmt.Errorf("failed to sign response: %w", err).Error()))
			return
		}

		samlResponse.Signature = sig
	case "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect":
		respStr, err := xml.Marshal(samlResponse)
		if err != nil {
			logging.Log("SAML-95jw3k").Error(err)
			response.sendBackResponse(r, w, response.makeResponderFailResponse(fmt.Errorf("failed to marshal response: %w", err).Error()))
			return
		}

		respData, err := xml.DeflateAndBase64([]byte(respStr))
		if err != nil {
			logging.Log("SAML-po2n1s").Error(err)
			response.sendBackResponse(r, w, response.makeResponderFailResponse(fmt.Errorf("failed to deflate response: %w", err).Error()))
			return
		}

		query := "SAMLResponse=" + url.QueryEscape(string(respData))
		if response.RelayState != "" {
			query += "&RelayState=" + response.RelayState
		}
		if p.signingContext.GetSignatureMethodIdentifier() != "" {
			query += "&SigAlg=" + url.QueryEscape(p.signingContext.GetSignatureMethodIdentifier())
		}

		sig, err := p.signingContext.SignString(query)
		if err != nil {
			logging.Log("SAML-jwnu2i").Error(err)
			response.sendBackResponse(r, w, response.makeResponderFailResponse(fmt.Errorf("failed to sign response: %w", err).Error()))
			return
		}

		response.Signature = url.QueryEscape(base64.StdEncoding.EncodeToString(sig))
		response.SigAlg = url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(p.signingContext.GetSignatureMethodIdentifier())))
	}

	response.sendBackResponse(r, w, samlResponse)
	return
}

func getIP(request *http.Request) net.IP {
	return httpapi.RemoteIPFromRequest(request)
}

const postTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN"
"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">
<body onload="document.getElementById('samlpost').submit()">
<noscript>
<p>
<strong>Note:</strong> Since your browser does not support JavaScript,
you must press the Continue button once to proceed.
</p>
</noscript>
<form action="{{ .AssertionConsumerServiceURL }}" method="post" id="samlpost">
<div>
<input type="hidden" name="RelayState"
value="{{ .RelayState }}"/>
<input type="hidden" name="SAMLResponse"
value="{{ .SAMLResponse }}"/>
</div>
<noscript>
<div>
<input type="submit" value="Continue"/>
</div>
</noscript>
</form>
</body>
</html>`

const logoutTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN"
"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">
<body onload="document.getElementById('samlpost').submit()">
<noscript>
<p>
<strong>Note:</strong> Since your browser does not support JavaScript,
you must press the Continue button once to proceed.
</p>
</noscript>
<form action="{{ .LogoutURL }}" method="post" id="samlpost">
<div>
<input type="hidden" name="RelayState"
value="{{ .RelayState }}"/>
<input type="hidden" name="SAMLResponse"
value="{{ .SAMLResponse }}"/>
</div>
<noscript>
<div>
<input type="submit" value="Continue"/>
</div>
</noscript>
</form>
</body>
</html>`
