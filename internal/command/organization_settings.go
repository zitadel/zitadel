package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SetOrganizationSettings struct {
	OrganizationID string

	OrganizationScopedUsernames *bool
}

func (e *SetOrganizationSettings) IsValid() error {
	if e.OrganizationID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-TODO", "Errors.Org.Settings.Invalid")
	}
	return nil
}

func (c *Commands) SetOrganizationSettings(ctx context.Context, set *SetOrganizationSettings) (_ *domain.ObjectDetails, err error) {
	if err := set.IsValid(); err != nil {
		return nil, err
	}
	wm, err := c.getOrganizationSettingsWriteModelByID(ctx, set.OrganizationID)
	if err != nil {
		return nil, err
	}
	if !wm.OrganizationState.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-TODO", "Errors.NotFound")
	}

	domainPolicy, err := c.domainPolicyWriteModel(ctx, wm.AggregateID)
	if err != nil {
		return nil, err
	}

	events, err := wm.NewSet(ctx,
		set.OrganizationScopedUsernames,
		domainPolicy.UserLoginMustBeDomain,
		c.getOrganizationScopedUsernames,
	)
	if err != nil {
		return nil, err
	}

	return c.pushAppendAndReduceDetails(ctx, wm, events...)
}

func (c *Commands) DeleteOrganizationSettings(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-TODO", "Errors.IDMissing")
	}
	wm, err := c.getOrganizationSettingsWriteModelByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !wm.State.Exists() {
		return writeModelToObjectDetails(wm.GetWriteModel()), nil
	}

	domainPolicy, err := c.domainPolicyWriteModel(ctx, wm.AggregateID)
	if err != nil {
		return nil, err
	}

	events, err := wm.NewRemoved(ctx,
		domainPolicy.UserLoginMustBeDomain,
		c.getOrganizationScopedUsernames,
	)
	if err != nil {
		return nil, err
	}

	return c.pushAppendAndReduceDetails(ctx, wm, events...)
}

func checkOrganizationScopedUsernames(ctx context.Context, filter preparation.FilterToQueryReducer, id string) (_ bool, err error) {
	wm := NewOrganizationSettingsWriteModel(id, nil)
	events, err := filter(ctx, wm.Query())
	if err != nil {
		return false, err
	}
	if len(events) == 0 {
		return false, nil
	}
	wm.AppendEvents(events...)
	err = wm.Reduce()

	if wm.State.Exists() && wm.OrganizationScopedUsernames {
		return true, nil
	}
	return false, nil
}

func (c *Commands) getOrganizationSettingsWriteModelByID(ctx context.Context, id string) (*OrganizationSettingsWriteModel, error) {
	wm := NewOrganizationSettingsWriteModel(id, c.checkPermission)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

func (c *Commands) checkOrganizationScopedUsernames(ctx context.Context, orgID string) (_ bool, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	wm, err := c.getOrganizationSettingsWriteModelByID(ctx, orgID)
	if err != nil {
		return false, err
	}

	if wm.State.Exists() && wm.OrganizationScopedUsernames {
		return true, nil
	}
	return false, nil
}

func (c *Commands) getOrganizationScopedUsernamesWriteModelByID(ctx context.Context, id string) (*OrganizationScopedUsernamesWriteModel, error) {
	wm := NewOrganizationScopedUsernamesWriteModel(id)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

func (c *Commands) getOrganizationScopedUsernames(ctx context.Context, id string) ([]string, error) {
	wm, err := c.getOrganizationScopedUsernamesWriteModelByID(ctx, id)
	if err != nil {
		return nil, err
	}
	usernames := make([]string, len(wm.Users))
	for i, user := range wm.Users {
		usernames[i] = user.username
	}
	return usernames, nil
}
