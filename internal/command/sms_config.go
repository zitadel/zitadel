package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/channels/twilio"
	"github.com/caos/zitadel/internal/repository/iam"
)

func (c *Commands) AddSMSConfigTwilio(ctx context.Context, config *twilio.TwilioConfig) (string, *domain.ObjectDetails, error) {
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	smsConfigWriteModel, err := c.getSMSConfig(ctx, id)
	if err != nil {
		return "", nil, err
	}

	iamAgg := IAMAggregateFromWriteModel(&smsConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam.NewSMSConfigTwilioAddedEvent(
		ctx,
		iamAgg,
		id,
		config.SID,
		config.Token,
		config.From))
	if err != nil {
		return "", nil, err
	}
	err = AppendAndReduce(smsConfigWriteModel, pushedEvents...)
	if err != nil {
		return "", nil, err
	}
	return id, writeModelToObjectDetails(&smsConfigWriteModel.WriteModel), nil
}

func (c *Commands) ChangeSMSConfigTwilio(ctx context.Context, id string, config *twilio.TwilioConfig) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "SMS-e9jwf", "Errors.IDMissing")
	}
	smsConfigWriteModel, err := c.getSMSConfig(ctx, id)
	if err != nil {
		return nil, err
	}
	if !smsConfigWriteModel.State.Exists() || smsConfigWriteModel.Twilio == nil {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2m9fw", "Errors.SMSConfig.NotFound")
	}
	iamAgg := IAMAggregateFromWriteModel(&smsConfigWriteModel.WriteModel)

	changedEvent, hasChanged, err := smsConfigWriteModel.NewChangedEvent(
		ctx,
		iamAgg,
		id,
		config.SID,
		config.Token,
		config.From)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-jf9wk", "Errors.NoChangesFound")
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smsConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smsConfigWriteModel.WriteModel), nil
}

func (c *Commands) ActivateSMSConfigTwilio(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "SMS-dn93n", "Errors.IDMissing")
	}
	smsConfigWriteModel, err := c.getSMSConfig(ctx, id)
	if err != nil {
		return nil, err
	}

	if !smsConfigWriteModel.State.Exists() || smsConfigWriteModel.Twilio == nil {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-sn9we", "Errors.SMSConfig.NotFound")
	}
	if smsConfigWriteModel.State == domain.SMSConfigStateActive {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-sn9we", "Errors.SMSConfig.AlreadyActive")
	}
	iamAgg := IAMAggregateFromWriteModel(&smsConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam.NewSMSConfigTwilioActivatedEvent(
		ctx,
		iamAgg,
		id))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smsConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smsConfigWriteModel.WriteModel), nil
}

func (c *Commands) DeactivateSMSConfigTwilio(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "SMS-frkwf", "Errors.IDMissing")
	}
	smsConfigWriteModel, err := c.getSMSConfig(ctx, id)
	if err != nil {
		return nil, err
	}
	if !smsConfigWriteModel.State.Exists() || smsConfigWriteModel.Twilio == nil {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-s39Kg", "Errors.SMSConfig.NotFound")
	}
	if smsConfigWriteModel.State == domain.SMSConfigStateInactive {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-dm9e3", "Errors.SMSConfig.AlreadyDeactivated")
	}

	iamAgg := IAMAggregateFromWriteModel(&smsConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam.NewSMSConfigTwilioDeactivatedEvent(
		ctx,
		iamAgg,
		id))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smsConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smsConfigWriteModel.WriteModel), nil
}

func (c *Commands) RemoveSMSConfigTwilio(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "SMS-3j9fs", "Errors.IDMissing")
	}
	smsConfigWriteModel, err := c.getSMSConfig(ctx, id)
	if err != nil {
		return nil, err
	}
	if !smsConfigWriteModel.State.Exists() || smsConfigWriteModel.Twilio == nil {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-sn9we", "Errors.SMSConfig.NotFound")
	}

	iamAgg := IAMAggregateFromWriteModel(&smsConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam.NewSMSConfigTwilioRemovedEvent(
		ctx,
		iamAgg,
		id))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smsConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smsConfigWriteModel.WriteModel), nil
}
func (c *Commands) getSMSConfig(ctx context.Context, id string) (_ *IAMSMSConfigWriteModel, err error) {
	writeModel := NewIAMSMSConfigWriteModel(id)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}
