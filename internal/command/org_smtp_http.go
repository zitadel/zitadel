package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddOrgSMTPConfigHTTP struct {
	Details       *domain.ObjectDetails
	ResourceOwner string
	ID            string

	Description string
	Endpoint    string
	SigningKey  string
}

func (c *Commands) AddOrgSMTPConfigHTTP(ctx context.Context, config *AddOrgSMTPConfigHTTP) (err error) {
	if config.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "ORG-FTNDXc1", "Errors.ResourceOwnerMissing")
	}
	if config.ID == "" {
		config.ID, err = c.idGenerator.Next()
		if err != nil {
			return err
		}
	}

	smtpConfigWriteModel, err := c.getOrgSMTPConfig(ctx, config.ResourceOwner, config.ID)
	if err != nil {
		return err
	}

	code, err := c.newSigningKey(ctx, c.eventstore.Filter, c.smtpEncryption) //nolint
	if err != nil {
		return err
	}
	config.SigningKey = code.PlainCode()

	err = c.pushAppendAndReduce(ctx, smtpConfigWriteModel, org.NewOrgSMTPConfigHTTPAddedEvent(
		ctx,
		OrgAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
		config.ID,
		config.Description,
		config.Endpoint,
		code.Crypted,
	))
	if err != nil {
		return err
	}
	config.Details = writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel)
	return nil
}

type ChangeOrgSMTPConfigHTTP struct {
	Details       *domain.ObjectDetails
	ResourceOwner string
	ID            string

	Description          string
	Endpoint             string
	ExpirationSigningKey bool
	SigningKey           *string
}

func (c *Commands) ChangeOrgSMTPConfigHTTP(ctx context.Context, config *ChangeOrgSMTPConfigHTTP) (err error) {
	if config.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "ORG-k7QCG1", "Errors.ResourceOwnerMissing")
	}
	if config.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "ORG-2MHkV1", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getOrgSMTPConfig(ctx, config.ResourceOwner, config.ID)
	if err != nil {
		return err
	}

	if !smtpConfigWriteModel.State.Exists() || smtpConfigWriteModel.HTTPConfig == nil {
		return zerrors.ThrowNotFound(nil, "ORG-xIrdl1", "Errors.SMTPConfig.NotFound")
	}

	var changedSigningKey *crypto.CryptoValue
	if config.ExpirationSigningKey {
		code, err := c.newSigningKey(ctx, c.eventstore.Filter, c.smtpEncryption) //nolint
		if err != nil {
			return err
		}
		changedSigningKey = code.Crypted
		config.SigningKey = &code.Plain
	}

	changedEvent, hasChanged, err := smtpConfigWriteModel.NewHTTPChangedEvent(
		ctx,
		OrgAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
		config.ID,
		config.Description,
		config.Endpoint,
		changedSigningKey,
	)
	if err != nil {
		return err
	}
	if !hasChanged {
		config.Details = writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel)
		return nil
	}

	err = c.pushAppendAndReduce(ctx, smtpConfigWriteModel, changedEvent)
	if err != nil {
		return err
	}
	config.Details = writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel)
	return nil
}

