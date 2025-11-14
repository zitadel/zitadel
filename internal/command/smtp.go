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
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddSMTPConfig struct {
	Details       *domain.ObjectDetails
	ResourceOwner string
	ID            string

	Description    string
	Host           string
	User           string
	Password       string
	Tls            bool
	From           string
	FromName       string
	ReplyToAddress string
}

func (c *Commands) AddSMTPConfig(ctx context.Context, config *AddSMTPConfig) (err error) {
	if config.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-PQN0wsqSyi", "Errors.ResourceOwnerMissing")
	}
	if config.ID == "" {
		config.ID, err = c.idGenerator.Next()
		if err != nil {
			return err
		}
	}

	from := strings.TrimSpace(config.From)
	if from == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-SAAFpV8VKV", "Errors.Invalid.Argument")
	}
	fromSplitted := strings.Split(from, "@")
	senderDomain := fromSplitted[len(fromSplitted)-1]
	description := strings.TrimSpace(config.Description)
	replyTo := strings.TrimSpace(config.ReplyToAddress)
	hostAndPort := strings.TrimSpace(config.Host)

	if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-EvAtufIinh", "Errors.Invalid.Argument")
	}

	var smtpPassword *crypto.CryptoValue
	if config.Password != "" {
		smtpPassword, err = crypto.Encrypt([]byte(config.Password), c.smtpEncryption)
		if err != nil {
			return err
		}
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, config.ResourceOwner, config.ID, senderDomain)
	if err != nil {
		return err
	}

	err = checkSenderAddress(smtpConfigWriteModel)
	if err != nil {
		return err
	}

	err = c.pushAppendAndReduce(ctx,
		smtpConfigWriteModel,
		instance.NewSMTPConfigAddedEvent(
			ctx,
			InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
			config.ID,
			description,
			config.Tls,
			config.From,
			config.FromName,
			replyTo,
			hostAndPort,
			config.User,
			smtpPassword,
		),
	)
	if err != nil {
		return err
	}
	config.Details = writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel)
	return nil
}

type ChangeSMTPConfig struct {
	Details       *domain.ObjectDetails
	ResourceOwner string
	ID            string

	Description    string
	Host           string
	User           string
	Password       string
	Tls            bool
	From           string
	FromName       string
	ReplyToAddress string
}

