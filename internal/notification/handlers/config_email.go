package handlers

import (
	"context"
	"net/http"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/notification/channels/email"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// GetActiveEmailConfig reads the SMTP provider config for the given org.
// Uses org-level SMTP when available. When orgID is empty, uses instance-level
// SMTP (backward compatibility). When orgID is set but org has no config,
// behavior depends on OrgSMTPFallbackToInstance: if true, falls back to
// instance-level; if false, returns an error.
func (n *NotificationQueries) GetActiveEmailConfig(ctx context.Context, orgID string) (*email.Config, error) {
	if orgID != "" {
		config, err := n.OrgSMTPConfigActive(ctx, orgID)
		if err == nil {
			return n.smtpConfigToEmailConfig(config)
		}
		// Only fall back on "not found" — propagate real errors (DB failures, etc.)
		if !zerrors.IsNotFound(err) {
			return nil, err
		}
		if !n.OrgSMTPFallbackToInstance {
			return nil, err
		}
		// fall through to instance-level
	}
	config, err := n.SMTPConfigActive(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return n.smtpConfigToEmailConfig(config)
}

// smtpConfigToEmailConfig converts a query.SMTPConfig to an email.Config,
// handling SMTP, HTTP webhook, and crypto decryption.
func (n *NotificationQueries) smtpConfigToEmailConfig(config *query.SMTPConfig) (*email.Config, error) {
	provider := &email.Provider{
		ID:          config.ID,
		Description: config.Description,
	}

	if config.SMTPConfig != nil {
		return smtpToEmailConfig(config.SMTPConfig, provider, n.SMTPPasswordCrypto)
	}

	if config.HTTPConfig != nil {
		return &email.Config{
			ProviderConfig: provider,
			WebhookConfig: &webhook.Config{
				CallURL:    config.HTTPConfig.Endpoint,
				Method:     http.MethodPost,
				Headers:    nil,
				SigningKey: config.HTTPConfig.SigningKey,
			},
		}, nil
	}
	return nil, zerrors.ThrowNotFound(nil, "QUERY-KPQleOckOV", "Errors.SMTPConfig.NotFound")
}

func smtpToEmailConfig(qs *query.SMTP, provider *email.Provider, passCrypto crypto.EncryptionAlgorithm) (*email.Config, error) {
	config := &email.Config{
		ProviderConfig: provider,
		SMTPConfig: &smtp.Config{
			From:           qs.SenderAddress,
			FromName:       qs.SenderName,
			ReplyToAddress: qs.ReplyToAddress,
			Tls:            qs.TLS,
			SMTP: smtp.SMTP{
				Host: qs.Host,
			},
		},
	}

	if qs.XOAuth2Auth != nil {
		config.SMTPConfig.SMTP.XOAuth2Auth = &smtp.XOAuth2AuthConfig{
			User:          qs.User,
			TokenEndpoint: qs.XOAuth2Auth.TokenEndpoint,
			Scopes:        qs.XOAuth2Auth.Scopes,
		}
		if qs.XOAuth2Auth.ClientCredentials != nil {
			clientSecret, err := crypto.DecryptString(qs.XOAuth2Auth.ClientCredentials.ClientSecret, passCrypto)
			if err != nil {
				return nil, err
			}
			config.SMTPConfig.SMTP.XOAuth2Auth.ClientCredentialsAuth = &smtp.OAuth2ClientCredentials{
				ClientId:     qs.XOAuth2Auth.ClientCredentials.ClientId,
				ClientSecret: clientSecret,
			}
		}
	}

	if qs.PlainAuth != nil {
		config.SMTPConfig.SMTP.PlainAuth = &smtp.PlainAuthConfig{
			User: qs.User,
		}

		if qs.PlainAuth.Password != nil {
			password, err := crypto.DecryptString(qs.PlainAuth.Password, passCrypto)
			if err != nil {
				return nil, err
			}
			config.SMTPConfig.SMTP.PlainAuth.Password = password
		}
	}

	// if no auth is configured but there is a user, use plain auth without a password
	if qs.User != "" &&
		config.SMTPConfig.SMTP.PlainAuth == nil &&
		config.SMTPConfig.SMTP.XOAuth2Auth == nil {
		config.SMTPConfig.SMTP.PlainAuth = &smtp.PlainAuthConfig{
			User: qs.User,
		}
	}

	return config, nil
}
