package command

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) ChangeHumanPhone(ctx context.Context, phone *domain.Phone, resourceOwner string, phoneCodeGenerator crypto.Generator) (*domain.Phone, error) {
	if err := phone.Normalize(); err != nil {
		return nil, err
	}
	existingPhone, err := c.phoneWriteModelByID(ctx, phone.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingPhone.UserState.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-3M0fs", "Errors.User.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	changedEvent, hasChanged := existingPhone.NewChangedEvent(ctx, userAgg, phone.PhoneNumber)

	// only continue if there were changes or there were no changes and the phone should be set to verified
	if !hasChanged && !(phone.IsPhoneVerified && existingPhone.IsPhoneVerified != phone.IsPhoneVerified) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-wF94r", "Errors.User.Phone.NotChanged")
	}

	events := make([]eventstore.Command, 0)
	if hasChanged {
		events = append(events, changedEvent)
	}
	if phone.IsPhoneVerified {
		events = append(events, user.NewHumanPhoneVerifiedEvent(ctx, userAgg))
	} else {
		phoneCode, err := domain.NewPhoneCode(phoneCodeGenerator)
		if err != nil {
			return nil, err
		}
		events = append(events, user.NewHumanPhoneCodeAddedEvent(ctx, userAgg, phoneCode.Code, phoneCode.Expiry))
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPhone, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToPhone(existingPhone), nil
}

func (c *Commands) VerifyHumanPhone(ctx context.Context, userID, code, resourceowner string, phoneCodeGenerator crypto.Generator) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Km9ds", "Errors.User.UserIDMissing")
	}
	if code == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-wMe9f", "Errors.User.Code.Empty")
	}

	existingCode, err := c.phoneWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	if !existingCode.UserState.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Rsj8c", "Errors.User.NotFound")
	}
	if !existingCode.State.Exists() || existingCode.Code == nil {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Rsj8c", "Errors.User.Code.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	err = crypto.VerifyCode(existingCode.CodeCreationDate, existingCode.CodeExpiry, existingCode.Code, code, phoneCodeGenerator)
	if err == nil {
		pushedEvents, err := c.eventstore.Push(ctx, user.NewHumanPhoneVerifiedEvent(ctx, userAgg))
		if err != nil {
			return nil, err
		}
		err = AppendAndReduce(existingCode, pushedEvents...)
		if err != nil {
			return nil, err
		}
		return writeModelToObjectDetails(&existingCode.WriteModel), nil
	}
	_, err = c.eventstore.Push(ctx, user.NewHumanPhoneVerificationFailedEvent(ctx, userAgg))
	logging.LogWithFields("COMMAND-5M9ds", "userID", userAgg.ID).OnError(err).Error("NewHumanPhoneVerificationFailedEvent push failed")
	return nil, zerrors.ThrowInvalidArgument(err, "COMMAND-sM0cs", "Errors.User.Code.Invalid")
}

func (c *Commands) CreateHumanPhoneVerificationCode(ctx context.Context, userID, resourceowner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4M0ds", "Errors.User.UserIDMissing")
	}

	existingPhone, err := c.phoneWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}

	if !existingPhone.UserState.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-2M0fs", "Errors.User.NotFound")
	}
	if !existingPhone.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-2b7Hf", "Errors.User.Phone.NotFound")
	}
	if existingPhone.IsPhoneVerified {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-2M9sf", "Errors.User.Phone.AlreadyVerified")
	}
	config, err := secretGeneratorConfig(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyPhoneCode)
	if err != nil {
		return nil, err
	}
	phoneCode, err := domain.NewPhoneCode(crypto.NewEncryptionGenerator(*config, c.userEncryption))
	if err != nil {
		return nil, err
	}

	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	if err = c.pushAppendAndReduce(ctx, existingPhone, user.NewHumanPhoneCodeAddedEvent(ctx, userAgg, phoneCode.Code, phoneCode.Expiry)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPhone.WriteModel), nil
}

func (c *Commands) HumanPhoneVerificationCodeSent(ctx context.Context, orgID, userID string) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-3m9Fs", "Errors.User.UserIDMissing")
	}

	existingPhone, err := c.phoneWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if !existingPhone.UserState.Exists() {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-3M9fs", "Errors.User.NotFound")
	}
	if !existingPhone.State.Exists() {
		return zerrors.ThrowNotFound(nil, "COMMAND-66n8J", "Errors.User.Phone.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	_, err = c.eventstore.Push(ctx, user.NewHumanPhoneCodeSentEvent(ctx, userAgg))
	return err
}

func (c *Commands) RemoveHumanPhone(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-6M0ds", "Errors.User.UserIDMissing")
	}

	existingPhone, err := c.phoneWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingPhone.UserState.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-3M9fs", "Errors.User.NotFound")
	}
	if !existingPhone.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-p6rsc", "Errors.User.Phone.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, user.NewHumanPhoneRemovedEvent(ctx, userAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPhone, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPhone.WriteModel), nil
}

func (c *Commands) phoneWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *HumanPhoneWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanPhoneWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
