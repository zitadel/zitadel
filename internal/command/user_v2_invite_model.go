package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type UserV2InviteWriteModel struct {
	eventstore.WriteModel

	InviteCode             *crypto.CryptoValue
	InviteCodeCreationDate time.Time
	InviteCodeExpiry       time.Duration

	ApplicationName string
	AuthRequestID   string

	InitCode             *crypto.CryptoValue
	InitCodeCreationDate time.Time
	InitCodeExpiry       time.Duration

	//PasswordEncodedHash      string
	//PasswordChangeRequired   bool
	//PasswordCode             *crypto.CryptoValue
	//PasswordCodeCreationDate time.Time
	//PasswordCodeExpiry       time.Duration

	//Email                 domain.EmailAddress
	//IsEmailVerified       bool
	//EmailCode             *crypto.CryptoValue
	//EmailCodeCreationDate time.Time
	//EmailCodeExpiry       time.Duration

	//Phone                 domain.PhoneNumber
	//IsPhoneVerified       bool
	//PhoneCode             *crypto.CryptoValue
	//PhoneCodeCreationDate time.Time
	//PhoneCodeExpiry       time.Duration

	UserState domain.UserState

	//IDPLinks []*domain.UserIDPLink
}

func newUserV2InviteWriteModel(userID, orgID string) *UserV2InviteWriteModel {
	return &UserV2InviteWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: orgID,
		},
	}
}

func (wm *UserV2InviteWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent,
			*user.HumanRegisteredEvent:
			wm.UserState = domain.UserStateActive
		case *user.HumanInitialCodeAddedEvent:
			wm.UserState = domain.UserStateInitial
			wm.SetInitCode(e.Code, e.Expiry, e.CreationDate())
		case *user.HumanInitializedCheckSucceededEvent:
			wm.UserState = domain.UserStateActive
			wm.EmptyInitCode()
		case *user.HumanInviteCodeAddedEvent:
			wm.SetInviteCode(e.Code, e.Expiry, e.CreationDate())
		//case *user.HumanInviteCheckSucceededEvent:
		//	wm.UserState = domain.UserStateActive
		//	wm.EmptyInitCode()

		////case *user.MachineAddedEvent:
		////	//wm.UserName = e.UserName
		////	//wm.Name = e.Name
		////	//wm.Description = e.Description
		////	//wm.AccessTokenType = e.AccessTokenType
		////	wm.UserState = domain.UserStateActive
		//
		//case *user.HumanEmailChangedEvent:
		//	wm.Email = e.EmailAddress
		//	wm.IsEmailVerified = false
		//	wm.EmptyEmailCode()
		//case *user.HumanEmailCodeAddedEvent:
		//	wm.IsEmailVerified = false
		//	wm.SetEMailCode(e.Code, e.Expiry, e.CreationDate())
		//case *user.HumanEmailVerifiedEvent:
		//	wm.IsEmailVerified = true
		//	wm.EmptyEmailCode()
		////case *user.HumanEmailVerificationFailedEvent:
		////	wm.EmailCheckFailedCount += 1
		//
		//case *user.HumanPhoneChangedEvent:
		//	wm.IsPhoneVerified = false
		//	wm.Phone = e.PhoneNumber
		//	wm.EmptyPhoneCode()
		//case *user.HumanPhoneCodeAddedEvent:
		//	wm.IsPhoneVerified = false
		//	wm.SetPhoneCode(e.Code, e.Expiry, e.CreationDate())
		//case *user.HumanPhoneVerifiedEvent:
		//	wm.IsPhoneVerified = true
		//	wm.EmptyPhoneCode()
		////case *user.HumanPhoneVerificationFailedEvent:
		////	wm.PhoneCheckFailedCount += 1
		//case *user.HumanPhoneRemovedEvent:
		//	wm.EmptyPhoneCode()
		//	wm.Phone = ""
		//	wm.IsPhoneVerified = false

		case *user.UserLockedEvent:
			wm.UserState = domain.UserStateLocked
		case *user.UserUnlockedEvent:
			//wm.PasswordCheckFailedCount = 0
			wm.UserState = domain.UserStateActive

		case *user.UserDeactivatedEvent:
			wm.UserState = domain.UserStateInactive
		case *user.UserReactivatedEvent:
			wm.UserState = domain.UserStateActive

		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		//
		//case *user.HumanPasswordHashUpdatedEvent:
		//	wm.PasswordEncodedHash = e.EncodedHash
		//case *user.HumanPasswordCheckFailedEvent:
		//	wm.PasswordCheckFailedCount += 1
		//case *user.HumanPasswordCheckSucceededEvent:
		//	wm.PasswordCheckFailedCount = 0
		//case *user.HumanPasswordChangedEvent:
		//	wm.PasswordEncodedHash = crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash)
		//	wm.PasswordChangeRequired = e.ChangeRequired
		//	wm.EmptyPasswordCode()
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

func (wm *UserV2InviteWriteModel) AddIDPLink(configID, displayName, externalUserID string) {
	//wm.IDPLinks = append(wm.IDPLinks, &domain.UserIDPLink{IDPConfigID: configID, DisplayName: displayName, ExternalUserID: externalUserID})
}

func (wm *UserV2InviteWriteModel) RemoveIDPLink(configID, externalUserID string) {
	//idx, _ := wm.IDPLinkByID(configID, externalUserID)
	//if idx < 0 {
	//	return
	//}
	//copy(wm.IDPLinks[idx:], wm.IDPLinks[idx+1:])
	//wm.IDPLinks[len(wm.IDPLinks)-1] = nil
	//wm.IDPLinks = wm.IDPLinks[:len(wm.IDPLinks)-1]
}

func (wm *UserV2InviteWriteModel) EmptyInitCode() {
	wm.InitCode = nil
	wm.InitCodeExpiry = 0
	wm.InitCodeCreationDate = time.Time{}
	//wm.InitCheckFailedCount = 0
}
func (wm *UserV2InviteWriteModel) SetInitCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
	wm.InitCode = code
	wm.InitCodeExpiry = expiry
	wm.InitCodeCreationDate = creationDate
	//wm.InitCheckFailedCount = 0
}

