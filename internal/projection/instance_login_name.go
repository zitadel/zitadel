package projection

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

var _ Projection = (*InstanceLoginNames)(nil)

func NewInstanceLoginNames(instance string) *InstanceLoginNames {
	return &InstanceLoginNames{
		instance: instance,
	}
}

func NewInstanceLoginNamesWithOwner(instance, owner, loginName string) *InstanceLoginNames {
	return &InstanceLoginNames{
		instance:  instance,
		loginName: loginName,
		Orgs: []*OrgLoginNames{
			{
				org:        owner,
				LoginNames: make(map[string]*UserLoginNames),
			},
		},
	}
}

type InstanceLoginNames struct {
	instance  string
	loginName string
	owner     string

	policy  loginNamePolicy
	removed bool

	Orgs []*OrgLoginNames
}

type OrgLoginNames struct {
	org string
	// LoginNames per user
	LoginNames map[string]*UserLoginNames

	policy  loginNamePolicy
	domains []*loginNameDomain
}

func (ln *InstanceLoginNames) Build(ctx context.Context, es *eventstore.Eventstore) ([]*UserLoginNames, error) {
	usernameQuery := ln.usernameQuery(ctx)
	events, err := es.Filter(ctx, usernameQuery)
	if err != nil {
		return nil, err
	}
	ln.reduceUsernameEvents(events)

	instanceQuery := ln.instanceQuery(ctx)
	events, err = es.Filter(ctx, instanceQuery)
	if err != nil {
		return nil, err
	}
	ln.reduceInstanceEvents(events)

	orgQuery := ln.orgQuery(ctx)
	events, err = es.Filter(ctx, orgQuery)
	if err != nil {
		return nil, err
	}
	ln.reduceOrgEvents(events)
	return ln.generate(), nil
}

func (ln *InstanceLoginNames) Reduce(events []eventstore.Event) {}

func (ln *InstanceLoginNames) SearchQuery(ctx context.Context) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(ln.instance).
		OrderAsc().
		// ResourceOwner(ln.owner).
		AddQuery().
		AggregateTypes(
			user.AggregateType,
		).
		EventTypes(
			user.UserV1AddedType,
			user.HumanAddedType,
			user.HumanRegisteredType,
			user.UserV1RegisteredType,
			user.MachineAddedEventType,
			user.UserRemovedType,
			user.UserUserNameChangedType,
			user.UserDomainClaimedType,
		).
		Or().
		AggregateTypes(org.AggregateType).
		// AggregateIDs(ln.owner).
		EventTypes(
			org.DomainPolicyAddedEventType,
			org.DomainPolicyChangedEventType,
			org.DomainPolicyRemovedEventType,
			org.OrgDomainPrimarySetEventType,
			org.OrgDomainRemovedEventType,
			org.OrgDomainVerifiedEventType,
		).
		Or().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(ln.instance).
		EventTypes(
			instance.DomainPolicyAddedEventType,
			instance.DomainPolicyChangedEventType,
			instance.InstanceRemovedEventType,
		).
		Builder()
}

func (ln *InstanceLoginNames) instanceQuery(ctx context.Context) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(ln.instance).
		OrderAsc().
		// ResourceOwner(ln.owner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(ln.instance).
		EventTypes(
			instance.DomainPolicyAddedEventType,
			instance.DomainPolicyChangedEventType,
			instance.InstanceRemovedEventType,
		).
		Builder()
}

func (ln *InstanceLoginNames) reduceInstanceEvents(events []eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.DomainPolicyAddedEvent:
			ln.reduceInstanceDomainPolicyAddedEvent(e)
		case *instance.DomainPolicyChangedEvent:
			ln.reduceInstanceDomainPolicyChangedEvent(e)
		case *instance.InstanceRemovedEvent:
			ln.reduceInstanceRemovedEvent(e)
		}
	}
}

func (ln *InstanceLoginNames) usernameQuery(ctx context.Context) *eventstore.SearchQueryBuilder {
	var username map[string]interface{}
	if ln.loginName != "" {
		username = map[string]interface{}{
			"userName": strings.Split(ln.loginName, "@")[0],
		}
	}
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(ln.instance).
		OrderAsc().
		ResourceOwner(ln.owner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		EventTypes(
			user.UserV1AddedType,
			user.HumanAddedType,
			user.HumanRegisteredType,
			user.UserV1RegisteredType,
			user.MachineAddedEventType,
			user.UserRemovedType,
			user.UserUserNameChangedType,
			user.UserDomainClaimedType,
		).EventData(username).
		Builder()
}

func (ln *InstanceLoginNames) reduceUsernameEvents(events []eventstore.Event) {
	for _, event := range events {
		org := ln.org(event)
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			org.LoginNames[e.Aggregate().ID] = &UserLoginNames{
				UserID:     e.Aggregate().ID,
				InstanceID: ln.instance,
				OwnerID:    org.org,
				username:   e.UserName,
			}
		case *user.HumanRegisteredEvent:
			org.LoginNames[e.Aggregate().ID] = &UserLoginNames{
				UserID:     e.Aggregate().ID,
				InstanceID: ln.instance,
				OwnerID:    org.org,
				username:   e.UserName,
			}
		case *user.MachineAddedEvent:
			org.LoginNames[e.Aggregate().ID] = &UserLoginNames{
				UserID:     e.Aggregate().ID,
				InstanceID: ln.instance,
				OwnerID:    org.org,
				username:   e.UserName,
			}
		case *user.UserRemovedEvent:
			org.LoginNames[e.Aggregate().ID] = nil
			delete(org.LoginNames, e.Aggregate().ID)
		case *user.UsernameChangedEvent:
			if _, ok := org.LoginNames[e.Aggregate().ID]; !ok {
				org.LoginNames[e.Aggregate().ID] = new(UserLoginNames)
			}
			org.LoginNames[e.Aggregate().ID].username = e.UserName
		case *user.DomainClaimedEvent:
			if _, ok := org.LoginNames[e.Aggregate().ID]; !ok {
				org.LoginNames[e.Aggregate().ID] = new(UserLoginNames)
			}
			org.LoginNames[e.Aggregate().ID].username = e.UserName
		}
	}
}

