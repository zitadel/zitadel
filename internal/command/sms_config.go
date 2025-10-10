package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddTwilioConfig struct {
	Details       *domain.ObjectDetails
	ResourceOwner string
	ID            string

	Description      string
	SID              string
	Token            string
	SenderNumber     string
	VerifyServiceSID string
}

func (c *Commands) AddSMSConfigTwilio(ctx context.Context, config *AddTwilioConfig) (err error) {
	if config.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-ZLrZhKSKq0", "Errors.ResourceOwnerMissing")
	}
	if config.ID == "" {
		config.ID, err = c.idGenerator.Next()
		if err != nil {
			return err
		}
	}
	smsConfigWriteModel, err := c.getSMSConfig(ctx, config.ResourceOwner, config.ID)
	if err != nil {
		return err
	}

	var token *crypto.CryptoValue
	if config.Token != "" {
		token, err = crypto.Encrypt([]byte(config.Token), c.smsEncryption)
		if err != nil {
			return err
		}
	}
	err = c.pushAppendAndReduce(ctx,
		smsConfigWriteModel,
		instance.NewSMSConfigTwilioAddedEvent(
			ctx,
			InstanceAggregateFromWriteModel(&smsConfigWriteModel.WriteModel),
			config.ID,
			config.Description,
			config.SID,
			config.SenderNumber,
			token,
			config.VerifyServiceSID,
		),
	)
	if err != nil {
		return err
	}
	config.Details = writeModelToObjectDetails(&smsConfigWriteModel.WriteModel)
	return nil
}

type ChangeTwilioConfig struct {
	Details       *domain.ObjectDetails
	ResourceOwner string
	ID            string

	Description      *string
	SID              *string
	Token            *string
	SenderNumber     *string
	VerifyServiceSID *string
}

func (c *Commands) ChangeSMSConfigTwilio(ctx context.Context, config *ChangeTwilioConfig) (err error) {
	if config.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-RHXryJwmFG", "Errors.ResourceOwnerMissing")
	}
	if config.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-gMr93iNhTR", "Errors.IDMissing")
	}
	smsConfigWriteModel, err := c.getSMSConfig(ctx, config.ResourceOwner, config.ID)
	if err != nil {
		return err
	}
	if !smsConfigWriteModel.State.Exists() || smsConfigWriteModel.Twilio == nil {
		return zerrors.ThrowNotFound(nil, "COMMAND-MUY0IFAf8O", "Errors.SMSConfig.NotFound")
	}
	changedEvent, hasChanged, err := smsConfigWriteModel.NewTwilioChangedEvent(
		ctx,
		InstanceAggregateFromWriteModel(&smsConfigWriteModel.WriteModel),
		config.ID,
		config.Description,
		config.SID,
		config.SenderNumber,
		config.VerifyServiceSID,
	)
	if err != nil {
		return err
	}
	if !hasChanged {
		config.Details = writeModelToObjectDetails(&smsConfigWriteModel.WriteModel)
		return nil
	}
	err = c.pushAppendAndReduce(ctx,
		smsConfigWriteModel,
		changedEvent,
	)
	if err != nil {
		return err
	}
	config.Details = writeModelToObjectDetails(&smsConfigWriteModel.WriteModel)
	return nil
}

func (c *Commands) ChangeSMSConfigTwilioToken(ctx context.Context, resourceOwner, id, token string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-sLLA1HnMzj", "Errors.ResourceOwnerMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "SMS-PeNaqbC0r0", "Errors.IDMissing")
	}

	smsConfigWriteModel, err := c.getSMSConfig(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if !smsConfigWriteModel.State.Exists() || smsConfigWriteModel.Twilio == nil {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-ij3NhEHATp", "Errors.SMSConfig.NotFound")
	}
	newtoken, err := crypto.Encrypt([]byte(token), c.smsEncryption)
	if err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx,
		smsConfigWriteModel,
		instance.NewSMSConfigTokenChangedEvent(
			ctx,
			InstanceAggregateFromWriteModel(&smsConfigWriteModel.WriteModel),
			id,
			newtoken,
		),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smsConfigWriteModel.WriteModel), nil
}

type AddSMSHTTP struct {
	Details       *domain.ObjectDetails
	ResourceOwner string
	ID            string

	Description string
	Endpoint    string
	SigningKey  string
}

func (c *Commands) AddSMSConfigHTTP(ctx context.Context, config *AddSMSHTTP) (err error) {
	if config.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-huy99qWjX4", "Errors.ResourceOwnerMissing")
	}
	if config.ID == "" {
		config.ID, err = c.idGenerator.Next()
		if err != nil {
			return err
		}
	}
	smsConfigWriteModel, err := c.getSMSConfig(ctx, config.ResourceOwner, config.ID)
	if err != nil {
		return err
	}

	code, err := c.newSigningKey(ctx, c.eventstore.Filter, c.smsEncryption) //nolint
	if err != nil {
		return err
	}
	config.SigningKey = code.PlainCode()

	err = c.pushAppendAndReduce(ctx,
		smsConfigWriteModel,
		instance.NewSMSConfigHTTPAddedEvent(
			ctx,
			InstanceAggregateFromWriteModel(&smsConfigWriteModel.WriteModel),
			config.ID,
			config.Description,
			config.Endpoint,
			code.Crypted,
		),
	)
	if err != nil {
		return err
	}
	config.Details = writeModelToObjectDetails(&smsConfigWriteModel.WriteModel)
	return nil
}

