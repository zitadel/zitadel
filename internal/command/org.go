package command

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command/preparation"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/repository/org"
	user_repo "github.com/caos/zitadel/internal/repository/user"
)

type OrgSetup struct {
	Name  string
	Human AddHuman
}

func (c *Commands) SetUpOrg(ctx context.Context, o *OrgSetup) (*domain.ObjectDetails, error) {
	orgID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	userID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	orgAgg := org.NewAggregate(orgID)
	userAgg := user_repo.NewAggregate(userID, orgID)

	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter,
		AddOrgCommand(ctx, orgAgg, o.Name),
		AddHumanCommand(userAgg, &o.Human, c.userPasswordAlg, c.smsEncryption, c.smtpEncryption, c.userEncryption),
		c.AddOrgMemberCommand(orgAgg, userID, domain.RoleOrgOwner),
	)
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
		ResourceOwner: orgID,
	}, nil
}

//AddOrgCommand defines the commands to create a new org,
// this includes the verified default domain
func AddOrgCommand(ctx context.Context, a *org.Aggregate, name string) preparation.Validation {
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
		return nil, caos_errs.ThrowInternal(err, "COMMAND-4M9sf", "Errors.Org.NotFound")
	}
	return orgWriteModelToOrg(writeModel), nil
}

func (c *Commands) checkOrgExists(ctx context.Context, orgID string) error {
	orgWriteModel, err := c.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return err
	}
	if orgWriteModel.State == domain.OrgStateUnspecified || orgWriteModel.State == domain.OrgStateRemoved {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0fs", "Errors.Org.NotFound")
	}
	return nil
}

func (c *Commands) AddOrg(ctx context.Context, name, userID, resourceOwner string, claimedUserIDs []string) (*domain.Org, error) {
	orgAgg, addedOrg, events, err := c.addOrg(ctx, &domain.Org{Name: name}, claimedUserIDs)
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
	if orgID == "" || name == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "EVENT-Mf9sd", "Errors.Org.Invalid")
	}
	orgWriteModel, err := c.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if orgWriteModel.State == domain.OrgStateUnspecified || orgWriteModel.State == domain.OrgStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-1MRds", "Errors.Org.NotFound")
	}
	if orgWriteModel.Name == name {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-4VSdf", "Errors.Org.NotChanged")
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
		return nil, caos_errs.ThrowNotFound(nil, "ORG-oL9nT", "Errors.Org.NotFound")
	}
	if orgWriteModel.State == domain.OrgStateInactive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Dbs2g", "Errors.Org.AlreadyDeactivated")
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
		return nil, caos_errs.ThrowNotFound(nil, "ORG-Dgf3g", "Errors.Org.NotFound")
	}
	if orgWriteModel.State == domain.OrgStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-bfnrh", "Errors.Org.AlreadyActive")
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

func (c *Commands) setUpOrg(
	ctx context.Context,
	organisation *domain.Org,
	admin *domain.Human,
	loginPolicy *domain.DomainPolicy,
	pwPolicy *domain.PasswordComplexityPolicy,
	initCodeGenerator crypto.Generator,
	phoneCodeGenerator crypto.Generator,
	claimedUserIDs []string,
	selfregistered bool,
) (orgAgg *eventstore.Aggregate, org *OrgWriteModel, human *HumanWriteModel, orgMember *OrgMemberWriteModel, events []eventstore.Command, err error) {
	orgAgg, orgWriteModel, addOrgEvents, err := c.addOrg(ctx, organisation, claimedUserIDs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	var userEvents []eventstore.Command
	if selfregistered {
		userEvents, human, err = c.registerHuman(ctx, orgAgg.ID, admin, nil, loginPolicy, pwPolicy, initCodeGenerator, phoneCodeGenerator)
	} else {
		userEvents, human, err = c.addHuman(ctx, orgAgg.ID, admin, loginPolicy, pwPolicy, initCodeGenerator, phoneCodeGenerator)
	}
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	addOrgEvents = append(addOrgEvents, userEvents...)

	addedMember := NewOrgMemberWriteModel(orgAgg.ID, human.AggregateID)
	orgMemberAgg := OrgAggregateFromWriteModel(&addedMember.WriteModel)
	orgMemberEvent, err := c.addOrgMember(ctx, orgMemberAgg, addedMember, domain.NewMember(orgMemberAgg.ID, human.AggregateID, domain.RoleOrgOwner))
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	addOrgEvents = append(addOrgEvents, orgMemberEvent)
	return orgAgg, orgWriteModel, human, addedMember, addOrgEvents, nil
}

func (c *Commands) addOrg(ctx context.Context, organisation *domain.Org, claimedUserIDs []string) (_ *eventstore.Aggregate, _ *OrgWriteModel, _ []eventstore.Command, err error) {
	if !organisation.IsValid() {
		return nil, nil, nil, caos_errs.ThrowInvalidArgument(nil, "COMM-deLSk", "Errors.Org.Invalid")
	}

	organisation.AggregateID, err = c.idGenerator.Next()
	if err != nil {
		return nil, nil, nil, caos_errs.ThrowInternal(err, "COMMA-OwciI", "Errors.Internal")
	}
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
