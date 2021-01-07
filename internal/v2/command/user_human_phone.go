package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) ChangeHumanPhone(ctx context.Context, phone *domain.Phone) (*domain.Phone, error) {
	if !phone.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M0ds", "Errors.Phone.Invalid")
	}

	existingPhone, err := r.phoneWriteModel(ctx, phone.AggregateID)
	if err != nil {
		return nil, err
	}
	if existingPhone.State == domain.PhoneStateUnspecified || existingPhone.State == domain.PhoneStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-5M0ds", "Errors.User.Phone.NotFound")
	}
	changedEvent, hasChanged := existingPhone.NewChangedEvent(ctx, phone.PhoneNumber)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-wF94r", "Errors.User.Phone.NotChanged")
	}
	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	userAgg.PushEvents(changedEvent)

	if phone.IsPhoneVerified {
		userAgg.PushEvents(user.NewHumanPhoneVerifiedEvent(ctx))
	} else {
		phoneCode, err := domain.NewPhoneCode(r.phoneVerificationCode)
		if err != nil {
			return nil, err
		}
		userAgg.PushEvents(user.NewHumanPhoneCodeAddedEvent(ctx, phoneCode.Code, phoneCode.Expiry))
	}

	err = r.eventstore.PushAggregate(ctx, existingPhone, userAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPhone(existingPhone), nil
}

func (r *CommandSide) CreateHumanPhoneVerificationCode(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0ds", "Errors.User.UserIDMissing")
	}

	existingPhone, err := r.phoneWriteModel(ctx, userID)
	if err != nil {
		return err
	}
	if existingPhone.State == domain.PhoneStateUnspecified || existingPhone.State == domain.PhoneStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-2M9fs", "Errors.User.Phone.NotFound")
	}
	if existingPhone.IsPhoneVerified {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9sf", "Errors.User.Phone.AlreadyVerified")
	}
	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	phoneCode, err := domain.NewPhoneCode(r.phoneVerificationCode)
	if err != nil {
		return err
	}
	userAgg.PushEvents(user.NewHumanPhoneCodeAddedEvent(ctx, phoneCode.Code, phoneCode.Expiry))
	return r.eventstore.PushAggregate(ctx, existingPhone, userAgg)
}

func (r *CommandSide) RemoveHumanPhone(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M0ds", "Errors.User.UserIDMissing")
	}

	existingPhone, err := r.phoneWriteModel(ctx, userID)
	if err != nil {
		return err
	}
	if existingPhone.State == domain.PhoneStateUnspecified || existingPhone.State == domain.PhoneStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5M0ds", "Errors.User.Phone.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	userAgg.PushEvents(
		user.NewHumanPhoneRemovedEvent(ctx),
	)
	return r.eventstore.PushAggregate(ctx, existingPhone, userAgg)
}

func (r *CommandSide) phoneWriteModel(ctx context.Context, userID string) (writeModel *HumanPhoneWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanPhoneWriteModel(userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
