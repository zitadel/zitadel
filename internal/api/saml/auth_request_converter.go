package saml

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/zitadel/saml/pkg/provider/models"
	"github.com/zitadel/saml/pkg/provider/xml/samlp"
	"time"
)

var _ models.AuthRequestInt = &AuthRequest{}

type AuthRequest struct {
	*domain.AuthRequest
}

func (a *AuthRequest) GetApplicationID() string {
	return a.ApplicationID
}

func (a *AuthRequest) GetID() string {
	return a.ID
}
func (a *AuthRequest) GetRelayState() string {
	return a.TransferState
}
func (a *AuthRequest) GetAccessConsumerServiceURL() string {
	return a.CallbackURI
}

func (a *AuthRequest) GetNameID() string {
	return a.UserName
}

func (a *AuthRequest) saml() *domain.AuthRequestSAML {
	return a.Request.(*domain.AuthRequestSAML)
}
func (a *AuthRequest) GetAuthRequestID() string {
	return a.saml().ID
}
func (a *AuthRequest) GetBindingType() string {
	return a.saml().BindingType
}
func (a *AuthRequest) GetIssuer() string {
	return a.saml().Issuer
}
func (a *AuthRequest) GetIssuerName() string {
	return a.saml().IssuerName
}
func (a *AuthRequest) GetDestination() string {
	return a.saml().Destination
}
func (a *AuthRequest) GetCode() string {
	return a.saml().Code
}
func (a *AuthRequest) GetUserID() string {
	return a.UserID
}
func (a *AuthRequest) GetUserName() string {
	return a.UserName
}
func (a *AuthRequest) Done() bool {
	for _, step := range a.PossibleSteps {
		if step.Type() == domain.NextStepRedirectToCallback {
			return true
		}
	}
	return false
}

func AuthRequestFromBusiness(authReq *domain.AuthRequest) (_ models.AuthRequestInt, err error) {
	if _, ok := authReq.Request.(*domain.AuthRequestSAML); !ok {
		return nil, errors.ThrowInvalidArgument(nil, "OIDC-Hbz7A", "auth request is not of type saml")
	}
	return &AuthRequest{authReq}, nil
}

func CreateAuthRequestToBusiness(authReq *samlp.AuthnRequestType, acsUrl, protocolBinding, applicationID, relayState, userAgentID string) *domain.AuthRequest {
	return &domain.AuthRequest{
		CreationDate:  time.Now(),
		AgentID:       userAgentID,
		ApplicationID: applicationID,
		CallbackURI:   acsUrl,
		TransferState: relayState,
		Request: &domain.AuthRequestSAML{
			ID:          authReq.Id,
			BindingType: protocolBinding,
			Code:        "",
			Issuer:      authReq.Issuer.Text,
			IssuerName:  authReq.Issuer.SPProvidedID,
			Destination: authReq.Destination,
		},
	}
}
