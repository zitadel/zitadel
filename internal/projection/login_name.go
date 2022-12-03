package projection

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func NewUserLoginNames(userID, instance string) *UserLoginNames {
	return &UserLoginNames{
		UserID:     userID,
		InstanceID: instance,
	}
}

func NewUserLoginNamesWithOwner(userID, instance, owner string) *UserLoginNames {
	return &UserLoginNames{
		UserID:     userID,
		InstanceID: instance,
		OwnerID:    owner,
	}
}

var _ Projection = (*UserLoginNames)(nil)

type UserLoginNames struct {
	UserID     string
	InstanceID string
	OwnerID    string

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

func (ln *UserLoginNames) Reduce(events []eventstore.Event) {
	// user events are reduced before the others
	// to ensure all the ids are set
	ln.reduceUserEvents(events)
	if ln.removed {
		return
	}

	for _, event := range events {
		// only apply events from the instance or owner of the user
		if event.Aggregate().ResourceOwner != ln.OwnerID && event.Aggregate().ResourceOwner != ln.InstanceID {
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

func (ln *UserLoginNames) reduceUserEvents(events []eventstore.Event) {
	for _, event := range events {
		if event.Aggregate().Type != user.AggregateType {
			continue
		}
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
		default:
			logging.WithFields("type", e.Type()).Debug("event not handeled")
		}
	}
}

func (ln *UserLoginNames) SearchQuery(ctx context.Context) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(ln.InstanceID).
		OrderAsc().
		// ResourceOwner(ln.ownerID).
		AddQuery().
		AggregateTypes(
			user.AggregateType,
		).
		AggregateIDs(ln.UserID).
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
		AggregateIDs(ln.OwnerID).
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
		AggregateIDs(ln.InstanceID).
		EventTypes(
			instance.DomainPolicyAddedEventType,
			instance.DomainPolicyChangedEventType,
			instance.InstanceRemovedEventType,
		).
		Builder()
}

func (ln *UserLoginNames) generate() {
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

func (ln *UserLoginNames) reduceHumanAdded(event *user.HumanAddedEvent) {
	ln.username = event.UserName
	ln.OwnerID = event.Aggregate().ResourceOwner
}

func (ln *UserLoginNames) reduceHumanRegistered(event *user.HumanRegisteredEvent) {
	ln.username = event.UserName
	ln.OwnerID = event.Aggregate().ResourceOwner
}

func (ln *UserLoginNames) reduceMachineAddedEvent(event *user.MachineAddedEvent) {
	ln.username = event.UserName
	ln.OwnerID = event.Aggregate().ResourceOwner
}

func (ln *UserLoginNames) reduceUserRemoved(event *user.UserRemovedEvent) {
	ln.removed = true
}

func (ln *UserLoginNames) reduceUsernameChanged(event *user.UsernameChangedEvent) {
	ln.username = event.UserName
}

func (ln *UserLoginNames) reduceUserDomainClaimed(event *user.DomainClaimedEvent) {
	ln.username = event.UserName
}

func (ln *UserLoginNames) reduceOrgDomainPolicyAddedEvent(event *org.DomainPolicyAddedEvent) {
	ln.ownerPolicy.mustBeDomain = event.UserLoginMustBeDomain
	ln.ownerPolicy.active = true
}

func (ln *UserLoginNames) reduceOrgDomainPolicyChangedEvent(event *org.DomainPolicyChangedEvent) {
	if event.UserLoginMustBeDomain == nil {
		return
	}
	ln.ownerPolicy.mustBeDomain = *event.UserLoginMustBeDomain
}

func (ln *UserLoginNames) reduceDomainPolicyRemovedEvent(event *org.DomainPolicyRemovedEvent) {
	ln.ownerPolicy.active = false
}

func (ln *UserLoginNames) reduceOrgDomainPrimarySetEvent(event *org.DomainPrimarySetEvent) {
	for _, domain := range ln.domains {
		domain.isPrimary = domain.name == event.Domain
	}
}

func (ln *UserLoginNames) reduceOrgDomainRemovedEvent(event *org.DomainRemovedEvent) {
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

func (ln *UserLoginNames) reduceOrgDomainVerifiedEvent(event *org.DomainVerifiedEvent) {
	ln.domains = append(ln.domains, &loginNameDomain{name: event.Domain})
}

func (ln *UserLoginNames) reduceInstanceDomainPolicyAddedEvent(event *instance.DomainPolicyAddedEvent) {
	ln.instancePolicy.mustBeDomain = event.UserLoginMustBeDomain
}

func (ln *UserLoginNames) reduceInstanceDomainPolicyChangedEvent(event *instance.DomainPolicyChangedEvent) {
	if event.UserLoginMustBeDomain == nil {
		return
	}
	ln.instancePolicy.mustBeDomain = *event.UserLoginMustBeDomain
}

func (ln *UserLoginNames) reduceInstanceRemovedEvent(event *instance.InstanceRemovedEvent) {
	ln.removed = true
}
