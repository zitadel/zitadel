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
	InitCheckFailedCount uint64

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
	EmailCheckFailedCount uint64

	PhoneWriteModel       bool
	Phone                 domain.PhoneNumber
	IsPhoneVerified       bool
	PhoneCode             *crypto.CryptoValue
	PhoneCodeCreationDate time.Time
	PhoneCodeExpiry       time.Duration
	PhoneCheckFailedCount uint64

	StateWriteModel bool
	UserState       domain.UserState
}

func NewUserHumanAllWriteModel(userID, resourceOwner string) *UserHumanWriteModel {
	return &UserHumanWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		AvatarWriteModel:   true,
		EmailWriteModel:    true,
		PhoneWriteModel:    true,
		PasswordWriteModel: true,
		StateWriteModel:    true,
	}
}

func NewUserHumanWriteModel(userID, resourceOwner string, profileWM, emailWM, phoneWM, passwordWM, stateWM, avatarWM bool) *UserHumanWriteModel {
	return &UserHumanWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		ProfileWriteModel:  profileWM,
		EmailWriteModel:    emailWM,
		PhoneWriteModel:    phoneWM,
		PasswordWriteModel: passwordWM,
		StateWriteModel:    stateWM,
		AvatarWriteModel:   avatarWM,
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
			wm.UserState = domain.UserStateInitial
			wm.SetInitCode(e.Code, e.Expiry, e.CreationDate())
		case *user.HumanInitializedCheckSucceededEvent:
			wm.UserState = domain.UserStateActive
			wm.EmptyInitCode()
		case *user.HumanInitializedCheckFailedEvent:
			wm.InitCheckFailedCount += 1

		case *user.UsernameChangedEvent:
			wm.UserName = e.UserName
		case *user.HumanProfileChangedEvent:
			wm.reduceHumanProfileChangedEvent(e)

		case *user.HumanEmailChangedEvent:
			wm.Email = e.EmailAddress
			wm.IsEmailVerified = false
			wm.EmptyEmailCode()
		case *user.HumanEmailCodeAddedEvent:
			wm.IsEmailVerified = false
			wm.SetEMailCode(e.Code, e.Expiry, e.CreationDate())
		case *user.HumanEmailVerifiedEvent:
			wm.IsEmailVerified = true
			wm.EmptyEmailCode()
		case *user.HumanEmailVerificationFailedEvent:
			wm.EmailCheckFailedCount += 1

		case *user.HumanPhoneChangedEvent:
			wm.IsPhoneVerified = false
			wm.Phone = e.PhoneNumber
			wm.EmptyPhoneCode()
		case *user.HumanPhoneCodeAddedEvent:
			wm.IsPhoneVerified = false
			wm.SetPhoneCode(e.Code, e.Expiry, e.CreationDate())
		case *user.HumanPhoneVerifiedEvent:
			wm.IsPhoneVerified = true
			wm.EmptyPhoneCode()
		case *user.HumanPhoneVerificationFailedEvent:
			wm.PhoneCheckFailedCount += 1
		case *user.HumanPhoneRemovedEvent:
			wm.EmptyPhoneCode()
			wm.Phone = ""
			wm.IsPhoneVerified = false

		case *user.HumanAvatarAddedEvent:
			wm.Avatar = e.StoreKey
		case *user.HumanAvatarRemovedEvent:
			wm.Avatar = ""

		case *user.UserLockedEvent:
			wm.UserState = domain.UserStateLocked
		case *user.UserUnlockedEvent:
			wm.PasswordCheckFailedCount = 0
			wm.UserState = domain.UserStateActive

		case *user.UserDeactivatedEvent:
			wm.UserState = domain.UserStateInactive
		case *user.UserReactivatedEvent:
			wm.UserState = domain.UserStateActive

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
			wm.EmptyPasswordCode()
		case *user.HumanPasswordCodeAddedEvent:
			wm.SetPasswordCode(e.Code, e.Expiry, e.CreationDate())
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UserHumanWriteModel) EmptyInitCode() {
	wm.InitCode = nil
	wm.InitCodeExpiry = 0
	wm.InitCodeCreationDate = time.Time{}
	wm.InitCheckFailedCount = 0
}
func (wm *UserHumanWriteModel) SetInitCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
	wm.InitCode = code
	wm.InitCodeExpiry = expiry
	wm.InitCodeCreationDate = creationDate
	wm.InitCheckFailedCount = 0
}
func (wm *UserHumanWriteModel) EmptyEmailCode() {
	wm.EmailCode = nil
	wm.EmailCodeExpiry = 0
	wm.EmailCodeCreationDate = time.Time{}
	wm.EmailCheckFailedCount = 0
}
func (wm *UserHumanWriteModel) SetEMailCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
	wm.EmailCode = code
	wm.EmailCodeExpiry = expiry
	wm.EmailCodeCreationDate = creationDate
	wm.EmailCheckFailedCount = 0
}
func (wm *UserHumanWriteModel) EmptyPhoneCode() {
	wm.PhoneCode = nil
	wm.PhoneCodeExpiry = 0
	wm.PhoneCodeCreationDate = time.Time{}
	wm.PhoneCheckFailedCount = 0
}
func (wm *UserHumanWriteModel) SetPhoneCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
	wm.PhoneCode = code
	wm.PhoneCodeExpiry = expiry
	wm.PhoneCodeCreationDate = creationDate
	wm.PhoneCheckFailedCount = 0
}
func (wm *UserHumanWriteModel) EmptyPasswordCode() {
	wm.PasswordCode = nil
	wm.PasswordCodeExpiry = 0
	wm.PasswordCodeCreationDate = time.Time{}
}
func (wm *UserHumanWriteModel) SetPasswordCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
	wm.PasswordCode = code
	wm.PasswordCodeExpiry = expiry
	wm.PasswordCodeCreationDate = creationDate
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
		user.HumanInitializedCheckFailedType,
		user.UserV1InitializedCheckFailedType,

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
			user.HumanEmailVerificationFailedType,
			user.UserV1EmailVerificationFailedType,
		)
	}
	if wm.PhoneWriteModel {
		eventTypes = append(eventTypes,
			user.UserV1PhoneChangedType,
			user.HumanPhoneChangedType,
			user.UserV1PhoneCodeAddedType,
			user.HumanPhoneCodeAddedType,

			user.UserV1PhoneVerifiedType,
			user.HumanPhoneVerifiedType,
			user.HumanPhoneVerificationFailedType,
			user.UserV1PhoneVerificationFailedType,

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
			user.HumanPasswordHashUpdatedType,

			user.HumanPasswordChangedType,
			user.UserV1PasswordChangedType,
			user.HumanPasswordCodeAddedType,
			user.UserV1PasswordCodeAddedType,

			user.HumanPasswordCheckFailedType,
			user.UserV1PasswordCheckFailedType,
			user.HumanPasswordCheckSucceededType,
			user.UserV1PasswordCheckSucceededType,
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
	wm.UserState = domain.UserStateInitial
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
	wm.UserState = domain.UserStateInitial
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
	if !isVerified || wm.IsEmailVerified == isVerified {
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
	if !isVerified || wm.IsPhoneVerified == isVerified {
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
