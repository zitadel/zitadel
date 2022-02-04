package command

import (
	"context"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/notification/channels/smtp"
	"github.com/caos/zitadel/internal/notification/channels/twilio"
	"github.com/caos/zitadel/internal/repository/project"
)

func (c *Commands) AddSMTPConfig(ctx context.Context, config *smtp.EmailConfig) (*domain.ObjectDetails, error) {
	smtpConfigWriteModel, err := c.getSMTPConfig(ctx)
	if err != nil {
		return nil, err
	}
	if smtpConfigWriteModel.State == domain.SMTPConfigStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-en9lw", "Errors.SMTPConfig.AlreadyExists")
	}
	var smtpPassword *crypto.CryptoValue
	if config.SMTP.Password != "" {
		smtpPassword, err = crypto.Encrypt([]byte(config.SMTP.Password), c.smtpPasswordCrypto)
		if err != nil {
			return nil, err
		}
	}

	iamAgg := IAMAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam.NewSMTPConfigAddedEvent(
		ctx,
		iamAgg,
		config.Tls,
		config.From,
		config.FromName,
		config.SMTP.Host,
		config.SMTP.User,
		smtpPassword))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smtpConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) ChangeSMTPConfig(ctx context.Context, config *smtp.EmailConfig) (*domain.ObjectDetails, error) {
	smtpConfigWriteModel, err := c.getSMTPConfig(ctx)
	if err != nil {
		return nil, err
	}
	if smtpConfigWriteModel.State == domain.SMTPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-3n9ls", "Errors.SMTPConfig.NotFound")
	}
	iamAgg := IAMAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel)

	changedEvent, hasChanged, err := smtpConfigWriteModel.NewChangedEvent(
		ctx,
		iamAgg,
		config.Tls,
		config.From,
		config.FromName,
		config.SMTP.Host,
		config.SMTP.User)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-m0o3f", "Errors.NoChangesFound")
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smtpConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) getSMSConfig(ctx context.Context) (_ *IAMSMTPConfigWriteModel, err error) {
	writeModel := NewIAMSMTPConfigWriteModel()
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}
