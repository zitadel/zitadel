package projection

import (
	"golang.org/x/text/language"

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

type Machine struct {
	projection

	id              string
	Name            string
	Description     string
	AccessTokenType domain.OIDCTokenType

	// TODO: separate projection?
	Secret *string
}

var _ eventstore.Reducer = (*Machine)(nil)

func NewMachineProjection(id string) *Machine {
	return &Machine{
		id: id,
	}
}

// Reduce implements eventstore.Reducer.
func (m *Machine) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		if !m.projection.shouldReduce(event) {
			continue
		}
		switch event.Type {
		case "user.machine.added":
			e, err := user.MachineAddedEventFromStorage(event)
			if err != nil {
				return err
			}

			m.Name = e.Payload.Name
			m.Description = e.Payload.Description
			m.AccessTokenType = e.Payload.AccessTokenType
		case "user.machine.changed":
			e, err := user.MachineChangedEventFromStorage(event)
			if err != nil {
				return err
			}

			if e.Payload.Name != nil {
				m.Name = *e.Payload.Name
			}
			if e.Payload.Description != nil {
				m.Description = *e.Payload.Description
			}
			if e.Payload.AccessTokenType != nil {
				m.AccessTokenType = *e.Payload.AccessTokenType
			}
		case "user.machine.secret.set":
			e, err := user.MachineSecretSetEventFromStorage(event)
			if err != nil {
				return err
			}
			m.Secret = &e.Payload.HashedSecret
		case "user.machine.secret.updated":
			e, err := user.MachineSecretHashUpdatedEventFromStorage(event)
			if err != nil {
				return err
			}
			m.Secret = &e.Payload.HashedSecret
		case "user.machine.secret.removed":
			e, err := user.MachineSecretHashUpdatedEventFromStorage(event)
			if err != nil {
				return err
			}
			m.Secret = &e.Payload.HashedSecret
		default:
			continue
		}
		m.projection.reduce(event)
	}
	return nil
}

func (m *Machine) Filter() []*eventstore.Filter {
	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.FilterPagination(
				eventstore.GlobalPositionGreater(&m.position),
			),
			eventstore.AppendAggregateFilter(
				user.AggregateType,
				eventstore.AggregateID(m.id),
				eventstore.AppendEvent(
					eventstore.EventTypes(
						"user.machine.added",
						"user.machine.changed",
						"user.machine.secret.set",
						"user.machine.secret.updated",
						"user.machine.secret.removed",
					),
				),
			),
		),
	}
}

type Human struct {
	projection
	id string

	FirstName              string
	LastName               string
	NickName               string
	DisplayName            string
	AvatarKey              *string
	PreferredLanguage      language.Tag
	Gender                 domain.Gender
	Email                  Email
	Phone                  *Phone
	PasswordChangeRequired bool
}

type Phone struct {
	Number     domain.PhoneNumber
	IsVerified bool
}

type Email struct {
	Address    domain.EmailAddress
	IsVerified bool
}

var _ eventstore.Reducer = (*Human)(nil)

func NewHumanProjection(id string) *Human {
	return &Human{
		id: id,
	}
}

