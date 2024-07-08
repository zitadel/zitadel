package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// InstanceOrgSetup is used for the first organisation in the instance setup.
// It used to be called OrgSetup, which now allows multiple Users, but it's used in the config.yaml and therefore
// a breaking change was not possible.
type InstanceOrgSetup struct {
	Name         string
	CustomDomain string
	Human        *AddHuman
	Machine      *AddMachine
	Roles        []string
}

type OrgSetup struct {
	Name         string
	CustomDomain string
	Admins       []*OrgSetupAdmin
}

// OrgSetupAdmin describes a user to be created (Human / Machine) or an existing (ID) to be used for an org setup.
type OrgSetupAdmin struct {
	ID      string
	Human   *AddHuman
	Machine *AddMachine
	Roles   []string
}

type orgSetupCommands struct {
	validations []preparation.Validation
	aggregate   *org.Aggregate
	commands    *Commands

	admins      []*OrgSetupAdmin
	pats        []*PersonalAccessToken
	machineKeys []*MachineKey
}

type CreatedOrg struct {
	ObjectDetails *domain.ObjectDetails
	CreatedAdmins []*CreatedOrgAdmin
}

type CreatedOrgAdmin struct {
	ID         string
	EmailCode  *string
	PhoneCode  *string
	PAT        *PersonalAccessToken
	MachineKey *MachineKey
}

func (c *Commands) setUpOrgWithIDs(ctx context.Context, o *OrgSetup, orgID string, allowInitialMail bool, userIDs ...string) (_ *CreatedOrg, err error) {
	cmds := c.newOrgSetupCommands(ctx, orgID, o)
	for _, admin := range o.Admins {
		if err = cmds.setupOrgAdmin(admin, allowInitialMail); err != nil {
			return nil, err
		}
	}
	if err = cmds.addCustomDomain(o.CustomDomain, userIDs); err != nil {
		return nil, err
	}

	return cmds.push(ctx)
}

func (c *Commands) newOrgSetupCommands(ctx context.Context, orgID string, orgSetup *OrgSetup) *orgSetupCommands {
	orgAgg := org.NewAggregate(orgID)
	validations := []preparation.Validation{
		AddOrgCommand(ctx, orgAgg, orgSetup.Name),
	}
	return &orgSetupCommands{
		validations: validations,
		aggregate:   orgAgg,
		commands:    c,
		admins:      orgSetup.Admins,
	}
}

func (c *orgSetupCommands) setupOrgAdmin(admin *OrgSetupAdmin, allowInitialMail bool) error {
	if admin.ID != "" {
		c.validations = append(c.validations, c.commands.AddOrgMemberCommand(c.aggregate, admin.ID, orgAdminRoles(admin.Roles)...))
		return nil
	}
	userID, err := id_generator.Next()
	if err != nil {
		return err
	}
	if admin.Human != nil {
		admin.Human.ID = userID
		c.validations = append(c.validations, c.commands.AddHumanCommand(admin.Human, c.aggregate.ID, c.commands.userPasswordHasher, c.commands.userEncryption, allowInitialMail))
	} else if admin.Machine != nil {
		admin.Machine.Machine.AggregateID = userID
		if err = c.setupOrgAdminMachine(c.aggregate, admin.Machine); err != nil {
			return err
		}
	}
	c.validations = append(c.validations, c.commands.AddOrgMemberCommand(c.aggregate, userID, orgAdminRoles(admin.Roles)...))
	return nil
}