func (c *Commands) TestOrgSMTPConfig(
	ctx context.Context, orgID, id, email string, config *smtp.Config,
) error {
	if email == "" {
		return zerrors.ThrowInvalidArgument(nil, "ORG-p9uy1", "Errors.SMTPConfig.TestEmailNotFound")
	}

	if id == "" && config.SMTP.PlainAuth != nil && config.SMTP.PlainAuth.Password == "" {
		return zerrors.ThrowInvalidArgument(nil, "ORG-p9kj1", "Errors.SMTPConfig.TestPassword")
	}

	if id == "" && config.SMTP.XOAuth2Auth != nil &&
		config.SMTP.XOAuth2Auth.ClientCredentialsAuth != nil &&
		config.SMTP.XOAuth2Auth.ClientCredentialsAuth.ClientSecret == "" {
		return zerrors.ThrowInvalidArgument(nil, "ORG-9OP961", "Errors.SMTPConfig.TestClientSecret")
	}

	if config.SMTP.PlainAuth != nil && config.SMTP.PlainAuth.Password == "" {
		smtpConfigWriteModel, err := c.getOrgSMTPConfig(ctx, orgID, id)
		if err != nil {
			return err
		}
		if !smtpConfigWriteModel.State.Exists() ||
			smtpConfigWriteModel.SMTPConfig == nil ||
			smtpConfigWriteModel.SMTPConfig.PlainAuth == nil {
			return zerrors.ThrowNotFound(nil, "ORG-p9cc1", "Errors.SMTPConfig.NotFound")
		}
		config.SMTP.PlainAuth.Password, err = crypto.DecryptString(
			smtpConfigWriteModel.SMTPConfig.PlainAuth.Password, c.smtpEncryption,
		)
		if err != nil {
			return err
		}
	}

	if config.SMTP.XOAuth2Auth != nil &&
		config.SMTP.XOAuth2Auth.ClientCredentialsAuth != nil &&
		config.SMTP.XOAuth2Auth.ClientCredentialsAuth.ClientSecret == "" {
		smtpConfigWriteModel, err := c.getOrgSMTPConfig(ctx, orgID, id)
		if err != nil {
			return err
		}
		if !smtpConfigWriteModel.State.Exists() ||
			smtpConfigWriteModel.SMTPConfig == nil ||
			smtpConfigWriteModel.SMTPConfig.XOAuth2Auth == nil {
			return zerrors.ThrowNotFound(nil, "ORG-p9c21", "Errors.SMTPConfig.NotFound")
		}
		config.SMTP.XOAuth2Auth.ClientCredentialsAuth.ClientSecret, err = crypto.DecryptString(
			smtpConfigWriteModel.SMTPConfig.XOAuth2Auth.ClientCredentials.ClientSecret,
			c.smtpEncryption,
		)
		if err != nil {
			return err
		}
	}

	return smtp.TestConfiguration(config, email)
}

func (c *Commands) TestOrgSMTPConfigById(ctx context.Context, orgID, id, email string) error {
	if id == "" {
		return zerrors.ThrowInvalidArgument(nil, "ORG-99oki1", "Errors.IDMissing")
	}
	if email == "" {
		return zerrors.ThrowInvalidArgument(nil, "ORG-99yth1", "Errors.SMTPConfig.TestEmailNotFound")
	}

	smtpConfigWriteModel, err := c.getOrgSMTPConfig(ctx, orgID, id)
	if err != nil {
		return err
	}
	if !smtpConfigWriteModel.State.Exists() || smtpConfigWriteModel.SMTPConfig == nil {
		return zerrors.ThrowNotFound(nil, "ORG-99klw1", "Errors.SMTPConfig.NotFound")
	}

	var plainAuth *smtp.PlainAuthConfig
	var xoauth2Auth *smtp.XOAuth2AuthConfig

	if smtpConfigWriteModel.SMTPConfig.PlainAuth != nil {
		password, err := crypto.DecryptString(
			smtpConfigWriteModel.SMTPConfig.PlainAuth.Password, c.smtpEncryption,
		)
		if err != nil {
			return err
		}
		plainAuth = &smtp.PlainAuthConfig{
			User:     smtpConfigWriteModel.SMTPConfig.User,
			Password: password,
		}
	}

	if smtpConfigWriteModel.SMTPConfig.XOAuth2Auth != nil {
		xoauth2Auth = &smtp.XOAuth2AuthConfig{
			User:          smtpConfigWriteModel.SMTPConfig.User,
			TokenEndpoint: smtpConfigWriteModel.SMTPConfig.XOAuth2Auth.TokenEndpoint,
			Scopes:        smtpConfigWriteModel.SMTPConfig.XOAuth2Auth.Scopes,
		}
		if smtpConfigWriteModel.SMTPConfig.XOAuth2Auth.ClientCredentials != nil {
			clientSecret, err := crypto.DecryptString(
				smtpConfigWriteModel.SMTPConfig.XOAuth2Auth.ClientCredentials.ClientSecret,
				c.smtpEncryption,
			)
			if err != nil {
				return err
			}
			xoauth2Auth.ClientCredentialsAuth = &smtp.OAuth2ClientCredentials{
				ClientId:     smtpConfigWriteModel.SMTPConfig.XOAuth2Auth.ClientCredentials.ClientId,
				ClientSecret: clientSecret,
			}
		}
	}

	smtpConfig := &smtp.Config{
		Tls:      smtpConfigWriteModel.SMTPConfig.TLS,
		From:     smtpConfigWriteModel.SMTPConfig.SenderAddress,
		FromName: smtpConfigWriteModel.SMTPConfig.SenderName,
		SMTP: smtp.SMTP{
			Host:        smtpConfigWriteModel.SMTPConfig.Host,
			PlainAuth:   plainAuth,
			XOAuth2Auth: xoauth2Auth,
		},
	}

	return smtp.TestConfiguration(smtpConfig, email)
}
