package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
)

type OrgDependencies struct {
	Users    []string
	Projects []string
}

type OrgSetup struct {
	Name         string
	CustomDomain string
	Human        AddHuman
	Roles        []string
}

func (c *Commands) SetUpOrgWithIDs(ctx context.Context, o *OrgSetup, orgID, userID string, userIDs ...string) (string, *domain.ObjectDetails, error) {
	existingOrg, err := c.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return "", nil, err
	}
	if existingOrg != nil {
		return "", nil, errors.ThrowPreconditionFailed(nil, "COMMAND-poaj2", "Errors.Org.AlreadyExisting")
	}

	return c.setUpOrgWithIDs(ctx, o, orgID, userID, userIDs...)
}

func (c *Commands) setUpOrgWithIDs(ctx context.Context, o *OrgSetup, orgID, userID string, userIDs ...string) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(orgID)
	userAgg := user_repo.NewAggregate(userID, orgID)

	roles := []string{domain.RoleOrgOwner}
	if len(o.Roles) > 0 {
		roles = o.Roles
	}

	validations := []preparation.Validation{
		AddOrgCommand(ctx, orgAgg, o.Name, userIDs...),
		AddHumanCommand(userAgg, &o.Human, c.userPasswordAlg, c.userEncryption),
		c.AddOrgMemberCommand(orgAgg, userID, roles...),
	}
	if o.CustomDomain != "" {
		validations = append(validations, c.prepareAddOrgDomain(orgAgg, o.CustomDomain, userIDs))
	}

	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validations...)
	if err != nil {
		return "", nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return userID, &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: orgID,
	}, nil
}

func (c *Commands) SetUpOrg(ctx context.Context, o *OrgSetup, userIDs ...string) (string, *domain.ObjectDetails, error) {
	orgID, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}

	userID, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}

	return c.setUpOrgWithIDs(ctx, o, orgID, userID, userIDs...)
}

// AddOrgCommand defines the commands to create a new org,
// this includes the verified default domain
func AddOrgCommand(ctx context.Context, a *org.Aggregate, name string, userIDs ...string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if name = strings.TrimSpace(name); name == "" {
			return nil, errors.ThrowInvalidArgument(nil, "ORG-mruNY", "Errors.Invalid.Argument")
		}
		defaultDomain := domain.NewIAMDomainName(name, authz.GetInstance(ctx).RequestedDomain())
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				org.NewOrgAddedEvent(ctx, &a.Aggregate, name),
				org.NewDomainAddedEvent(ctx, &a.Aggregate, defaultDomain),
				org.NewDomainVerifiedEvent(ctx, &a.Aggregate, defaultDomain),
				org.NewDomainPrimarySetEvent(ctx, &a.Aggregate, defaultDomain),
			}, nil
		}, nil
	}
}

func (c *Commands) getOrg(ctx context.Context, orgID string) (*domain.Org, error) {
	writeModel, err := c.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if writeModel.State == domain.OrgStateUnspecified || writeModel.State == domain.OrgStateRemoved {
		return nil, errors.ThrowInternal(err, "COMMAND-4M9sf", "Errors.Org.NotFound")
	}
	return orgWriteModelToOrg(writeModel), nil
}

func (c *Commands) checkOrgExists(ctx context.Context, orgID string) error {
	orgWriteModel, err := c.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return err
	}
	if orgWriteModel.State == domain.OrgStateUnspecified || orgWriteModel.State == domain.OrgStateRemoved {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-QXPGs", "Errors.Org.NotFound")
	}
	return nil
}

func (c *Commands) AddOrgWithID(ctx context.Context, name, userID, resourceOwner, orgID string, claimedUserIDs []string) (*domain.Org, error) {
	existingOrg, err := c.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if existingOrg.State != domain.OrgStateUnspecified {
		return nil, errors.ThrowNotFound(nil, "ORG-lapo2m", "Errors.Org.AlreadyExisting")
	}

	return c.addOrgWithIDAndMember(ctx, name, userID, resourceOwner, orgID, claimedUserIDs)
}

func (c *Commands) AddOrg(ctx context.Context, name, userID, resourceOwner string, claimedUserIDs []string) (*domain.Org, error) {
	if name = strings.TrimSpace(name); name == "" {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-Mf9sd", "Errors.Org.Invalid")
	}

	orgID, err := c.idGenerator.Next()
	if err != nil {
		return nil, errors.ThrowInternal(err, "COMMA-OwciI", "Errors.Internal")
	}

	return c.addOrgWithIDAndMember(ctx, name, userID, resourceOwner, orgID, claimedUserIDs)
}