// Reduce implements eventstore.Reducer.
func (h *Human) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		if !h.projection.shouldReduce(event) {
			continue
		}

		switch event.Type {
		case "user.human.added", "user.added":
			e, err := user.HumanAddedEventFromStorage(event)
			if err != nil {
				return err
			}

			h.FirstName = e.Payload.FirstName
			h.LastName = e.Payload.LastName
			h.NickName = e.Payload.NickName
			h.DisplayName = e.Payload.DisplayName
			h.PreferredLanguage = e.Payload.PreferredLanguage
			h.Gender = e.Payload.Gender
			h.Email.Address = e.Payload.EmailAddress
			if e.Payload.PhoneNumber != "" {
				h.Phone = &Phone{
					Number: e.Payload.PhoneNumber,
				}
			}
			h.PasswordChangeRequired = e.Payload.PasswordChangeRequired
		case "user.human.selfregistered":
			e, err := user.HumanRegisteredEventFromStorage(event)
			if err != nil {
				return err
			}

			h.FirstName = e.Payload.FirstName
			h.LastName = e.Payload.LastName
			h.NickName = e.Payload.NickName
			h.DisplayName = e.Payload.DisplayName
			h.PreferredLanguage = e.Payload.PreferredLanguage
			h.Gender = e.Payload.Gender
			h.Email.Address = e.Payload.EmailAddress
			if e.Payload.PhoneNumber != "" {
				h.Phone = &Phone{
					Number: e.Payload.PhoneNumber,
				}
			}
			h.PasswordChangeRequired = e.Payload.PasswordChangeRequired
		case "user.human.profile.changed":
			e, err := user.HumanProfileChangedEventFromStorage(event)
			if err != nil {
				return err
			}
			if e.Payload.FirstName != "" {
				h.FirstName = e.Payload.FirstName
			}
			if e.Payload.LastName != "" {
				h.LastName = e.Payload.LastName
			}
			if e.Payload.NickName != nil {
				h.NickName = *e.Payload.NickName
			}
			if e.Payload.DisplayName != nil {
				h.DisplayName = *e.Payload.DisplayName
			}
			if e.Payload.PreferredLanguage != nil {
				h.PreferredLanguage = *e.Payload.PreferredLanguage
			}
			if e.Payload.Gender != nil {
				h.Gender = *e.Payload.Gender
			}
		case "user.human.phone.changed":
			h.Phone = new(Phone)
			e, err := user.HumanPhoneChangedEventFromStorage(event)
			if err != nil {
				return err
			}
			h.Phone.Number = e.Payload.PhoneNumber
		case "user.human.phone.removed":
			h.Phone = nil
		case "user.human.phone.verified":
			h.Phone.IsVerified = true
		case "user.human.email.changed":
			h.Email.IsVerified = false
			e, err := user.HumanEmailChangedEventFromStorage(event)
			if err != nil {
				return err
			}
			h.Email.Address = e.Payload.Address
		case "user.human.email.verified":
			h.Email.IsVerified = true
		case "user.human.avatar.added":
			e, err := user.HumanAvatarAddedEventFromStorage(event)
			if err != nil {
				return err
			}
			h.AvatarKey = &e.Payload.StoreKey
		case "user.human.avatar.removed":
			h.AvatarKey = nil
		case "user.human.password.changed":
			e, err := user.HumanPasswordChangedEventFromStorage(event)
			if err != nil {
				return err
			}
			h.PasswordChangeRequired = e.Payload.ChangeRequired
		default:
			continue
		}
		h.projection.reduce(event)
	}
	return nil
}

func (h *Human) Filter() []*eventstore.Filter {
	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.FilterPagination(
				eventstore.GlobalPositionGreater(&h.position),
			),
			eventstore.AppendAggregateFilter(
				user.AggregateType,
				eventstore.AggregateID(h.id),
				eventstore.AppendEvent(
					eventstore.EventTypes(
						"user.human.added",
						"user.added",
						"user.human.selfregistered",
						"user.human.profile.changed",
						"user.human.phone.changed",
						"user.human.phone.removed",
						"user.human.phone.verified",
						"user.human.email.changed",
						"user.human.email.verified",
						"user.human.avatar.added",
						"user.human.avatar.removed",
						"user.human.password.changed",
					),
				),
			),
		),
	}
}

// import (
// 	"github.com/zitadel/zitadel/internal/v2/eventstore"
// 	"github.com/zitadel/zitadel/internal/v2/org"
// )

// type OrgState struct {
// 	projection

// 	id string

// 	org.State
// }

// func NewStateProjection(id string) *OrgState {
// 	// TODO: check buffer for id and return from buffer if exists
// 	return &OrgState{
// 		id: id,
// 	}
// }

// func (p *OrgState) Filter() []*eventstore.Filter {
// 	return []*eventstore.Filter{
// 		eventstore.NewFilter(
// 			eventstore.FilterPagination(
// 				eventstore.Descending(),
// 				eventstore.GlobalPositionGreater(&p.position),
// 			),
// 			eventstore.AppendAggregateFilter(
// 				org.AggregateType,
// 				eventstore.AggregateID(p.id),
// 				eventstore.AppendEvent(
// 					eventstore.EventType("org.added"),
// 				),
// 				eventstore.AppendEvent(
// 					eventstore.EventType("org.deactivated"),
// 				),
// 				eventstore.AppendEvent(
// 					eventstore.EventType("org.reactivated"),
// 				),
// 				eventstore.AppendEvent(
// 					eventstore.EventType("org.removed"),
// 				),
// 			),
// 		),
// 	}
// }

// func (p *OrgState) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
// 	for _, event := range events {
// 		if !p.shouldReduce(event) {
// 			continue
// 		}

// 		switch {
// 		case org.Added.IsType(event.Type):
// 			p.State = org.ActiveState
// 		case org.Deactivated.IsType(event.Type):
// 			p.State = org.InactiveState
// 		case org.Reactivated.IsType(event.Type):
// 			p.State = org.ActiveState
// 		case org.Removed.IsType(event.Type):
// 			p.State = org.RemovedState
// 		default:
// 			continue
// 		}
// 		p.position = event.Position
// 	}

// 	// TODO: if more than x events store state

// 	return nil
// }
