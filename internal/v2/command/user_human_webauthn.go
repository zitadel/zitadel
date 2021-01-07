package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) RemoveHumanU2F(ctx context.Context, userID, webAuthNID string) error {
	event := user.NewHumanU2FRemovedEvent(ctx, webAuthNID)
	return r.removeHumanWebAuthN(ctx, userID, webAuthNID, event)
}

func (r *CommandSide) RemoveHumanPasswordless(ctx context.Context, userID, webAuthNID string) error {
	event := user.NewHumanPasswordlessRemovedEvent(ctx, webAuthNID)
	return r.removeHumanWebAuthN(ctx, userID, webAuthNID, event)
}

func (r *CommandSide) removeHumanWebAuthN(ctx context.Context, userID, webAuthNID string, event eventstore.EventPusher) error {
	if userID == "" || webAuthNID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M9de", "Errors.IDMissing")
	}

	existingWebAuthN, err := r.webauthNWriteModelByID(ctx, userID, webAuthNID)
	if err != nil {
		return err
	}
	if existingWebAuthN.UserState == domain.UserStateUnspecified || existingWebAuthN.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowAlreadyExists(nil, "COMMAND-5M0ds", "Errors.User.ExternalIDP.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingWebAuthN.WriteModel)
	userAgg.PushEvents(event)

	return r.eventstore.PushAggregate(ctx, existingWebAuthN, userAgg)
}

func (r *CommandSide) webauthNWriteModelByID(ctx context.Context, userID, webAuthNID string) (writeModel *HumanWebAuthNWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanWebAuthNWriteModel(userID, webAuthNID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
