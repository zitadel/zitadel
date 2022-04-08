package saml

import (
	"encoding/base64"
	"fmt"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/saml/checker"
	"github.com/caos/zitadel/internal/api/saml/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/md"
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
		certificateCheckNecessary(
			func() *xml_dsig.SignatureType { return authNRequest.Signature },
			func() *md.EntityDescriptorType { return sp.metadata },
		),
		checkCertificate(
			func() *xml_dsig.SignatureType { return authNRequest.Signature },
			func() *md.EntityDescriptorType { return sp.metadata },
		),
		"SAML-b17d9a",
		func() {
			response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("failed to validate certificate from request: %w", err).Error()))
		},
	)

	/*
		// get signature out of request if POST-binding
		checker.WithConditionalLogicStep(
			signatureGetNecessary(
				func() *md.IDPSSODescriptorType { return p.Metadata },
				func() *md.EntityDescriptorType { return sp.metadata },
				func() string { return authRequestForm.Sig },
				func() *xml_dsig.SignatureType { return authNRequest.Signature },
				func() string { return authNRequest.ProtocolBinding },
			),
			getSignatureFromAuthRequest(
				func() *xml_dsig.SignatureType { return authNRequest.Signature },
				func() string { return authRequestForm.AuthRequest },
				func(request string) { authRequestForm.AuthRequest = request },
				func(sig string) { authRequestForm.Sig = sig },
				func(sigAlg string) { authRequestForm.SigAlg = sigAlg },
				func(errF error) { err = errF },
			),
			"SAML-i1o2mh",
			func() {
				response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("failed to extract signature from request: %w", err).Error()))
			},
		)*/

	// verify signature if necessary
	checker.WithConditionalLogicStep(
		signatureRedirectVerificationNecessary(
			func() *md.IDPSSODescriptorType { return p.Metadata },
			func() *md.EntityDescriptorType { return sp.metadata },
			func() string { return authRequestForm.Sig },
			func() string { return authNRequest.ProtocolBinding },
		),
		verifyRedirectSignature(
			func() string { return authRequestForm.AuthRequest },
			func() string { return authRequestForm.RelayState },
			func() string { return authRequestForm.Sig },
			func() string { return authRequestForm.SigAlg },
			func() *ServiceProvider { return sp },
			func(errF error) { err = errF },
		),
		"SAML-817n2s",
		func() {
			response.sendBackResponse(r, w, response.makeDeniedResponse(fmt.Errorf("failed to verify signature: %w", err).Error()))
		},
	)

	// verify signature if necessary
	checker.WithConditionalLogicStep(
		signaturePostVerificationNecessary(
			func() *md.IDPSSODescriptorType { return p.Metadata },
			func() *md.EntityDescriptorType { return sp.metadata },
			func() *xml_dsig.SignatureType { return authNRequest.Signature },
			func() string { return authNRequest.ProtocolBinding },
		),
		verifyPostSignature(
			func() string { return authRequestForm.AuthRequest },
			func() *ServiceProvider { return sp },
			func(errF error) { err = errF },
		),
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

	regexSignValue := regexp.MustCompile(`(<)(.?)(.?)(:?)(SignatureValue)(.|\n|\t|\r|\f)*(</)(.?)(.?)(:?)(SignatureValue>)`)
	authRequestWithoutSignValue := regexSignValue.ReplaceAll(reqBytes, []byte(""))

	regexKeyInfo := regexp.MustCompile(`(<)(.?)(.?)(:?)(KeyInfo)(.|\n|\t|\r|\f)*(</)(.?)(.?)(:?)(KeyInfo>)`)
	authRequest := regexKeyInfo.ReplaceAll(authRequestWithoutSignValue, []byte(""))

	return base64.StdEncoding.EncodeToString(authRequest), nil

}

func certificateCheckNecessary(
	authRequestSignatureF func() *xml_dsig.SignatureType,
	spMetadataF func() *md.EntityDescriptorType,
) func() bool {
	return func() bool {
		sig := authRequestSignatureF()
		spMetadata := spMetadataF()
		return sig != nil && sig.KeyInfo != nil &&
			spMetadata.SPSSODescriptor.KeyDescriptor != nil && len(spMetadata.SPSSODescriptor.KeyDescriptor) > 0
	}
}

func checkCertificate(
	authRequestSignatureF func() *xml_dsig.SignatureType,
	spMetadataF func() *md.EntityDescriptorType,
) func() error {
	return func() error {
		for _, keyDesc := range spMetadataF().SPSSODescriptor.KeyDescriptor {
			for _, spX509Data := range keyDesc.KeyInfo.X509Data {
				for _, reqX509Data := range authRequestSignatureF().KeyInfo.X509Data {
					if spX509Data.X509Certificate == reqX509Data.X509Certificate {
						return nil
					}
				}

			}
		}
		return fmt.Errorf("unknown certificate used to sign request")
	}
}

func signatureRedirectVerificationNecessary(
	idpMetadataF func() *md.IDPSSODescriptorType,
	spMetadataF func() *md.EntityDescriptorType,
	signatureF func() string,
	protocolBinding func() string,
) func() bool {
	return func() bool {

		return (spMetadataF().SPSSODescriptor.AuthnRequestsSigned == "true" ||
			idpMetadataF().WantAuthnRequestsSigned == "true" ||
			signatureF() != "") &&
			protocolBinding() == "urn:oasis:names:ts:SAML:2.0:bindings:HTTP-Redirect"
	}
}

func verifyRedirectSignature(
	authRequest func() string,
	relayState func() string,
	sig func() string,
	sigAlg func() string,
	sp func() *ServiceProvider,
	errF func(error),
) func() error {
	return func() error {
		if authRequest() == "" {
			return fmt.Errorf("no authrequest provided but required")
		}
		if relayState() == "" {
			return fmt.Errorf("no relaystate provided but required")
		}
		if sig() == "" {
			return fmt.Errorf("no signature provided but required")
		}
		if sigAlg() == "" {
			return fmt.Errorf("no signature algorithm provided but required")
		}

		spInstance := sp()
		err := spInstance.verifyRedirectSignature(
			authRequest(),
			relayState(),
			sigAlg(),
			sig(),
		)
		errF(err)
		return err
	}
}

func signaturePostVerificationNecessary(
	idpMetadataF func() *md.IDPSSODescriptorType,
	spMetadataF func() *md.EntityDescriptorType,
	authRequestSignatureF func() *xml_dsig.SignatureType,
	protocolBinding func() string,
) func() bool {
	return func() bool {
		authRequestSignature := authRequestSignatureF()

		return spMetadataF().SPSSODescriptor.AuthnRequestsSigned == "true" ||
			idpMetadataF().WantAuthnRequestsSigned == "true" ||
			(authRequestSignature != nil &&
				!reflect.DeepEqual(authRequestSignature.SignatureValue, xml_dsig.SignatureValueType{}) &&
				authRequestSignature.SignatureValue.Text != "") &&
				protocolBinding() == "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
	}
}

func verifyPostSignature(
	authRequestF func() string,
	spF func() *ServiceProvider,
	errF func(error),
) func() error {
	return func() error {
		sp := spF()

		data, err := base64.StdEncoding.DecodeString(authRequestF())
		if err != nil {
			errF(err)
			return err
		}

		if err := sp.verifyPostSignature(string(data)); err != nil {
			errF(err)
			return err
		}
		return nil
	}
}