func (c *orgSetupCommands) setupOrgAdminMachine(orgAgg *org.Aggregate, machine *AddMachine) error {
	userAgg := user.NewAggregate(machine.Machine.AggregateID, orgAgg.ID)
	c.validations = append(c.validations, AddMachineCommand(userAgg, machine.Machine))
	var pat *PersonalAccessToken
	var machineKey *MachineKey
	if machine.Pat != nil {
		pat = NewPersonalAccessToken(orgAgg.ID, machine.Machine.AggregateID, machine.Pat.ExpirationDate, machine.Pat.Scopes, domain.UserTypeMachine)
		tokenID, err := id_generator.Next()
		if err != nil {
			return err
		}
		pat.TokenID = tokenID
		c.pats = append(c.pats, pat)
		c.validations = append(c.validations, prepareAddPersonalAccessToken(pat, c.commands.keyAlgorithm))
	}
	if machine.MachineKey != nil {
		machineKey = NewMachineKey(orgAgg.ID, machine.Machine.AggregateID, machine.MachineKey.ExpirationDate, machine.MachineKey.Type)
		keyID, err := id_generator.Next()
		if err != nil {
			return err
		}
		machineKey.KeyID = keyID
		c.machineKeys = append(c.machineKeys, machineKey)
		c.validations = append(c.validations, prepareAddUserMachineKey(machineKey, c.commands.keySize))
	}
	return nil
}

func (c *orgSetupCommands) addCustomDomain(domain string, userIDs []string) error {
	if domain != "" {
		c.validations = append(c.validations, c.commands.prepareAddOrgDomain(c.aggregate, domain, userIDs))
	}
	return nil
}

func orgAdminRoles(roles []string) []string {
	if len(roles) > 0 {
		return roles
	}
	return []string{domain.RoleOrgOwner}
}

func (c *orgSetupCommands) push(ctx context.Context) (_ *CreatedOrg, err error) {
	cmds, err := preparation.PrepareCommands(ctx, c.commands.eventstore.Filter, c.validations...)
	if err != nil {
		return nil, err
	}

	events, err := c.commands.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}

	return &CreatedOrg{
		ObjectDetails: &domain.ObjectDetails{
			Sequence:      events[len(events)-1].Sequence(),
			EventDate:     events[len(events)-1].CreatedAt(),
			ResourceOwner: c.aggregate.ID,
		},
		CreatedAdmins: c.createdAdmins(),
	}, nil
}

func (c *orgSetupCommands) createdAdmins() []*CreatedOrgAdmin {
	users := make([]*CreatedOrgAdmin, 0, len(c.admins))
	for _, admin := range c.admins {
		if admin.ID != "" {
			continue
		}
		if admin.Human != nil {
			users = append(users, c.createdHumanAdmin(admin))
			continue
		}
		if admin.Machine != nil {
			users = append(users, c.createdMachineAdmin(admin))
		}
	}
	return users
}

func (c *orgSetupCommands) createdHumanAdmin(admin *OrgSetupAdmin) *CreatedOrgAdmin {
	createdAdmin := &CreatedOrgAdmin{
		ID: admin.Human.ID,
	}
	if admin.Human.EmailCode != nil {
		createdAdmin.EmailCode = admin.Human.EmailCode
	}
	return createdAdmin
}

func (c *orgSetupCommands) createdMachineAdmin(admin *OrgSetupAdmin) *CreatedOrgAdmin {
	createdAdmin := &CreatedOrgAdmin{
		ID: admin.Machine.Machine.AggregateID,
	}
	if admin.Machine.Pat != nil {
		for _, pat := range c.pats {
			if pat.AggregateID == createdAdmin.ID {
				createdAdmin.PAT = pat
			}
		}
	}
	if admin.Machine.MachineKey != nil {
		for _, key := range c.machineKeys {
			if key.AggregateID == createdAdmin.ID {
				createdAdmin.MachineKey = key
			}
		}
	}
	return createdAdmin
}

func (c *Commands) SetUpOrg(ctx context.Context, o *OrgSetup, allowInitialMail bool, userIDs ...string) (*CreatedOrg, error) {
	orgID, err := id_generator.Next()
	if err != nil {
		return nil, err
	}

	return c.setUpOrgWithIDs(ctx, o, orgID, allowInitialMail, userIDs...)
}

