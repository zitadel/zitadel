package saml

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/samlp"
	"net/http"
)

const (
	EncodingDeflate = "urn:oasis:names:tc:SAML:2.0:bindings:URL-Encoding:DEFLATE"
)

type AuthRequestForm struct {
	AuthRequest string
	Encoding    string
	RelayState  string
	SigAlg      string
	Sig         string
}

func (p *IdentityProvider) ssoHandleFunc(w http.ResponseWriter, r *http.Request) {
	authRequestForm, err := getAuthRequestFromRequest(r)
	if err != nil {
		logging.Log("SAML-837n2s").Error(err)
		http.Error(w, fmt.Errorf("failed to parse form: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	if err := verifyForm(authRequestForm); err != nil {
		logging.Log("SAML-827n2s").Error(err)
		if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeDeniedResponse(
			"",
			"",
			p.EntityID,
			fmt.Errorf("failed to validate form: %w", err).Error(),
		)); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
		return
	}
	authRequestForm.AuthRequest = "PHNhbWxwOkF1dGhuUmVxdWVzdCB4bWxuczpzYW1scD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIiBEZXN0aW5hdGlvbj0iaHR0cHM6Ly9hY2NvdW50cy5vcmJvcy5pby9zYW1sL1NTTyIgSUQ9Il81ZWI0MjM5NWE2ODBiYmI3MDM4ODY5MDg1NGVhZGJmOSIgSXNzdWVJbnN0YW50PSIyMDIyLTAzLTAxVDEyOjE3OjMwWiIgUHJvdG9jb2xCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1QT1NUIiBWZXJzaW9uPSIyLjAiPjxzYW1sOklzc3VlciB4bWxuczpzYW1sPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YXNzZXJ0aW9uIj51cm46YXV0aDA6b3Jib3MtaW86b3Jib3M8L3NhbWw6SXNzdWVyPjxTaWduYXR1cmUgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyMiPjxTaWduZWRJbmZvPjxDYW5vbmljYWxpemF0aW9uTWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS8xMC94bWwtZXhjLWMxNG4jIi8+PFNpZ25hdHVyZU1ldGhvZCBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvMDQveG1sZHNpZy1tb3JlI3JzYS1zaGEyNTYiLz48UmVmZXJlbmNlIFVSST0iI181ZWI0MjM5NWE2ODBiYmI3MDM4ODY5MDg1NGVhZGJmOSI+PFRyYW5zZm9ybXM+PFRyYW5zZm9ybSBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNlbnZlbG9wZWQtc2lnbmF0dXJlIi8+PFRyYW5zZm9ybSBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvMTAveG1sLWV4Yy1jMTRuIyIvPjwvVHJhbnNmb3Jtcz48RGlnZXN0TWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS8wNC94bWxlbmMjc2hhMjU2Ii8+PERpZ2VzdFZhbHVlPlFqdVpUYWgrQm5jM1JObCs2N0JFZlBPclBkSGhyZUpmK0hRdEhQWkoyTkU9PC9EaWdlc3RWYWx1ZT48L1JlZmVyZW5jZT48L1NpZ25lZEluZm8+PFNpZ25hdHVyZVZhbHVlPmtnY1RSdVZUbEFZU2ovYVUrRytlb3FjeExBT0VMYzNFNXhLRUZlK3hHc2tJVGJtWlVScnU3bEFhOWFRVlFQTGI0bFk1S3VRTU9KekttNko0MHVEYVB3U1lpMTlaSWtQZFJ1TVRmYlFFbXNtdDRiNmxBRkQzYWt3SXdnZ2dNNTg4aHAwY09ReEdaNEI4azBWM1RJcDNmKytGQ3F3TkVEaTFmZndHSlNUWm5oYmtmU3BNTlFIbUxyZVJ3ZDBaeGovR3hhblBMSUUvWHRVb0RkZ1M2dGwxemVnYXRVaG5nb1BTZ2RjaXEzNEh5RnB2eVRrTldLcGt0U1R6VjNTWTNzVGl3djRkK3FibzhpT0MwVGI2clZSTXlTaDN1OEZWWXJMWGk3OGdWa01kdmRtZ2hFMDhFSXlOMlFwWFI5aUJyMEVBYklqb3h3cER1SVR6VzI0ZHhSQ0dMQT09PC9TaWduYXR1cmVWYWx1ZT48S2V5SW5mbz48WDUwOURhdGE+PFg1MDlDZXJ0aWZpY2F0ZT5NSUlEQlRDQ0FlMmdBd0lCQWdJSmJ1QU12QVNoZzBrY01BMEdDU3FHU0liM0RRRUJDd1VBTUNBeEhqQWNCZ05WQkFNVEZXOXlZbTl6TFdsdkxtVjFMbUYxZEdnd0xtTnZiVEFlRncweU1qQXpNREV3T1RRMk1EQmFGdzB6TlRFeE1EZ3dPVFEyTURCYU1DQXhIakFjQmdOVkJBTVRGVzl5WW05ekxXbHZMbVYxTG1GMWRHZ3dMbU52YlRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBSzl2MFY2WldLY2FGR0FJcHdUaGVXMkNDVmhOUFVGNFNRRmQybzRTcDdEOWFaTUY1QzR0SUZGejZJRGZ0T3piV0o2RWYxdWhhbDRxdEsvN08vK3ZOaWJrK2VQNzFORlR2cTVMTFBiN2ovcW42WElkbncwZ3VSNW5zRVlyejljaWR4OUpLZGpacjQzRm9JcWs1enM3dms5cnEwNENLU09LVVpqT2pXbmxiS2xsNGdES2NnZ2p6TStjUHZHVDBPWlIyU3laV21kSTFZQlM0ZnUvT2FodEJmR2JKN2xGUjRjdDlXQVR2MEkyNlB1VGFpUjJyYVd5MWd1T1ZsaDZQOHdTeDBHQjJ5S0lGekFUZUVyZmFla1JZcGd0Q1ZPSTA2WExKSmp2MmRKY0MvdWc3S1V1UnVRTmwvRVdSNjRHdVpMbWZCai8vQXgzZXhLbzQ2eGZwN21IQ3VjQ0F3RUFBYU5DTUVBd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVW5LczN4YndyTjBvS3g5cWVqV3Bpd1Qvb3g3NHdEZ1lEVlIwUEFRSC9CQVFEQWdLRU1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRQmpmeld5YkRneTZDdGpZSy9iVUF3cU92TG1IR0pubHQ0UWI1cVJZb1o0Qm5LWVo0cnlZTjM1NWRSdXhJU0p6T1FuRWZNS3hoS1A0T1pRVzRHdG5rdEQwT1huOVBGdkRwc1lyR2hQUGRzMjdOT2o1RUJMenRKWjMvbk9RV0NBNEdvUkM1WktlNVlKNW1qZGthVVAyTmE1eE9YSGo5OHhPTU1VVXgwVEFBanNnM0Nub0VsWFlVcXc1UC9aRkw0WFVrUFFERU5hUEhIVjE4bEJXSnc2SHVSZ3FwR3pTeWdDVjRrQ21HOW10cXRzeUlBdzJCK0ZXM25XdjRGRGl3aUxlMDNCejR6SFdUUWlXc2hhOEtoaEhFYzB6aDRoU3VkSTZjNkxWYWRjMjJhSytucVlDSjh1bnNKMEVicHR1bFdoaUU0R09qQmJTZVFtbEU3b2ZFMnRkUW1EPC9YNTA5Q2VydGlmaWNhdGU+PC9YNTA5RGF0YT48L0tleUluZm8+PC9TaWduYXR1cmU+PC9zYW1scDpBdXRoblJlcXVlc3Q+"

	authNRequest, err := decodeAuthNRequest(authRequestForm.Encoding, authRequestForm.AuthRequest)
	if err != nil {
		logging.Log("SAML-837s2s").Error(err)
		if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeDeniedResponse(
			authNRequest.Id,
			authNRequest.AssertionConsumerServiceURL,
			p.EntityID,
			fmt.Errorf("failed to decode request: %w", err).Error(),
		)); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
		return
	}

	sp, err := p.GetServiceProvider(r.Context(), authNRequest.Issuer.Text)
	if err != nil {
		logging.Log("SAML-317s2s").Error(err)
		http.Error(w, fmt.Errorf("failed to find registered serviceprovider: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	if sp == nil {
		logging.Log("SAML-837nas").Error(err)
		if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeDeniedResponse(
			authNRequest.Id,
			authNRequest.AssertionConsumerServiceURL,
			p.EntityID,
			fmt.Errorf("unknown service provider: %w", err).Error(),
		)); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
		return
	}

	if sp.metadata.SPSSODescriptor.AuthnRequestsSigned == "true" ||
		p.Metadata.WantAuthnRequestsSigned == "true" ||
		authRequestForm.Sig != "" ||
		(authNRequest.Signature != nil && authNRequest.Signature.SignatureValue.Text != "") {
		//DEFLATE encoding sends signature information with the parameters, other encodings include the signed value inside the request
		if authRequestForm.Encoding != EncodingDeflate {
			authRequestForm.SigAlg = authNRequest.Signature.SignedInfo.SignatureMethod.Algorithm
			authRequestForm.Sig = authNRequest.Signature.SignatureValue.Text
		}

		if err := sp.verifySignature(
			authRequestForm.AuthRequest,
			authRequestForm.RelayState,
			authRequestForm.SigAlg,
			authRequestForm.Sig,
		); err != nil {
			logging.Log("SAML-817n2s").Error(err)
			if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeDeniedResponse(
				authNRequest.Id,
				authNRequest.AssertionConsumerServiceURL,
				p.EntityID,
				fmt.Errorf("failed to verify signature: %w", err).Error(),
			)); err != nil {
				http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	if err := p.verifyRequestDestination(authNRequest); err != nil {
		logging.Log("SAML-83722s").Error(err)
		if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeDeniedResponse(
			authNRequest.Id,
			authNRequest.AssertionConsumerServiceURL,
			p.EntityID,
			fmt.Errorf("failed to verify request destination: %w", err).Error(),
		)); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
		return
	}

	acsURL := ""
	for _, acs := range sp.metadata.SPSSODescriptor.AssertionConsumerService {
		if acs.Binding == authNRequest.ProtocolBinding {
			acsURL = acs.Location
			break
		}
	}
	if acsURL == "" {
		logging.Log("SAML-83711s").Error(err)
		if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeUnsupportedBindingResponse(
			authNRequest.Id,
			authNRequest.AssertionConsumerServiceURL,
			p.EntityID,
			fmt.Errorf("unsupported binding").Error(),
		)); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := verifyRequestContent(
		authNRequest,
		string(sp.metadata.EntityID),
		acsURL,
	); err != nil {
		logging.Log("SAML-8kj22s").Error(err)
		if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeDeniedResponse(
			authNRequest.Id,
			authNRequest.AssertionConsumerServiceURL,
			p.EntityID,
			fmt.Errorf("failed to verify request content: %w", err).Error(),
		)); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
		return
	}

	authRequest, err := p.storage.CreateAuthRequest(r.Context(), authNRequest, authRequestForm.RelayState, sp.ID)
	if err != nil {
		logging.Log("SAML-8opi22s").Error(err)
		if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeResponderFailResponse(
			authNRequest.Id,
			authNRequest.AssertionConsumerServiceURL,
			p.EntityID,
			fmt.Errorf("failed to persist request %w", err).Error(),
		)); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
		return
	}

	switch authNRequest.ProtocolBinding {
	case "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect":
		http.Redirect(w, r, sp.LoginURL(authRequest.GetID()), http.StatusTemporaryRedirect)
	case "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST":
		http.Redirect(w, r, sp.LoginURL(authRequest.GetID()), http.StatusTemporaryRedirect)
	default:
		logging.Log("SAML-67722s").Error(err)
		if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeUnsupportedBindingResponse(
			authNRequest.Id,
			authNRequest.AssertionConsumerServiceURL,
			p.EntityID,
			fmt.Errorf("unsupported binding: %s", authNRequest.ProtocolBinding).Error(),
		)); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
	}
	return
}

