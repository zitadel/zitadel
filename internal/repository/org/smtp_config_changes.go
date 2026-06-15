package org

import (
	"github.com/zitadel/zitadel/internal/crypto"
)

// Change helper functions for OrgSMTPConfigChangedEvent.
// These mirror the instance-level change helpers.

func ChangeOrgSMTPConfigID(id string) func(event *OrgSMTPConfigChangedEvent) {
	return func(e *OrgSMTPConfigChangedEvent) {
		e.ID = id
	}
}

func ChangeOrgSMTPConfigDescription(description string) func(event *OrgSMTPConfigChangedEvent) {
	return func(e *OrgSMTPConfigChangedEvent) {
		e.Description = &description
	}
}

func ChangeOrgSMTPConfigTLS(tls bool) func(event *OrgSMTPConfigChangedEvent) {
	return func(e *OrgSMTPConfigChangedEvent) {
		e.TLS = &tls
	}
}

func ChangeOrgSMTPConfigFromAddress(senderAddress string) func(event *OrgSMTPConfigChangedEvent) {
	return func(e *OrgSMTPConfigChangedEvent) {
		e.FromAddress = &senderAddress
	}
}

func ChangeOrgSMTPConfigFromName(senderName string) func(event *OrgSMTPConfigChangedEvent) {
	return func(e *OrgSMTPConfigChangedEvent) {
		e.FromName = &senderName
	}
}

func ChangeOrgSMTPConfigReplyToAddress(replyToAddress string) func(event *OrgSMTPConfigChangedEvent) {
	return func(e *OrgSMTPConfigChangedEvent) {
		e.ReplyToAddress = &replyToAddress
	}
}

func ChangeOrgSMTPConfigSMTPHost(smtpHost string) func(event *OrgSMTPConfigChangedEvent) {
	return func(e *OrgSMTPConfigChangedEvent) {
		e.Host = &smtpHost
	}
}

func ChangeOrgSMTPConfigSMTPUser(smtpUser string) func(event *OrgSMTPConfigChangedEvent) {
	return func(e *OrgSMTPConfigChangedEvent) {
		e.User = &smtpUser
	}
}

func ChangeOrgSMTPConfigSMTPPassword(password *crypto.CryptoValue) func(event *OrgSMTPConfigChangedEvent) {
	return func(e *OrgSMTPConfigChangedEvent) {
		e.Password = password
		e.PlainAuth.Password = password
	}
}

func ChangeOrgSMTPConfigXOAuth2ClientCredentialsClientId(clientId string) func(event *OrgSMTPConfigChangedEvent) {
	return func(e *OrgSMTPConfigChangedEvent) {
		e.XOAuth2Auth.ClientCredentials.ClientId = &clientId
	}
}

func ChangeOrgSMTPConfigXOAuth2ClientCredentialsClientSecret(clientSecret *crypto.CryptoValue) func(event *OrgSMTPConfigChangedEvent) {
	return func(e *OrgSMTPConfigChangedEvent) {
		e.XOAuth2Auth.ClientCredentials.ClientSecret = clientSecret
	}
}

func ChangeOrgSMTPConfigXOAuth2TokenEndpoint(tokenEndpoint string) func(event *OrgSMTPConfigChangedEvent) {
	return func(e *OrgSMTPConfigChangedEvent) {
		e.XOAuth2Auth.TokenEndpoint = &tokenEndpoint
	}
}

func ChangeOrgSMTPConfigXOAuth2Scopes(scopes []string) func(event *OrgSMTPConfigChangedEvent) {
	return func(e *OrgSMTPConfigChangedEvent) {
		e.XOAuth2Auth.Scopes = scopes
	}
}

// Change helper functions for OrgSMTPConfigHTTPChangedEvent.

func ChangeOrgSMTPConfigHTTPID(id string) func(event *OrgSMTPConfigHTTPChangedEvent) {
	return func(e *OrgSMTPConfigHTTPChangedEvent) {
		e.ID = id
	}
}

func ChangeOrgSMTPConfigHTTPDescription(description string) func(event *OrgSMTPConfigHTTPChangedEvent) {
	return func(e *OrgSMTPConfigHTTPChangedEvent) {
		e.Description = &description
	}
}

func ChangeOrgSMTPConfigHTTPEndpoint(endpoint string) func(event *OrgSMTPConfigHTTPChangedEvent) {
	return func(e *OrgSMTPConfigHTTPChangedEvent) {
		e.Endpoint = &endpoint
	}
}

func ChangeOrgSMTPConfigHTTPSigningKey(signingKey *crypto.CryptoValue) func(event *OrgSMTPConfigHTTPChangedEvent) {
	return func(e *OrgSMTPConfigHTTPChangedEvent) {
		e.SigningKey = signingKey
	}
}
