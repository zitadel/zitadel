package command

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type UserHumanWriteModel struct {
	eventstore.WriteModel

	UserName string

	ProfileWriteModel bool
	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	PreferredLanguage language.Tag
	Gender            domain.Gender

	AvatarWriteModel bool
	Avatar           string

	InitCode             *crypto.CryptoValue
	InitCodeCreationDate time.Time
	InitCodeExpiry       time.Duration

	PasswordWriteModel       bool
	PasswordEncodedHash      string
	PasswordChangeRequired   bool
	PasswordCode             *crypto.CryptoValue
	PasswordCodeCreationDate time.Time
	PasswordCodeExpiry       time.Duration
	PasswordCheckFailedCount uint64

	EmailWriteModel       bool
	Email                 domain.EmailAddress
	IsEmailVerified       bool
	EmailCode             *crypto.CryptoValue
	EmailCodeCreationDate time.Time
	EmailCodeExpiry       time.Duration

	PhoneWriteModel bool
	Phone           domain.PhoneNumber
	IsPhoneVerified bool

	StateWriteModel bool
	UserState       domain.UserState
}

func NewUserHumanAllWriteModel(userID, resourceOwner string) *UserHumanWriteModel {
	return &UserHumanWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		AvatarWriteModel: true,
		EmailWriteModel:  true,
		PhoneWriteModel:  true,
		StateWriteModel:  true,
	}
}

func NewUserHumanWriteModel(userID, resourceOwner string, profileWM, emailWM, phoneWM, passwordWM bool) *UserHumanWriteModel {
	return &UserHumanWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		ProfileWriteModel:  profileWM,
		EmailWriteModel:    emailWM,
		PhoneWriteModel:    phoneWM,
		PasswordWriteModel: passwordWM,
	}
}

func NewUserHumanStateWriteModel(userID, resourceOwner string) *UserHumanWriteModel {
	return &UserHumanWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		StateWriteModel: true,
	}
}

func NewUserHumanEmailWriteModel(userID, resourceOwner string) *UserHumanWriteModel {
	return &UserHumanWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		EmailWriteModel: true,
	}
}

func NewUserHumanPhoneWriteModel(userID, resourceOwner string) *UserHumanWriteModel {
	return &UserHumanWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		PhoneWriteModel: true,
	}
}

