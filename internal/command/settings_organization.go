package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SetSettingsOrganization struct {
	OrganizationID string

	UserUniqueness *bool
}

func (e *SetSettingsOrganization) IsValid() error {
	if e.OrganizationID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-TODO", "Errors.Org.Settings.Invalid")
	}
	return nil
}

func (c *Commands) SetSettingsOrganization(ctx context.Context, set *SetSettingsOrganization) (_ *domain.ObjectDetails, err error) {
	if err := set.IsValid(); err != nil {
		return nil, err
	}
	wm, err := c.getSettingsOrganizationWriteModelByID(ctx, set.OrganizationID)
	if err != nil {
		return nil, err
	}
	if !wm.OrganizationState.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-TODO", "Errors.NotFound")
	}

	events, err := wm.NewSet(ctx,
		set.UserUniqueness,
	)
	if err != nil {
		return nil, err
	}

	return c.pushAppendAndReduceDetails(ctx, wm, events...)
}

func (c *Commands) DeleteSettingsOrganization(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-TODO", "Errors.IDMissing")
	}
	wm, err := c.getSettingsOrganizationWriteModelByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !wm.State.Exists() {
		return writeModelToObjectDetails(wm.GetWriteModel()), nil
	}

	events, err := wm.NewRemoved(ctx)
	if err != nil {
		return nil, err
	}

	return c.pushAppendAndReduceDetails(ctx, wm, events...)
}

func (c *Commands) getSettingsOrganizationWriteModelByID(ctx context.Context, id string) (*SettingsOrganizationWriteModel, error) {
	wm := NewSettingsOrganizationWriteModel(id, c.checkPermission)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	return wm, nil
}
