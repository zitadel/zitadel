package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (r *CommandSide) BulkAddedHumanExternalIDP(ctx context.Context, userID, resourceOwner string, externalIDPs []*domain.ExternalIDP) (err error) {
	if len(externalIDPs) == 0 {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Ek9s", "Errors.User.ExternalIDP.MinimumExternalIDPNeeded")
	}

	events := make([]eventstore.EventPusher, len(externalIDPs))
	for i, externalIDP := range externalIDPs {
		externalIDPWriteModel := NewHumanExternalIDPWriteModel(userID, externalIDP.IDPConfigID, externalIDP.ExternalUserID, resourceOwner)
		userAgg := UserAggregateFromWriteModel(&externalIDPWriteModel.WriteModel)

		events[i], err = r.addHumanExternalIDP(ctx, userAgg, externalIDP)
		if err != nil {
			return err
		}
	}

	_, err = r.eventstore.PushEvents(ctx, events...)
	return err
}

func (r *CommandSide) addHumanExternalIDP(ctx context.Context, aggregate *eventstore.Aggregate, externalIDP *domain.ExternalIDP) (eventstore.EventPusher, error) {
	if !externalIDP.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6m9Kd", "Errors.User.ExternalIDP.Invalid")
	}
	//TODO: check if idpconfig exists
	return user.NewHumanExternalIDPAddedEvent(ctx, aggregate, externalIDP.IDPConfigID, externalIDP.DisplayName, externalIDP.ExternalUserID), nil
}

func (r *CommandSide) RemoveHumanExternalIDP(ctx context.Context, externalIDP *domain.ExternalIDP) error {
	event, err := r.removeHumanExternalIDP(ctx, externalIDP, false)
	if err != nil {
		return err
	}
	_, err = r.eventstore.PushEvents(ctx, event)
	return err
}

func (r *CommandSide) removeHumanExternalIDP(ctx context.Context, externalIDP *domain.ExternalIDP, cascade bool) (eventstore.EventPusher, error) {
	if externalIDP.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M9ds", "Errors.IDMissing")
	}

	existingExternalIDP, err := r.externalIDPWriteModelByID(ctx, externalIDP.AggregateID, externalIDP.IDPConfigID, externalIDP.ExternalUserID, externalIDP.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if existingExternalIDP.State == domain.ExternalIDPStateUnspecified || existingExternalIDP.State == domain.ExternalIDPStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-1M9xR", "Errors.User.ExternalIDP.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingExternalIDP.WriteModel)
	if cascade {
		return user.NewHumanExternalIDPCascadeRemovedEvent(ctx, userAgg, externalIDP.IDPConfigID, externalIDP.ExternalUserID), nil
	}
	return user.NewHumanExternalIDPRemovedEvent(ctx, userAgg, externalIDP.IDPConfigID, externalIDP.ExternalUserID), nil
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
	_, err = r.eventstore.PushEvents(ctx, user.NewHumanExternalIDPCheckSucceededEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
	return err
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