// AddOrgCommand defines the commands to create a new org,
// this includes the verified default domain
func AddOrgCommand(ctx context.Context, a *org.Aggregate, name string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if name = strings.TrimSpace(name); name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-mruNY", "Errors.Invalid.Argument")
		}
		defaultDomain, err := domain.NewIAMDomainName(name, authz.GetInstance(ctx).RequestedDomain())
		if err != nil {
			return nil, err
		}
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
	if !isOrgStateExists(writeModel.State) {
		return nil, zerrors.ThrowInternal(err, "COMMAND-4M9sf", "Errors.Org.NotFound")
	}
	return orgWriteModelToOrg(writeModel), nil
}

func (c *Commands) checkOrgExists(ctx context.Context, orgID string) error {
	orgWriteModel, err := c.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return err
	}
	if !isOrgStateExists(orgWriteModel.State) {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-QXPGs", "Errors.Org.NotFound")
	}
	return nil
}

func (c *Commands) AddOrgWithID(ctx context.Context, name, userID, resourceOwner, orgID string, claimedUserIDs []string) (_ *domain.Org, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	existingOrg, err := c.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if existingOrg.State != domain.OrgStateUnspecified {
		return nil, zerrors.ThrowNotFound(nil, "ORG-lapo2m", "Errors.Org.AlreadyExisting")
	}

	return c.addOrgWithIDAndMember(ctx, name, userID, resourceOwner, orgID, claimedUserIDs)
}

func (c *Commands) AddOrg(ctx context.Context, name, userID, resourceOwner string, claimedUserIDs []string) (*domain.Org, error) {
	if name = strings.TrimSpace(name); name == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "EVENT-Mf9sd", "Errors.Org.Invalid")
	}

	orgID, err := id_generator.Next()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "COMMA-OwciI", "Errors.Internal")
	}

	return c.addOrgWithIDAndMember(ctx, name, userID, resourceOwner, orgID, claimedUserIDs)
}

func (c *Commands) addOrgWithIDAndMember(ctx context.Context, name, userID, resourceOwner, orgID string, claimedUserIDs []string) (_ *domain.Org, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

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
		return nil, zerrors.ThrowInvalidArgument(nil, "EVENT-Mf9sd", "Errors.Org.Invalid")
	}

	orgWriteModel, err := c.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if !isOrgStateExists(orgWriteModel.State) {
		return nil, zerrors.ThrowNotFound(nil, "ORG-1MRds", "Errors.Org.NotFound")
	}
	if orgWriteModel.Name == name {
		return nil, zerrors.ThrowPreconditionFailed(nil, "ORG-4VSdf", "Errors.Org.NotChanged")
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
	if !isOrgStateExists(orgWriteModel.State) {
		return nil, zerrors.ThrowNotFound(nil, "ORG-oL9nT", "Errors.Org.NotFound")
	}
	if orgWriteModel.State == domain.OrgStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "EVENT-Dbs2g", "Errors.Org.AlreadyDeactivated")
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
	if !isOrgStateExists(orgWriteModel.State) {
		return nil, zerrors.ThrowNotFound(nil, "ORG-Dgf3g", "Errors.Org.NotFound")
	}
	if orgWriteModel.State == domain.OrgStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "EVENT-bfnrh", "Errors.Org.AlreadyActive")
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

