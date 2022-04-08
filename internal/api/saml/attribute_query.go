package saml

import (
	"encoding/base64"
	"fmt"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/saml/checker"
	"github.com/caos/zitadel/internal/api/saml/signature"
	"github.com/caos/zitadel/internal/api/saml/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/saml"
	"github.com/caos/zitadel/internal/api/saml/xml/samlp"
	"github.com/caos/zitadel/internal/api/saml/xml/soap"
	"github.com/caos/zitadel/internal/api/saml/xml/xml_dsig"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
)

func (p *IdentityProvider) attributeQueryHandleFunc(w http.ResponseWriter, r *http.Request) {
	checker := checker.Checker{}
	var attrQueryRequest string
	var err error
	var sp *ServiceProvider
	var attrQuery *samlp.AttributeQueryType
	var sigAlg, sig string
	var response *soap.ResponseEnvelope

	//parse body to string
	checker.WithLogicStep(
		func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return err
			}
			attrQueryRequest = string(b)
			return nil
		},
		"SAML-ap2j3n1",
		func() {
			http.Error(w, fmt.Errorf("failed to parse body: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	// decode request from xml into golang struct
	checker.WithLogicStep(
		func() error {
			attrQuery, err = xml.DecodeAttributeQuery(attrQueryRequest)
			if err != nil {
				return err
			}
			return nil
		},
		"SAML-qpoin2a",
		func() {
			http.Error(w, fmt.Errorf("failed to decode request: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	// get persistet service provider from issuer out of the request
	checker.WithLogicStep(
		func() error {
			sp, err = p.GetServiceProvider(r.Context(), attrQuery.Issuer.Text)
			if err != nil {
				return err
			}
			return nil
		},
		" SAML-asdi1n",
		func() {
			http.Error(w, fmt.Errorf("failed to find registered serviceprovider: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	//validate used certificate for signing the request
	checker.WithConditionalLogicStep(
		func() bool {
			return attrQuery.Signature != nil && attrQuery.Signature.KeyInfo != nil &&
				sp.metadata.SPSSODescriptor.KeyDescriptor != nil && len(sp.metadata.SPSSODescriptor.KeyDescriptor) > 0
		},
		func() error {
			for _, keyDesc := range sp.metadata.SPSSODescriptor.KeyDescriptor {
				for _, spX509Data := range keyDesc.KeyInfo.X509Data {
					for _, reqX509Data := range attrQuery.Signature.KeyInfo.X509Data {
						if spX509Data.X509Certificate == reqX509Data.X509Certificate {
							return nil
						}
					}

				}
			}
			return fmt.Errorf("unknown certificate used to sign request")
		},
		"SAML-bxi3n5",
		func() {
			http.Error(w, fmt.Errorf("failed to validate certificate from request: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	// get signature out of request if POST-binding
	checker.WithConditionalLogicStep(
		func() bool {
			return attrQuery.Signature != nil &&
				!reflect.DeepEqual(attrQuery.Signature.SignatureValue, xml_dsig.SignatureValueType{}) &&
				attrQuery.Signature.SignatureValue.Text != ""
		},
		func() error {
			if attrQuery.Signature != nil &&
				!reflect.DeepEqual(attrQuery.Signature.SignedInfo, xml_dsig.SignedInfoType{}) &&
				!reflect.DeepEqual(attrQuery.Signature.SignedInfo.SignatureMethod, xml_dsig.SignatureMethodType{}) &&
				attrQuery.Signature.SignedInfo.SignatureMethod.Algorithm == "" {
				sigAlg = attrQuery.Signature.SignedInfo.SignatureMethod.Algorithm
			}

			if attrQuery.Signature != nil &&
				!reflect.DeepEqual(attrQuery.Signature.SignatureValue, xml_dsig.SignatureValueType{}) &&
				attrQuery.Signature.SignatureValue.Text == "" {
				sig = attrQuery.Signature.SignatureValue.Text
			}

			attrQueryRequest, err = attrRequestIntoStringWithoutSignature(attrQueryRequest)
			if err != nil {
				return err
			}
			return nil
		},
		"SAML-ao1n2ps",
		func() {
			http.Error(w, fmt.Errorf("failed to extract signature from request: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	// verify signature if necessary
	checker.WithLogicStep(
		func() error {
			if attrQueryRequest == "" {
				return fmt.Errorf("no authrequest provided but required")
			}
			if sig == "" {
				return fmt.Errorf("no signature provided but required")
			}
			if sigAlg == "" {
				return fmt.Errorf("no signature algorithm provided but required")
			}

			err = sp.verifyRedirectSignature(
				attrQueryRequest,
				"",
				sigAlg,
				sig,
			)
			return err
		},
		"SAML-owm2b4",
		func() {
			http.Error(w, fmt.Errorf("failed to verify signature: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	// verify that destination in request is this IDP
	checker.WithLogicStep(
		func() error { err = p.verifyRequestDestinationOfAttrQuery(attrQuery); return err },
		"SAML-ap2n1a",
		func() {
			http.Error(w, fmt.Errorf("failed to verify request destination: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	attrs := &Attributes{}
	checker.WithLogicStep(
		func() error {
			if err := p.storage.SetUserinfoWithLoginName(r.Context(), attrs, attrQuery.Subject.NameID.Text, []int{}); err != nil {
				return err
			}

			queriedAttrs := []saml.AttributeType{}
			if attrQuery.Attribute != nil {
				for _, queriedAttr := range attrQuery.Attribute {
					queriedAttrs = append(queriedAttrs, queriedAttr)
				}
			}

			response = getAttributeQueryResponse(attrQuery.Id, p.EntityID, sp.GetEntityID(), attrs, queriedAttrs)
			return nil
		},
		"SAML-wosm22",
		func() {
			http.Error(w, fmt.Errorf("failed to get userinfo: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	checker.WithLogicStep(
		func() error {
			signature, err := signature.Create(p.signingContext, response.Body.Response)
			if err != nil {
				return err
			}
			response.Body.Response.Signature = signature
			return nil
		},
		"SAML-p012sa",
		func() {
			http.Error(w, fmt.Errorf("failed to sign response: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	//check and log errors if necessary
	if checker.CheckFailed() {
		return
	}

	if err := xml.WriteXML(w, response); err != nil {
		logging.Log("SAML-91j12bk").Error(err)
		http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
	}
}

func attrRequestIntoStringWithoutSignature(message string) (string, error) {
	reqBytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode: %w", err)
	}

	regex := regexp.MustCompile(`(<)(.?)(.?)(:?)(Signature)(.|\n|\t|\r|\f)*(</)(.?)(.?)(:?)(Signature>)`)
	request := regex.ReplaceAll(reqBytes, []byte(""))

	return base64.StdEncoding.EncodeToString(request), nil
}

func getAttributeQueryResponse(
	requestID string,
	issuer string,
	entityID string,
	attributes *Attributes,
	queriedAttrs []saml.AttributeType,
) *soap.ResponseEnvelope {
	resp := makeAttributeQueryResponse(requestID, issuer, entityID, attributes, queriedAttrs)

	return &soap.ResponseEnvelope{
		Body: soap.ResponseBody{
			Response: resp,
		},
	}
}
