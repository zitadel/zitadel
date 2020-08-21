package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	usr_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
)

func UserGrantByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ols34", "id should be filled")
	}
	return UserGrantQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func UserGrantQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.UserGrantAggregate).
		LatestSequenceFilter(latestSequence)
}

func UserGrantUniqueQuery(resourceOwner, projectID, userID string) *es_models.SearchQuery {
	grantID := resourceOwner + projectID + userID
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.UserGrantUniqueAggregate).
		AggregateIDFilter(grantID).
		OrderDesc().
		SetLimit(1)
}

func UserGrantAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, grant *model.UserGrant) (*es_models.Aggregate, error) {
	if grant == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dis83", "existing grant should not be nil")
	}
	return aggCreator.NewAggregate(ctx, grant.AggregateID, model.UserGrantAggregate, model.UserGrantVersion, grant.Sequence)
}

func UserGrantAddedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, grant *model.UserGrant) ([]*es_models.Aggregate, error) {
	agg, err := UserGrantAggregate(ctx, aggCreator, grant)
	if err != nil {
		return nil, err
	}
	validationQuery := es_models.NewSearchQuery().
		AggregateTypeFilter(usr_model.UserAggregate, org_es_model.OrgAggregate, proj_es_model.ProjectAggregate).
		AggregateIDsFilter(grant.UserID, authz.GetCtxData(ctx).OrgID, grant.ProjectID)

	validation := addUserGrantValidation(authz.GetCtxData(ctx).OrgID, grant)
	agg, err = agg.SetPrecondition(validationQuery, validation).AppendEvent(model.UserGrantAdded, grant)
	if err != nil {
		return nil, err
	}

	uniqueAggregate, err := reservedUniqueUserGrantAggregate(ctx, aggCreator, grant)
	if err != nil {
		return nil, err
	}
	return []*es_models.Aggregate{
		agg,
		uniqueAggregate,
	}, nil
}

func reservedUniqueUserGrantAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, grant *model.UserGrant) (*es_models.Aggregate, error) {
	grantID := authz.GetCtxData(ctx).OrgID + grant.ProjectID + grant.UserID
	aggregate, err := aggCreator.NewAggregate(ctx, grantID, model.UserGrantUniqueAggregate, model.UserGrantVersion, 0)
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.UserGrantReserved, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(UserGrantUniqueQuery(authz.GetCtxData(ctx).OrgID, grant.ProjectID, grant.UserID), isEventValidation(aggregate, model.UserGrantReserved)), nil
}

func releasedUniqueUserGrantAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, grant *model.UserGrant) (aggregate *es_models.Aggregate, err error) {
	grantID := grant.ResourceOwner + grant.ProjectID + grant.UserID
	aggregate, err = aggCreator.NewAggregate(ctx, grantID, model.UserGrantUniqueAggregate, model.UserGrantVersion, 0)
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.UserGrantReleased, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(UserGrantUniqueQuery(grant.ResourceOwner, grant.ProjectID, grant.UserID), isEventValidation(aggregate, model.UserGrantReleased)), nil
}

func UserGrantChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.UserGrant, grant *model.UserGrant, cascade bool) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-osl8x", "grant should not be nil")
		}
		agg, err := UserGrantAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.Changes(grant)
		if !cascade {
			return agg.AppendEvent(model.UserGrantChanged, changes)
		}
		return agg.AppendEvent(model.UserGrantCascadeChanged, changes)
	}
}

func UserGrantDeactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.UserGrant, grant *model.UserGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-lo21s", "grant should not be nil")
		}
		agg, err := UserGrantAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserGrantDeactivated, nil)
	}
}

func UserGrantReactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.UserGrant, grant *model.UserGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-mks34", "grant should not be nil")
		}
		agg, err := UserGrantAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserGrantReactivated, nil)
	}
}

func UserGrantRemovedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.UserGrant, grant *model.UserGrant, cascade bool) ([]*es_models.Aggregate, error) {
	if grant == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-lo21s", "grant should not be nil")
	}
	agg, err := UserGrantAggregate(ctx, aggCreator, existing)
	if err != nil {
		return nil, err
	}
	eventType := model.UserGrantRemoved
	if cascade {
		eventType = model.UserGrantCascadeRemoved
	}
	agg, err = agg.AppendEvent(eventType, nil)
	if err != nil {
		return nil, err
	}
	uniqueAggregate, err := releasedUniqueUserGrantAggregate(ctx, aggCreator, existing)
	if err != nil {
		return nil, err
	}
	return []*es_models.Aggregate{
		agg,
		uniqueAggregate,
	}, nil
}

func isEventValidation(aggregate *es_models.Aggregate, eventType es_models.EventType) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		if len(events) == 0 {
			aggregate.PreviousSequence = 0
			return nil
		}
		if events[0].Type == eventType {
			return errors.ThrowPreconditionFailedf(nil, "EVENT-eJQqe", "user_grant is already %v", eventType)
		}
		aggregate.PreviousSequence = events[0].Sequence
		return nil
	}
}

func addUserGrantValidation(resourceOwner string, grant *model.UserGrant) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		existsOrg := false
		existsUser := false
		project := new(proj_es_model.Project)
		for _, event := range events {
			switch event.AggregateType {
			case usr_model.UserAggregate:
				switch event.Type {
				case usr_model.UserAdded, usr_model.UserRegistered, usr_model.HumanAdded, usr_model.MachineAdded:
					existsUser = true
				case usr_model.UserRemoved:
					existsUser = false
				}
			case org_es_model.OrgAggregate:
				switch event.Type {
				case org_es_model.OrgAdded:
					existsOrg = true
				case org_es_model.OrgRemoved:
					existsOrg = false
				}
			case proj_es_model.ProjectAggregate:
				project.AppendEvent(event)
			}
		}
		if !existsOrg {
			return errors.ThrowPreconditionFailed(nil, "EVENT-3OfIm", "org doesn't exist")
		}
		if !existsUser {
			return errors.ThrowPreconditionFailed(nil, "EVENT-Sl8uS", "user doesn't exist")
		}
		return checkProjectConditions(resourceOwner, grant, project)
	}
}

//TODO: rethink this function i know it's ugly.
func checkProjectConditions(resourceOwner string, grant *model.UserGrant, project *proj_es_model.Project) error {
	if grant.ProjectID != project.AggregateID {
		return errors.ThrowInvalidArgument(nil, "EVENT-ixlMx", "project doesn't exist")
	}
	if project.State == int32(proj_model.ProjectStateRemoved) {
		return errors.ThrowPreconditionFailed(nil, "EVENT-Lxp0s", "project doesn't exist")
	}
	if resourceOwner == project.ResourceOwner {
		return checkIfProjectHasRoles(grant.RoleKeys, project.Roles)
	}

	if _, projectGrant := proj_es_model.GetProjectGrantByOrgID(project.Grants, resourceOwner); projectGrant != nil {
		return checkIfProjectGrantHasRoles(grant.RoleKeys, projectGrant.RoleKeys)
	}
	return nil
}

func checkIfProjectHasRoles(roles []string, existing []*proj_es_model.ProjectRole) error {
	for _, roleKey := range roles {
		if _, role := proj_es_model.GetProjectRole(existing, roleKey); role == nil {
			return errors.ThrowPreconditionFailedf(nil, "EVENT-Lxp0s", "project doesn't have role %v", roleKey)
		}
	}
	return nil
}

func checkIfProjectGrantHasRoles(roles []string, existing []string) error {
	roleExists := false
	for _, roleKey := range roles {
		for _, existingRoleKey := range existing {
			if roleKey == existingRoleKey {
				roleExists = true
				continue
			}
		}
		if !roleExists {
			return errors.ThrowPreconditionFailed(nil, "EVENT-LSpwi", "project grant doesn't have role")
		}
	}
	return nil
}
