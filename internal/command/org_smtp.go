package command

import (
	"context"
	"net"
	"strings"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddOrgSMTPConfig struct {
	Details       *domain.ObjectDetails
	ResourceOwner string
	ID            string

	Description    string
	Host           string
	User           string
	Tls            bool
	From           string
	FromName       string
	ReplyToAddress string
	PlainAuth      *PlainAuth
	XOAuth2Auth    *XOAuth2Auth
}

func (c *Commands) AddOrgSMTPConfig(ctx context.Context, config *AddOrgSMTPConfig) (err error) {
	if config.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "ORG-sn93Ss", "Errors.ResourceOwnerMissing")
	}
	if config.ID == "" {
		config.ID, err = c.idGenerator.Next()
		if err != nil {
			return err
		}
	}

	from := strings.TrimSpace(config.From)
	if from == "" {
		return zerrors.ThrowInvalidArgument(nil, "ORG-ASF3g2", "Errors.Invalid.Argument")
	}
	description := strings.TrimSpace(config.Description)
	replyTo := strings.TrimSpace(config.ReplyToAddress)
	hostAndPort := strings.TrimSpace(config.Host)

	if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
		return zerrors.ThrowInvalidArgument(nil, "ORG-gK9RE2", "Errors.Invalid.Argument")
	}

	smtpConfigWriteModel, err := c.getOrgSMTPConfig(ctx, config.ResourceOwner, config.ID)
	if err != nil {
		return err
	}

	var plainAuth *instance.PlainAuth
	var xoauth2Auth *instance.XOAuth2Auth

	if config.XOAuth2Auth != nil {
		xoauth2Auth = &instance.XOAuth2Auth{
			TokenEndpoint: config.XOAuth2Auth.TokenEndpoint,
			Scopes:        config.XOAuth2Auth.Scopes,
		}
		if config.XOAuth2Auth.ClientCredentialsAuth != nil {
			clientSecret, err := crypto.Encrypt(
				[]byte(config.XOAuth2Auth.ClientCredentialsAuth.ClientSecret),
				c.smtpEncryption,
			)
			if err != nil {
				return err
			}
			xoauth2Auth.ClientCredentials = &instance.XOAuth2ClientCredentials{
				ClientId:     config.XOAuth2Auth.ClientCredentialsAuth.ClientId,
				ClientSecret: clientSecret,
			}
		}
	}
	if config.PlainAuth != nil {
		plainAuth = &instance.PlainAuth{}
		if config.PlainAuth.Password != "" {
			plainAuth.Password, err = crypto.Encrypt([]byte(config.PlainAuth.Password), c.smtpEncryption)
			if err != nil {
				return err
			}
		}
	}

	err = c.pushAppendAndReduce(ctx,
		smtpConfigWriteModel,
		org.NewOrgSMTPConfigAddedEvent(
			ctx,
			OrgAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
			config.ID,
			description,
			config.Tls,
			config.From,
			config.FromName,
			replyTo,
			hostAndPort,
			config.User,
			plainAuth,
			xoauth2Auth,
		),
	)
	if err != nil {
		return err
	}
	config.Details = writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel)
	return nil
}

type ChangeOrgSMTPConfig struct {
	Details       *domain.ObjectDetails
	ResourceOwner string
	ID            string

	Description    string
	Host           string
	User           string
	Tls            bool
	From           string
	FromName       string
	ReplyToAddress string
	PlainAuth      *PlainAuth
	XOAuth2Auth    *XOAuth2Auth
}