func (ln *InstanceLoginNames) orgQuery(ctx context.Context) *eventstore.SearchQueryBuilder {
	ids := make([]string, len(ln.Orgs))
	for i, org := range ln.Orgs {
		ids[i] = org.org
	}
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(ln.instance).
		OrderAsc().
		// ResourceOwner(ln.owner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(ids...).
		EventTypes(
			org.DomainPolicyAddedEventType,
			org.DomainPolicyChangedEventType,
			org.DomainPolicyRemovedEventType,
			org.OrgDomainPrimarySetEventType,
			org.OrgDomainRemovedEventType,
			org.OrgDomainVerifiedEventType,
		).
		Builder()
}

func (ln *InstanceLoginNames) reduceOrgEvents(events []eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.DomainPolicyAddedEvent:
			ln.reduceOrgDomainPolicyAddedEvent(e)
		case *org.DomainPolicyChangedEvent:
			ln.reduceOrgDomainPolicyChangedEvent(e)
		case *org.DomainPolicyRemovedEvent:
			ln.reduceOrgRemovedEvent(e)
		case *org.DomainPrimarySetEvent:
			ln.reduceOrgDomainPrimarySetEvent(e)
		case *org.DomainRemovedEvent:
			ln.reduceOrgDomainRemovedEvent(e)
		case *org.DomainVerifiedEvent:
			ln.reduceOrgDomainVerifiedEvent(e)
		}
	}
}

func (ln *InstanceLoginNames) reduceInstanceDomainPolicyAddedEvent(event *instance.DomainPolicyAddedEvent) {
	ln.policy.mustBeDomain = event.UserLoginMustBeDomain
}

func (ln *InstanceLoginNames) reduceInstanceDomainPolicyChangedEvent(event *instance.DomainPolicyChangedEvent) {
	if event.UserLoginMustBeDomain == nil {
		return
	}
	ln.policy.mustBeDomain = *event.UserLoginMustBeDomain
}

func (ln *InstanceLoginNames) reduceInstanceRemovedEvent(event *instance.InstanceRemovedEvent) {
	ln.removed = true
}

func (ln *InstanceLoginNames) org(event eventstore.Event) *OrgLoginNames {
	for _, org := range ln.Orgs {
		if org.org == event.Aggregate().ResourceOwner {
			return org
		}
	}
	org := &OrgLoginNames{
		org:        event.Aggregate().ResourceOwner,
		LoginNames: make(map[string]*UserLoginNames),
	}

	ln.Orgs = append(ln.Orgs, org)

	return org
}

func (ln *InstanceLoginNames) reduceOrgDomainPolicyAddedEvent(event *org.DomainPolicyAddedEvent) {
	org := ln.org(event)
	org.policy = loginNamePolicy{
		mustBeDomain: event.UserLoginMustBeDomain,
	}
}

func (ln *InstanceLoginNames) reduceOrgDomainPolicyChangedEvent(event *org.DomainPolicyChangedEvent) {
	if event.UserLoginMustBeDomain == nil {
		return
	}
	org := ln.org(event)
	org.policy = loginNamePolicy{
		mustBeDomain: *event.UserLoginMustBeDomain,
	}
}

func (ln *InstanceLoginNames) reduceOrgRemovedEvent(event *org.DomainPolicyRemovedEvent) {
	for i, org := range ln.Orgs {
		if org.org != event.Aggregate().ID {
			continue
		}
		ln.Orgs[i] = ln.Orgs[len(ln.Orgs)-1]
		ln.Orgs[len(ln.Orgs)-1] = nil
		ln.Orgs = ln.Orgs[:len(ln.Orgs)-1]
		return
	}
}

func (ln *InstanceLoginNames) reduceOrgDomainPrimarySetEvent(event *org.DomainPrimarySetEvent) {
	org := ln.org(event)

	for _, domain := range org.domains {
		domain.isPrimary = domain.name == event.Domain
	}
}

func (ln *InstanceLoginNames) reduceOrgDomainRemovedEvent(event *org.DomainRemovedEvent) {
	org := ln.org(event)

	for i, domain := range org.domains {
		if domain.name != event.Domain {
			continue
		}
		org.domains[i] = org.domains[len(org.domains)-1]
		org.domains[len(org.domains)-1] = nil
		org.domains = org.domains[:len(org.domains)-1]
	}
}

func (ln *InstanceLoginNames) reduceOrgDomainVerifiedEvent(event *org.DomainVerifiedEvent) {
	org := ln.org(event)
	org.domains = append(org.domains, &loginNameDomain{name: event.Domain})
}

func (ln *InstanceLoginNames) generate() (loginNames []*UserLoginNames) {
	for _, org := range ln.Orgs {
		for userID, loginName := range org.LoginNames {
			loginName.UserID = userID
			loginName.InstanceID = ln.instance
			loginName.instancePolicy = ln.policy
			loginName.ownerPolicy = org.policy
			loginName.domains = org.domains
			loginName.removed = loginName.removed || ln.removed

			loginName.generate()
			loginNames = append(loginNames, loginName)
		}
	}

	return loginNames
}