func (wm *UserV2InviteWriteModel) SetInviteCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
	wm.InviteCode = code
	wm.InviteCodeExpiry = expiry
	wm.InviteCodeCreationDate = creationDate
	//wm.InitCheckFailedCount = 0
}

//	func (wm *UserV2InviteWriteModel) EmptyEmailCode() {
//		wm.EmailCode = nil
//		wm.EmailCodeExpiry = 0
//		wm.EmailCodeCreationDate = time.Time{}
//		wm.EmailCheckFailedCount = 0
//	}
//
//	func (wm *UserV2InviteWriteModel) SetEMailCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
//		wm.EmailCode = code
//		wm.EmailCodeExpiry = expiry
//		wm.EmailCodeCreationDate = creationDate
//		wm.EmailCheckFailedCount = 0
//	}
//
//	func (wm *UserV2InviteWriteModel) EmptyPhoneCode() {
//		wm.PhoneCode = nil
//		wm.PhoneCodeExpiry = 0
//		wm.PhoneCodeCreationDate = time.Time{}
//		wm.PhoneCheckFailedCount = 0
//	}
//
//	func (wm *UserV2InviteWriteModel) SetPhoneCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
//		wm.PhoneCode = code
//		wm.PhoneCodeExpiry = expiry
//		wm.PhoneCodeCreationDate = creationDate
//		wm.PhoneCheckFailedCount = 0
//	}
//
//	func (wm *UserV2InviteWriteModel) EmptyPasswordCode() {
//		wm.PasswordCode = nil
//		wm.PasswordCodeExpiry = 0
//		wm.PasswordCodeCreationDate = time.Time{}
//	}
func (wm *UserV2InviteWriteModel) SetPasswordCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
	//wm.PasswordCode = code
	//wm.PasswordCodeExpiry = expiry
	//wm.PasswordCodeCreationDate = creationDate
}

func (wm *UserV2InviteWriteModel) Query() *eventstore.SearchQueryBuilder {
	// remove events are always processed
	// and username is based for machine and human
	eventTypes := []eventstore.EventType{
		user.UserRemovedType,
		user.UserUserNameChangedType,
		user.UserV1AddedType,
		user.HumanAddedType,
		user.UserV1RegisteredType,
		user.HumanRegisteredType,
		user.MachineChangedEventType,
		user.MachineAddedEventType,
		user.UserV1EmailChangedType,
		user.HumanEmailChangedType,
		user.UserV1EmailCodeAddedType,
		user.HumanEmailCodeAddedType,
		user.UserV1EmailVerifiedType,
		user.HumanEmailVerifiedType,
		user.HumanEmailVerificationFailedType,
		user.UserV1EmailVerificationFailedType,
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
		user.HumanAvatarAddedType,
		user.HumanAvatarRemovedType,
		user.HumanPasswordHashUpdatedType,
		user.HumanPasswordChangedType,
		user.UserV1PasswordChangedType,
		user.HumanPasswordCodeAddedType,
		user.UserV1PasswordCodeAddedType,
		user.HumanPasswordCheckFailedType,
		user.UserV1PasswordCheckFailedType,
		user.HumanPasswordCheckSucceededType,
		user.UserV1PasswordCheckSucceededType,
		user.UserIDPLinkAddedType,
		user.UserIDPLinkRemovedType,
		user.UserIDPLinkCascadeRemovedType,
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

func (wm *UserV2InviteWriteModel) Aggregate() *user.Aggregate {
	return user.NewAggregate(wm.AggregateID, wm.ResourceOwner)
}

//func (wm *UserV2InviteWriteModel) IDPLinkByID(idpID, externalUserID string) (idx int, idp *domain.UserIDPLink) {
//	for idx, idp = range wm.IDPLinks {
//		if idp.IDPConfigID == idpID && idp.ExternalUserID == externalUserID {
//			return idx, idp
//		}
//	}
//	return -1, nil
//}
