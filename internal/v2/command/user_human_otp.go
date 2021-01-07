package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) RemoveHumanOTP(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0sd", "Errors.User.UserIDMissing")
	}

	existingOTP, err := r.otpWriteModelByID(ctx, userID)
	if err != nil {
		return err
	}
	if existingOTP.UserState == domain.UserStateUnspecified || existingOTP.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowAlreadyExists(nil, "COMMAND-5M0ds", "Errors.User.OTP.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	userAgg.PushEvents(
		user.NewHumanOTPRemovedEvent(ctx),
	)

	return r.eventstore.PushAggregate(ctx, existingOTP, userAgg)
}

func (r *CommandSide) otpWriteModelByID(ctx context.Context, userID string) (writeModel *HumanOTPWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanOTPWriteModel(userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
