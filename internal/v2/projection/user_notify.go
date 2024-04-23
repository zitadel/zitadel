package projection

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/user"
)

type NotifyUser struct {
	projection
	id string

	Username string
	// LoginNames         database.TextArray[string]
	// PreferredLoginName string
	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	AvatarKey         *string
	PreferredLanguage language.Tag
	Gender            domain.Gender
	LastEmail         *string
	VerifiedEmail     *string
	LastPhone         *string
	VerifiedPhone     *string
	PasswordSet       bool
}

func NewNotifyUser(id string) *NotifyUser {
	return &NotifyUser{
		id: id,
	}
}

func (n *NotifyUser) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) (err error) {
	for _, event := range events {
		if !n.shouldReduce(event) {
			continue
		}

		switch event.Type {
		case "user.human.added", "user.added":
			err = n.reduceAdded(event)
		case "user.human.selfregistered":
			err = n.reduceRegistered(event)
		case "user.human.phone.changed":
			err = n.reducePhoneChanged(event)
		case "user.human.phone.removed":
			err = n.reducePhoneRemoved()
		case "user.human.phone.verified":
			err = n.reducePhoneVerified()
		case "user.human.email.changed":
			err = n.reduceEmailChanged(event)
		case "user.human.email.verified":
			err = n.reduceEmailVerified()
		case "user.human.avatar.added":
			err = n.reduceAvatarAdded(event)
		case "user.human.avatar.removed":
			err = n.reduceAvatarRemoved()
		case "user.human.password.changed":
			err = n.reducePasswordChanged()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *NotifyUser) Filter() []*eventstore.Filter {
	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.FilterPagination(
				eventstore.GlobalPositionGreater(&n.position),
			),
			eventstore.AppendAggregateFilter(
				user.AggregateType,
				eventstore.AggregateID(n.id),
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

func (n *NotifyUser) reduceAdded(event *eventstore.Event[eventstore.StoragePayload]) error {
	e, err := user.HumanAddedEventFromStorage(event)
	if err != nil {
		return err
	}

	n.Username = e.Payload.Username
	n.FirstName = e.Payload.FirstName
	n.LastName = e.Payload.LastName
	n.NickName = e.Payload.NickName
	n.DisplayName = e.Payload.DisplayName
	n.PreferredLanguage = e.Payload.PreferredLanguage
	n.Gender = e.Payload.Gender
	if e.Payload.EmailAddress != "" {
		n.LastEmail = (*string)(&e.Payload.EmailAddress)
	}
	if e.Payload.PhoneNumber != "" {
		n.LastPhone = (*string)(&e.Payload.PhoneNumber)
	}
	n.PasswordSet = crypto.SecretOrEncodedHash(e.Payload.Secret, e.Payload.EncodedHash) != ""
	return nil
}

func (n *NotifyUser) reduceRegistered(event *eventstore.Event[eventstore.StoragePayload]) error {
	e, err := user.HumanRegisteredEventFromStorage(event)
	if err != nil {
		return err
	}

	n.Username = e.Payload.Username
	n.FirstName = e.Payload.FirstName
	n.LastName = e.Payload.LastName
	n.NickName = e.Payload.NickName
	n.DisplayName = e.Payload.DisplayName
	n.PreferredLanguage = e.Payload.PreferredLanguage
	n.Gender = e.Payload.Gender
	if e.Payload.EmailAddress != "" {
		n.LastEmail = (*string)(&e.Payload.EmailAddress)
	}
	if e.Payload.PhoneNumber != "" {
		n.LastPhone = (*string)(&e.Payload.PhoneNumber)
	}
	n.PasswordSet = crypto.SecretOrEncodedHash(e.Payload.Secret, e.Payload.EncodedHash) != ""
	return nil
}

func (n *NotifyUser) reducePhoneChanged(event *eventstore.Event[eventstore.StoragePayload]) error {
	e, err := user.HumanPhoneChangedEventFromStorage(event)
	if err != nil {
		return err
	}
	n.LastPhone = (*string)(&e.Payload.PhoneNumber)
	return nil
}

func (n *NotifyUser) reducePhoneRemoved() error {
	n.LastPhone = nil
	n.VerifiedPhone = nil
	return nil
}

func (n *NotifyUser) reducePhoneVerified() error {
	n.VerifiedPhone = n.LastPhone
	return nil
}

func (n *NotifyUser) reduceEmailChanged(event *eventstore.Event[eventstore.StoragePayload]) error {
	e, err := user.HumanEmailChangedEventFromStorage(event)
	if err != nil {
		return err
	}
	n.LastEmail = (*string)(&e.Payload.Address)
	if e.Payload.Address == "" {
		n.LastEmail = nil
	}
	return nil
}

func (n *NotifyUser) reduceEmailVerified() error {
	n.VerifiedEmail = n.LastEmail
	return nil
}

func (n *NotifyUser) reducePasswordChanged() error {
	n.PasswordSet = true
	return nil
}

func (n *NotifyUser) reduceAvatarAdded(event *eventstore.Event[eventstore.StoragePayload]) error {
	e, err := user.HumanAvatarAddedEventFromStorage(event)
	if err != nil {
		return err
	}

	n.AvatarKey = &e.Payload.StoreKey
	return nil
}

func (n *NotifyUser) reduceAvatarRemoved() error {
	n.AvatarKey = nil
	return nil
}
