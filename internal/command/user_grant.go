package command

import (
	"context"
	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// AddUserGrant authorizes a user for a project with the given role keys.
// The project must be owned by or granted to the resourceOwner.
// If the resourceOwner is nil, the project must be owned by the project that belongs to usergrant.ProjectID.
func (c *Commands) AddUserGrant(ctx context.Context, usergrant *domain.UserGrant, resourceOwner *string, check PermissionCheck) (_ *domain.UserGrant, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	event, addedUserGrant, err := c.addUserGrant(ctx, usergrant, resourceOwner, check)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(addedUserGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return userGrantWriteModelToUserGrant(addedUserGrant), nil
}

func (c *Commands) addUserGrant(ctx context.Context, userGrant *domain.UserGrant, resourceOwner *string, check PermissionCheck) (command eventstore.Command, _ *UserGrantWriteModel, err error) {
	if !userGrant.IsValid() {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-kVfMa", "Errors.UserGrant.Invalid")
	}
	err = c.checkUserGrantPreCondition(ctx, userGrant, resourceOwner, check)
	if resourceOwner == nil {
		resourceOwner = &userGrant.ResourceOwner
	}
	if err != nil {
		return nil, nil, err
	}
	userGrant.AggregateID, err = c.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}

	addedUserGrant := NewUserGrantWriteModel(userGrant.AggregateID, *resourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&addedUserGrant.WriteModel)
	command = usergrant.NewUserGrantAddedEvent(
		ctx,
		userGrantAgg,
		userGrant.UserID,
		userGrant.ProjectID,
		userGrant.ProjectGrantID,
		userGrant.RoleKeys,
	)
	return command, addedUserGrant, nil
}

func (c *Commands) ChangeUserGrant(ctx context.Context, userGrant *domain.UserGrant, resourceOwner *string, ignoreUnchanged bool, check PermissionCheck) (_ *domain.UserGrant, err error) {
	event, changedUserGrant, err := c.changeUserGrant(ctx, userGrant, resourceOwner, false, ignoreUnchanged, check)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return userGrantWriteModelToUserGrant(changedUserGrant), nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(changedUserGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return userGrantWriteModelToUserGrant(changedUserGrant), nil
}

func (c *Commands) changeUserGrant(ctx context.Context, userGrant *domain.UserGrant, resourceOwner *string, cascade, ignoreUnchanged bool, check PermissionCheck) (_ eventstore.Command, _ *UserGrantWriteModel, err error) {
	if userGrant.AggregateID == "" {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-3M0sd", "Errors.UserGrant.Invalid")
	}
	existingUserGrant, err := c.userGrantWriteModelByID(ctx, userGrant.AggregateID, gu.Value(resourceOwner))
	if err != nil {
		return nil, nil, err
	}
	if check == nil {
		// checkUserGrantPreCondition checks the permission later when we have the necessary context information
		err = checkExplicitProjectPermission(ctx, existingUserGrant.ProjectGrantID, existingUserGrant.ProjectID)
		if err != nil {
			return nil, nil, err
		}
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.UserGrant.NotFound")
	}

	grantUnchanged := len(existingUserGrant.RoleKeys) == len(userGrant.RoleKeys)
	if grantUnchanged {
	outer:
		for _, reqKey := range userGrant.RoleKeys {
			for _, existingKey := range existingUserGrant.RoleKeys {
				if reqKey == existingKey {
					continue outer
				}
			}
			grantUnchanged = false
			break
		}
	}

	if grantUnchanged {
		if ignoreUnchanged {
			return nil, existingUserGrant, nil
		} else {
			return nil, nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Rs8fy", "Errors.UserGrant.NotChanged")
		}
	}
	userGrant.UserID = existingUserGrant.UserID
	userGrant.ProjectID = existingUserGrant.ProjectID
	userGrant.ProjectGrantID = existingUserGrant.ProjectGrantID
	userGrant.ResourceOwner = existingUserGrant.ResourceOwner

	err = c.checkUserGrantPreCondition(ctx, userGrant, &existingUserGrant.ResourceOwner, check)
	if err != nil {
		return nil, nil, err
	}

	changedUserGrant := NewUserGrantWriteModel(userGrant.AggregateID, userGrant.ResourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&changedUserGrant.WriteModel)

	if cascade {
		return usergrant.NewUserGrantCascadeChangedEvent(ctx, userGrantAgg, userGrant.RoleKeys), existingUserGrant, nil
	}
	return usergrant.NewUserGrantChangedEvent(ctx, userGrantAgg, existingUserGrant.UserID, userGrant.RoleKeys), existingUserGrant, nil
}

func (c *Commands) removeRoleFromUserGrant(ctx context.Context, userGrantID string, roleKeys []string, cascade bool) (_ eventstore.Command, err error) {
	existingUserGrant, err := c.userGrantWriteModelByID(ctx, userGrantID, "")
	if err != nil {
		return nil, err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.UserGrant.NotFound")
	}
	keyExists := false
	for i, key := range existingUserGrant.RoleKeys {
		for _, roleKey := range roleKeys {
			if key == roleKey {
				keyExists = true
				copy(existingUserGrant.RoleKeys[i:], existingUserGrant.RoleKeys[i+1:])
				existingUserGrant.RoleKeys[len(existingUserGrant.RoleKeys)-1] = ""
				existingUserGrant.RoleKeys = existingUserGrant.RoleKeys[:len(existingUserGrant.RoleKeys)-1]
				continue
			}
		}
	}
	if !keyExists {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5m8g9", "Errors.UserGrant.RoleKeyNotFound")
	}
	changedUserGrant := NewUserGrantWriteModel(userGrantID, existingUserGrant.ResourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&changedUserGrant.WriteModel)

	if cascade {
		return usergrant.NewUserGrantCascadeChangedEvent(ctx, userGrantAgg, existingUserGrant.RoleKeys), nil
	}

	return usergrant.NewUserGrantChangedEvent(ctx, userGrantAgg, existingUserGrant.UserID, existingUserGrant.RoleKeys), nil
}

func (c *Commands) DeactivateUserGrant(ctx context.Context, grantID, resourceOwner string) (objectDetails *domain.ObjectDetails, err error) {
	if grantID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-N2OhG", "Errors.UserGrant.IDMissing")
	}

	existingUserGrant, err := c.userGrantWriteModelByID(ctx, grantID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.UserGrant.NotFound")
	}
	if existingUserGrant.State != domain.UserGrantStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-1S9gx", "Errors.UserGrant.NotActive")
	}
	err = checkExplicitProjectPermission(ctx, existingUserGrant.ProjectGrantID, existingUserGrant.ProjectID)
	if err != nil {
		return nil, err
	}

	deactivateUserGrant := NewUserGrantWriteModel(grantID, resourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&deactivateUserGrant.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, usergrant.NewUserGrantDeactivatedEvent(ctx, userGrantAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUserGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUserGrant.WriteModel), nil
}

func (c *Commands) ReactivateUserGrant(ctx context.Context, grantID, resourceOwner string) (objectDetails *domain.ObjectDetails, err error) {
	if grantID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Qxy8v", "Errors.UserGrant.IDMissing")
	}

	existingUserGrant, err := c.userGrantWriteModelByID(ctx, grantID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Lp0gs", "Errors.UserGrant.NotFound")
	}
	if existingUserGrant.State != domain.UserGrantStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-1ML0v", "Errors.UserGrant.NotInactive")
	}
	err = checkExplicitProjectPermission(ctx, existingUserGrant.ProjectGrantID, existingUserGrant.ProjectID)
	if err != nil {
		return nil, err
	}
	deactivateUserGrant := NewUserGrantWriteModel(grantID, resourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&deactivateUserGrant.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, usergrant.NewUserGrantReactivatedEvent(ctx, userGrantAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUserGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUserGrant.WriteModel), nil
}

