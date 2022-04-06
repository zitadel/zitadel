package saml

import (
	"encoding/base64"
	"fmt"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/saml/checker"
	"github.com/caos/zitadel/internal/api/saml/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/samlp"
	"github.com/caos/zitadel/internal/api/saml/xml/xml_dsig"
	"net/http"
	"reflect"
	"regexp"
)

type AuthRequestForm struct {
	AuthRequest string
	Encoding    string
	RelayState  string
	SigAlg      string
	Sig         string
}

func (p *IdentityProvider) ssoHandleFunc(w http.ResponseWriter, r *http.Request) {
	checker := checker.Checker{}
	var authRequestForm *AuthRequestForm
	var authNRequest *samlp.AuthnRequestType
	var sp *ServiceProvider
	var authRequest AuthRequestInt
	var err error

	response := &Response{
		PostTemplate: p.postTemplate,
		ErrorFunc: func(err error) {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		},
		Issuer: p.EntityID,
	}

	// parse form to cover POST and REDIRECT binding
	checker.WithLogicStep(
		func() error {
			authRequestForm, err = getAuthRequestFromRequest(r)
			if err != nil {
				return err
			}
			response.SigAlg = authRequestForm.SigAlg
			response.RelayState = authRequestForm.RelayState
			return nil
		},
		"SAML-837n2s",
		func() {
			http.Error(w, fmt.Errorf("failed to parse form: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	// verify that relayState is provided
	checker.WithValueNotEmptyCheck(
		"relayState",
		func() string { return authRequestForm.RelayState },
		"SAML-86272s",
		func() {
			response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("empty relaystate").Error()))
		},
	)

	// verify that request is not empty
	checker.WithValueNotEmptyCheck(
		"SAMLRequest",
		func() string { return authRequestForm.AuthRequest },
		"SAML-nu32kq",
		func() {
			response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("no auth request provided").Error()))
		},
	)

	// verify that there is a signature provided if signature algorithm is provided
	checker.WithConditionalValueNotEmpty(
		func() bool { return authRequestForm.SigAlg != "" },
		"Signature",
		func() string { return authRequestForm.Sig },
		"SAML-827n2s",
		func() {
			response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("signature algorith provided but no signature").Error()))
		},
	)

	// decode request from xml into golang struct
	checker.WithLogicStep(
		func() error {
			authNRequest, err = xml.DecodeAuthNRequest(authRequestForm.Encoding, authRequestForm.AuthRequest)
			if err != nil {
				return err
			}
			response.RequestID = authNRequest.Id
			return nil
		},
		"SAML-837s2s",
		func() {
			response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("failed to decode request").Error()))
		},
	)

	// get persistet service provider from issuer out of the request
	checker.WithLogicStep(
		func() error {
			sp, err = p.GetServiceProvider(r.Context(), authNRequest.Issuer.Text)
			if err != nil {
				return err
			}
			response.Audience = sp.GetEntityID()
			return nil
		},
		" SAML-317s2s",
		func() {
			response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("failed to find registered serviceprovider: %w", err).Error()))
		},
	)

	//validate used certificate for signing the request
	checker.WithConditionalLogicStep(
		func() bool {
			return authNRequest.Signature != nil && authNRequest.Signature.KeyInfo != nil &&
				sp.metadata.SPSSODescriptor.KeyDescriptor != nil && len(sp.metadata.SPSSODescriptor.KeyDescriptor) > 0
		},
		func() error {
			for _, keyDesc := range sp.metadata.SPSSODescriptor.KeyDescriptor {
				for _, spX509Data := range keyDesc.KeyInfo.X509Data {
					for _, reqX509Data := range authNRequest.Signature.KeyInfo.X509Data {
						if spX509Data.X509Certificate == reqX509Data.X509Certificate {
							return nil
						}
					}

				}
			}
			return fmt.Errorf("unknown certificate used to sign request")
		},
		"SAML-b17d9a",
		func() {
			response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("failed to validate certificate from request: %w", err).Error()))
		},
	)

	// get signature out of request if POST-binding
	checker.WithConditionalLogicStep(
		func() bool {
			return sp.metadata.SPSSODescriptor.AuthnRequestsSigned == "true" ||
				p.Metadata.WantAuthnRequestsSigned == "true" ||
				authRequestForm.Sig != "" ||
				(authNRequest.Signature != nil &&
					!reflect.DeepEqual(authNRequest.Signature.SignatureValue, xml_dsig.SignatureValueType{}) &&
					authNRequest.Signature.SignatureValue.Text != "") &&
					authNRequest.ProtocolBinding == "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
		},
		func() error {
			if authNRequest.Signature != nil &&
				!reflect.DeepEqual(authNRequest.Signature.SignedInfo, xml_dsig.SignedInfoType{}) &&
				!reflect.DeepEqual(authNRequest.Signature.SignedInfo.SignatureMethod, xml_dsig.SignatureMethodType{}) &&
				authNRequest.Signature.SignedInfo.SignatureMethod.Algorithm == "" {
				authRequestForm.SigAlg = authNRequest.Signature.SignedInfo.SignatureMethod.Algorithm
			}

			if authNRequest.Signature != nil &&
				!reflect.DeepEqual(authNRequest.Signature.SignatureValue, xml_dsig.SignatureValueType{}) &&
				authNRequest.Signature.SignatureValue.Text == "" {
				authRequestForm.Sig = authNRequest.Signature.SignatureValue.Text
			}

			authRequestForm.AuthRequest, err = authNRequestIntoStringWithoutSignature(authRequestForm.AuthRequest)
			if err != nil {
				return err
			}
			return nil
		},
		"SAML-i1o2mh",
		func() {
			response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("failed to extract signature from request: %w", err).Error()))
		},
	)

	// verify signature if necessary
	checker.WithConditionalLogicStep(
		func() bool {
			return sp.metadata.SPSSODescriptor.AuthnRequestsSigned == "true" ||
				p.Metadata.WantAuthnRequestsSigned == "true" ||
				authRequestForm.Sig != "" ||
				(authNRequest.Signature != nil &&
					!reflect.DeepEqual(authNRequest.Signature.SignatureValue, xml_dsig.SignatureValueType{}) &&
					authNRequest.Signature.SignatureValue.Text != "")
		},
		func() error {
			if authRequestForm.AuthRequest == "" {
				return fmt.Errorf("no authrequest provided but required")
			}
			if authRequestForm.RelayState == "" {
				return fmt.Errorf("no relaystate provided but required")
			}
			if authRequestForm.Sig == "" {
				return fmt.Errorf("no signature provided but required")
			}
			if authRequestForm.SigAlg == "" {
				return fmt.Errorf("no signature algorithm provided but required")
			}

			err = sp.verifySignature(
				authRequestForm.AuthRequest,
				authRequestForm.RelayState,
				authRequestForm.SigAlg,
				authRequestForm.Sig,
			)
			return err
		},
		"SAML-817n2s",
		func() {
			response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("failed to verify signature: %w", err).Error()))
		},
	)

	// verify that destination in request is this IDP
	checker.WithLogicStep(
		func() error { err = p.verifyRequestDestinationOfAuthRequest(authNRequest); return err },
		"SAML-83722s",
		func() {
			response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("failed to verify request destination: %w", err).Error()))
		},
	)

	// work out used acs url and protocolbinding for response
	checker.WithValueStep(
		func() {
			for _, acs := range sp.metadata.SPSSODescriptor.AssertionConsumerService {
				if acs.Binding == authNRequest.ProtocolBinding {
					response.AcsUrl = acs.Location
					response.ProtocolBinding = acs.Binding
					break
				}
			}
			if response.AcsUrl == "" {
				for _, acs := range sp.metadata.SPSSODescriptor.AssertionConsumerService {
					response.AcsUrl = acs.Location
					response.ProtocolBinding = acs.Binding
					break
				}
			}
		},
	)

	// check if supported acs url is provided
	checker.WithValueNotEmptyCheck(
		"acsUrl",
		func() string { return response.AcsUrl },
		"SAML-83712s",
		func() {
			response.sendBackResponse(r, w, response.makeUnsupportedBindingResponse(fmt.Errorf("missing usable assertion consumer url").Error()))
		},
	)

	// check if supported protocolbinding is provided
	checker.WithValueNotEmptyCheck(
		"protocol binding",
		func() string { return response.ProtocolBinding },
		"SAML-83711s",
		func() {
			response.sendBackResponse(r, w, response.makeUnsupportedBindingResponse(fmt.Errorf("missing usable protocol binding").Error()))
		},
	)

	//check if authrequest has required attributes
	checker.WithValuesNotEmptyCheck(
		func() []string { return []string{authNRequest.Id, authNRequest.Version, authNRequest.Issuer.Text} },
		"SAML-8kj22s",
		func() {
			response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("request is missing requiered attributes").Error()))
		},
	)

	//check if entityId used in the request and serviceprovider is equal
	checker.WithValueEqualsCheck(
		"entityID",
		func() string { return authNRequest.Issuer.Text },
		func() string { return string(sp.metadata.EntityID) },
		"SAML-7qj22s",
		func() {
			response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("provided issuer is not equal to known service provider").Error()))
		},
	)

	// persist authrequest
	checker.WithLogicStep(
		func() error {
			authRequest, err = p.storage.CreateAuthRequest(
				r.Context(),
				authNRequest,
				response.AcsUrl,
				response.ProtocolBinding,
				authRequestForm.RelayState,
				sp.ID,
			)
			return err
		},
		"SAML-8opi22s",
		func() {
			response.sendBackResponse(r, w, response.makeResponderFailResponse(fmt.Errorf("failed to persist request: %w", err).Error()))
		},
	)

	//check and log errors if necessary
	if checker.CheckFailed() {
		return
	}

	switch response.ProtocolBinding {
	case "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect", "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST":
		http.Redirect(w, r, sp.LoginURL(authRequest.GetID()), http.StatusSeeOther)
	default:
		logging.Log("SAML-67722s").Error(err)
		response.sendBackResponse(r, w, response.makeUnsupportedBindingResponse(fmt.Errorf("unsupported binding: %s", response.ProtocolBinding).Error()))
	}
	return
}

func getAuthRequestFromRequest(r *http.Request) (*AuthRequestForm, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("failed to parse form: %w", err)
	}

	request := &AuthRequestForm{
		AuthRequest: r.FormValue("SAMLRequest"),
		Encoding:    r.FormValue("SAMLEncoding"),
		RelayState:  r.FormValue("RelayState"),
		SigAlg:      r.FormValue("SigAlg"),
		Sig:         r.FormValue("Signature"),
	}

	return request, nil
}

func authNRequestIntoStringWithoutSignature(message string) (string, error) {
	reqBytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode: %w", err)
	}

	regex := regexp.MustCompile(`(<)(.?)(.?)(:?)(Signature)(.|\n|\t|\r|\f)*(</)(.?)(.?)(:?)(Signature>)`)
	authRequest := regex.ReplaceAll(reqBytes, []byte(""))

	return base64.StdEncoding.EncodeToString(authRequest), nil

}
