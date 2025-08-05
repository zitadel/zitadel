package saml

import (
	"github.com/zitadel/saml/pkg/provider/models"

	"github.com/zitadel/zitadel/internal/command"
)

var _ models.AuthRequestInt = &AuthRequestV2{}

type AuthRequestV2 struct {
	*command.CurrentSAMLRequest
}

func (a *AuthRequestV2) GetApplicationID() string {
	return a.ApplicationID
}

func (a *AuthRequestV2) GetID() string {
	return a.ID
}
func (a *AuthRequestV2) GetRelayState() string {
	return a.RelayState
}
func (a *AuthRequestV2) GetAccessConsumerServiceURL() string {
	return a.ACSURL
}
func (a *AuthRequestV2) GetAuthRequestID() string {
	return a.RequestID
}
func (a *AuthRequestV2) GetBindingType() string {
	return a.Binding
}
func (a *AuthRequestV2) GetIssuer() string {
	return a.Issuer
}
func (a *AuthRequestV2) GetDestination() string {
	return a.Destination
}
func (a *AuthRequestV2) GetUserID() string {
	return a.UserID
}
func (a *AuthRequestV2) Done() bool {
	return a.UserID != "" && a.SessionID != ""
}
