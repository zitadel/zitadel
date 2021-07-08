package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) ChangeIDPAuthConnectorConfig(ctx context.Context, config *domain.AuthConnectorIDPConfig, resourceOwner string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-Dbwwg", "Errors.ResourceOwnerMissing")
	}
	if config.IDPConfigID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-GVbhg", "Errors.IDMissing")
	}

	existingConfig := NewOrgIDPAuthConnectorConfigWriteModel(config.IDPConfigID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, existingConfig)
	if err != nil {
		return nil, err
	}

	if existingConfig.State == domain.IDPConfigStateRemoved || existingConfig.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "Org-GFdg2", "Errors.Org.IDPConfig.NotFound")
	}

	machine, err := c.machineWriteModelByID(ctx, config.MachineID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(machine.UserState) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-GEgh2", "Errors.User.NotFound")
	}

	iamAgg := OrgAggregateFromWriteModel(&existingConfig.WriteModel)
	changedEvent, hasChanged, err := existingConfig.NewChangedEvent(
		ctx,
		iamAgg,
		config.IDPConfigID,
		config.BaseURL,
		config.ProviderID,
		config.MachineID,
	)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-H5u2j", "Errors.NoChanges")
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingConfig, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingConfig.AuthConnectorConfigWriteModel.WriteModel), nil
}

func (c *Commands) removeMachineUserFromAuthConnector(ctx context.Context, idpConfigID string, resourceOwner string) ([]eventstore.EventPusher, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-GBj52", "Errors.ResourceOwnerMissing")
	}
	if idpConfigID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-Grgh3", "Errors.IDMissing")
	}

	existingConfig := NewOrgIDPAuthConnectorConfigWriteModel(idpConfigID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, existingConfig)
	if err != nil {
		return nil, err
	}

	if existingConfig.State == domain.IDPConfigStateRemoved || existingConfig.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "Org-ADghn", "Errors.Org.IDPConfig.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingConfig.WriteModel)
	events := []eventstore.EventPusher{
		org.NewIDPAuthConnectorMachineUserRemovedEvent(ctx, orgAgg, existingConfig.IDPConfigID),
	}
	if existingConfig.State == domain.IDPConfigStateActive {
		events = append(events, org.NewIDPConfigDeactivatedEvent(ctx, orgAgg, existingConfig.IDPConfigID))
	}
	return events, nil
}
