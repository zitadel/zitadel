package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

func (c *Commands) ChangeDefaultIDPAuthConnectorConfig(ctx context.Context, config *domain.AuthConnectorIDPConfig) (*domain.ObjectDetails, error) {
	if config.IDPConfigID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-Bfsfg", "Errors.IDMissing")
	}
	existingConfig := NewIAMIDPAuthConnectorConfigWriteModel(config.IDPConfigID)
	err := c.eventstore.FilterToQueryReducer(ctx, existingConfig)
	if err != nil {
		return nil, err
	}

	if existingConfig.State == domain.IDPConfigStateRemoved || existingConfig.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-BVfwd", "Errors.IAM.IDPConfig.NotFound")
	}
	machine, err := c.machineWriteModelByID(ctx, config.MachineID, "")
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(machine.UserState) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-BGf31", "Errors.User.NotFound")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingConfig.WriteModel)
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-DVbvf", "Errors.IAM.LabelPolicy.NotChanged")
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