func (c *Commands) RemoveOrg(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(id)

	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareRemoveOrg(orgAgg))
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}

	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) prepareRemoveOrg(a *org.Aggregate) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			instance := authz.GetInstance(ctx)
			if a.ID == instance.DefaultOrganisationID() {
				return nil, zerrors.ThrowPreconditionFailed(nil, "COMMA-wG9p1", "Errors.Org.DefaultOrgNotDeletable")
			}

			err := c.checkProjectExists(ctx, instance.ProjectID(), a.ID)
			// if there is no error, the ZITADEL project was found on the org to be deleted
			if err == nil {
				return nil, zerrors.ThrowPreconditionFailed(err, "COMMA-AF3JW", "Errors.Org.ZitadelOrgNotDeletable")
			}
			// "precondition failed" error means the project does not exist, return other errors
			if !zerrors.IsPreconditionFailed(err) {
				return nil, err
			}
			writeModel, err := c.getOrgWriteModelByID(ctx, a.ID)
			if err != nil {
				return nil, zerrors.ThrowPreconditionFailed(err, "COMMA-wG9p1", "Errors.Org.NotFound")
			}
			if !isOrgStateExists(writeModel.State) {
				return nil, zerrors.ThrowNotFound(nil, "COMMA-aps2n", "Errors.Org.NotFound")
			}

			domainPolicy, err := c.domainPolicyWriteModel(ctx, a.ID)
			if err != nil {
				return nil, err
			}
			usernames, err := OrgUsers(ctx, filter, a.ID)
			if err != nil {
				return nil, err
			}
			domains, err := OrgDomains(ctx, filter, a.ID)
			if err != nil {
				return nil, err
			}
			links, err := OrgUserIDPLinks(ctx, filter, a.ID)
			if err != nil {
				return nil, err
			}
			entityIds, err := OrgSamlEntityIDs(ctx, filter, a.ID)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{org.NewOrgRemovedEvent(ctx, &a.Aggregate, writeModel.Name, usernames, domainPolicy.UserLoginMustBeDomain, domains, links, entityIds)}, nil
		}, nil
	}
}

func OrgUserIDPLinks(ctx context.Context, filter preparation.FilterToQueryReducer, orgID string) ([]*domain.UserIDPLink, error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(orgID).
		OrderAsc().
		AddQuery().
		AggregateTypes(user.AggregateType).
		EventTypes(
			user.UserIDPLinkAddedType, user.UserIDPLinkRemovedType, user.UserIDPLinkCascadeRemovedType,
		).Builder())
	if err != nil {
		return nil, err
	}
	links := make([]*domain.UserIDPLink, 0)
	for _, event := range events {
		switch eventTyped := event.(type) {
		case *user.UserIDPLinkAddedEvent:
			links = append(links, &domain.UserIDPLink{
				IDPConfigID:    eventTyped.IDPConfigID,
				ExternalUserID: eventTyped.ExternalUserID,
				DisplayName:    eventTyped.DisplayName,
			})
		case *user.UserIDPLinkRemovedEvent:
			for i := range links {
				if links[i].ExternalUserID == eventTyped.ExternalUserID &&
					links[i].IDPConfigID == eventTyped.IDPConfigID {
					links[i] = links[len(links)-1]
					links[len(links)-1] = nil
					links = links[:len(links)-1]
					break
				}
			}

		case *user.UserIDPLinkCascadeRemovedEvent:
			for i := range links {
				if links[i].ExternalUserID == eventTyped.ExternalUserID &&
					links[i].IDPConfigID == eventTyped.IDPConfigID {
					links[i] = links[len(links)-1]
					links[len(links)-1] = nil
					links = links[:len(links)-1]
					break
				}
			}
		}
	}
	return links, nil
}

type samlEntityID struct {
	appID    string
	entityID string
}

func OrgSamlEntityIDs(ctx context.Context, filter preparation.FilterToQueryReducer, orgID string) ([]string, error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(orgID).
		OrderAsc().
		AddQuery().
		AggregateTypes(project.AggregateType).
		EventTypes(
			project.SAMLConfigAddedType, project.SAMLConfigChangedType, project.ApplicationRemovedType,
		).Builder())
	if err != nil {
		return nil, err
	}
	entityIDs := make([]samlEntityID, 0)
	for _, event := range events {
		switch eventTyped := event.(type) {
		case *project.SAMLConfigAddedEvent:
			entityIDs = append(entityIDs, samlEntityID{appID: eventTyped.AppID, entityID: eventTyped.EntityID})
		case *project.SAMLConfigChangedEvent:
			for i := range entityIDs {
				if entityIDs[i].appID == eventTyped.AppID {
					entityIDs[i].entityID = eventTyped.EntityID
					break
				}
			}
		case *project.ApplicationRemovedEvent:
			for i := range entityIDs {
				if entityIDs[i].appID == eventTyped.AppID {
					entityIDs[i] = entityIDs[len(entityIDs)-1]
					entityIDs = entityIDs[:len(entityIDs)-1]
					break
				}
			}
		}
	}
	ids := make([]string, len(entityIDs))
	for i := range entityIDs {
		ids[i] = entityIDs[i].entityID
	}
	return ids, nil
}

