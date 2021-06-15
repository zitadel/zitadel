package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
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
		return nil, caos_errs.ThrowNotFound(nil, "Org-GFdg2", "Errors.Org.IDPConfig.AlreadyExists")
	}

	iamAgg := OrgAggregateFromWriteModel(&existingConfig.WriteModel)
	changedEvent, hasChanged, err := existingConfig.NewChangedEvent(
		ctx,
		iamAgg,
		config.IDPConfigID,
		config.BaseURL,
		config.BackendConnectorID,
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
