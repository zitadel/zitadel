package saml

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/caos/zitadel/internal/api/saml/xml/saml"
	"github.com/caos/zitadel/internal/api/saml/xml/samlp"
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
	StatusCodePartialLogout          = "urn:oasis:names:tc:SAML:2.0:status:PartialLogout"
)

type Response struct {
	PostTemplate    *template.Template
	ProtocolBinding string
	RelayState      string
	AcsUrl          string
	Signature       string
	SigAlg          string
	ErrorFunc       func(err error)

	RequestID string
	Issuer    string
	Audience  string
	SendIP    string
}

func (r *Response) doResponse(request *http.Request, w http.ResponseWriter, response string) {
	if r.AcsUrl == "" {
		if _, err := w.Write([]byte(response)); err != nil {
			r.ErrorFunc(err)
			return
		}
	}

	switch r.ProtocolBinding {
	case "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST":
		data := AuthResponseForm{
			url.QueryEscape(r.RelayState),
			response,
			r.AcsUrl,
		}

		if err := r.PostTemplate.Execute(w, data); err != nil {
			r.ErrorFunc(err)
			return
		}
	case "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect":
		redirectURL := fmt.Sprintf("%s?SAMLResponse=%s&RelayState= %s", r.AcsUrl, url.QueryEscape(response), url.QueryEscape(r.RelayState))
		if r.Signature != "" && r.SigAlg != "" {
			redirectURL = fmt.Sprintf("%s&Signature=%s&SigAlg=%s", redirectURL, url.QueryEscape(r.Signature), url.QueryEscape(r.SigAlg))
		}
		http.Redirect(w, request, redirectURL, http.StatusFound)
		return
	default:
		//TODO: no binding
	}
}

type AuthResponseForm struct {
	RelayState                  string
	SAMLResponse                string
	AssertionConsumerServiceURL string
}

func (r *Response) sendBackResponse(
	req *http.Request,
	w http.ResponseWriter,
	resp *samlp.ResponseType,
) {
	var xmlbuff bytes.Buffer

	memWriter := bufio.NewWriter(&xmlbuff)
	_, err := memWriter.Write([]byte(xml.Header))
	if err != nil {
		r.ErrorFunc(err)
		return
	}

	encoder := xml.NewEncoder(memWriter)
	err = encoder.Encode(resp)
	if err != nil {
		r.ErrorFunc(err)
		return
	}

	err = memWriter.Flush()
	if err != nil {
		r.ErrorFunc(err)
		return
	}

	samlMessage := base64.StdEncoding.EncodeToString(xmlbuff.Bytes())

	r.doResponse(req, w, samlMessage)
}

func (r *Response) makeUnsupportedBindingResponse(
	message string,
) *samlp.ResponseType {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)
	return makeResponse(
		r.RequestID,
		r.AcsUrl,
		nowStr,
		StatusCodeUnsupportedBinding,
		message,
		r.Issuer,
	)
}

func (r *Response) makeResponderFailResponse(
	message string,
) *samlp.ResponseType {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)
	return makeResponse(
		r.RequestID,
		r.AcsUrl,
		nowStr,
		StatusCodeResponder,
		message,
		r.Issuer,
	)
}

func (r *Response) makeDeniedResponse(
	message string,
) *samlp.ResponseType {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)
	return makeResponse(
		r.RequestID,
		r.AcsUrl,
		nowStr,
		StatusCodeRequestDenied,
		message,
		r.Issuer,
	)
}

func (r *Response) makeFailedResponse(
	message string,
) *samlp.ResponseType {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)
	return makeResponse(
		r.RequestID,
		r.AcsUrl,
		nowStr,
		StatusCodeAuthNFailed,
		message,
		r.Issuer,
	)
}

func (r *Response) makeSuccessfulResponse(
	attributes *Attributes,
) *samlp.ResponseType {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)
	fiveFromNowStr := now.Add(5 * time.Minute).Format(DefaultTimeFormat)

	return r.makeAssertionResponse(
		nowStr,
		fiveFromNowStr,
		attributes,
	)
}

