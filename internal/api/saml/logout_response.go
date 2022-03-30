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

type LogoutResponse struct {
	LogoutTemplate *template.Template
	RelayState     string
	SAMLResponse   string
	LogoutURL      string

	RequestID string
	Issuer    string
	ErrorFunc func(err error)
}

type LogoutResponseForm struct {
	RelayState   string
	SAMLResponse string
	LogoutURL    string
}

func (r *LogoutResponse) sendBackLogoutResponse(w http.ResponseWriter, resp *samlp.LogoutResponse) {
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

	data := LogoutResponseForm{
		RelayState:   url.QueryEscape(r.RelayState),
		SAMLResponse: samlMessage,
		LogoutURL:    r.LogoutURL,
	}

	if err := r.LogoutTemplate.Execute(w, data); err != nil {
		r.ErrorFunc(err)
		return
	}
}

func (r *LogoutResponse) makeSuccessfulLogoutResponse() *samlp.LogoutResponse {
	return makeLogoutResponse(
		r.RequestID,
		r.LogoutURL,
		time.Now().UTC().Format(DefaultTimeFormat),
		StatusCodeSuccess,
		"",
		getIssuer(r.Issuer),
	)
}

func (r *LogoutResponse) makeUnsupportedlLogoutResponse(
	message string,
) *samlp.LogoutResponse {
	return makeLogoutResponse(
		r.RequestID,
		r.LogoutURL,
		time.Now().UTC().Format(DefaultTimeFormat),
		StatusCodeRequestUnsupported,
		message,
		getIssuer(r.Issuer),
	)
}

func (r *LogoutResponse) makePartialLogoutResponse(
	message string,
) *samlp.LogoutResponse {
	return makeLogoutResponse(
		r.RequestID,
		r.LogoutURL,
		time.Now().UTC().Format(DefaultTimeFormat),
		StatusCodePartialLogout,
		message,
		getIssuer(r.Issuer),
	)
}

func (r *LogoutResponse) makeDeniedLogoutResponse(
	message string,
) *samlp.LogoutResponse {
	return makeLogoutResponse(
		r.RequestID,
		r.LogoutURL,
		time.Now().UTC().Format(DefaultTimeFormat),
		StatusCodeRequestDenied,
		message,
		getIssuer(r.Issuer),
	)
}

func makeLogoutResponse(
	requestID string,
	logoutURL string,
	issueInstant string,
	status string,
	message string,
	issuer *saml.NameIDType,
) *samlp.LogoutResponse {
	return &samlp.LogoutResponse{
		Id:           NewID(),
		InResponseTo: requestID,
		Version:      "2.0",
		IssueInstant: issueInstant,
		Destination:  logoutURL,
		Issuer:       issuer,
		Status: samlp.StatusType{
			StatusCode: samlp.StatusCodeType{
				Value: status,
			},
			StatusMessage: message,
		},
	}
}