func (c *Commands) RemoveUserGrant(ctx context.Context, grantID, resourceOwner string) (objectDetails *domain.ObjectDetails, err error) {
	event, existingUserGrant, err := c.removeUserGrant(ctx, grantID, resourceOwner, false)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingUserGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingUserGrant.WriteModel), nil
}

func (c *Commands) BulkRemoveUserGrant(ctx context.Context, grantIDs []string, resourceOwner string) (err error) {
	if len(grantIDs) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-5M0sd", "Errors.UserGrant.IDMissing")
	}
	events := make([]eventstore.Command, len(grantIDs))
	for i, grantID := range grantIDs {
		event, _, err := c.removeUserGrant(ctx, grantID, resourceOwner, false)
		if err != nil {
			return err
		}
		events[i] = event
	}
	_, err = c.eventstore.Push(ctx, events...)
	return err
}

func (c *Commands) removeUserGrant(ctx context.Context, grantID, resourceOwner string, cascade bool) (_ eventstore.Command, writeModel *UserGrantWriteModel, err error) {
	if grantID == "" {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-J9sc5", "Errors.UserGrant.IDMissing")
	}

	existingUserGrant, err := c.userGrantWriteModelByID(ctx, grantID, resourceOwner)
	if err != nil {
		return nil, nil, err
	}
	if existingUserGrant.State == domain.UserGrantStateUnspecified || existingUserGrant.State == domain.UserGrantStateRemoved {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-1My0t", "Errors.UserGrant.NotFound")
	}
	if !cascade {
		err = checkExplicitProjectPermission(ctx, existingUserGrant.ProjectGrantID, existingUserGrant.ProjectID)
		if err != nil {
			return nil, nil, err
		}
	}

	removeUserGrant := NewUserGrantWriteModel(grantID, existingUserGrant.ResourceOwner)
	userGrantAgg := UserGrantAggregateFromWriteModel(&removeUserGrant.WriteModel)
	if cascade {
		return usergrant.NewUserGrantCascadeRemovedEvent(
			ctx,
			userGrantAgg,
			existingUserGrant.UserID,
			existingUserGrant.ProjectID,
			existingUserGrant.ProjectGrantID), existingUserGrant, nil
	}
	return usergrant.NewUserGrantRemovedEvent(
		ctx,
		userGrantAgg,
		existingUserGrant.UserID,
		existingUserGrant.ProjectID,
		existingUserGrant.ProjectGrantID), existingUserGrant, nil
}

