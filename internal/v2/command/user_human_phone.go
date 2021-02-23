package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) ChangeHumanPhone(ctx context.Context, phone *domain.Phone) (*domain.Phone, error) {
	if !phone.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M0ds", "Errors.Phone.Invalid")
	}

	existingPhone, err := r.phoneWriteModelByID(ctx, phone.AggregateID, phone.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if existingPhone.State == domain.PhoneStateUnspecified || existingPhone.State == domain.PhoneStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-aM9cs", "Errors.User.Phone.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	changedEvent, hasChanged := existingPhone.NewChangedEvent(ctx, userAgg, phone.PhoneNumber)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-wF94r", "Errors.User.Phone.NotChanged")
	}

	events := []eventstore.EventPusher{changedEvent}
	if phone.IsPhoneVerified {
		events = append(events, user.NewHumanPhoneVerifiedEvent(ctx, userAgg))
	} else {
		phoneCode, err := domain.NewPhoneCode(r.phoneVerificationCode)
		if err != nil {
			return nil, err
		}
		events = append(events, user.NewHumanPhoneCodeAddedEvent(ctx, userAgg, phoneCode.Code, phoneCode.Expiry))
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPhone, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToPhone(existingPhone), nil
}

func (r *CommandSide) VerifyHumanPhone(ctx context.Context, userID, code, resourceowner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Km9ds", "Errors.User.UserIDMissing")
	}
	if code == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-wMe9f", "Errors.User.Code.Empty")
	}

	existingCode, err := r.phoneWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	if existingCode.Code == nil || existingCode.State == domain.PhoneStateUnspecified || existingCode.State == domain.PhoneStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-Rsj8c", "Errors.User.Code.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	err = crypto.VerifyCode(existingCode.CodeCreationDate, existingCode.CodeExpiry, existingCode.Code, code, r.phoneVerificationCode)
	if err == nil {
		_, err = r.eventstore.PushEvents(ctx, user.NewHumanPhoneVerifiedEvent(ctx, userAgg))
		return err
	}
	_, err = r.eventstore.PushEvents(ctx, user.NewHumanPhoneVerificationFailedEvent(ctx, userAgg))

	logging.LogWithFields("COMMAND-5M9ds", "userID", userAgg.ID).OnError(err).Error("NewHumanPhoneVerificationFailedEvent push failed")
	return caos_errs.ThrowInvalidArgument(err, "COMMAND-sM0cs", "Errors.User.Code.Invalid")
}

func (r *CommandSide) CreateHumanPhoneVerificationCode(ctx context.Context, userID, resourceowner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0ds", "Errors.User.UserIDMissing")
	}

	existingPhone, err := r.phoneWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return err
	}

	//TODO: code like the following if is written many times find way to simplify
	if existingPhone.State == domain.PhoneStateUnspecified || existingPhone.State == domain.PhoneStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-2b7Hf", "Errors.User.Phone.NotFound")
	}
	if existingPhone.IsPhoneVerified {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9sf", "Errors.User.Phone.AlreadyVerified")
	}

	phoneCode, err := domain.NewPhoneCode(r.phoneVerificationCode)
	if err != nil {
		return err
	}

	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, user.NewHumanPhoneCodeAddedEvent(ctx, userAgg, phoneCode.Code, phoneCode.Expiry))
	return err
}

func (r *CommandSide) HumanPhoneVerificationCodeSent(ctx context.Context, orgID, userID string) (err error) {
	existingPhone, err := r.phoneWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingPhone.State == domain.PhoneStateUnspecified || existingPhone.State == domain.PhoneStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-66n8J", "Errors.User.Phone.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, user.NewHumanPhoneCodeSentEvent(ctx, userAgg))
	return err
}

func (r *CommandSide) RemoveHumanPhone(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M0ds", "Errors.User.UserIDMissing")
	}

	existingPhone, err := r.phoneWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingPhone.State == domain.PhoneStateUnspecified || existingPhone.State == domain.PhoneStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-p6rsc", "Errors.User.Phone.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, user.NewHumanPhoneRemovedEvent(ctx, userAgg))
	return err
}

func (r *CommandSide) phoneWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *HumanPhoneWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanPhoneWriteModel(userID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
