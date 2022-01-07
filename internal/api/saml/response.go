package saml

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/saml"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/samlp"
	"net/http"
	"net/url"
	"text/template"
	"time"
)

const (
	DefaultTimeFormat                = "2006-01-02T15:04:05.999999Z"
	StatusCodeSuccess                = "urn:oasis:names:tc:SAML:2.0:status:Success"
	StatusCodeVersionMissmatch       = "urn:oasis:names:tc:SAML:2.0:status:VersionMismatch"
	StatusCodeAuthNFailed            = "urn:oasis:names:tc:SAML:2.0:status:AuthnFailed"
	StatusCodeInvalidAttrNameOrValue = "urn:oasis:names:tc:SAML:2.0:status:InvalidAttrNameOrValue"
	StatusCodeInvalidNameIDPolicy    = "urn:oasis:names:tc:SAML:2.0:status:InvalidNameIDPolicy"
	StatusCodeRequestDenied          = "urn:oasis:names:tc:SAML:2.0:status:RequestDenied"
	StatusCodeRequestUnsupported     = "urn:oasis:names:tc:SAML:2.0:status:RequestUnsupported"
	StatusCodeUnsupportedBinding     = "urn:oasis:names:tc:SAML:2.0:status:UnsupportedBinding"
	StatusCodeResponder              = "urn:oasis:names:tc:SAML:2.0:status:Responder"
)

type AuthResponseForm struct {
	RelayState                  string
	SAMLResponse                string
	AssertionConsumerServiceURL string
}

func sendBackResponse(template *template.Template, w http.ResponseWriter, relayState string, acsURL string, resp *samlp.Response) error {
	var xmlbuff bytes.Buffer

	memWriter := bufio.NewWriter(&xmlbuff)
	_, err := memWriter.Write([]byte(xml.Header))
	if err != nil {
		return err
	}

	encoder := xml.NewEncoder(memWriter)
	err = encoder.Encode(resp)
	if err != nil {
		return err
	}

	err = memWriter.Flush()
	if err != nil {
		return err
	}

	samlMessage := base64.StdEncoding.EncodeToString(xmlbuff.Bytes())

	data := AuthResponseForm{
		url.QueryEscape(relayState),
		samlMessage,
		acsURL,
	}

	return template.Execute(w, data)
}

func makeUnsupportedBindingResponse(
	requestID string,
	acsURL string,
	issuer string,
	message string,
) *samlp.Response {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)
	return makeResponse(
		requestID,
		acsURL,
		nowStr,
		StatusCodeUnsupportedBinding,
		message,
		&saml.NameIDType{
			Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
			Text:   issuer,
		},
	)
}

func makeResponderFailResponse(
	requestID string,
	acsURL string,
	issuer string,
	message string,
) *samlp.Response {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)
	return makeResponse(
		requestID,
		acsURL,
		nowStr,
		StatusCodeResponder,
		message,
		&saml.NameIDType{
			Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
			Text:   issuer,
		},
	)
}

func makeDeniedResponse(
	requestID string,
	acsURL string,
	issuer string,
	message string,
) *samlp.Response {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)
	return makeResponse(
		requestID,
		acsURL,
		nowStr,
		StatusCodeRequestDenied,
		message,
		&saml.NameIDType{
			Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
			Text:   issuer,
		},
	)
}

func makeFailedResponse(
	requestID string,
	acsURL string,
	issuer string,
	message string,
) *samlp.Response {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)
	return makeResponse(
		requestID,
		acsURL,
		nowStr,
		StatusCodeAuthNFailed,
		message,
		&saml.NameIDType{
			Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
			Text:   issuer,
		},
	)
}

func makeSuccessfulResponse(
	request AuthRequestInt,
	sendIP string,
	nameID string,
	attributes []*saml.AttributeType,
) *samlp.Response {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)
	fiveFromNowStr := now.Add(5 * time.Minute).Format(DefaultTimeFormat)
	issuer := &saml.NameIDType{
		Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
		Text:   request.GetDestination(),
	}

	return makeAssertionResponse(
		request.GetAuthRequestID(),
		request.GetAccessConsumerServiceURL(),
		sendIP,
		nowStr,
		fiveFromNowStr,
		issuer,
		nameID,
		attributes,
	)
}

func makeAssertionResponse(
	requestID string,
	acsURL string,
	sendIP string,
	issueInstant string,
	untilInstant string,
	issuer *saml.NameIDType,
	nameID string,
	attributes []*saml.AttributeType,
) *samlp.Response {
	response := makeResponse(requestID, acsURL, issueInstant, StatusCodeSuccess, "", issuer)
	assertion := makeAssertion(requestID, acsURL, sendIP, issueInstant, untilInstant, issuer, nameID, attributes)
	response.Assertion = *assertion
	return response
}

func makeAssertion(
	requestID string,
	sendIP string,
	acsURL string,
	issueInstant string,
	untilInstant string,
	issuer *saml.NameIDType,
	nameID string,
	attributes []*saml.AttributeType,
) *saml.Assertion {
	id := NewID()
	return &saml.Assertion{
		Version:      "2.0",
		Id:           id,
		IssueInstant: issueInstant,
		Issuer:       *issuer,
		Subject: &saml.SubjectType{
			NameID: &saml.NameIDType{
				Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
				Text:   nameID,
			},
			SubjectConfirmation: []saml.SubjectConfirmationType{
				{
					Method: "urn:oasis:names:tc:SAML:2.0:cm:sender-vouches",
				},
				{
					Method: "urn:oasis:names:tc:SAML:2.0:cm:bearer",
					SubjectConfirmationData: &saml.SubjectConfirmationDataType{
						Address:      sendIP,
						InResponseTo: requestID,
						Recipient:    acsURL,
						NotOnOrAfter: untilInstant,
					},
				},
			},
		},
		Conditions: &saml.ConditionsType{
			NotBefore:    issueInstant,
			NotOnOrAfter: untilInstant,
			AudienceRestriction: []saml.AudienceRestrictionType{
				{Audience: []string{issuer.Text}},
			},
		},
		AttributeStatement: []saml.AttributeStatementType{
			{Attribute: attributes},
		},
		AuthnStatement: []saml.AuthnStatementType{
			{
				AuthnInstant: issueInstant,
				SessionIndex: id,
				AuthnContext: saml.AuthnContextType{
					AuthnContextClassRef: "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport",
				},
			},
		},
	}
}

func makeResponse(
	requestID string,
	acsURL string,
	issueInstant string,
	status string,
	message string,
	issuer *saml.NameIDType,
) *samlp.Response {
	return &samlp.Response{
		Version:      "2.0",
		Id:           NewID(),
		IssueInstant: issueInstant,
		Status: samlp.StatusType{
			StatusCode: samlp.StatusCodeType{
				Value: status,
			},
			StatusMessage: message,
		},
		InResponseTo: requestID,
		Destination:  acsURL,
		Issuer:       issuer,
	}
}