func (c *Commands) ChangeSMTPConfig(ctx context.Context, config *ChangeSMTPConfig) error {
	if config.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-jwA8gxldy3", "Errors.ResourceOwnerMissing")
	}
	if config.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-2JPlSRzuHy", "Errors.IDMissing")
	}

	from := strings.TrimSpace(config.From)
	if from == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-gyPUXOTA4N", "Errors.Invalid.Argument")
	}
	fromSplitted := strings.Split(from, "@")
	senderDomain := fromSplitted[len(fromSplitted)-1]
	description := strings.TrimSpace(config.Description)
	replyTo := strings.TrimSpace(config.ReplyToAddress)
	hostAndPort := strings.TrimSpace(config.Host)
	if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-kZNVkuL32L", "Errors.Invalid.Argument")
	}

	var smtpPassword *crypto.CryptoValue
	var err error
	if config.Password != "" {
		smtpPassword, err = crypto.Encrypt([]byte(config.Password), c.smtpEncryption)
		if err != nil {
			return err
		}
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, config.ResourceOwner, config.ID, senderDomain)
	if err != nil {
		return err
	}

	if !smtpConfigWriteModel.State.Exists() {
		return zerrors.ThrowNotFound(nil, "COMMAND-j5IDFtt3T1", "Errors.SMTPConfig.NotFound")
	}

	err = checkSenderAddress(smtpConfigWriteModel)
	if err != nil {
		return err
	}

	changedEvent, hasChanged, err := smtpConfigWriteModel.NewChangedEvent(
		ctx,
		InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
		config.ID,
		description,
		config.Tls,
		from,
		config.FromName,
		replyTo,
		hostAndPort,
		config.User,
		smtpPassword,
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

func (c *Commands) ChangeSMTPConfigPassword(ctx context.Context, resourceOwner, id string, password string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-gHAyvUXCAF", "Errors.ResourceOwnerMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-BCkAf7LcJA", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, resourceOwner, id, "")
	if err != nil {
		return nil, err
	}
	if smtpConfigWriteModel.State != domain.SMTPConfigStateActive {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-rDHzqjGuKQ", "Errors.SMTPConfig.NotFound")
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
		instance.NewSMTPConfigPasswordChangedEvent(
			ctx,
			InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
			id,
			smtpPassword,
		),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

type AddSMTPConfigHTTP struct {
	Details       *domain.ObjectDetails
	ResourceOwner string
	ID            string

	Description string
	Endpoint    string
	SigningKey  string
}

func (c *Commands) AddSMTPConfigHTTP(ctx context.Context, config *AddSMTPConfigHTTP) (err error) {
	if config.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-FTNDXc8ACS", "Errors.ResourceOwnerMissing")
	}
	if config.ID == "" {
		config.ID, err = c.idGenerator.Next()
		if err != nil {
			return err
		}
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, config.ResourceOwner, config.ID, "")
	if err != nil {
		return err
	}

	code, err := c.newSigningKey(ctx, c.eventstore.Filter, c.smtpEncryption) //nolint
	if err != nil {
		return err
	}
	config.SigningKey = code.PlainCode()

	err = c.pushAppendAndReduce(ctx, smtpConfigWriteModel, instance.NewSMTPConfigHTTPAddedEvent(
		ctx,
		InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
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

type ChangeSMTPConfigHTTP struct {
	Details       *domain.ObjectDetails
	ResourceOwner string
	ID            string

	Description          string
	Endpoint             string
	ExpirationSigningKey bool
	SigningKey           *string
}

func (c *Commands) ChangeSMTPConfigHTTP(ctx context.Context, config *ChangeSMTPConfigHTTP) (err error) {
	if config.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-k7QCGOWyJA", "Errors.ResourceOwnerMissing")
	}
	if config.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-2MHkV8ObWo", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, config.ResourceOwner, config.ID, "")
	if err != nil {
		return err
	}

	if !smtpConfigWriteModel.State.Exists() || smtpConfigWriteModel.HTTPConfig == nil {
		return zerrors.ThrowNotFound(nil, "COMMAND-xIrdledqv4", "Errors.SMTPConfig.NotFound")
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
		InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
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

func (c *Commands) ActivateSMTPConfig(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-h5htMCebv3", "Errors.ResourceOwnerMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-1hPl6oVMJa", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, resourceOwner, id, "")
	if err != nil {
		return nil, err
	}

	if !smtpConfigWriteModel.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-E9K20hxOS9", "Errors.SMTPConfig.NotFound")
	}
	if smtpConfigWriteModel.State == domain.SMTPConfigStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-vUHBSmBzaw", "Errors.SMTPConfig.AlreadyActive")
	}

	err = c.pushAppendAndReduce(ctx,
		smtpConfigWriteModel,
		instance.NewSMTPConfigActivatedEvent(
			ctx,
			InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
			id,
		),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) DeactivateSMTPConfig(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-pvNHou89Tw", "Errors.ResourceOwnerMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-jLTIMrtApO", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, resourceOwner, id, "")
	if err != nil {
		return nil, err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-k39PJ", "Errors.SMTPConfig.NotFound")
	}
	if smtpConfigWriteModel.State == domain.SMTPConfigStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-km8g3", "Errors.SMTPConfig.AlreadyDeactivated")
	}

	err = c.pushAppendAndReduce(ctx,
		smtpConfigWriteModel,
		instance.NewSMTPConfigDeactivatedEvent(
			ctx,
			InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
			id,
		),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) RemoveSMTPConfig(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-t2WsPRgGaK", "Errors.ResourceOwnerMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-0ZV5whuUfu", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, resourceOwner, id, "")
	if err != nil {
		return nil, err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-09CXlTDL6w", "Errors.SMTPConfig.NotFound")
	}

	err = c.pushAppendAndReduce(ctx,
		smtpConfigWriteModel,
		instance.NewSMTPConfigRemovedEvent(
			ctx,
			InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel),
			id,
		),
	)
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
		if !smtpConfigWriteModel.State.Exists() || smtpConfigWriteModel.SMTPConfig == nil {
			return zerrors.ThrowNotFound(nil, "SMTP-p9cc", "Errors.SMTPConfig.NotFound")
		}

		password, err = crypto.DecryptString(smtpConfigWriteModel.SMTPConfig.Password, c.smtpEncryption)
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

	if !smtpConfigWriteModel.State.Exists() || smtpConfigWriteModel.SMTPConfig == nil {
		return zerrors.ThrowNotFound(nil, "SMTP-99klw", "Errors.SMTPConfig.NotFound")
	}

	password, err := crypto.DecryptString(smtpConfigWriteModel.SMTPConfig.Password, c.smtpEncryption)
	if err != nil {
		return err
	}

	smtpConfig := &smtp.Config{
		Tls:      smtpConfigWriteModel.SMTPConfig.TLS,
		From:     smtpConfigWriteModel.SMTPConfig.SenderAddress,
		FromName: smtpConfigWriteModel.SMTPConfig.SenderName,
		SMTP: smtp.SMTP{
			Host:     smtpConfigWriteModel.SMTPConfig.Host,
			User:     smtpConfigWriteModel.SMTPConfig.User,
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
		return zerrors.ThrowInvalidArgument(nil, "INST-xtWIiR2ZbR", "Errors.SMTPConfig.SenderAdressNotCustomDomain")
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
			id, err := c.idGenerator.Next()
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
