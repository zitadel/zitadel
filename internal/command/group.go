package command

import (
	"context"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// Mostly required this one
func (c *Commands) AddGroupWithID(ctx context.Context, group *domain.Group, resourceOwner, groupID string) (_ *domain.Group, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-w8tnSoJxtn", "Errors.ResourceOwnerMissing")
	}
	if groupID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-nDXf5vXoUj", "Errors.IDMissing")
	}

	existingGroup, err := c.getGroupWriteModelByID(ctx, groupID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingGroup.State != domain.GroupStateUnspecified {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-opamwu", "Errors.Group.AlreadyExists")
	}
	group, err = c.addGroupWithID(ctx, group, resourceOwner, groupID)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (c *Commands) AddGroup(ctx context.Context, group *domain.Group, resourceOwner, ownerUserID string) (_ *domain.Group, err error) {
	if !group.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "GROUP-IOGDD", "Errors.Group.Invalid")
	}
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-emq7bpQY2s", "Errors.ResourceOwnerMissing")
	}
	// Check if the org exists
	orgWriteModel, err := c.getOrgWriteModelByID(ctx, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isOrgStateExists(orgWriteModel.State) {
		return nil, zerrors.ThrowNotFound(nil, "ORG-1MRds", "Errors.Org.NotFound")
	}
	if ownerUserID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-yf75Gl3Drp", "Errors.OwnerIDMissing")
	}

	groupID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}

	group, err = c.addGroupWithIDWithOwner(ctx, group, resourceOwner, ownerUserID, groupID)
	if err != nil {
		if zerrors.IsErrorAlreadyExists(err) {
			return nil, zerrors.ThrowAlreadyExists(err, "COMMAND-9f8sJ", "Errors.Group.AlreadyExists")
		}
		return nil, err
	}
	return group, nil
}

func (c *Commands) addGroupWithID(ctx context.Context, groupAdd *domain.Group, resourceOwner, groupID string) (_ *domain.Group, err error) {
	groupAdd.AggregateID = groupID
	addedGroup := NewGroupWriteModel(groupAdd.AggregateID, resourceOwner)
	groupAgg := GroupAggregateFromWriteModel(&addedGroup.WriteModel)

	events := []eventstore.Command{
		group.NewGroupAddedEvent(
			ctx,
			groupAgg,
			groupAdd.Name,
			groupAdd.Description),
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedGroup, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return groupWriteModelToGroup(addedGroup), nil
}

func (c *Commands) addGroupWithIDWithOwner(ctx context.Context, groupAdd *domain.Group, resourceOwner, ownerUserID, groupID string) (_ *domain.Group, err error) {
	if !groupAdd.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "GROUP-IOGDD", "Errors.Group.Invalid")
	}
	groupAdd.AggregateID = groupID
	addedGroup := NewGroupWriteModel(groupAdd.AggregateID, resourceOwner)
	groupAgg := GroupAggregateFromWriteModel(&addedGroup.WriteModel)

	events := []eventstore.Command{
		group.NewGroupAddedEvent(
			ctx,
			groupAgg,
			groupAdd.Name,
			groupAdd.Description),
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedGroup, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return groupWriteModelToGroup(addedGroup), nil
}

func AddGroupCommand(
	a *group.Aggregate,
	name string,
	owner string,
	description string,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if name = strings.TrimSpace(name); name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "GROUP-D02yo", "Errors.Invalid.Argument")
		}
		if owner == "" {
			return nil, zerrors.ThrowPreconditionFailed(nil, "GROUP-hxyoq", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				group.NewGroupAddedEvent(ctx, &a.Aggregate,
					name,
					description,
				),
			}, nil
		}, nil
	}
}

func groupWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer, groupID, resourceOwner string) (group *GroupWriteModel, err error) {
	group = NewGroupWriteModel(groupID, resourceOwner)
	events, err := filter(ctx, group.Query())
	if err != nil {
		return nil, err
	}

	group.AppendEvents(events...)
	if err := group.Reduce(); err != nil {
		return nil, err
	}

	return group, nil
}