type ChangeSMSHTTP struct {
	Details       *domain.ObjectDetails
	ResourceOwner string
	ID            string

	Description          *string
	Endpoint             *string
	ExpirationSigningKey bool
	SigningKey           *string
}

func (c *Commands) ChangeSMSConfigHTTP(ctx context.Context, config *ChangeSMSHTTP) (err error) {
	if config.ResourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-M622CFQnwK", "Errors.ResourceOwnerMissing")
	}
	if config.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-phyb2e4Kll", "Errors.IDMissing")
	}
	smsConfigWriteModel, err := c.getSMSConfig(ctx, config.ResourceOwner, config.ID)
	if err != nil {
		return err
	}
	if !smsConfigWriteModel.State.Exists() || smsConfigWriteModel.HTTP == nil {
		return zerrors.ThrowNotFound(nil, "COMMAND-6NW4I5Kqzj", "Errors.SMSConfig.NotFound")
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

	changedEvent, hasChanged, err := smsConfigWriteModel.NewHTTPChangedEvent(
		ctx,
		InstanceAggregateFromWriteModel(&smsConfigWriteModel.WriteModel),
		config.ID,
		config.Description,
		config.Endpoint,
		changedSigningKey)
	if err != nil {
		return err
	}
	if !hasChanged {
		config.Details = writeModelToObjectDetails(&smsConfigWriteModel.WriteModel)
		return nil
	}
	err = c.pushAppendAndReduce(ctx, smsConfigWriteModel, changedEvent)
	if err != nil {
		return err
	}
	config.Details = writeModelToObjectDetails(&smsConfigWriteModel.WriteModel)
	return nil
}

func (c *Commands) ActivateSMSConfig(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-EFgoOg997V", "Errors.ResourceOwnerMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-jJ6TVqzvjp", "Errors.IDMissing")
	}

	smsConfigWriteModel, err := c.getSMSConfig(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}

	if !smsConfigWriteModel.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-9ULtp9PH5E", "Errors.SMSConfig.NotFound")
	}
	if smsConfigWriteModel.State == domain.SMSConfigStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-B25GFeIvRi", "Errors.SMSConfig.AlreadyActive")
	}
	err = c.pushAppendAndReduce(ctx, smsConfigWriteModel,
		instance.NewSMSConfigActivatedEvent(
			ctx,
			InstanceAggregateFromWriteModel(&smsConfigWriteModel.WriteModel),
			id,
		),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smsConfigWriteModel.WriteModel), nil
}

func (c *Commands) DeactivateSMSConfig(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-V9NWOZj8Gi", "Errors.ResourceOwnerMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-xs1ah1v1CL", "Errors.IDMissing")
	}
	smsConfigWriteModel, err := c.getSMSConfig(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if !smsConfigWriteModel.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-La91dGNhbM", "Errors.SMSConfig.NotFound")
	}
	if smsConfigWriteModel.State == domain.SMSConfigStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-OSZAEkYvk7", "Errors.SMSConfig.AlreadyDeactivated")
	}
	err = c.pushAppendAndReduce(ctx,
		smsConfigWriteModel,
		instance.NewSMSConfigDeactivatedEvent(
			ctx,
			InstanceAggregateFromWriteModel(&smsConfigWriteModel.WriteModel),
			id,
		),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smsConfigWriteModel.WriteModel), nil
}

func (c *Commands) RemoveSMSConfig(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-cw0NSJsn1v", "Errors.ResourceOwnerMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Qrz7lvdC4c", "Errors.IDMissing")
	}
	smsConfigWriteModel, err := c.getSMSConfig(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if !smsConfigWriteModel.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-povEVHPCkV", "Errors.SMSConfig.NotFound")
	}

	err = c.pushAppendAndReduce(ctx,
		smsConfigWriteModel,
		instance.NewSMSConfigRemovedEvent(
			ctx,
			InstanceAggregateFromWriteModel(&smsConfigWriteModel.WriteModel),
			id,
		),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smsConfigWriteModel.WriteModel), nil
}

func (c *Commands) getSMSConfig(ctx context.Context, instanceID, id string) (_ *IAMSMSConfigWriteModel, err error) {
	writeModel := NewIAMSMSConfigWriteModel(instanceID, id)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

// getActiveSMSConfig returns the last activated SMS configuration
func (c *Commands) getActiveSMSConfig(ctx context.Context, instanceID string) (_ *IAMSMSConfigWriteModel, err error) {
	writeModel := NewIAMSMSLastActivatedConfigWriteModel(instanceID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return c.getSMSConfig(ctx, instanceID, writeModel.activeID)
}
