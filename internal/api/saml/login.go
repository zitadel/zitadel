package saml

import (
	"fmt"
	httpapi "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/saml"
	"net"
	"net/http"
)

func (p *IdentityProvider) loginHandleFunc(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Errorf("failed to parse form: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	requestID := r.Form.Get("requestId")
	if requestID == "" {
		http.Error(w, fmt.Errorf("no requestID provided").Error(), http.StatusInternalServerError)
		return
	}

	authRequest, err := p.storage.AuthRequestByID(r.Context(), requestID)

	if err != nil {
		if err := sendBackResponse(p.postTemplate, w, authRequest.GetRelayState(), "", makeDeniedResponse(
			"",
			"",
			p.EntityID,
			fmt.Errorf("failed to get request: %w", err).Error(),
		)); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
		return
	}

	entityID, err := p.storage.GetEntityIDByAppID(r.Context(), authRequest.GetApplicationID())
	if err != nil {
		http.Error(w, fmt.Errorf("failed to get entityID: %w", err).Error(), http.StatusInternalServerError)
	}

	nameID := authRequest.GetNameID()
	nameID = "stefan"
	ip := getIP(r).String()
	ip = "192.168.0.1"
	//TODO get userinfo
	attributes := []*saml.AttributeType{
		{
			Name:           "FirstName",
			NameFormat:     "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
			AttributeValue: []string{"stefan"},
		},
		{
			Name:           "SurName",
			NameFormat:     "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
			AttributeValue: []string{"benz"},
		},
		{
			Name:           "Email",
			NameFormat:     "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
			AttributeValue: []string{"stefan@caos.ch"},
		},
	}

	resp := makeSuccessfulResponse(
		authRequest,
		p.EntityID,
		ip,
		nameID,
		attributes,
		entityID,
	)

	signature, err := createSignatureP(p.signer, resp.Assertion)
	if err != nil {
		if err := sendBackResponse(p.postTemplate, w, authRequest.GetRelayState(), "", makeResponderFailResponse(
			"",
			"",
			p.EntityID,
			fmt.Errorf("failed to sign response: %w", err).Error(),
		)); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
	}
	resp.Assertion.Signature = signature

	if err := sendBackResponse(p.postTemplate, w, authRequest.GetRelayState(), authRequest.GetAccessConsumerServiceURL(), resp); err != nil {
		http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
	}
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
