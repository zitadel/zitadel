package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	repo "github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// SetGroupManagerRoles sets the ZITADEL manager roles every member of the group
// receives for the group's organization. An empty role list removes them.
// The caller needs the permission to manage organization members.
func (c *Commands) SetGroupManagerRoles(ctx context.Context, groupID string, roles []string) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if groupID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "CMDGMR-3jWqLs", "Errors.Group.Invalid")
	}
	if len(roles) > 0 && len(domain.CheckForInvalidRoles(roles, domain.OrgRolePrefix, c.zitadelRoles)) > 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "CMDGMR-9pXdVt", "Errors.Org.MemberInvalid")
	}

	group, err := c.checkGroupExists(ctx, groupID, nil)
	if err != nil {
		return nil, err
	}
	if err = c.checkPermissionUpdateOrgMember(ctx, group.ResourceOwner, group.ResourceOwner); err != nil {
		return nil, err
	}

	existingRoles, err := c.groupManagerRoles(ctx, groupID, group.ResourceOwner)
	if err != nil {
		return nil, err
	}
	sortedRoles := slices.Clone(roles)
	slices.Sort(sortedRoles)
	sortedExisting := slices.Clone(existingRoles)
	slices.Sort(sortedExisting)
	if slices.Equal(sortedRoles, sortedExisting) {
		return writeModelToObjectDetails(&group.WriteModel), nil
	}

	err = c.pushAppendAndReduce(ctx,
		group,
		repo.NewGroupManagerRolesSetEvent(ctx,
			GroupAggregateFromWriteModel(ctx, &group.WriteModel),
			roles,
		))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&group.WriteModel), nil
}

func (c *Commands) groupManagerRoles(ctx context.Context, groupID, resourceOwner string) ([]string, error) {
	wm := newGroupManagerRolesWriteModel(groupID, resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, wm); err != nil {
		return nil, err
	}
	return wm.roles, nil
}

type groupManagerRolesWriteModel struct {
	eventstore.WriteModel

	roles []string
}

func newGroupManagerRolesWriteModel(groupID, resourceOwner string) *groupManagerRolesWriteModel {
	return &groupManagerRolesWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   groupID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *groupManagerRolesWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(repo.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(repo.GroupManagerRolesSetEventType).
		Builder()
}

func (wm *groupManagerRolesWriteModel) Reduce() error {
	for _, event := range wm.Events {
		if e, ok := event.(*repo.GroupManagerRolesSetEvent); ok {
			wm.roles = e.Roles
		}
	}
	return wm.WriteModel.Reduce()
}
