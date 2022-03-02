package saml

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/samlp"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/xml_dsig"
	"net/http"
	"regexp"
	"strings"
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
	authRequestForm.AuthRequest = "PHNhbWxwOkF1dGhuUmVxdWVzdCB4bWxuczpzYW1scD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIiAgICAJRGVzdGluYXRpb249Imh0dHBzOi8vYWNjb3VudHMub3Jib3MuaW8vc2FtbC9TU08iICAgIElEPSJfZmFiYmZmMzNkNDc3ZjAxZGY5Mzk2NTU1NmM2NTU0MDIiICAgIElzc3VlSW5zdGFudD0iMjAyMi0wMy0wMVQxNzo1MzozMFoiICAgIFByb3RvY29sQmluZGluZz0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmJpbmRpbmdzOkhUVFAtUE9TVCIgVmVyc2lvbj0iMi4wIj48c2FtbDpJc3N1ZXIgeG1sbnM6c2FtbD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmFzc2VydGlvbiI+dXJuOmF1dGgwOm9yYm9zLWlvOm9yYm9zPC9zYW1sOklzc3Vlcj48L3NhbWxwOkF1dGhuUmVxdWVzdD4="
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
		(authNRequest.Signature != nil && authNRequest.Signature.SignatureValue != xml_dsig.SignatureValueType{} && authNRequest.Signature.SignatureValue.Text != "") {

		switch authNRequest.ProtocolBinding {
		case "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST":
			authRequestForm.SigAlg = authNRequest.Signature.SignedInfo.SignatureMethod.Algorithm
			authRequestForm.Sig = authNRequest.Signature.SignatureValue.Text

			authRequestForm.AuthRequest, err = decodeAuthNRequestIntoStringWithoutSignature(authRequestForm.Encoding, authRequestForm.AuthRequest)
			if err != nil {
				logging.Log("SAML-i1o2mh").Error(err)
				logging.Log("SAML-817n2s").Error(err)
				if err := sendBackResponse(p.postTemplate, w, authRequestForm.RelayState, "", makeDeniedResponse(
					authNRequest.Id,
					authNRequest.AssertionConsumerServiceURL,
					p.EntityID,
					fmt.Errorf("failed to handle signature in request: %w", err).Error(),
				)); err != nil {
					http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
				}
				return
			}
		case "redirect":
			//do nothing? as everything should be in the form
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
		return nil, fmt.Errorf("failed to base64 decode: %w", err)
	}

	req := &samlp.AuthnRequest{}
	switch encoding {
	case EncodingDeflate:
		reader := flate.NewReader(bytes.NewReader(reqBytes))
		decoder := xml.NewDecoder(reader)
		if err = decoder.Decode(req); err != nil {
			return nil, fmt.Errorf("failed to defalte decode: %w", err)
		}
	default:
		if err := xml.Unmarshal(reqBytes, req); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %w", err)
		}
	}

	return req, nil
}

func decodeAuthNRequestIntoStringWithoutSignature(encoding string, message string) (string, error) {
	reqBytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode: %w", err)
	}

	regex := regexp.MustCompile(`(<)(.?)(.?)(:?)(Signature)(.|\n|\t|\r|\f)*(</)(.?)(.?)(:?)(Signature>)`)
	return regex.ReplaceAllString(string(reqBytes), ""), nil

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

	if !strings.Contains(entityID, "auth0.com") {
		if request.AssertionConsumerServiceURL == "" {
			return fmt.Errorf("request with empty assertionConsumerServiceURL")
		}

		if request.AssertionConsumerServiceURL != acsURL {
			return fmt.Errorf("request with unknown assertionConsumerServiceURL")
		}
	}

	return nil
}
