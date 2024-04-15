package projection

import (
	"slices"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/instance"
	"github.com/zitadel/zitadel/internal/v2/org"
	"github.com/zitadel/zitadel/internal/v2/user"
)

var _ eventstore.Reducer = (*LoginNames)(nil)

type LoginNames struct {
	projection
	UserID string

	Owner string

	LoginNames []*LoginName

	username       string
	ownerPolicy    loginNamePolicy
	instancePolicy loginNamePolicy
	domains        []*loginNameDomain
	removed        bool
}

func NewLoginNamesWithOwner(userID, instance, owner string) *LoginNames {
	return &LoginNames{
		projection: projection{
			instance: instance,
		},
		UserID: userID,
		Owner:  owner,
	}
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

// Reduce implements eventstore.Reducer.
func (ln *LoginNames) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	// user events are reduced before the others
	// to ensure all the ids are set
	if err := ln.reduceUserEvents(events); err != nil {
		return err
	}
	if ln.removed {
		return nil
	}

	for _, event := range events {
		// only apply events from the instance or owner of the user
		if event.Aggregate.Owner != ln.Owner && event.Aggregate.Owner != ln.instance {
			continue
		}
		switch event.Type {
		case "org.policy.domain.added":
			ln.ownerPolicy.active = true

			e, err := org.DomainPolicyAddedEventFromStorage(event)
			if err != nil {
				return err
			}
			ln.ownerPolicy.mustBeDomain = e.Payload.UserLoginMustBeDomain
		case "org.policy.domain.changed":
			e, err := org.DomainPolicyChangedEventFromStorage(event)
			if err != nil {
				return err
			}
			if e.Payload.UserLoginMustBeDomain == nil {
				continue
			}
			ln.ownerPolicy.mustBeDomain = *e.Payload.UserLoginMustBeDomain

		case "org.policy.domain.removed":
			ln.ownerPolicy.active = false
		case "org.domain.primary.set":
			e, err := org.DomainPrimarySetEventFromStorage(event)
			if err != nil {
				return err
			}
			for _, domain := range ln.domains {
				domain.isPrimary = domain.name == e.Payload.Name
			}
		case "org.domain.removed":
			e, err := org.DomainRemovedEventFromStorage(event)
			if err != nil {
				return err
			}
			ln.domains = slices.DeleteFunc(ln.domains, func(domain *loginNameDomain) bool {
				return domain.name == e.Payload.Name
			})
		case "org.domain.verified":
			e, err := org.DomainVerifiedEventFromStorage(event)
			if err != nil {
				return err
			}
			ln.domains = append(ln.domains, &loginNameDomain{name: e.Payload.Name})
		case "org.removed":
			ln.removed = true
		case "instance.policy.domain.added":
			e, err := instance.DomainPolicyAddedEventFromStorage(event)
			if err != nil {
				return err
			}
			ln.instancePolicy.mustBeDomain = e.Payload.UserLoginMustBeDomain
		case "instance.policy.domain.changed":
			e, err := instance.DomainPolicyChangedEventFromStorage(event)
			if err != nil {
				return err
			}
			if e.Payload.UserLoginMustBeDomain == nil {
				continue
			}
			ln.instancePolicy.mustBeDomain = *e.Payload.UserLoginMustBeDomain
		case "instance.removed":
			ln.removed = true
		}
	}
	return nil
}

func (ln *LoginNames) Filter() []*eventstore.Filter {
	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.AppendAggregateFilters(
				eventstore.NewAggregateFilter(
					user.AggregateType,
					eventstore.AggregateID(ln.UserID),
					eventstore.AppendEvent(
						eventstore.EventTypes(
							"user.added",
							"user.human.added",
							"user.human.selfregistered",
							"user.machine.added",
							"user.removed",
							"user.username.changed",
							"ser.domain.claimed.sent",
						),
					),
				),
				eventstore.NewAggregateFilter(
					org.AggregateType,
					eventstore.AggregateID(ln.Owner),
					eventstore.AppendEvent(
						eventstore.EventTypes(
							"org.policy.domain.added",
							"org.policy.domain.changed",
							"org.policy.domain.removed",
							"org.domain.primary.set",
							"org.domain.removed",
							"org.domain.verified",
							"org.removed",
						),
					),
				),
				eventstore.NewAggregateFilter(
					instance.AggregateType,
					eventstore.AggregateID(ln.instance),
					eventstore.AppendEvent(
						eventstore.EventTypes(
							"instance.policy.domain.added",
							"instance.policy.domain.changed",
							"instance.removed",
						),
					),
				),
			),
		),
	}
}

func (ln *LoginNames) Generate() {
	if ln.removed {
		return
	}

	mustBeDomain := ln.instancePolicy.mustBeDomain
	if ln.ownerPolicy.active {
		mustBeDomain = ln.ownerPolicy.mustBeDomain
	}
	if !mustBeDomain {
		ln.LoginNames = append(ln.LoginNames, &LoginName{
			Name:      ln.username,
			IsPrimary: true,
		})
		return
	}

	for _, domain := range ln.domains {
		ln.LoginNames = append(ln.LoginNames, &LoginName{
			Name:      strings.Join([]string{ln.username, domain.name}, "@"),
			IsPrimary: domain.isPrimary,
		})
	}
}

func (ln *LoginNames) reduceUserEvents(events []*eventstore.Event[eventstore.StoragePayload]) error {
	if ln.removed {
		return nil
	}

	for _, event := range events {
		if event.Aggregate.Type != user.AggregateType {
			continue
		}
		switch event.Type {
		case "user.added", "user.human.added":
			ln.Owner = event.Aggregate.Owner

			e, err := user.HumanAddedEventFromStorage(event)
			if err != nil {
				return err
			}
			ln.username = e.Payload.Username
		case "user.human.selfregistered":
			ln.Owner = event.Aggregate.Owner

			e, err := user.HumanRegisteredEventFromStorage(event)
			if err != nil {
				return err
			}
			ln.username = e.Payload.Username
		case "user.machine.added":
			ln.Owner = event.Aggregate.Owner

			e, err := user.MachineAddedEventFromStorage(event)
			if err != nil {
				return err
			}
			ln.username = e.Payload.Username
		case "user.removed":
			ln.removed = true
			return nil
		case "user.username.changed":
			e, err := user.UsernameChangedEventFromStorage(event)
			if err != nil {
				return err
			}
			ln.username = e.Payload.Username
		case "user.domain.claimed.sent":
			e, err := user.DomainClaimedEventFromStorage(event)
			if err != nil {
				return err
			}
			ln.username = e.Payload.Username
		default:
			logging.WithFields("type", event.Type).Debug("event not handled")
		}
	}

	return nil
}
