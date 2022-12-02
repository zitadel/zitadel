package projection

import (
	"context"

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

func NewInstanceLoginNamesWithOwner(instance, owner string) *InstanceLoginNames {
	return &InstanceLoginNames{
		instance: instance,
		LoginNames: []*OrgLoginNames{
			{
				org: owner,
			},
		},
	}
}

type InstanceLoginNames struct {
	instance string
	policy   loginNamePolicy
	removed  bool

	LoginNames []*OrgLoginNames
}

type OrgLoginNames struct {
	org string
	// LoginNames per user
	LoginNames map[string]*UserLoginNames

	policy  loginNamePolicy
	domains []*loginNameDomain
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

func (ln *InstanceLoginNames) reduceInstanceEvents(events []eventstore.Event) {
	for _, event := range events {
		if event.Aggregate().Type != instance.AggregateType {
			continue
		}
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
