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

type UserV2WriteModel struct {
	eventstore.WriteModel

	UserName string

	MachineWriteModel bool
	Name              string
	Description       string
	AccessTokenType   domain.OIDCTokenType

	MachineSecretWriteModel bool
	ClientSecret            *crypto.CryptoValue

	ProfileWriteModel bool
	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	PreferredLanguage language.Tag
	Gender            domain.Gender

	AvatarWriteModel bool
	Avatar           string

	HumanWriteModel      bool
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

	IDPLinkWriteModel bool
	IDPLinks          []*domain.UserIDPLink
}

func NewUserExistsWriteModel(userID, resourceOwner string) *UserV2WriteModel {
	return newUserV2WriteModel(userID, resourceOwner, WithHuman(), WithMachine())
}

func NewUserStateWriteModel(userID, resourceOwner string) *UserV2WriteModel {
	return newUserV2WriteModel(userID, resourceOwner, WithHuman(), WithMachine(), WithState())
}

func NewUserRemoveWriteModel(userID, resourceOwner string) *UserV2WriteModel {
	return newUserV2WriteModel(userID, resourceOwner, WithHuman(), WithMachine(), WithState(), WithIDPLinks())
}

func NewUserHumanWriteModel(userID, resourceOwner string, profileWM, emailWM, phoneWM, passwordWM, avatarWM, idpLinks bool) *UserV2WriteModel {
	opts := []UserV2WMOption{WithHuman(), WithState()}
	if profileWM {
		opts = append(opts, WithProfile())
	}
	if emailWM {
		opts = append(opts, WithEmail())
	}
	if phoneWM {
		opts = append(opts, WithPhone())
	}
	if passwordWM {
		opts = append(opts, WithPassword())
	}
	if avatarWM {
		opts = append(opts, WithAvatar())
	}
	if idpLinks {
		opts = append(opts, WithIDPLinks())
	}
	return newUserV2WriteModel(userID, resourceOwner, opts...)
}

func newUserV2WriteModel(userID, resourceOwner string, opts ...UserV2WMOption) *UserV2WriteModel {
	wm := &UserV2WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}

	for _, optFunc := range opts {
		optFunc(wm)
	}
	return wm
}

type UserV2WMOption func(o *UserV2WriteModel)

func WithHuman() UserV2WMOption {
	return func(o *UserV2WriteModel) {
		o.HumanWriteModel = true
	}
}
func WithMachine() UserV2WMOption {
	return func(o *UserV2WriteModel) {
		o.MachineWriteModel = true
	}
}
func WithProfile() UserV2WMOption {
	return func(o *UserV2WriteModel) {
		o.ProfileWriteModel = true
	}
}
func WithEmail() UserV2WMOption {
	return func(o *UserV2WriteModel) {
		o.EmailWriteModel = true
	}
}
func WithPhone() UserV2WMOption {
	return func(o *UserV2WriteModel) {
		o.PhoneWriteModel = true
	}
}
func WithPassword() UserV2WMOption {
	return func(o *UserV2WriteModel) {
		o.PasswordWriteModel = true
	}
}
func WithState() UserV2WMOption {
	return func(o *UserV2WriteModel) {
		o.StateWriteModel = true
	}
}
func WithAvatar() UserV2WMOption {
	return func(o *UserV2WriteModel) {
		o.AvatarWriteModel = true
	}
}
func WithIDPLinks() UserV2WMOption {
	return func(o *UserV2WriteModel) {
		o.IDPLinkWriteModel = true
	}
}

