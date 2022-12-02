package projection

import (
	"context"
	"time"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

var _ Projection = (*User)(nil)

func NewUser(id, instance string) *User {
	return &User{
		ID:         id,
		instanceID: instance,
	}
}

func NewUserWithOwner(id, instance, owner string) *User {
	return &User{
		ID:            id,
		instanceID:    instance,
		ResourceOwner: owner,
	}
}

type User struct {
	instanceID string

	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64
	State         domain.UserState
	Type          domain.UserType
	Username      string
	Human         *Human
	Machine       *Machine
}

type Human struct {
	Profile Profile
	Email   Email
	Phone   Phone
}

type Profile struct {
	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	AvatarKey         string
	PreferredLanguage language.Tag
	Gender            domain.Gender
}

type Phone struct {
	Number     string
	IsVerified bool
}

type Email struct {
	Address    string
	IsVerified bool
}

type Machine struct {
	Name        string
	Description string
}

func (u *User) Reduce(events []eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			u.reduceHumanAdded(e)
		case *user.HumanRegisteredEvent:
			u.reduceHumanRegistered(e)
		case *user.HumanInitialCodeAddedEvent:
			u.reduceStateChange(domain.UserStateInitial)
		case *user.HumanInitializedCheckSucceededEvent:
			u.reduceStateChange(domain.UserStateActive)
		case *user.UserLockedEvent:
			u.reduceStateChange(domain.UserStateLocked)
		case *user.UserUnlockedEvent:
			u.reduceStateChange(domain.UserStateActive)
		case *user.UserDeactivatedEvent:
			u.reduceStateChange(domain.UserStateInactive)
		case *user.UserReactivatedEvent:
			u.reduceStateChange(domain.UserStateActive)
		case *user.UserRemovedEvent:
			u.reduceStateChange(domain.UserStateDeleted)
		case *user.UsernameChangedEvent:
			u.reduceUserNameChanged(e)
		case *user.DomainClaimedEvent:
			u.reduceDomainClaimed(e)
		case *user.HumanProfileChangedEvent:
			u.reduceHumanProfileChanged(e)
		case *user.HumanPhoneChangedEvent:
			u.reduceHumanPhoneChanged(e)
		case *user.HumanPhoneRemovedEvent:
			u.reduceHumanPhoneRemoved(e)
		case *user.HumanPhoneVerifiedEvent:
			u.reduceHumanPhoneVerified(e)
		case *user.HumanEmailChangedEvent:
			u.reduceHumanEmailChanged(e)
		case *user.HumanEmailVerifiedEvent:
			u.reduceHumanEmailVerified(e)
		case *user.HumanAvatarAddedEvent:
			u.reduceHumanAvatarAdded(e)
		case *user.HumanAvatarRemovedEvent:
			u.reduceHumanAvatarRemoved(e)
		case *user.MachineAddedEvent:
			u.reduceMachineAdded(e)
		case *user.MachineChangedEvent:
			u.reduceMachineChanged(e)
		// case *user.HumanPasswordChangedEvent:
		// 	u.reduceHumanPasswordChanged(e)
		default:
			logging.WithFields("type", e.Type()).Debug("event not handeled")
		}
		u.ChangeDate = event.CreationDate()
		u.Sequence = event.Sequence()
	}
}

func (u *User) SearchQuery(context.Context) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(u.instanceID).
		OrderAsc().
		ResourceOwner(u.ResourceOwner).
		AddQuery().
		AggregateTypes(
			user.AggregateType,
		).
		AggregateIDs(u.ID).
		EventTypes(
			user.UserV1AddedType,
			user.HumanAddedType,
			user.UserV1RegisteredType,
			user.HumanRegisteredType,
			user.HumanInitialCodeAddedType,
			user.UserV1InitialCodeAddedType,
			user.HumanInitializedCheckSucceededType,
			user.UserV1InitializedCheckSucceededType,
			user.UserLockedType,
			user.UserUnlockedType,
			user.UserDeactivatedType,
			user.UserReactivatedType,
			user.UserRemovedType,
			user.UserUserNameChangedType,
			user.UserDomainClaimedType,
			user.HumanProfileChangedType,
			user.UserV1ProfileChangedType,
			user.HumanPhoneChangedType,
			user.UserV1PhoneChangedType,
			user.HumanPhoneRemovedType,
			user.UserV1PhoneRemovedType,
			user.HumanPhoneVerifiedType,
			user.UserV1PhoneVerifiedType,
			user.HumanEmailChangedType,
			user.UserV1EmailChangedType,
			user.HumanEmailVerifiedType,
			user.UserV1EmailVerifiedType,
			user.HumanAvatarAddedType,
			user.HumanAvatarRemovedType,
			user.MachineAddedEventType,
			user.MachineChangedEventType,
			user.HumanPasswordChangedType,
		).
		Builder()
}

