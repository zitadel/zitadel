package projection

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/user"
)

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
func (h *Human) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) (err error) {
	for _, event := range events {
		if !h.projection.shouldReduce(event) {
			continue
		}

		switch event.Type {
		case "user.human.added", "user.added":
			err = h.reduceAdded(event)
		case "user.human.selfregistered":
			err = h.reduceRegistered(event)
		case "user.human.profile.changed":
			err = h.reduceProfileChanged(event)
		case "user.human.phone.changed":
			err = h.reducePhoneChanged(event)
		case "user.human.phone.removed":
			h.Phone = nil
		case "user.human.phone.verified":
			h.Phone.IsVerified = true
		case "user.human.email.changed":
			err = h.reduceEmailChanged(event)
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
		if err != nil {
			return err
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

func (h *Human) reduceAdded(event *eventstore.Event[eventstore.StoragePayload]) error {
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

	return nil
}

func (h *Human) reduceRegistered(event *eventstore.Event[eventstore.StoragePayload]) error {
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

	return nil
}

func (h *Human) reduceProfileChanged(event *eventstore.Event[eventstore.StoragePayload]) error {
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
	return nil
}

func (h *Human) reducePhoneChanged(event *eventstore.Event[eventstore.StoragePayload]) error {
	h.Phone = new(Phone)
	e, err := user.HumanPhoneChangedEventFromStorage(event)
	if err != nil {
		return err
	}
	h.Phone.Number = e.Payload.PhoneNumber
	return nil
}

func (h *Human) reduceEmailChanged(event *eventstore.Event[eventstore.StoragePayload]) error {
	h.Email.IsVerified = false
	e, err := user.HumanEmailChangedEventFromStorage(event)
	if err != nil {
		return err
	}
	h.Email.Address = e.Payload.Address
	return nil
}