func (wm *UserV2WriteModel) Reduce() error {
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

		case *user.MachineChangedEvent:
			if e.Name != nil {
				wm.Name = *e.Name
			}
			if e.Description != nil {
				wm.Description = *e.Description
			}
			if e.AccessTokenType != nil {
				wm.AccessTokenType = *e.AccessTokenType
			}

		case *user.MachineAddedEvent:
			wm.UserName = e.UserName
			wm.Name = e.Name
			wm.Description = e.Description
			wm.AccessTokenType = e.AccessTokenType
			wm.UserState = domain.UserStateActive

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
		case *user.UserIDPLinkAddedEvent:
			wm.AddIDPLink(e.IDPConfigID, e.DisplayName, e.ExternalUserID)
		case *user.UserIDPLinkRemovedEvent:
			wm.RemoveIDPLink(e.IDPConfigID, e.ExternalUserID)
		case *user.UserIDPLinkCascadeRemovedEvent:
			wm.RemoveIDPLink(e.IDPConfigID, e.ExternalUserID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UserV2WriteModel) AddIDPLink(configID, displayName, externalUserID string) {
	wm.IDPLinks = append(wm.IDPLinks, &domain.UserIDPLink{IDPConfigID: configID, DisplayName: displayName, ExternalUserID: externalUserID})
}

func (wm *UserV2WriteModel) RemoveIDPLink(configID, externalUserID string) {
	idx, _ := wm.IDPLinkByID(configID, externalUserID)
	if idx < 0 {
		return
	}
	copy(wm.IDPLinks[idx:], wm.IDPLinks[idx+1:])
	wm.IDPLinks[len(wm.IDPLinks)-1] = nil
	wm.IDPLinks = wm.IDPLinks[:len(wm.IDPLinks)-1]
}

func (wm *UserV2WriteModel) EmptyInitCode() {
	wm.InitCode = nil
	wm.InitCodeExpiry = 0
	wm.InitCodeCreationDate = time.Time{}
	wm.InitCheckFailedCount = 0
}
func (wm *UserV2WriteModel) SetInitCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
	wm.InitCode = code
	wm.InitCodeExpiry = expiry
	wm.InitCodeCreationDate = creationDate
	wm.InitCheckFailedCount = 0
}
func (wm *UserV2WriteModel) EmptyEmailCode() {
	wm.EmailCode = nil
	wm.EmailCodeExpiry = 0
	wm.EmailCodeCreationDate = time.Time{}
	wm.EmailCheckFailedCount = 0
}
func (wm *UserV2WriteModel) SetEMailCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
	wm.EmailCode = code
	wm.EmailCodeExpiry = expiry
	wm.EmailCodeCreationDate = creationDate
	wm.EmailCheckFailedCount = 0
}
func (wm *UserV2WriteModel) EmptyPhoneCode() {
	wm.PhoneCode = nil
	wm.PhoneCodeExpiry = 0
	wm.PhoneCodeCreationDate = time.Time{}
	wm.PhoneCheckFailedCount = 0
}
func (wm *UserV2WriteModel) SetPhoneCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
	wm.PhoneCode = code
	wm.PhoneCodeExpiry = expiry
	wm.PhoneCodeCreationDate = creationDate
	wm.PhoneCheckFailedCount = 0
}
func (wm *UserV2WriteModel) EmptyPasswordCode() {
	wm.PasswordCode = nil
	wm.PasswordCodeExpiry = 0
	wm.PasswordCodeCreationDate = time.Time{}
}
func (wm *UserV2WriteModel) SetPasswordCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
	wm.PasswordCode = code
	wm.PasswordCodeExpiry = expiry
	wm.PasswordCodeCreationDate = creationDate
}

func (wm *UserV2WriteModel) Query() *eventstore.SearchQueryBuilder {
	// remove events are always processed
	// and username is based for machine and human
	eventTypes := []eventstore.EventType{
		user.UserRemovedType,
		user.UserUserNameChangedType,
	}

	if wm.HumanWriteModel {
		eventTypes = append(eventTypes,
			user.UserV1AddedType,
			user.HumanAddedType,
			user.UserV1RegisteredType,
			user.HumanRegisteredType,
		)
	}

	if wm.MachineWriteModel {
		eventTypes = append(eventTypes,
			user.MachineChangedEventType,
			user.MachineAddedEventType,
		)
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
			user.UserV1InitialCodeAddedType,
			user.HumanInitialCodeAddedType,

			user.UserV1InitializedCheckSucceededType,
			user.HumanInitializedCheckSucceededType,
			user.HumanInitializedCheckFailedType,
			user.UserV1InitializedCheckFailedType,

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
	if wm.IDPLinkWriteModel {
		eventTypes = append(eventTypes,
			user.UserIDPLinkAddedType,
			user.UserIDPLinkRemovedType,
			user.UserIDPLinkCascadeRemovedType,
		)
	}

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
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

func (wm *UserV2WriteModel) reduceHumanAddedEvent(e *user.HumanAddedEvent) {
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

func (wm *UserV2WriteModel) reduceHumanRegisteredEvent(e *user.HumanRegisteredEvent) {
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

func (wm *UserV2WriteModel) reduceHumanProfileChangedEvent(e *user.HumanProfileChangedEvent) {
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

func (wm *UserV2WriteModel) Aggregate() *user.Aggregate {
	return user.NewAggregate(wm.AggregateID, wm.ResourceOwner)
}

func (wm *UserV2WriteModel) NewProfileChangedEvent(
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

func (wm *UserV2WriteModel) IDPLinkByID(idpID, externalUserID string) (idx int, idp *domain.UserIDPLink) {
	for idx, idp = range wm.IDPLinks {
		if idp.IDPConfigID == idpID && idp.ExternalUserID == externalUserID {
			return idx, idp
		}
	}
	return -1, nil
}
