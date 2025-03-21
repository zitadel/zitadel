package saml

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/zitadel/saml/pkg/provider"
	"github.com/zitadel/saml/pkg/provider/models"
	"github.com/zitadel/saml/pkg/provider/xml"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
)

func (p *Provider) CreateErrorResponse(authReq models.AuthRequestInt, reason domain.SAMLErrorReason, description string) (string, string, error) {
	resp := &provider.Response{
		ProtocolBinding: authReq.GetBindingType(),
		RelayState:      authReq.GetRelayState(),
		AcsUrl:          authReq.GetAccessConsumerServiceURL(),
		RequestID:       authReq.GetAuthRequestID(),
		Issuer:          authReq.GetDestination(),
		Audience:        authReq.GetIssuer(),
	}
	return createResponse(p.AuthCallbackErrorResponse(resp, domain.SAMLErrorReasonToString(reason), description), authReq.GetBindingType(), authReq.GetAccessConsumerServiceURL(), resp.RelayState, resp.SigAlg, resp.Signature)
}

func (p *Provider) CreateResponse(ctx context.Context, authReq models.AuthRequestInt) (string, string, error) {
	resp := &provider.Response{
		ProtocolBinding: authReq.GetBindingType(),
		RelayState:      authReq.GetRelayState(),
		AcsUrl:          authReq.GetAccessConsumerServiceURL(),
		RequestID:       authReq.GetAuthRequestID(),
		Audience:        authReq.GetIssuer(),
	}

	issuer := ContextToIssuer(ctx)
	req, err := http.NewRequestWithContext(provider.ContextWithIssuer(ctx, issuer), http.MethodGet, issuer, nil)
	if err != nil {
		return "", "", err
	}
	resp.Issuer = p.GetEntityID(req)

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

	return createResponse(samlResponse, authReq.GetBindingType(), authReq.GetAccessConsumerServiceURL(), resp.RelayState, resp.SigAlg, resp.Signature)
}

func createResponse(samlResponse interface{}, binding, acs, relayState, sigAlg, sig string) (string, string, error) {
	respData, err := xml.Marshal(samlResponse)
	if err != nil {
		return "", "", err
	}

	switch binding {
	case provider.PostBinding:
		return acs, base64.StdEncoding.EncodeToString(respData), nil
	case provider.RedirectBinding:
		respData, err := xml.DeflateAndBase64(respData)
		if err != nil {
			return "", "", err
		}
		parsed, err := url.Parse(acs)
		if err != nil {
			return "", "", err
		}
		values := parsed.Query()
		values.Add("SAMLResponse", string(respData))
		values.Add("RelayState", relayState)
		values.Add("SigAlg", sigAlg)
		values.Add("Signature", sig)
		parsed.RawQuery = values.Encode()
		return parsed.String(), "", nil
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
