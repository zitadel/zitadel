package command

import (
	"context"
	"net"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func (c *Commands) AddSMTPConfig(ctx context.Context, instanceID string, config *smtp.Config) (string, *domain.ObjectDetails, error) {
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}

	from := strings.TrimSpace(config.From)
	if from == "" {
		return "", nil, errors.ThrowInvalidArgument(nil, "INST-ASv2d", "Errors.Invalid.Argument")
	}
	fromSplitted := strings.Split(from, "@")
	senderDomain := fromSplitted[len(fromSplitted)-1]

	replyTo := strings.TrimSpace(config.ReplyToAddress)
	hostAndPort := strings.TrimSpace(config.SMTP.Host)

	if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
		return "", nil, errors.ThrowInvalidArgument(nil, "INST-9JdRe", "Errors.Invalid.Argument")
	}

	var smtpPassword *crypto.CryptoValue
	if config.SMTP.Password != "" {
		smtpPassword, err = crypto.Encrypt([]byte(config.SMTP.Password), c.smtpEncryption)
		if err != nil {
			return "", nil, err
		}
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id, senderDomain)
	if err != nil {
		return "", nil, err
	}

	err = checkSenderAddress(smtpConfigWriteModel)
	if err != nil {
		return "", nil, err
	}

	iamAgg := InstanceAggregateFromWriteModel(ctx, &smtpConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewSMTPConfigAddedEvent(
		ctx,
		iamAgg,
		id,
		config.Tls,
		config.From,
		config.FromName,
		replyTo,
		hostAndPort,
		config.SMTP.User,
		smtpPassword,
		config.SMTP.ProviderType,
	))
	if err != nil {
		return "", nil, err
	}

	err = AppendAndReduce(smtpConfigWriteModel, pushedEvents...)
	if err != nil {
		return "", nil, err
	}
	return id, writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) ChangeSMTPConfig(ctx context.Context, instanceID string, id string, config *smtp.Config) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "SMTP-x8vo9", "Errors.IDMissing")
	}

	from := strings.TrimSpace(config.From)
	if from == "" {
		return nil, errors.ThrowInvalidArgument(nil, "INST-HSv2d", "Errors.Invalid.Argument")
	}
	fromSplitted := strings.Split(from, "@")
	senderDomain := fromSplitted[len(fromSplitted)-1]

	replyTo := strings.TrimSpace(config.ReplyToAddress)
	hostAndPort := strings.TrimSpace(config.SMTP.Host)
	if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
		return nil, errors.ThrowInvalidArgument(nil, "INST-Kv875", "Errors.Invalid.Argument")
	}

	var smtpPassword *crypto.CryptoValue
	var err error
	if config.SMTP.Password != "" {
		smtpPassword, err = crypto.Encrypt([]byte(config.SMTP.Password), c.smtpEncryption)
		if err != nil {
			return nil, err
		}
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id, senderDomain)
	if err != nil {
		return nil, err
	}

	if !smtpConfigWriteModel.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-7j8gv", "Errors.SMTPConfig.NotFound")
	}

	err = checkSenderAddress(smtpConfigWriteModel)
	if err != nil {
		return nil, err
	}

	iamAgg := InstanceAggregateFromWriteModel(ctx, &smtpConfigWriteModel.WriteModel)

	changedEvent, hasChanged, err := smtpConfigWriteModel.NewChangedEvent(
		ctx,
		iamAgg,
		id,
		config.Tls,
		from,
		config.FromName,
		replyTo,
		hostAndPort,
		config.SMTP.User,
		smtpPassword,
		config.SMTP.ProviderType,
	)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-lh3op", "Errors.NoChangesFound")
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

