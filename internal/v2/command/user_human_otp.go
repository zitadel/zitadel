package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) RemoveHumanOTP(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0sd", "Errors.User.UserIDMissing")
	}

	existingOTP, err := r.otpWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingOTP.State == domain.OTPStateUnspecified || existingOTP.State == domain.OTPStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5M0ds", "Errors.User.OTP.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	userAgg.PushEvents(
		user.NewHumanOTPRemovedEvent(ctx),
	)

	return r.eventstore.PushAggregate(ctx, existingOTP, userAgg)
}

func (r *CommandSide) otpWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *HumanOTPWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanOTPWriteModel(userID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