func (r *Response) makeAssertionResponse(
	issueInstant string,
	untilInstant string,
	attributes *Attributes,
) *samlp.ResponseType {

	response := makeResponse(r.RequestID, r.AcsUrl, issueInstant, StatusCodeSuccess, "", r.Issuer)
	assertion := makeAssertion(r.RequestID, r.AcsUrl, r.SendIP, issueInstant, untilInstant, r.Issuer, attributes.GetNameID(), attributes.GetSAML(), r.Audience, true)
	response.Assertion = *assertion
	return response
}

func getIssuer(entityID string) *saml.NameIDType {
	return &saml.NameIDType{
		Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
		Text:   entityID,
	}
}

func makeAttributeQueryResponse(
	requestID string,
	issuer string,
	entityID string,
	attributes *Attributes,
	queriedAttrs []saml.AttributeType,
) *samlp.ResponseType {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)
	fiveMinutes, _ := time.ParseDuration("5m")
	fiveFromNow := now.Add(fiveMinutes)
	fiveFromNowStr := fiveFromNow.Format(DefaultTimeFormat)

	providedAttrs := []*saml.AttributeType{}
	attrsSaml := attributes.GetSAML()
	if queriedAttrs == nil || len(queriedAttrs) == 0 {
		for _, attrSaml := range attrsSaml {
			providedAttrs = append(providedAttrs, attrSaml)
		}
	} else {
		for _, attrSaml := range attrsSaml {
			for _, queriedAttr := range queriedAttrs {
				if attrSaml.Name == queriedAttr.Name && attrSaml.NameFormat == queriedAttr.NameFormat {
					providedAttrs = append(providedAttrs, attrSaml)
				}
			}
		}
	}

	response := makeResponse(requestID, "", nowStr, StatusCodeSuccess, "", issuer)
	assertion := makeAssertion(requestID, "", "", nowStr, fiveFromNowStr, issuer, attributes.GetNameID(), providedAttrs, entityID, false)
	response.Assertion = *assertion
	return response
}

func makeAssertion(
	requestID string,
	acsURL string,
	sendIP string,
	issueInstant string,
	untilInstant string,
	issuer string,
	nameID *saml.NameIDType,
	attributes []*saml.AttributeType,
	audience string,
	authN bool,
) *saml.AssertionType {
	id := NewID()
	issuerP := getIssuer(issuer)

	ret := &saml.AssertionType{
		Version:      "2.0",
		Id:           id,
		IssueInstant: issueInstant,
		Issuer:       *issuerP,
		Subject: &saml.SubjectType{
			NameID: nameID,
			SubjectConfirmation: []saml.SubjectConfirmationType{
				{
					Method: "urn:oasis:names:tc:SAML:2.0:cm:bearer",
					SubjectConfirmationData: &saml.SubjectConfirmationDataType{
						InResponseTo: requestID,
						NotBefore:    issueInstant,
						NotOnOrAfter: untilInstant,
					},
				},
			},
		},
		Conditions: &saml.ConditionsType{
			NotBefore:    issueInstant,
			NotOnOrAfter: untilInstant,
			AudienceRestriction: []saml.AudienceRestrictionType{
				{Audience: []string{audience}},
			},
		},
		AttributeStatement: []saml.AttributeStatementType{
			{Attribute: attributes},
		},
	}
	if acsURL != "" {
		ret.Subject.SubjectConfirmation[0].SubjectConfirmationData.Recipient = acsURL
	}
	if sendIP != "" {
		ret.Subject.SubjectConfirmation[0].SubjectConfirmationData.Address = sendIP
	}
	if authN {
		ret.AuthnStatement = []saml.AuthnStatementType{
			{
				AuthnInstant: issueInstant,
				SessionIndex: id,
				AuthnContext: saml.AuthnContextType{
					AuthnContextClassRef: "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport",
				},
			},
		}
	}
	return ret
}

func makeResponse(
	requestID string,
	acsURL string,
	issueInstant string,
	status string,
	message string,
	issuer string,
) *samlp.ResponseType {
	resp := &samlp.ResponseType{
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
		Issuer:       getIssuer(issuer),
	}

	if acsURL != "" {
		resp.Destination = acsURL
	}
	return resp
}
