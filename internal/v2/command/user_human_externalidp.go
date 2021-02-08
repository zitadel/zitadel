package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) BulkAddedHumanExternalIDP(ctx context.Context, userID, resourceOwner string, externalIDPs []*domain.ExternalIDP) error {
	if len(externalIDPs) == 0 {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Ek9s", "Errors.User.ExternalIDP.MinimumExternalIDPNeeded")
	}
	aggregates := make([]eventstore.Aggregater, len(externalIDPs))

	for i, externalIDP := range externalIDPs {
		externalIDPWriteModel := NewHumanExternalIDPWriteModel(userID, externalIDP.IDPConfigID, externalIDP.ExternalUserID, resourceOwner)
		userAgg := UserAggregateFromWriteModel(&externalIDPWriteModel.WriteModel)
		err := r.addHumanExternalIDP(ctx, userAgg, externalIDP)
		if err != nil {
			return err
		}
		aggregates[i] = userAgg
	}
	_, err := r.eventstore.PushAggregates(ctx, aggregates...)
	return err
}

func (r *CommandSide) addHumanExternalIDP(ctx context.Context, userAgg *user.Aggregate, externalIDP *domain.ExternalIDP) error {
	if !externalIDP.IsValid() {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6m9Kd", "Errors.User.ExternalIDP.Invalid")
	}
	//TODO: check if idpconfig exists
	userAgg.PushEvents(user.NewHumanExternalIDPAddedEvent(ctx, externalIDP.IDPConfigID, externalIDP.DisplayName, externalIDP.ExternalUserID))
	return nil
}

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
	return r.eventstore.PushAggregate(ctx, existingExternalIDP, userAgg)
}

func (r *CommandSide) HumanExternalLoginChecked(ctx context.Context, orgID, userID string, authRequest *domain.AuthRequest) (err error) {
	if userID == "" {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5n8sM", "Errors.IDMissing")
	}

	existingHuman, err := r.getHumanWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingHuman.UserState == domain.UserStateUnspecified || existingHuman.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-dn88J", "Errors.User.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingHuman.WriteModel)
	userAgg.PushEvents(user.NewHumanExternalIDPCheckSucceededEvent(ctx, authRequestDomainToAuthRequestInfo(authRequest)))
	return r.eventstore.PushAggregate(ctx, existingHuman, userAgg)
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
