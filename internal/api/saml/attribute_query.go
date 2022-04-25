package saml

import (
	"fmt"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/saml/checker"
	"github.com/caos/zitadel/internal/api/saml/serviceprovider"
	"github.com/caos/zitadel/internal/api/saml/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/md"
	"github.com/caos/zitadel/internal/api/saml/xml/saml"
	"github.com/caos/zitadel/internal/api/saml/xml/samlp"
	"github.com/caos/zitadel/internal/api/saml/xml/soap"
	"github.com/caos/zitadel/internal/api/saml/xml/xml_dsig"
	"io/ioutil"
	"net/http"
)

func (p *IdentityProvider) attributeQueryHandleFunc(w http.ResponseWriter, r *http.Request) {
	checkerInstance := checker.Checker{}
	var attrQueryRequest string
	var err error
	var sp *serviceprovider.ServiceProvider
	var attrQuery *samlp.AttributeQueryType
	var response *samlp.ResponseType

	//parse body to string
	checkerInstance.WithLogicStep(
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
	checkerInstance.WithLogicStep(
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

	// get persisted service provider from issuer out of the request
	checkerInstance.WithLogicStep(
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
	checkerInstance.WithConditionalLogicStep(
		certificateCheckNecessary(
			func() *xml_dsig.SignatureType { return attrQuery.Signature },
			func() *md.EntityDescriptorType { return sp.Metadata },
		),
		checkCertificate(
			func() *xml_dsig.SignatureType { return attrQuery.Signature },
			func() *md.EntityDescriptorType { return sp.Metadata },
		),
		"SAML-bxi3n5",
		func() {
			http.Error(w, fmt.Errorf("failed to validate certificate from request: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	// get signature out of request if POST-binding
	checkerInstance.WithConditionalLogicStep(
		signaturePostProvided(
			func() *xml_dsig.SignatureType { return attrQuery.Signature },
		),
		verifyPostSignature(
			func() string { return attrQueryRequest },
			func() *serviceprovider.ServiceProvider { return sp },
			func(errF error) { err = errF },
		),
		"SAML-ao1n2ps",
		func() {
			http.Error(w, fmt.Errorf("failed to extract signature from request: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	// verify that destination in request is this IDP
	checkerInstance.WithLogicStep(
		func() error { err = p.verifyRequestDestinationOfAttrQuery(attrQuery); return err },
		"SAML-ap2n1a",
		func() {
			http.Error(w, fmt.Errorf("failed to verify request destination: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	// read userinfo and fill queried attributes into reponse
	attrs := &Attributes{}
	checkerInstance.WithLogicStep(
		func() error {
			if err := p.storage.SetUserinfoWithLoginName(r.Context(), attrs, attrQuery.Subject.NameID.Text, []int{}); err != nil {
				return err
			}

			queriedAttrs := make([]saml.AttributeType, 0)
			if attrQuery.Attribute != nil {
				for _, queriedAttr := range attrQuery.Attribute {
					queriedAttrs = append(queriedAttrs, queriedAttr)
				}
			}
			response = makeAttributeQueryResponse(attrQuery.Id, p.EntityID, sp.GetEntityID(), attrs, queriedAttrs)
			return nil
		},
		"SAML-wosm22",
		func() {
			http.Error(w, fmt.Errorf("failed to get userinfo: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	// create enveloped signature
	checkerInstance.WithLogicStep(
		func() error {
			return createPostSignature(response, p)
		},
		"SAML-p012sa",
		func() {
			http.Error(w, fmt.Errorf("failed to sign response: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	//check and log errors if necessary
	if checkerInstance.CheckFailed() {
		return
	}

	soapResponse := &soap.ResponseEnvelope{
		Body: soap.ResponseBody{
			Response: response,
		},
	}

	if err := xml.WriteXMLMarshalled(w, soapResponse); err != nil {
		logging.Log("SAML-91j12bk").Error(err)
		http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
	}
}