func getAuthRequestFromRequest(r *http.Request) (*AuthRequestForm, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("failed to parse form: %w", err)
	}

	request := &AuthRequestForm{
		AuthRequest: r.Form.Get("SAMLRequest"),
		Encoding:    r.Form.Get("SAMLEncoding"),
		RelayState:  r.Form.Get("RelayState"),
		SigAlg:      r.Form.Get("SigAlg"),
		Sig:         r.Form.Get("Signature"),
	}

	return request, nil
}

func verifyForm(r *AuthRequestForm) error {
	if r.AuthRequest == "" {
		return fmt.Errorf("empty SAMLRequest")
	}

	/*if r.Encoding == "" {
		r.Encoding = EncodingDeflate
	}*/

	if r.RelayState == "" {
		return fmt.Errorf("empty RelayState")
	}
	//should be 80, but google / SNOW implement it wrong
	if len(r.RelayState) > 300 {
		return fmt.Errorf("relaystate should not be longer than 300")
	}

	if r.SigAlg != "" {
		if r.Sig == "" {
			return fmt.Errorf("empty Signature")
		}
		return fmt.Errorf("signature algorithm is empty")
	}

	return nil
}

func decodeAuthNRequest(encoding string, message string) (*samlp.AuthnRequest, error) {
	reqBytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return nil, err
	}

	req := &samlp.AuthnRequest{}
	switch encoding {
	case EncodingDeflate:
		reader := flate.NewReader(bytes.NewReader(reqBytes))
		decoder := xml.NewDecoder(reader)
		if err = decoder.Decode(req); err != nil {
			return nil, err
		}
	default:
		if err := xml.Unmarshal(reqBytes, req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

func verifyRequestContent(request *samlp.AuthnRequest, entityID, acsURL string) error {
	if request.Id == "" {
		return fmt.Errorf("request with empty id")
	}

	if request.Version == "" {
		return fmt.Errorf("request with empty version")
	}

	if request.Issuer.Text == "" {
		return fmt.Errorf("request with empty issuer")
	}

	if request.Issuer.Text != entityID {
		return fmt.Errorf("request with unknown issuer")
	}

	if request.AssertionConsumerServiceURL == "" {
		return fmt.Errorf("request with empty assertionConsumerServiceURL")
	}

	if request.AssertionConsumerServiceURL != acsURL {
		return fmt.Errorf("request with unknown assertionConsumerServiceURL")
	}

	return nil
}