func (c *Commands) getGroupByID(ctx context.Context, groupID, resourceOwner string) (_ *domain.Group, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	groupWriteModel, err := c.getGroupWriteModelByID(ctx, groupID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if groupWriteModel.State == domain.GroupStateUnspecified || groupWriteModel.State == domain.GroupStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "GROUP-Pc4ee", "Errors.Group.NotFound")
	}
	return groupWriteModelToGroup(groupWriteModel), nil
}

func (c *Commands) groupAggregateByID(ctx context.Context, groupID, resourceOwner string) (*eventstore.Aggregate, domain.GroupState, error) {
	result, err := c.groupState(ctx, groupID, resourceOwner)
	if err != nil {
		return nil, domain.GroupStateUnspecified, zerrors.ThrowNotFound(err, "COMMA-MDQqF", "Errors.Group.NotFound")
	}

	if len(result) == 0 {
		// _ = projection.GroupGrantFields.Trigger(ctx)
		result, err = c.groupState(ctx, groupID, resourceOwner)
		if err != nil || len(result) == 0 {
			return nil, domain.GroupStateUnspecified, zerrors.ThrowNotFound(err, "COMMA-V2Mza", "Errors.Group.NotFound")
		}
	}

	var state domain.GroupState
	err = result[0].Value.Unmarshal(&state)
	if err != nil {
		return nil, state, zerrors.ThrowNotFound(err, "COMMA-p1n7e", "Errors.Group.NotFound")
	}
	return &result[0].Aggregate, state, nil
}

func (c *Commands) groupState(ctx context.Context, groupID, resourceOwner string) ([]*eventstore.SearchResult, error) {
	return c.eventstore.Search(
		ctx,
		map[eventstore.FieldType]any{
			eventstore.FieldTypeObjectType:     group.GroupSearchType,
			eventstore.FieldTypeObjectID:       groupID,
			eventstore.FieldTypeObjectRevision: group.GroupObjectRevision,
			eventstore.FieldTypeFieldName:      group.GroupStateSearchField,
			eventstore.FieldTypeResourceOwner:  resourceOwner,
		},
	)
}

func (c *Commands) checkGroupExists(ctx context.Context, groupID, resourceOwner string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	_, state, err := c.groupAggregateByID(ctx, groupID, resourceOwner)
	if err != nil || !state.Valid() {
		// return zerrors.ThrowPreconditionFailed(err, "COMMA-UDnxD", "Errors.Group.NotFound")
		return zerrors.ThrowInvalidArgumentf(nil, "GROUP-UDnxD", "Errors.Group.NotFound %s,ResourceOwner: %s, GroupID: %s, Revision: %s", state, resourceOwner, groupID, group.GroupObjectRevision)

	}
	return nil
}

func (c *Commands) ChangeGroup(ctx context.Context, groupChange *domain.Group, resourceOwner string) (*domain.Group, error) {
	if !groupChange.IsValid() || groupChange.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-5n8vQ", "Errors.Group.Invalid")
	}

	existingGroup, err := c.getGroupWriteModelByID(ctx, groupChange.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingGroup.State == domain.GroupStateUnspecified || existingGroup.State == domain.GroupStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-4m9vS", "Errors.Group.NotFound")
	}

	groupAgg := GroupAggregateFromWriteModel(&existingGroup.WriteModel)
	changedEvent, hasChanged, err := existingGroup.NewChangedEvent(
		ctx,
		groupAgg,
		groupChange.Name,
		groupChange.Description)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-2M0fs", "Errors.NoChangesFound")
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGroup, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return groupWriteModelToGroup(existingGroup), nil
}