func OrgDomains(ctx context.Context, filter preparation.FilterToQueryReducer, orgID string) ([]string, error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(orgID).
		OrderAsc().
		AddQuery().
		AggregateTypes(org.AggregateType).
		EventTypes(
			org.OrgDomainVerifiedEventType,
			org.OrgDomainRemovedEventType,
		).Builder())
	if err != nil {
		return nil, err
	}
	names := make([]string, 0)
	for _, event := range events {
		switch eventTyped := event.(type) {
		case *org.DomainVerifiedEvent:
			names = append(names, eventTyped.Domain)
		case *org.DomainRemovedEvent:
			for i := range names {
				if names[i] == eventTyped.Domain {
					names[i] = names[len(names)-1]
					names = names[:len(names)-1]
					break
				}
			}
		}
	}
	return names, nil
}

type userIDName struct {
	name string
	id   string
}

func OrgUsers(ctx context.Context, filter preparation.FilterToQueryReducer, orgID string) ([]string, error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(authz.GetInstance(ctx).InstanceID()).
		ResourceOwner(orgID).
		OrderAsc().
		AddQuery().
		AggregateTypes(user.AggregateType).
		EventTypes(
			user.HumanAddedType,
			user.MachineAddedEventType,
			user.HumanRegisteredType,
			user.UserDomainClaimedType,
			user.UserUserNameChangedType,
			user.UserRemovedType,
		).Builder())
	if err != nil {
		return nil, err
	}

	users := make([]userIDName, 0)
	for _, event := range events {
		switch eventTyped := event.(type) {
		case *user.HumanAddedEvent:
			users = append(users, userIDName{eventTyped.UserName, eventTyped.Aggregate().ID})
		case *user.MachineAddedEvent:
			users = append(users, userIDName{eventTyped.UserName, eventTyped.Aggregate().ID})
		case *user.HumanRegisteredEvent:
			users = append(users, userIDName{eventTyped.UserName, eventTyped.Aggregate().ID})
		case *user.DomainClaimedEvent:
			for i := range users {
				if users[i].id == eventTyped.Aggregate().ID {
					users[i].name = eventTyped.UserName
				}
			}
		case *user.UsernameChangedEvent:
			for i := range users {
				if users[i].id == eventTyped.Aggregate().ID {
					users[i].name = eventTyped.UserName
				}
			}
		case *user.UserRemovedEvent:
			for i := range users {
				if users[i].id == eventTyped.Aggregate().ID {
					users[i] = users[len(users)-1]
					users = users[:len(users)-1]
					break
				}
			}
		}
	}
	names := make([]string, len(users))
	for i := range users {
		names[i] = users[i].name
	}
	return names, nil
}

func ExistsOrg(ctx context.Context, filter preparation.FilterToQueryReducer, id string) (exists bool, err error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(authz.GetInstance(ctx).InstanceID()).
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
		return nil, nil, nil, zerrors.ThrowInvalidArgument(nil, "COMM-deLSk", "Errors.Org.Invalid")
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

func (c *Commands) getOrgWriteModelByID(ctx context.Context, orgID string) (_ *OrgWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	orgWriteModel := NewOrgWriteModel(orgID)
	err = c.eventstore.FilterToQueryReducer(ctx, orgWriteModel)
	if err != nil {
		return nil, err
	}
	return orgWriteModel, nil
}

func isOrgStateExists(state domain.OrgState) bool {
	return !hasOrgState(state, domain.OrgStateRemoved, domain.OrgStateUnspecified)
}

func hasOrgState(check domain.OrgState, states ...domain.OrgState) bool {
	for _, state := range states {
		if check == state {
			return true
		}
	}
	return false
}