func (u *User) reduceHumanAdded(event *user.HumanAddedEvent) {
	u.Username = event.UserName
	u.Human.Profile.FirstName = event.FirstName
	u.Human.Profile.LastName = event.LastName
	u.Human.Profile.NickName = event.NickName
	u.Human.Profile.DisplayName = event.DisplayName
	u.Human.Profile.PreferredLanguage = event.PreferredLanguage
	u.Human.Profile.Gender = event.Gender
	u.Human.Email.Address = event.EmailAddress
	u.Human.Phone.Number = event.PhoneNumber
	u.CreationDate = event.CreationDate()
	// u.Human. = event.Secret
	// u.Human.ChangeRequired = event.ChangeRequired
}

func (u *User) reduceHumanRegistered(event *user.HumanRegisteredEvent) {
	u.Username = event.UserName
	u.Human.Profile.FirstName = event.FirstName
	u.Human.Profile.LastName = event.LastName
	u.Human.Profile.NickName = event.NickName
	u.Human.Profile.DisplayName = event.DisplayName
	u.Human.Profile.PreferredLanguage = event.PreferredLanguage
	u.Human.Profile.Gender = event.Gender
	u.Human.Email.Address = event.EmailAddress
	u.Human.Phone.Number = event.PhoneNumber
}

func (u *User) reduceStateChange(state domain.UserState) {
	u.State = state
}

func (u *User) reduceUserNameChanged(event *user.UsernameChangedEvent) {
	u.Username = event.UserName
}

func (u *User) reduceDomainClaimed(event *user.DomainClaimedEvent) {
	u.Username = event.UserName
}

func (u *User) reduceHumanProfileChanged(event *user.HumanProfileChangedEvent) {
	if event.FirstName != "" {
		u.Human.Profile.FirstName = event.FirstName
	}
	if event.LastName != "" {
		u.Human.Profile.LastName = event.LastName
	}
	if event.NickName != nil {
		u.Human.Profile.NickName = *event.NickName
	}
	if event.DisplayName != nil {
		u.Human.Profile.DisplayName = *event.DisplayName
	}
	if event.PreferredLanguage != nil {
		u.Human.Profile.PreferredLanguage = *event.PreferredLanguage
	}
	if event.Gender != nil {
		u.Human.Profile.Gender = *event.Gender
	}
}

func (u *User) reduceHumanPhoneChanged(event *user.HumanPhoneChangedEvent) {
	u.Human.Phone.Number = event.PhoneNumber
	u.Human.Phone.IsVerified = false
}

func (u *User) reduceHumanPhoneRemoved(event *user.HumanPhoneRemovedEvent) {
	u.Human.Phone.Number = ""
	u.Human.Phone.IsVerified = false
}

func (u *User) reduceHumanPhoneVerified(event *user.HumanPhoneVerifiedEvent) {
	u.Human.Phone.IsVerified = true
}

func (u *User) reduceHumanEmailChanged(event *user.HumanEmailChangedEvent) {
	u.Human.Email.Address = event.EmailAddress
	u.Human.Email.IsVerified = false
}

func (u *User) reduceHumanEmailVerified(event *user.HumanEmailVerifiedEvent) {
	u.Human.Email.IsVerified = true
}

func (u *User) reduceHumanAvatarAdded(event *user.HumanAvatarAddedEvent) {
	u.Human.Profile.AvatarKey = event.StoreKey
}

func (u *User) reduceHumanAvatarRemoved(event *user.HumanAvatarRemovedEvent) {
	u.Human.Profile.AvatarKey = ""
}

func (u *User) reduceMachineAdded(event *user.MachineAddedEvent) {
	u.Machine.Description = event.Description
	u.Machine.Name = event.Name
}

func (u *User) reduceMachineChanged(event *user.MachineChangedEvent) {
	if event.Description != nil {
		u.Machine.Description = *event.Description
	}

	if event.Name != nil {
		u.Machine.Name = *event.Name
	}
}

// func (u *User) reduceHumanPasswordChanged(event *user.HumanPasswordChangedEvent) {

// }