func (c *Commands) DeactivateGroup(ctx context.Context, groupID string, resourceOwner string) (*domain.ObjectDetails, error) {
	if groupID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-98iF1", "Errors.Group.GroupIDMissing")
	}

	groupAgg, state, err := c.groupAggregateByID(ctx, groupID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if state == domain.GroupStateUnspecified || state == domain.GroupStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-223N9", "Errors.Group.NotFound")
	}
	if state != domain.GroupStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-nle66", "Errors.Group.NotActive")
	}

	existingGroup, err := c.getGroupWriteModelByID(ctx, groupID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingGroup.State == domain.GroupStateUnspecified || existingGroup.State == domain.GroupStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-223N9", "Errors.Group.NotFound")
	}
	if existingGroup.State != domain.GroupStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-nlf66", "Errors.Group.NotActive")
	}

	pushedEvents, err := c.eventstore.Push(ctx, group.NewGroupDeactivatedEvent(ctx, groupAgg))
	if err != nil {
		return nil, err
	}

	return &domain.ObjectDetails{
		ResourceOwner: pushedEvents[0].Aggregate().ResourceOwner,
		Sequence:      pushedEvents[0].Sequence(),
		EventDate:     pushedEvents[0].CreatedAt(),
	}, nil
}

func (c *Commands) ReactivateGroup(ctx context.Context, groupID string, resourceOwner string) (*domain.ObjectDetails, error) {
	if groupID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4fjsG", "Errors.Group.GroupIDMissing")
	}

	groupAgg, state, err := c.groupAggregateByID(ctx, groupID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if state == domain.GroupStateUnspecified || state == domain.GroupStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-4N9sd", "Errors.Group.NotFound")
	}

	if state != domain.GroupStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-6n8cs", "Errors.Group.NotInactive")
	}

	existingGroup, err := c.getGroupWriteModelByID(ctx, groupID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingGroup.State == domain.GroupStateUnspecified || existingGroup.State == domain.GroupStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-4N9sd", "Errors.Group.NotFound")
	}
	if existingGroup.State != domain.GroupStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-6n8cs", "Errors.Group.NotInactive")
	}

	pushedEvents, err := c.eventstore.Push(ctx, group.NewGroupReactivatedEvent(ctx, groupAgg))
	if err != nil {
		return nil, err
	}

	return &domain.ObjectDetails{
		ResourceOwner: pushedEvents[0].Aggregate().ResourceOwner,
		Sequence:      pushedEvents[0].Sequence(),
		EventDate:     pushedEvents[0].CreatedAt(),
	}, nil
}

func (c *Commands) RemoveGroup(ctx context.Context, groupID, resourceOwner string, cascadingUserGrantIDs ...string) (*domain.ObjectDetails, error) {
	if groupID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-77Hn9", "Errors.Group.GroupIDMissing")
	}

	existingGroup, err := c.getGroupWriteModelByID(ctx, groupID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingGroup.State == domain.GroupStateUnspecified || existingGroup.State == domain.GroupStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-4N5sD", "Errors.Group.NotFound")
	}

	groupAgg := GroupAggregateFromWriteModel(&existingGroup.WriteModel)
	events := []eventstore.Command{
		group.NewGroupRemovedEvent(ctx, groupAgg, existingGroup.Name),
	}

	for _, grantID := range cascadingUserGrantIDs {
		event, _, err := c.removeUserGrant(ctx, grantID, "", true)
		if err != nil {
			logging.LogWithFields("COMMAND-b8Djf", "usergrantid", grantID).WithError(err).Warn("could not cascade remove user grant")
			continue
		}
		events = append(events, event)
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGroup, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGroup.WriteModel), nil
}

func (c *Commands) getGroupWriteModelByID(ctx context.Context, groupID, resourceOwner string) (_ *GroupWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	groupWriteModel := NewGroupWriteModel(groupID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, groupWriteModel)
	if err != nil {
		return nil, err
	}
	return groupWriteModel, nil
}
