package saml

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/zitadel/saml/pkg/provider"
	"github.com/zitadel/saml/pkg/provider/models"
	"github.com/zitadel/saml/pkg/provider/xml"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
)

func (p *Provider) CreateErrorResponse(authReq models.AuthRequestInt, reason, description string) (string, string, error) {
	resp := &provider.Response{
		ProtocolBinding: authReq.GetBindingType(),
		RelayState:      authReq.GetRelayState(),
		AcsUrl:          authReq.GetAccessConsumerServiceURL(),
		RequestID:       authReq.GetAuthRequestID(),
		Issuer:          authReq.GetDestination(),
		Audience:        authReq.GetIssuer(),
	}
	samlResponse := p.AuthCallbackErrorResponse(resp, reason, description)

	respData, err := xml.Marshal(samlResponse)
	if err != nil {
		return "", "", err
	}

	switch authReq.GetBindingType() {
	case provider.PostBinding:
		return authReq.GetAccessConsumerServiceURL(), base64.StdEncoding.EncodeToString(respData), nil
	case provider.RedirectBinding:
		respData, err := xml.DeflateAndBase64(respData)
		if err != nil {
			return "", "", err
		}
		return fmt.Sprintf("%s?%s", authReq.GetAccessConsumerServiceURL(), provider.BuildRedirectQuery(string(respData), resp.RelayState, resp.SigAlg, resp.Signature)), "", nil
	}
	return "", "", nil
}

func (p *Provider) CreateResponse(ctx context.Context, authReq models.AuthRequestInt) (string, string, error) {
	resp := &provider.Response{
		ProtocolBinding: authReq.GetBindingType(),
		RelayState:      authReq.GetRelayState(),
		AcsUrl:          authReq.GetAccessConsumerServiceURL(),
		RequestID:       authReq.GetAuthRequestID(),
		Issuer:          authReq.GetDestination(),
		Audience:        authReq.GetIssuer(),
	}
	samlResponse, err := p.AuthCallbackResponse(ctx, authReq, resp)
	if err != nil {
		return "", "", err
	}

	if err := p.command.CreateSAMLSessionFromSAMLRequest(
		setContextUserSystem(ctx),
		authReq.GetID(),
		samlComplianceChecker(),
		samlResponse.Id,
		p.Expiration(),
	); err != nil {
		return "", "", err
	}

	respData, err := xml.Marshal(samlResponse)
	if err != nil {
		return "", "", err
	}

	switch authReq.GetBindingType() {
	case provider.PostBinding:
		return authReq.GetAccessConsumerServiceURL(), base64.StdEncoding.EncodeToString(respData), nil
	case provider.RedirectBinding:
		respData, err := xml.DeflateAndBase64(respData)
		if err != nil {
			return "", "", err
		}
		return fmt.Sprintf("%s?%s", authReq.GetAccessConsumerServiceURL(), provider.BuildRedirectQuery(string(respData), resp.RelayState, resp.SigAlg, resp.Signature)), "", nil
	}
	return "", "", nil
}

func setContextUserSystem(ctx context.Context) context.Context {
	data := authz.CtxData{
		UserID: "SYSTEM",
	}
	return authz.SetCtxData(ctx, data)
}

func samlComplianceChecker() command.SAMLRequestComplianceChecker {
	return func(_ context.Context, samlReq *command.SAMLRequestWriteModel) error {
		if err := samlReq.CheckAuthenticated(); err != nil {
			return err
		}
		return nil
	}
}
