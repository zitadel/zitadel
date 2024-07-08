package command

import (
	"context"
	"net"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddSMTPConfig(ctx context.Context, instanceID string, config *smtp.Config) (string, *domain.ObjectDetails, error) {
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}

	from := strings.TrimSpace(config.From)
	if from == "" {
		return "", nil, zerrors.ThrowInvalidArgument(nil, "INST-ASv2d", "Errors.Invalid.Argument")
	}
	fromSplitted := strings.Split(from, "@")
	senderDomain := fromSplitted[len(fromSplitted)-1]
	description := strings.TrimSpace(config.Description)
	replyTo := strings.TrimSpace(config.ReplyToAddress)
	hostAndPort := strings.TrimSpace(config.SMTP.Host)

	if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
		return "", nil, zerrors.ThrowInvalidArgument(nil, "INST-9JdRe", "Errors.Invalid.Argument")
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

	iamAgg := InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewSMTPConfigAddedEvent(
		ctx,
		iamAgg,
		id,
		description,
		config.Tls,
		config.From,
		config.FromName,
		replyTo,
		hostAndPort,
		config.SMTP.User,
		smtpPassword,
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
		return nil, zerrors.ThrowInvalidArgument(nil, "SMTP-x8vo9", "Errors.IDMissing")
	}

	from := strings.TrimSpace(config.From)
	if from == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-HSv2d", "Errors.Invalid.Argument")
	}
	fromSplitted := strings.Split(from, "@")
	senderDomain := fromSplitted[len(fromSplitted)-1]
	description := strings.TrimSpace(config.Description)
	replyTo := strings.TrimSpace(config.ReplyToAddress)
	hostAndPort := strings.TrimSpace(config.SMTP.Host)
	if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-Kv875", "Errors.Invalid.Argument")
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
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-7j8gv", "Errors.SMTPConfig.NotFound")
	}

	err = checkSenderAddress(smtpConfigWriteModel)
	if err != nil {
		return nil, err
	}

	iamAgg := InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel)

	changedEvent, hasChanged, err := smtpConfigWriteModel.NewChangedEvent(
		ctx,
		iamAgg,
		id,
		description,
		config.Tls,
		from,
		config.FromName,
		replyTo,
		hostAndPort,
		config.SMTP.User,
		smtpPassword,
	)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-lh3op", "Errors.NoChangesFound")
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
	if smtpConfigWriteModel.State != domain.SMTPConfigStateActive {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3n9ls", "Errors.SMTPConfig.NotFound")
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
		return nil, zerrors.ThrowInvalidArgument(nil, "SMTP-nm56k", "Errors.IDMissing")
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
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-kg8yr", "Errors.SMTPConfig.NotFound")
	}

	if smtpConfigWriteModel.State == domain.SMTPConfigStateActive {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-ed3lr", "Errors.SMTPConfig.AlreadyActive")
	}

	iamAgg := InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel)
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
		return nil, zerrors.ThrowInvalidArgument(nil, "SMTP-98ikl", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id, "")
	if err != nil {
		return nil, err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-k39PJ", "Errors.SMTPConfig.NotFound")
	}
	if smtpConfigWriteModel.State == domain.SMTPConfigStateInactive {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-km8g3", "Errors.SMTPConfig.AlreadyDeactivated")
	}

	iamAgg := InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel)
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
		return nil, zerrors.ThrowInvalidArgument(nil, "SMTP-7f5cv", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id, "")
	if err != nil {
		return nil, err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-kg8rt", "Errors.SMTPConfig.NotFound")
	}

	iamAgg := InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel)
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