func (c *Commands) ChangeSMTPConfigPassword(ctx context.Context, instanceID, id string, password string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id, "")
	if err != nil {
		return nil, err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return nil, errors.ThrowNotFound(nil, "COMMAND-3n9ls", "Errors.SMTPConfig.NotFound")
	}

	var smtpPassword *crypto.CryptoValue
	if password != "" {
		smtpPassword, err = crypto.Encrypt([]byte(password), c.smtpEncryption)
		if err != nil {
			return nil, err
		}
	}

	pushedEvents, err := c.eventstore.Push(ctx, instance.NewSMTPConfigPasswordChangedEvent(
		ctx,
		&instanceAgg.Aggregate,
		id,
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

func (c *Commands) ActivateSMTPConfig(ctx context.Context, instanceID, id, activatedId string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "SMTP-nm56k", "Errors.IDMissing")
	}

	if len(activatedId) > 0 {
		_, err := c.DeactivateSMTPConfig(ctx, instanceID, activatedId)
		if err != nil {
			return nil, err
		}
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id, "")
	if err != nil {
		return nil, err
	}

	if !smtpConfigWriteModel.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-kg8yr", "Errors.SMTPConfig.NotFound")
	}

	if smtpConfigWriteModel.State == domain.SMTPConfigStateActive {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-ed3lr", "Errors.SMTPConfig.AlreadyActive")
	}

	iamAgg := InstanceAggregateFromWriteModel(ctx, &smtpConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewSMTPConfigActivatedEvent(
		ctx,
		iamAgg,
		id))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smtpConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) DeactivateSMTPConfig(ctx context.Context, instanceID, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "SMTP-98ikl", "Errors.IDMissing")
	}
	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id, "")
	if err != nil {
		return nil, err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-k39PJ", "Errors.SMTPConfig.NotFound")
	}
	if smtpConfigWriteModel.State == domain.SMTPConfigStateInactive {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-km8g3", "Errors.SMTPConfig.AlreadyDeactivated")
	}

	iamAgg := InstanceAggregateFromWriteModel(ctx, &smtpConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewSMTPConfigDeactivatedEvent(
		ctx,
		iamAgg,
		id))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smtpConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) RemoveSMTPConfig(ctx context.Context, instanceID, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "SMTP-7f5cv", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id, "")
	if err != nil {
		return nil, err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-kg8rt", "Errors.SMTPConfig.NotFound")
	}

	iamAgg := InstanceAggregateFromWriteModel(ctx, &smtpConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewSMTPConfigRemovedEvent(
		ctx,
		iamAgg,
		id))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smtpConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func checkSenderAddress(writeModel *IAMSMTPConfigWriteModel) error {
	if !writeModel.smtpSenderAddressMatchesInstanceDomain {
		return nil
	}
	if !writeModel.domainState.Exists() {
		return errors.ThrowInvalidArgument(nil, "INST-83nl8", "Errors.SMTPConfig.SenderAdressNotCustomDomain")
	}
	return nil
}

func (c *Commands) getSMTPConfig(ctx context.Context, instanceID, id, domain string) (_ *IAMSMTPConfigWriteModel, err error) {
	writeModel := NewIAMSMTPConfigWriteModel(instanceID, id, domain)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}

// TODO: SetUpInstance still uses this and would be removed as soon as deprecated PrepareCommands is removed
func (c *Commands) prepareAddSMTPConfig(a *instance.Aggregate, from, name, replyTo, hostAndPort, user string, password []byte, tls bool, providerType uint32) preparation.Validation {
	return func() (preparation.CreateCommands, error) {

		if from = strings.TrimSpace(from); from == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-mruNY", "Errors.Invalid.Argument")
		}

		replyTo = strings.TrimSpace(replyTo)

		hostAndPort = strings.TrimSpace(hostAndPort)
		if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
			return nil, errors.ThrowInvalidArgument(nil, "INST-9JdRe", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			id, err := c.idGenerator.Next()
			if err != nil {
				return nil, errors.ThrowInternal(nil, "INST-9JdRe", "Errors.Invalid.Argument")
			}

			fromSplitted := strings.Split(from, "@")
			senderDomain := fromSplitted[len(fromSplitted)-1]
			writeModel, err := getSMTPConfigWriteModel(ctx, filter, id, senderDomain)
			if err != nil {
				return nil, err
			}
			if writeModel.State == domain.SMTPConfigStateActive {
				return nil, errors.ThrowAlreadyExists(nil, "INST-W3VS2", "Errors.SMTPConfig.AlreadyExists")
			}
			err = checkSenderAddress(writeModel)
			if err != nil {
				return nil, err
			}
			var smtpPassword *crypto.CryptoValue
			if password != nil {
				smtpPassword, err = crypto.Encrypt(password, c.smtpEncryption)
				if err != nil {
					return nil, err
				}
			}
			return []eventstore.Command{
				instance.NewSMTPConfigAddedEvent(
					ctx,
					&a.Aggregate,
					id,
					tls,
					from,
					name,
					replyTo,
					hostAndPort,
					user,
					smtpPassword,
					providerType,
				),
			}, nil
		}, nil
	}
}

// TODO: SetUpInstance still uses this and would be removed as soon as deprecated PrepareCommands is removed
func getSMTPConfigWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer, id, domain string) (_ *IAMSMTPConfigWriteModel, err error) {
	writeModel := NewIAMSMTPConfigWriteModel(authz.GetInstance(ctx).InstanceID(), id, domain)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}