func (c *Commands) addOrgWithIDAndMember(ctx context.Context, name, userID, resourceOwner, orgID string, claimedUserIDs []string) (*domain.Org, error) {
	orgAgg, addedOrg, events, err := c.addOrgWithID(ctx, &domain.Org{Name: name}, orgID, claimedUserIDs)
	if err != nil {
		return nil, err
	}
	err = c.checkUserExists(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	addedMember := NewOrgMemberWriteModel(addedOrg.AggregateID, userID)
	orgMemberEvent, err := c.addOrgMember(ctx, orgAgg, addedMember, domain.NewMember(orgAgg.ID, userID, domain.RoleOrgOwner))
	if err != nil {
		return nil, err
	}
	events = append(events, orgMemberEvent)
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedOrg, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return orgWriteModelToOrg(addedOrg), nil
}

func (c *Commands) ChangeOrg(ctx context.Context, orgID, name string) (*domain.ObjectDetails, error) {
	name = strings.TrimSpace(name)
	if orgID == "" || name == "" {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-Mf9sd", "Errors.Org.Invalid")
	}

	orgWriteModel, err := c.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if orgWriteModel.State == domain.OrgStateUnspecified || orgWriteModel.State == domain.OrgStateRemoved {
		return nil, errors.ThrowNotFound(nil, "ORG-1MRds", "Errors.Org.NotFound")
	}
	if orgWriteModel.Name == name {
		return nil, errors.ThrowPreconditionFailed(nil, "ORG-4VSdf", "Errors.Org.NotChanged")
	}
	orgAgg := OrgAggregateFromWriteModel(&orgWriteModel.WriteModel)
	events := make([]eventstore.Command, 0)
	events = append(events, org.NewOrgChangedEvent(ctx, orgAgg, orgWriteModel.Name, name))
	changeDomainEvents, err := c.changeDefaultDomain(ctx, orgID, name)
	if err != nil {
		return nil, err
	}
	if len(changeDomainEvents) > 0 {
		events = append(events, changeDomainEvents...)
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(orgWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&orgWriteModel.WriteModel), nil
}

func (c *Commands) DeactivateOrg(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	orgWriteModel, err := c.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if orgWriteModel.State == domain.OrgStateUnspecified || orgWriteModel.State == domain.OrgStateRemoved {
		return nil, errors.ThrowNotFound(nil, "ORG-oL9nT", "Errors.Org.NotFound")
	}
	if orgWriteModel.State == domain.OrgStateInactive {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Dbs2g", "Errors.Org.AlreadyDeactivated")
	}
	orgAgg := OrgAggregateFromWriteModel(&orgWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewOrgDeactivatedEvent(ctx, orgAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(orgWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&orgWriteModel.WriteModel), nil
}

func (c *Commands) ReactivateOrg(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	orgWriteModel, err := c.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if orgWriteModel.State == domain.OrgStateUnspecified || orgWriteModel.State == domain.OrgStateRemoved {
		return nil, errors.ThrowNotFound(nil, "ORG-Dgf3g", "Errors.Org.NotFound")
	}
	if orgWriteModel.State == domain.OrgStateActive {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-bfnrh", "Errors.Org.AlreadyActive")
	}
	orgAgg := OrgAggregateFromWriteModel(&orgWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewOrgReactivatedEvent(ctx, orgAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(orgWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&orgWriteModel.WriteModel), nil
}

func (c *Commands) RemoveOrg(ctx context.Context, id string, dependencies *OrgDependencies) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(id)
	valdiations := []preparation.Validation{
		c.prepareRemoveOrg(orgAgg),
	}
	for _, user := range dependencies.Users {
		valdiations = append(valdiations, c.prepareUserOwnerRemoved(user_repo.NewAggregate(user, id)))
	}
	for _, depProject := range dependencies.Projects {
		valdiations = append(valdiations, c.prepareProjectOwnerRemoved(project.NewAggregate(depProject, id)))
	}

	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, valdiations...)
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}

	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) prepareRemoveOrg(a *org.Aggregate) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := c.getOrgWriteModelByID(ctx, a.ID)
			if err != nil {
				return nil, errors.ThrowPreconditionFailed(err, "COMMA-wG9p1", "Errors.Org.NotFound")
			}
			if writeModel.State == domain.OrgStateRemoved {
				return nil, errors.ThrowInvalidArgument(nil, "COMMA-pSAVZ", "Errors.NoChangesFound")
			}
			return []eventstore.Command{org.NewOrgRemovedEvent(ctx, &a.Aggregate, writeModel.Name)}, nil
		}, nil
	}
}

func ExistsOrg(ctx context.Context, filter preparation.FilterToQueryReducer, id string) (exists bool, err error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(id).
		OrderAsc().
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(id).
		EventTypes(
			org.OrgAddedEventType,
			org.OrgDeactivatedEventType,
			org.OrgReactivatedEventType,
			org.OrgRemovedEventType,
		).Builder())
	if err != nil {
		return false, err
	}

	for _, event := range events {
		switch event.(type) {
		case *org.OrgAddedEvent, *org.OrgReactivatedEvent:
			exists = true
		case *org.OrgDeactivatedEvent, *org.OrgRemovedEvent:
			exists = false
		}
	}

	return exists, nil
}

func (c *Commands) addOrgWithID(ctx context.Context, organisation *domain.Org, orgID string, claimedUserIDs []string) (_ *eventstore.Aggregate, _ *OrgWriteModel, _ []eventstore.Command, err error) {
	if !organisation.IsValid() {
		return nil, nil, nil, errors.ThrowInvalidArgument(nil, "COMM-deLSk", "Errors.Org.Invalid")
	}

	organisation.AggregateID = orgID
	organisation.AddIAMDomain(authz.GetInstance(ctx).RequestedDomain())
	addedOrg := NewOrgWriteModel(organisation.AggregateID)

	orgAgg := OrgAggregateFromWriteModel(&addedOrg.WriteModel)
	events := []eventstore.Command{
		org.NewOrgAddedEvent(ctx, orgAgg, organisation.Name),
	}
	for _, orgDomain := range organisation.Domains {
		orgDomainEvents, err := c.addOrgDomain(ctx, orgAgg, NewOrgDomainWriteModel(orgAgg.ID, orgDomain.Domain), orgDomain, claimedUserIDs)
		if err != nil {
			return nil, nil, nil, err
		}
		events = append(events, orgDomainEvents...)
	}
	return orgAgg, addedOrg, events, nil
}

func (c *Commands) getOrgWriteModelByID(ctx context.Context, orgID string) (*OrgWriteModel, error) {
	orgWriteModel := NewOrgWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, orgWriteModel)
	if err != nil {
		return nil, err
	}
	return orgWriteModel, nil
}
