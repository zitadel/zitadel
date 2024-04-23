package projection

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/user"
)

type User struct {
	projection

	ID       string
	State    domain.UserState
	Type     domain.UserType
	Username string
}

func NewUserProjection(id string) *User {
	return &User{
		ID: id,
	}
}

var _ eventstore.Reducer = (*User)(nil)

// Reduce implements eventstore.Reducer.
func (u *User) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		if !u.projection.shouldReduce(event) {
			continue
		}
		switch event.Type {
		case "user.added", "user.human.added":
			u.ID = event.Aggregate.ID
			u.State = domain.UserStateActive
			u.Type = domain.UserTypeHuman

			e, err := user.HumanAddedEventFromStorage(event)
			if err != nil {
				return err
			}
			u.Username = e.Payload.Username
		case "user.human.selfregistered":
			u.ID = event.Aggregate.ID
			u.State = domain.UserStateActive
			u.Type = domain.UserTypeHuman

			e, err := user.HumanRegisteredEventFromStorage(event)
			if err != nil {
				return err
			}
			u.Username = e.Payload.Username
		case "user.machine.added":
			u.ID = event.Aggregate.ID
			u.State = domain.UserStateActive
			u.Type = domain.UserTypeMachine

			e, err := user.MachineAddedEventFromStorage(event)
			if err != nil {
				return err
			}
			u.Username = e.Payload.Username
		case "user.locked":
			u.State = domain.UserStateLocked
		case "user.unlocked":
			u.State = domain.UserStateActive
		case "user.deactivated":
			u.State = domain.UserStateInactive
		case "user.reactivated":
			u.State = domain.UserStateActive
		case "user.removed":
			u.State = domain.UserStateDeleted
		case "user.human.initialization.code.added":
			u.State = domain.UserStateInitial
		case "user.human.initialization.check.succeeded":
			u.State = domain.UserStateActive
		case "user.username.changed":
			e, err := user.UsernameChangedEventFromStorage(event)
			if err != nil {
				return err
			}
			u.Username = e.Payload.Username
		case "user.domain.claimed.sent":
			e, err := user.DomainClaimedEventFromStorage(event)
			if err != nil {
				return err
			}
			u.Username = e.Payload.Username
		}
		u.projection.reduce(event)
	}
	return nil
}

func (u *User) Filter() []*eventstore.Filter {
	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.FilterPagination(
				eventstore.GlobalPositionGreater(&u.position),
			),
			eventstore.AppendAggregateFilter(
				user.AggregateType,
				eventstore.AggregateID(u.ID),
				eventstore.AppendEvent(
					eventstore.EventTypes(
						"user.added",
						"user.human.added",
						"user.human.selfregistered",
						"user.machine.added",
						"user.locked",
						"user.unlocked",
						"user.deactivated",
						"user.reactivated",
						"user.removed",
						"user.human.initialization.code.added",
						"user.human.initialization.check.succeeded",
						"user.username.changed",
						"ser.domain.claimed.sent",
					),
				),
			),
		),
	}
}