func (c *Commands) ChangeOrgSMTPConfig(ctx context.Context, config *ChangeOrgSMTPConfig) error {
	if config.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "ORG-wJk3s0", "Errors.ResourceOwnerMissing")
	}
	if config.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "ORG-2MHqR1", "Errors.IDMissing")
	}

	from := strings.TrimSpace(config.From)
	if from == "" {
		return zerrors.ThrowInvalidArgument(nil, "ORG-g9PXN4", "Errors.Invalid.Argument")
	}
	description := strings.TrimSpace(config.Description)
	replyTo := strings.TrimSpace(config.ReplyToAddress)
	hostAndPort := strings.TrimSpace(config.Host)
	if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
		return zerrors.ThrowInvalidArgument(nil, "ORG-kZ3Vk2", "Errors.Invalid.Argument")
	}

	smtpConfigWriteModel, err := c.getOrgSMTPConfig(ctx, config.ResourceOwner, config.ID)
	if err != nil {
		return err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return zerrors.ThrowNotFound(nil, "ORG-j5IDn3", "Errors.SMTPConfig.NotFound")
	}

	var plainAuth *instance.PlainAuth
	var xoauth2Auth *instance.XOAuth2Auth

	if config.PlainAuth != nil {
		plainAuth = &instance.PlainAuth{}
		if config.PlainAuth.Password != "" {
			plainAuth.Password, err = crypto.Encrypt([]byte(config.PlainAuth.Password), c.smtpEncryption)
			if err != nil {
				return err
			}
		}
	}
	if config.XOAuth2Auth != nil {
		xoauth2Auth = &instance.XOAuth2Auth{
			TokenEndpoint: config.XOAuth2Auth.TokenEndpoint,
			Scopes:        config.XOAuth2Auth.Scopes,
		}
		if config.XOAuth2Auth.ClientCredentialsAuth != nil {
			xoauth2Auth.ClientCredentials = &instance.XOAuth2ClientCredentials{
				ClientId: config.XOAuth2Auth.ClientCredentialsAuth.ClientId,
			}
			if config.XOAuth2Auth.ClientCredentialsAuth.ClientSecret != "" {
				xoauth2Auth.ClientCredentials.ClientSecret, err = crypto.Encrypt(
					[]byte(config.XOAuth2Auth.ClientCredentialsAuth.ClientSecret),
					c.smtpEncryption,
				)
				if err != nil {
					return err
				}
			}
		}
	}

	changedEvent, hasChanged, err := smtpConfigWriteModel.NewChangedEvent(
		ctx,
		OrgAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
		config.ID,
		description,
		config.Tls,
		from,
		config.FromName,
		replyTo,
		hostAndPort,
		config.User,
		plainAuth,
		xoauth2Auth,
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

func (c *Commands) RemoveOrgSMTPConfig(ctx context.Context, orgID, id string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-t2WsP1", "Errors.ResourceOwnerMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-0ZV5w1", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getOrgSMTPConfig(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "ORG-09CXl1", "Errors.SMTPConfig.NotFound")
	}

	err = c.pushAppendAndReduce(ctx,
		smtpConfigWriteModel,
		org.NewOrgSMTPConfigRemovedEvent(
			ctx,
			OrgAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
			id,
		),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) ChangeOrgSMTPConfigPassword(
	ctx context.Context, orgID, id string, password string,
) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-gHAyv1", "Errors.ResourceOwnerMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-BCkAf1", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getOrgSMTPConfig(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if smtpConfigWriteModel.State != domain.SMTPConfigStateActive {
		return nil, zerrors.ThrowNotFound(nil, "ORG-rDHzq1", "Errors.SMTPConfig.NotFound")
	}

	var smtpPassword *crypto.CryptoValue
	if password != "" {
		smtpPassword, err = crypto.Encrypt([]byte(password), c.smtpEncryption)
		if err != nil {
			return nil, err
		}
	}

	err = c.pushAppendAndReduce(ctx,
		smtpConfigWriteModel,
		org.NewOrgSMTPConfigPasswordChangedEvent(
			ctx,
			OrgAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
			id,
			smtpPassword,
		),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) ActivateOrgSMTPConfig(
	ctx context.Context, orgID, id string,
) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-h5htM1", "Errors.ResourceOwnerMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-1hPl61", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getOrgSMTPConfig(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "ORG-E9K201", "Errors.SMTPConfig.NotFound")
	}
	if smtpConfigWriteModel.State == domain.SMTPConfigStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "ORG-vUHBS1", "Errors.SMTPConfig.AlreadyActive")
	}

	err = c.pushAppendAndReduce(ctx,
		smtpConfigWriteModel,
		org.NewOrgSMTPConfigActivatedEvent(
			ctx,
			OrgAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
			id,
		),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) DeactivateOrgSMTPConfig(
	ctx context.Context, orgID, id string,
) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-pvNHo1", "Errors.ResourceOwnerMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-jLTIM1", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getOrgSMTPConfig(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "ORG-k39PJ1", "Errors.SMTPConfig.NotFound")
	}
	if smtpConfigWriteModel.State == domain.SMTPConfigStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "ORG-km8g31", "Errors.SMTPConfig.AlreadyDeactivated")
	}

	err = c.pushAppendAndReduce(ctx,
		smtpConfigWriteModel,
		org.NewOrgSMTPConfigDeactivatedEvent(
			ctx,
			OrgAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
			id,
		),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) getOrgSMTPConfig(
	ctx context.Context, orgID, id string,
) (*OrgSMTPConfigWriteModel, error) {
	writeModel := NewOrgSMTPConfigWriteModel(orgID, id)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