func (c *Commands) userGrantWriteModelByID(ctx context.Context, userGrantID, resourceOwner string) (writeModel *UserGrantWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserGrantWriteModel(userGrantID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) checkUserGrantPreCondition(ctx context.Context, usergrant *domain.UserGrant, resourceOwner *string, check PermissionCheck) (err error) {
	if !authz.GetFeatures(ctx).ShouldUseImprovedPerformance(feature.ImprovedPerformanceTypeUserGrant) {
		return c.checkUserGrantPreConditionOld(ctx, usergrant, resourceOwner, check)
	}

	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err := c.checkUserExists(ctx, usergrant.UserID, ""); err != nil {
		return err
	}
	if usergrant.ProjectGrantID != "" || resourceOwner == nil {
		projectOwner, grantID, err := c.searchProjectOwnerAndGrantID(ctx, usergrant.ProjectID, gu.Value(resourceOwner))
		if err != nil {
			return err
		}
		if resourceOwner == nil {
			resourceOwner = &projectOwner
		}
		if usergrant.ProjectGrantID == "" {
			usergrant.ProjectGrantID = grantID
		}
	}
	existingRoleKeys, err := c.searchUserGrantPreConditionState(ctx, usergrant, *resourceOwner)
	if check != nil {
		if usergrant.ProjectGrantID != "" {
			grantErr := check(*resourceOwner, usergrant.ProjectGrantID)
			if grantErr != nil {
				projectErr := check(*resourceOwner, usergrant.ProjectID)
				if projectErr != nil {
					return grantErr
				}
			}
		} else {
			projectErr := check(*resourceOwner, usergrant.ProjectID)
			if projectErr != nil {
				return projectErr
			}
		}
	}
	if err != nil {
		return err
	}
	if usergrant.HasInvalidRoles(existingRoleKeys) {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-mm9F4", "Errors.Project.Role.NotFound")
	}
	return nil
}

// this code needs to be rewritten anyways as soon as we improved the fields handling
//
//nolint:gocognit
func (c *Commands) searchUserGrantPreConditionState(ctx context.Context, userGrant *domain.UserGrant, resourceOwner string) (existingRoleKeys []string, err error) {
	criteria := []map[eventstore.FieldType]any{
		// project state query
		{
			eventstore.FieldTypeAggregateType: project.AggregateType,
			eventstore.FieldTypeAggregateID:   userGrant.ProjectID,
			eventstore.FieldTypeFieldName:     project.ProjectStateSearchField,
			eventstore.FieldTypeObjectType:    project.ProjectSearchType,
		},
		// granted org query
		{
			eventstore.FieldTypeAggregateType: org.AggregateType,
			eventstore.FieldTypeAggregateID:   resourceOwner,
			eventstore.FieldTypeFieldName:     org.OrgStateSearchField,
			eventstore.FieldTypeObjectType:    org.OrgSearchType,
		},
	}
	if userGrant.ProjectGrantID != "" {
		criteria = append(criteria, map[eventstore.FieldType]any{
			eventstore.FieldTypeAggregateType: project.AggregateType,
			eventstore.FieldTypeAggregateID:   userGrant.ProjectID,
			eventstore.FieldTypeObjectType:    project.ProjectGrantSearchType,
			eventstore.FieldTypeObjectID:      userGrant.ProjectGrantID,
		})
	} else {
		criteria = append(criteria, map[eventstore.FieldType]any{
			eventstore.FieldTypeAggregateType: project.AggregateType,
			eventstore.FieldTypeAggregateID:   userGrant.ProjectID,
			eventstore.FieldTypeObjectType:    project.ProjectRoleSearchType,
			eventstore.FieldTypeFieldName:     project.ProjectRoleKeySearchField,
		})
	}
	results, err := c.eventstore.Search(ctx, criteria...)
	if err != nil {
		return nil, err
	}

	var (
		existsProject    bool
		existsGrantedOrg bool
		existsGrant      bool
	)

	for _, result := range results {
		switch result.Object.Type {
		case project.ProjectRoleSearchType:
			var role string
			err := result.Value.Unmarshal(&role)
			if err != nil {
				return nil, err
			}
			existingRoleKeys = append(existingRoleKeys, role)
		case org.OrgSearchType:
			var state domain.OrgState
			err := result.Value.Unmarshal(&state)
			if err != nil {
				return nil, err
			}
			existsGrantedOrg = state.Valid() && state != domain.OrgStateRemoved
		case project.ProjectSearchType:
			var state domain.ProjectState
			err := result.Value.Unmarshal(&state)
			if err != nil {
				return nil, err
			}
			existsProject = state.Valid() && state != domain.ProjectStateRemoved
		case project.ProjectGrantSearchType:
			switch result.FieldName {
			case project.ProjectGrantGrantedOrgIDSearchField:
				var orgID string
				err := result.Value.Unmarshal(&orgID)
				if err != nil || orgID != resourceOwner {
					return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-3m9gg", "Errors.Org.NotFound")
				}
			case project.ProjectGrantStateSearchField:
				var state domain.ProjectGrantState
				err := result.Value.Unmarshal(&state)
				if err != nil {
					return nil, err
				}
				existsGrant = state.Valid() && state != domain.ProjectGrantStateRemoved
			case project.ProjectGrantRoleKeySearchField:
				var role string
				err := result.Value.Unmarshal(&role)
				if err != nil {
					return nil, err
				}
				existingRoleKeys = append(existingRoleKeys, role)
			case project.ProjectGrantGrantIDSearchField:
				var grantID string
				err := result.Value.Unmarshal(&grantID)
				if err != nil || grantID != userGrant.ProjectGrantID {
					return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-huvKF", "Errors.Project.Grant.NotFound")
				}
			}
		}
	}

	if !existsProject {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-m9gsd", "Errors.Project.NotFound")
	}
	if !existsGrantedOrg {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-3m9gg", "Errors.Org.NotFound")
	}
	if userGrant.ProjectGrantID != "" && !existsGrant {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-huvKF", "Errors.Project.Grant.NotFound")
	}
	return existingRoleKeys, nil
}

func (c *Commands) checkUserGrantPreConditionOld(ctx context.Context, usergrant *domain.UserGrant, resourceOwner *string, check PermissionCheck) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	preConditions := NewUserGrantPreConditionReadModel(usergrant.UserID, usergrant.ProjectID, usergrant.ProjectGrantID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, preConditions)
	if err != nil {
		return err
	}
	if usergrant.ResourceOwner == "" {
		usergrant.ResourceOwner = preConditions.ProjectResourceOwner
	}
	if usergrant.ProjectGrantID == "" {
		usergrant.ProjectGrantID = preConditions.ProjectGrantID
	}
	if !preConditions.UserExists {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-4f8sg", "Errors.User.NotFound")
	}
	projectIsOwned := resourceOwner == nil || *resourceOwner == preConditions.ProjectResourceOwner
	if projectIsOwned && !preConditions.ProjectExists {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-3n77S", "Errors.Project.NotFound")
	}
	if !projectIsOwned && !preConditions.ProjectGrantExists {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-4m9ff", "Errors.Project.Grant.NotFound")
	}
	if usergrant.HasInvalidRoles(preConditions.ExistingRoleKeys) {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-mm9F4", "Errors.Project.Role.NotFound")
	}
	if check != nil {
		if usergrant.ProjectGrantID != "" {
			grantErr := check(usergrant.ResourceOwner, usergrant.ProjectGrantID)
			if grantErr != nil {
				projectErr := check(usergrant.ResourceOwner, usergrant.ProjectID)
				if projectErr != nil {
					return grantErr
				}
			}
		} else {
			projectErr := check(usergrant.ResourceOwner, usergrant.ProjectID)
			if projectErr != nil {
				return projectErr
			}
		}
	}
	return nil
}

func (c *Commands) searchProjectOwnerAndGrantID(ctx context.Context, projectID string, grantedOrgID string) (projectOwner string, grantID string, err error) {
	grantIDQuery := map[eventstore.FieldType]any{
		eventstore.FieldTypeAggregateType: project.AggregateType,
		eventstore.FieldTypeAggregateID:   projectID,
		eventstore.FieldTypeFieldName:     project.ProjectGrantGrantedOrgIDSearchField,
	}
	if grantedOrgID != "" {
		grantIDQuery[eventstore.FieldTypeValue] = grantedOrgID
		grantIDQuery[eventstore.FieldTypeObjectType] = project.ProjectGrantSearchType

	}
	results, err := c.eventstore.Search(ctx, grantIDQuery)
	if err != nil {
		return "", "", err
	}
	for _, result := range results {
		projectOwner = result.Aggregate.ResourceOwner
		if grantedOrgID != "" && grantedOrgID == projectOwner {
			return projectOwner, "", nil
		}
		if result.Object.Type == project.ProjectGrantSearchType {
			return projectOwner, result.Object.ID, nil
		}
	}
	return projectOwner, grantID, err
}
