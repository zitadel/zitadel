package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

var _ Projection = (*LoginNames)(nil)

type LoginNames struct {
	userID     string
	instanceID string
	ownerID    string

	LoginNames []*LoginName

	username       string
	ownerPolicy    loginNamePolicy
	instancePolicy loginNamePolicy
	domains        []*loginNameDomain
	removed        bool
}

type LoginName struct {
	Name      string
	IsPrimary bool
}

type loginNameDomain struct {
	name      string
	isPrimary bool
}

type loginNamePolicy struct {
	mustBeDomain bool
	active       bool
}

func NewLoginNames(userID, instance string) *LoginNames {
	return &LoginNames{
		userID:     userID,
		instanceID: instance,
	}
}

func NewLoginNamesWithOwner(userID, instance, owner string) *LoginNames {
	return &LoginNames{
		userID:     userID,
		instanceID: instance,
		ownerID:    owner,
	}
}

func (ln *LoginNames) Reduce(events []eventstore.Event) {
	// user events are reduced before the others
	// to ensure all the ids are set
	ln.reduceUserEvents(events)
	if ln.removed {
		return
	}

	for _, event := range events {
		// only apply events from the instance or owner of the user
		if event.Aggregate().ResourceOwner != ln.ownerID && event.Aggregate().ResourceOwner != ln.instanceID {
			continue
		}

		switch e := event.(type) {
		case *org.DomainPolicyAddedEvent:
			ln.reduceOrgDomainPolicyAddedEvent(e)
		case *org.DomainPolicyChangedEvent:
			ln.reduceOrgDomainPolicyChangedEvent(e)
		case *org.DomainPolicyRemovedEvent:
			ln.reduceDomainPolicyRemovedEvent(e)
		case *org.DomainPrimarySetEvent:
			ln.reduceOrgDomainPrimarySetEvent(e)
		case *org.DomainRemovedEvent:
			ln.reduceOrgDomainRemovedEvent(e)
		case *org.DomainVerifiedEvent:
			ln.reduceOrgDomainVerifiedEvent(e)
		case *instance.DomainPolicyAddedEvent:
			ln.reduceInstanceDomainPolicyAddedEvent(e)
		case *instance.DomainPolicyChangedEvent:
			ln.reduceInstanceDomainPolicyChangedEvent(e)
		case *instance.InstanceRemovedEvent:
			ln.reduceInstanceRemovedEvent(e)
		}
	}

	ln.generate()
}

func (ln *LoginNames) reduceUserEvents(events []eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			ln.reduceHumanAdded(e)
		case *user.HumanRegisteredEvent:
			ln.reduceHumanRegistered(e)
		case *user.MachineAddedEvent:
			ln.reduceMachineAddedEvent(e)
		case *user.UserRemovedEvent:
			ln.reduceUserRemoved(e)
		case *user.UsernameChangedEvent:
			ln.reduceUsernameChanged(e)
		case *user.DomainClaimedEvent:
			ln.reduceUserDomainClaimed(e)
		}
	}
}

func (ln *LoginNames) SearchQuery(ctx context.Context) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(ln.instanceID).
		OrderAsc().
		// ResourceOwner(ln.ownerID).
		AddQuery().
		AggregateTypes(
			user.AggregateType,
		).
		AggregateIDs(ln.userID).
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
		AggregateIDs(ln.ownerID).
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
		AggregateIDs(ln.instanceID).
		EventTypes(
			instance.DomainPolicyAddedEventType,
			instance.DomainPolicyChangedEventType,
			instance.InstanceRemovedEventType,
		).
		Builder()
}

func (ln *LoginNames) generate() {
	if ln.removed {
		return
	}

	mustBeDomain := ln.instancePolicy.mustBeDomain
	if ln.ownerPolicy.active {
		mustBeDomain = ln.ownerPolicy.mustBeDomain
	}
	if !mustBeDomain {
		ln.LoginNames = append(ln.LoginNames, &LoginName{
			Name: ln.username,
		})
		return
	}

	for _, domain := range ln.domains {
		ln.LoginNames = append(ln.LoginNames, &LoginName{
			Name:      ln.username + "@" + domain.name,
			IsPrimary: domain.isPrimary,
		})
	}
}

func (ln *LoginNames) reduceHumanAdded(event *user.HumanAddedEvent) {
	ln.username = event.UserName
	ln.ownerID = event.Aggregate().ResourceOwner
}

func (ln *LoginNames) reduceHumanRegistered(event *user.HumanRegisteredEvent) {
	ln.username = event.UserName
	ln.ownerID = event.Aggregate().ResourceOwner
}

func (ln *LoginNames) reduceMachineAddedEvent(event *user.MachineAddedEvent) {
	ln.username = event.UserName
	ln.ownerID = event.Aggregate().ResourceOwner
}

func (ln *LoginNames) reduceUserRemoved(event *user.UserRemovedEvent) {
	ln.removed = true
}

func (ln *LoginNames) reduceUsernameChanged(event *user.UsernameChangedEvent) {
	ln.username = event.UserName
}

func (ln *LoginNames) reduceUserDomainClaimed(event *user.DomainClaimedEvent) {
	ln.username = event.UserName
}

func (ln *LoginNames) reduceOrgDomainPolicyAddedEvent(event *org.DomainPolicyAddedEvent) {
	ln.ownerPolicy.mustBeDomain = event.UserLoginMustBeDomain
	ln.ownerPolicy.active = true
}

func (ln *LoginNames) reduceOrgDomainPolicyChangedEvent(event *org.DomainPolicyChangedEvent) {
	if event.UserLoginMustBeDomain == nil {
		return
	}
	ln.ownerPolicy.mustBeDomain = *event.UserLoginMustBeDomain
}

func (ln *LoginNames) reduceDomainPolicyRemovedEvent(event *org.DomainPolicyRemovedEvent) {
	ln.ownerPolicy.active = false
}

func (ln *LoginNames) reduceOrgDomainPrimarySetEvent(event *org.DomainPrimarySetEvent) {
	for _, domain := range ln.domains {
		domain.isPrimary = domain.name == event.Domain
	}
}

func (ln *LoginNames) reduceOrgDomainRemovedEvent(event *org.DomainRemovedEvent) {
	for i, domain := range ln.domains {
		if domain.name != event.Domain {
			continue
		}
		ln.domains[i] = ln.domains[len(ln.domains)-1]
		ln.domains[len(ln.domains)-1] = nil
		ln.domains = ln.domains[:len(ln.domains)-1]
		return
	}
}

func (ln *LoginNames) reduceOrgDomainVerifiedEvent(event *org.DomainVerifiedEvent) {
	ln.domains = append(ln.domains, &loginNameDomain{name: event.Domain})
}

func (ln *LoginNames) reduceInstanceDomainPolicyAddedEvent(event *instance.DomainPolicyAddedEvent) {
	ln.instancePolicy.mustBeDomain = event.UserLoginMustBeDomain
}

func (ln *LoginNames) reduceInstanceDomainPolicyChangedEvent(event *instance.DomainPolicyChangedEvent) {
	if event.UserLoginMustBeDomain == nil {
		return
	}
	ln.instancePolicy.mustBeDomain = *event.UserLoginMustBeDomain
}

func (ln *LoginNames) reduceInstanceRemovedEvent(event *instance.InstanceRemovedEvent) {
	ln.removed = true
}
