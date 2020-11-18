package oidc

import "github.com/caos/zitadel/internal/crypto"

type AddedEvent struct {
	eventstore.BaseEvent

	IDPConfigID           string              `json:"idpConfigId"`
	ClientID              string              `json:"clientId"`
	Secret                *crypto.CryptoValue `json:"clientSecret"`
	Issuer                string              `json:"issuer"`
	Scopes                []string            `json:"scpoes"`
	IDPDisplayNameMapping int32               `json:"idpDisplayNameMapping,omitempty"`
	UsernameMapping       int32               `json:"usernameMapping,omitempty"`
}