func (wm *UserHumanWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.reduceHumanAddedEvent(e)
		case *user.HumanRegisteredEvent:
			wm.reduceHumanRegisteredEvent(e)
		case *user.HumanInitialCodeAddedEvent:
			wm.InitCode = e.Code
			wm.InitCodeCreationDate = e.CreationDate()
			wm.InitCodeExpiry = e.Expiry
			wm.UserState = domain.UserStateInitial
		case *user.HumanInitializedCheckSucceededEvent:
			wm.InitCode = nil
			wm.UserState = domain.UserStateActive
		case *user.UsernameChangedEvent:
			wm.UserName = e.UserName
		case *user.HumanProfileChangedEvent:
			wm.reduceHumanProfileChangedEvent(e)
		case *user.HumanEmailChangedEvent:
			wm.Email = e.EmailAddress
			wm.IsEmailVerified = false
			wm.EmailCode = nil
		case *user.HumanEmailCodeAddedEvent:
			wm.EmailCode = e.Code
			wm.EmailCodeCreationDate = e.CreationDate()
			wm.EmailCodeExpiry = e.Expiry
		case *user.HumanEmailVerifiedEvent:
			wm.IsEmailVerified = true
			wm.EmailCode = nil
		case *user.HumanPhoneChangedEvent:
			wm.reduceHumanPhoneChangedEvent(e)
		case *user.HumanPhoneVerifiedEvent:
			wm.reduceHumanPhoneVerifiedEvent()
		case *user.HumanPhoneRemovedEvent:
			wm.reduceHumanPhoneRemovedEvent()
		case *user.HumanAvatarAddedEvent:
			wm.Avatar = e.StoreKey
		case *user.HumanAvatarRemovedEvent:
			wm.Avatar = ""
		case *user.UserLockedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateLocked
			}
		case *user.UserUnlockedEvent:
			wm.PasswordCheckFailedCount = 0
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateActive
			}
		case *user.UserDeactivatedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateInactive
			}
		case *user.UserReactivatedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateActive
			}
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted

		case *user.HumanPasswordHashUpdatedEvent:
			wm.PasswordEncodedHash = e.EncodedHash
		case *user.HumanPasswordCheckFailedEvent:
			wm.PasswordCheckFailedCount += 1
		case *user.HumanPasswordCheckSucceededEvent:
			wm.PasswordCheckFailedCount = 0
		case *user.HumanPasswordChangedEvent:
			wm.PasswordEncodedHash = user.SecretOrEncodedHash(e.Secret, e.EncodedHash)
			wm.PasswordChangeRequired = e.ChangeRequired
			wm.PasswordCode = nil
			wm.PasswordCheckFailedCount = 0
		case *user.HumanPasswordCodeAddedEvent:
			wm.PasswordCode = e.Code
			wm.PasswordCodeCreationDate = e.CreationDate()
			wm.PasswordCodeExpiry = e.Expiry
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UserHumanWriteModel) Query() *eventstore.SearchQueryBuilder {
	eventTypes := []eventstore.EventType{
		user.UserV1AddedType,
		user.HumanAddedType,
		user.UserV1RegisteredType,
		user.HumanRegisteredType,
		user.UserV1InitialCodeAddedType,
		user.HumanInitialCodeAddedType,
		user.UserV1InitializedCheckSucceededType,
		user.HumanInitializedCheckSucceededType,
		user.UserRemovedType,

		user.UserUserNameChangedType,
	}

	if wm.EmailWriteModel {
		eventTypes = append(eventTypes,
			user.UserV1EmailChangedType,
			user.HumanEmailChangedType,
			user.UserV1EmailCodeAddedType,
			user.HumanEmailCodeAddedType,
			user.UserV1EmailVerifiedType,
			user.HumanEmailVerifiedType,
		)
	}
	if wm.PhoneWriteModel {
		eventTypes = append(eventTypes,
			user.UserV1PhoneChangedType,
			user.HumanPhoneChangedType,
			user.UserV1PhoneVerifiedType,
			user.HumanPhoneVerifiedType,
			user.UserV1PhoneRemovedType,
			user.HumanPhoneRemovedType,
		)
	}
	if wm.ProfileWriteModel {
		eventTypes = append(eventTypes,
			user.UserV1ProfileChangedType,
			user.HumanProfileChangedType,
		)
	}
	if wm.StateWriteModel {
		eventTypes = append(eventTypes,
			user.UserLockedType,
			user.UserUnlockedType,
			user.UserDeactivatedType,
			user.UserReactivatedType,
		)
	}
	if wm.AvatarWriteModel {
		eventTypes = append(eventTypes,
			user.HumanAvatarAddedType,
			user.HumanAvatarRemovedType,
		)
	}
	if wm.PasswordWriteModel {
		eventTypes = append(eventTypes,
			user.HumanPasswordChangedType,
			user.HumanPasswordCodeAddedType,
			user.HumanPasswordCheckFailedType,
			user.HumanPasswordCheckSucceededType,
			user.HumanPasswordHashUpdatedType,
			user.UserV1PasswordChangedType,
			user.UserV1PasswordCodeAddedType,
			user.UserV1PasswordCheckFailedType,
			user.UserV1PasswordCheckSucceededType,

			user.HumanEmailVerifiedType,
		)
	}

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(eventTypes...).
		Builder()
	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func (wm *UserHumanWriteModel) reduceHumanAddedEvent(e *user.HumanAddedEvent) {
	wm.UserName = e.UserName
	wm.FirstName = e.FirstName
	wm.LastName = e.LastName
	wm.NickName = e.NickName
	wm.DisplayName = e.DisplayName
	wm.PreferredLanguage = e.PreferredLanguage
	wm.Gender = e.Gender
	wm.Email = e.EmailAddress
	wm.Phone = e.PhoneNumber
	wm.UserState = domain.UserStateActive
	wm.PasswordEncodedHash = user.SecretOrEncodedHash(e.Secret, e.EncodedHash)
	wm.PasswordChangeRequired = e.ChangeRequired
}

func (wm *UserHumanWriteModel) reduceHumanRegisteredEvent(e *user.HumanRegisteredEvent) {
	wm.UserName = e.UserName
	wm.FirstName = e.FirstName
	wm.LastName = e.LastName
	wm.NickName = e.NickName
	wm.DisplayName = e.DisplayName
	wm.PreferredLanguage = e.PreferredLanguage
	wm.Gender = e.Gender
	wm.Email = e.EmailAddress
	wm.Phone = e.PhoneNumber
	wm.UserState = domain.UserStateActive
	wm.PasswordEncodedHash = user.SecretOrEncodedHash(e.Secret, e.EncodedHash)
	wm.PasswordChangeRequired = e.ChangeRequired
}

func (wm *UserHumanWriteModel) reduceHumanProfileChangedEvent(e *user.HumanProfileChangedEvent) {
	if e.FirstName != "" {
		wm.FirstName = e.FirstName
	}
	if e.LastName != "" {
		wm.LastName = e.LastName
	}
	if e.NickName != nil {
		wm.NickName = *e.NickName
	}
	if e.DisplayName != nil {
		wm.DisplayName = *e.DisplayName
	}
	if e.PreferredLanguage != nil {
		wm.PreferredLanguage = *e.PreferredLanguage
	}
	if e.Gender != nil {
		wm.Gender = *e.Gender
	}
}

func (wm *UserHumanWriteModel) reduceHumanEmailChangedEvent(e *user.HumanEmailChangedEvent) {
	wm.Email = e.EmailAddress
	wm.IsEmailVerified = false
	wm.EmailCode = nil
}

func (wm *UserHumanWriteModel) reduceHumanPhoneChangedEvent(e *user.HumanPhoneChangedEvent) {
	wm.Phone = e.PhoneNumber
	wm.IsPhoneVerified = false
}

func (wm *UserHumanWriteModel) reduceHumanPhoneVerifiedEvent() {
	wm.IsPhoneVerified = true
}

func (wm *UserHumanWriteModel) reduceHumanPhoneRemovedEvent() {
	wm.Phone = ""
	wm.IsPhoneVerified = false
}

func (wm *UserHumanWriteModel) NewEmailAddressChangedEvent(
	ctx context.Context,
	email domain.EmailAddress,
) *user.HumanEmailChangedEvent {
	if wm.Email == email {
		return nil
	}
	return user.NewHumanEmailChangedEvent(ctx, &wm.Aggregate().Aggregate, email)
}

func (wm *UserHumanWriteModel) NewEmailIsVerifiedEvent(
	ctx context.Context,
	isVerified bool,
) *user.HumanEmailVerifiedEvent {
	if wm.IsEmailVerified == isVerified || !wm.IsEmailVerified {
		return nil
	}
	return user.NewHumanEmailVerifiedEvent(ctx, &wm.Aggregate().Aggregate)
}

func (wm *UserHumanWriteModel) NewPhoneNumberChangedEvent(
	ctx context.Context,
	phone domain.PhoneNumber,
) *user.HumanPhoneChangedEvent {
	if wm.Phone == phone {
		return nil
	}
	return user.NewHumanPhoneChangedEvent(ctx, &wm.Aggregate().Aggregate, phone)
}

func (wm *UserHumanWriteModel) NewPhoneIsVerifiedEvent(
	ctx context.Context,
	isVerified bool,
) *user.HumanPhoneVerifiedEvent {
	if wm.IsPhoneVerified == isVerified || !wm.IsPhoneVerified {
		return nil
	}
	return user.NewHumanPhoneVerifiedEvent(ctx, &wm.Aggregate().Aggregate)
}

func (wm *UserHumanWriteModel) Aggregate() *user.Aggregate {
	return user.NewAggregate(wm.AggregateID, wm.ResourceOwner)
}

func (wm *UserHumanWriteModel) NewProfileChangedEvent(
	ctx context.Context,
	firstName,
	lastName,
	nickName,
	displayName *string,
	preferredLanguage *language.Tag,
	gender *domain.Gender,
) (*user.HumanProfileChangedEvent, error) {
	changes := make([]user.ProfileChanges, 0)
	if firstName != nil && wm.FirstName != *firstName {
		changes = append(changes, user.ChangeFirstName(*firstName))
	}
	if lastName != nil && wm.LastName != *lastName {
		changes = append(changes, user.ChangeLastName(*lastName))
	}
	if nickName != nil && wm.NickName != *nickName {
		changes = append(changes, user.ChangeNickName(*nickName))
	}
	if displayName != nil && wm.DisplayName != *displayName {
		changes = append(changes, user.ChangeDisplayName(*displayName))
	}
	if preferredLanguage != nil && wm.PreferredLanguage != *preferredLanguage {
		changes = append(changes, user.ChangePreferredLanguage(*preferredLanguage))
	}
	if gender != nil && wm.Gender != *gender {
		changes = append(changes, user.ChangeGender(*gender))
	}
	if len(changes) == 0 {
		return nil, nil
	}
	return user.NewHumanProfileChangedEvent(ctx, &wm.Aggregate().Aggregate, changes)
}