func (c *Commands) TestSMTPConfig(ctx context.Context, instanceID, id, email string, config *smtp.Config) error {
	password := config.SMTP.Password

	if email == "" {
		return zerrors.ThrowInvalidArgument(nil, "SMTP-p9uy", "Errors.SMTPConfig.TestEmailNotFound")
	}

	if id == "" && password == "" {
		return zerrors.ThrowInvalidArgument(nil, "SMTP-p9kj", "Errors.SMTPConfig.TestPassword")
	}

	// If the password is not sent it'd mean that the password hasn't been changed for
	// the stored configuration identified by its id so we can try to retrieve it
	if id != "" && password == "" {
		smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id, "")
		if err != nil {
			return err
		}
		if !smtpConfigWriteModel.State.Exists() {
			return zerrors.ThrowNotFound(nil, "SMTP-p9cc", "Errors.SMTPConfig.NotFound")
		}

		password, err = crypto.DecryptString(smtpConfigWriteModel.Password, c.smtpEncryption)
		if err != nil {
			return err
		}
	}

	config.SMTP.Password = password

	// Try to send an email
	err := smtp.TestConfiguration(config, email)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) TestSMTPConfigById(ctx context.Context, instanceID, id, email string) error {
	if id == "" {
		return zerrors.ThrowInvalidArgument(nil, "SMTP-99oki", "Errors.IDMissing")
	}

	if email == "" {
		return zerrors.ThrowInvalidArgument(nil, "SMTP-99yth", "Errors.SMTPConfig.TestEmailNotFound")
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id, "")
	if err != nil {
		return err
	}

	if !smtpConfigWriteModel.State.Exists() {
		return zerrors.ThrowNotFound(nil, "SMTP-99klw", "Errors.SMTPConfig.NotFound")
	}

	password, err := crypto.DecryptString(smtpConfigWriteModel.Password, c.smtpEncryption)
	if err != nil {
		return err
	}

	smtpConfig := &smtp.Config{
		Description: smtpConfigWriteModel.Description,
		Tls:         smtpConfigWriteModel.TLS,
		From:        smtpConfigWriteModel.SenderAddress,
		FromName:    smtpConfigWriteModel.SenderName,
		SMTP: smtp.SMTP{
			Host:     smtpConfigWriteModel.Host,
			User:     smtpConfigWriteModel.User,
			Password: password,
		},
	}

	// Try to send an email
	err = smtp.TestConfiguration(smtpConfig, email)
	if err != nil {
		return err
	}

	return nil
}

func checkSenderAddress(writeModel *IAMSMTPConfigWriteModel) error {
	if !writeModel.smtpSenderAddressMatchesInstanceDomain {
		return nil
	}
	if !writeModel.domainState.Exists() {
		return zerrors.ThrowInvalidArgument(nil, "INST-83nl8", "Errors.SMTPConfig.SenderAdressNotCustomDomain")
	}
	return nil
}

func (c *Commands) getSMTPConfig(ctx context.Context, instanceID, id, domain string) (writeModel *IAMSMTPConfigWriteModel, err error) {
	writeModel = NewIAMSMTPConfigWriteModel(instanceID, id, domain)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}

// TODO: SetUpInstance still uses this and would be removed as soon as deprecated PrepareCommands is removed
func (c *Commands) prepareAddAndActivateSMTPConfig(a *instance.Aggregate, description, from, name, replyTo, hostAndPort, user string, password []byte, tls bool) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if from = strings.TrimSpace(from); from == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-mruNY", "Errors.Invalid.Argument")
		}

		replyTo = strings.TrimSpace(replyTo)

		hostAndPort = strings.TrimSpace(hostAndPort)
		if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-9JdRe", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			id, err := id_generator.Next()
			if err != nil {
				return nil, zerrors.ThrowInternal(nil, "INST-9JdRe", "Errors.Invalid.Argument")
			}

			fromSplitted := strings.Split(from, "@")
			senderDomain := fromSplitted[len(fromSplitted)-1]
			writeModel, err := getSMTPConfigWriteModel(ctx, filter, id, senderDomain)
			if err != nil {
				return nil, err
			}
			if writeModel.State == domain.SMTPConfigStateActive {
				return nil, zerrors.ThrowAlreadyExists(nil, "INST-W3VS2", "Errors.SMTPConfig.AlreadyExists")
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
					description,
					tls,
					from,
					name,
					replyTo,
					hostAndPort,
					user,
					smtpPassword,
				),
				instance.NewSMTPConfigActivatedEvent(
					ctx,
					&a.Aggregate,
					id,
				),
			}, nil
		}, nil
	}
}

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
