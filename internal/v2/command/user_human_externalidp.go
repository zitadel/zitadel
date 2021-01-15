package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) RemoveHumanExternalIDP(ctx context.Context, externalIDP *domain.ExternalIDP) error {
	return r.removeHumanExternalIDP(ctx, externalIDP, false)
}

func (r *CommandSide) removeHumanExternalIDP(ctx context.Context, externalIDP *domain.ExternalIDP, cascade bool) error {
	if externalIDP.IsValid() {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M9ds", "Errors.IDMissing")
	}

	existingExternalIDP, err := r.externalIDPWriteModelByID(ctx, externalIDP.AggregateID, externalIDP.IDPConfigID, externalIDP.ExternalUserID, externalIDP.ResourceOwner)
	if err != nil {
		return err
	}
	if existingExternalIDP.State == domain.ExternalIDPStateUnspecified || existingExternalIDP.State == domain.ExternalIDPStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-1M9xR", "Errors.User.ExternalIDP.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingExternalIDP.WriteModel)
	if !cascade {
		userAgg.PushEvents(
			user.NewHumanExternalIDPRemovedEvent(ctx, externalIDP.IDPConfigID, externalIDP.ExternalUserID),
		)
	} else {
		userAgg.PushEvents(
			user.NewHumanExternalIDPCascadeRemovedEvent(ctx, externalIDP.IDPConfigID, externalIDP.ExternalUserID),
		)
	}

	//TODO: Release unique externalidp
	return r.eventstore.PushAggregate(ctx, existingExternalIDP, userAgg)
}

func (r *CommandSide) externalIDPWriteModelByID(ctx context.Context, userID, idpConfigID, externalUserID, resourceOwner string) (writeModel *HumanExternalIDPWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanExternalIDPWriteModel(userID, idpConfigID, externalUserID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
