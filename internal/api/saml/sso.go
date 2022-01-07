package saml

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/xml"
	"fmt"
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
		http.Error(w, fmt.Errorf("failed to parse form: %w", err).Error(), http.StatusInternalServerError)
	}

	if err := verifyForm(authRequestForm); err != nil {
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

	authNRequest, err := decodeAuthNRequest(authRequestForm.Encoding, authRequestForm.AuthRequest)
	if err != nil {
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

	sp := p.GetServiceProvider(authNRequest.Issuer.Text)
	if sp == nil {
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
			if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeDeniedResponse(
				authNRequest.Id,
				authNRequest.AssertionConsumerServiceURL,
				p.EntityID,
				fmt.Errorf("failed to verify signature: %w", err).Error(),
			)); err != nil {
				http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
			}
		}
	}

	if err := p.verifyRequestDestination(authNRequest); err != nil {
		if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeDeniedResponse(
			authNRequest.Id,
			authNRequest.AssertionConsumerServiceURL,
			p.EntityID,
			fmt.Errorf("failed to verify request destination: %w", err).Error(),
		)); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
	}

	if err := verifyRequestContent(
		authNRequest,
		p.EntityID,
		p.Metadata.SingleSignOnService[0].Location,
	); err != nil {
		if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeDeniedResponse(
			authNRequest.Id,
			authNRequest.AssertionConsumerServiceURL,
			p.EntityID,
			fmt.Errorf("failed to verify request content: %w", err).Error(),
		)); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
	}

	authRequest, err := p.storage.CreateAuthRequest(r.Context(), authNRequest, authRequestForm.RelayState, sp.ID)
	if err != nil {
		if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeResponderFailResponse(
			authNRequest.Id,
			authNRequest.AssertionConsumerServiceURL,
			p.EntityID,
			fmt.Errorf("failed to persist request %w", err).Error(),
		)); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
	}

	switch authNRequest.ProtocolBinding {
	case "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect":
		http.Redirect(w, r, p.GetRedirectURL(authRequest.GetID()), http.StatusTemporaryRedirect)
	case "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST":
		http.Redirect(w, r, p.GetRedirectURL(authRequest.GetID()), http.StatusTemporaryRedirect)
	default:
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

	if r.Encoding == "" {
		r.Encoding = EncodingDeflate
	}

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
		return nil, fmt.Errorf("unknown encoding")
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
