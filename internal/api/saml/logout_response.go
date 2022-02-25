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

type LogoutResponseForm struct {
	RelayState   string
	SAMLResponse string
	LogoutURL    string
}

func sendBackLogoutResponse(template *template.Template, w http.ResponseWriter, relayState string, logoutURL string, resp *samlp.LogoutResponse) error {
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

	data := LogoutResponseForm{
		RelayState:   url.QueryEscape(relayState),
		SAMLResponse: samlMessage,
		LogoutURL:    logoutURL,
	}

	return template.Execute(w, data)
}

func makeSuccessfulLogoutResponse(
	request *samlp.LogoutRequest,
	logoutURL string,
	entityID string,
) *samlp.LogoutResponse {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)

	issuer := &saml.NameIDType{
		Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
		Text:   entityID,
	}

	return makeLogoutResponse(
		request.Id,
		logoutURL,
		nowStr,
		StatusCodeSuccess,
		"",
		issuer,
	)
}

func makeUnsupportedlLogoutResponse(
	request *samlp.LogoutRequest,
	logoutURL string,
	entityID string,
	message string,
) *samlp.LogoutResponse {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)

	issuer := &saml.NameIDType{
		Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
		Text:   entityID,
	}

	return makeLogoutResponse(
		request.Id,
		logoutURL,
		nowStr,
		StatusCodeRequestUnsupported,
		message,
		issuer,
	)
}

func makePartialLogoutResponse(
	request *samlp.LogoutRequest,
	logoutURL string,
	entityID string,
	message string,
) *samlp.LogoutResponse {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)

	issuer := &saml.NameIDType{
		Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
		Text:   entityID,
	}

	return makeLogoutResponse(
		request.Id,
		logoutURL,
		nowStr,
		StatusCodePartialLogout,
		message,
		issuer,
	)
}

func makeDeniedLogoutResponse(
	request *samlp.LogoutRequest,
	logoutURL string,
	entityID string,
	message string,
) *samlp.LogoutResponse {
	now := time.Now().UTC()
	nowStr := now.Format(DefaultTimeFormat)

	issuer := &saml.NameIDType{
		Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
		Text:   entityID,
	}

	return makeLogoutResponse(
		request.Id,
		logoutURL,
		nowStr,
		StatusCodeRequestDenied,
		message,
		issuer,
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
