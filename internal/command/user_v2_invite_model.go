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

	InviteCode              *crypto.CryptoValue
	InviteCodeCreationDate  time.Time
	InviteCodeExpiry        time.Duration
	InviteCheckFailureCount uint8

	ApplicationName string
	AuthRequestID   string
	URLTemplate     string
	CodeReturned    bool
	EmailVerified   bool

	// authMethods holds the primary auth methods (password, passkeys, IDP links)
	// the user currently has, keyed uniquely per method. It is tracked as the
	// current set rather than a sticky flag, so that removing all auth methods
	// re-enables the invite flow, mirroring what ListUserAuthMethodTypes reports.
	authMethods map[authMethodKey]struct{}

	UserState domain.UserState
}

// authMethodKey uniquely identifies a single primary auth method within the
// authMethods set.
// ID components (e.g. an externalUserID supplied by an external IDP) cannot
// collide across method instances regardless of their contents.
type authMethodKey struct {
	kind           string
	idpConfigID    string
	externalUserID string
	passkeyTokenID string
}

func passwordAuthMethodKey() authMethodKey {
	return authMethodKey{kind: "password"}
}

func idpAuthMethodKey(idpConfigID, externalUserID string) authMethodKey {
	return authMethodKey{kind: "idp", idpConfigID: idpConfigID, externalUserID: externalUserID}
}

func passkeyAuthMethodKey(tokenID string) authMethodKey {
	return authMethodKey{kind: "passkey", passkeyTokenID: tokenID}
}

// AuthMethodSet reports whether the user currently has at least one primary
// auth method (password, passkey or IDP link) set up.
func (wm *UserV2InviteWriteModel) AuthMethodSet() bool {
	return len(wm.authMethods) > 0
}

func (wm *UserV2InviteWriteModel) CreationAllowed() bool {
	return !wm.AuthMethodSet()
}

func (wm *UserV2InviteWriteModel) addAuthMethod(key authMethodKey) {
	if wm.authMethods == nil {
		wm.authMethods = make(map[authMethodKey]struct{})
	}
	wm.authMethods[key] = struct{}{}
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
		case *user.HumanAddedEvent:
			wm.EmptyInviteCode()
			wm.ApplicationName = ""
			wm.AuthRequestID = ""
			wm.URLTemplate = ""
			wm.CodeReturned = false
			wm.EmailVerified = false
			wm.UserState = domain.UserStateActive
			wm.authMethods = nil
			if crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash) != "" {
				wm.addAuthMethod(passwordAuthMethodKey())
			}
		case *user.HumanRegisteredEvent:
			wm.EmptyInviteCode()
			wm.ApplicationName = ""
			wm.AuthRequestID = ""
			wm.URLTemplate = ""
			wm.CodeReturned = false
			wm.EmailVerified = false
			wm.UserState = domain.UserStateActive
			wm.authMethods = nil
			if crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash) != "" {
				wm.addAuthMethod(passwordAuthMethodKey())
			}
		case *user.MachineAddedEvent:
			wm.EmptyInviteCode()
			wm.ApplicationName = ""
			wm.AuthRequestID = ""
			wm.URLTemplate = ""
			wm.CodeReturned = false
			wm.EmailVerified = false
			wm.authMethods = nil
		case *user.HumanInviteCodeAddedEvent:
			wm.SetInviteCode(e.Code, e.Expiry, e.CreationDate())
			wm.URLTemplate = e.URLTemplate
			wm.CodeReturned = e.CodeReturned
			wm.ApplicationName = e.ApplicationName
			wm.AuthRequestID = e.AuthRequestID
		case *user.HumanInviteCheckSucceededEvent:
			wm.EmptyInviteCode()
		case *user.HumanInviteCheckFailedEvent:
			wm.InviteCheckFailureCount++
			if wm.InviteCheckFailureCount >= 3 || crypto.IsCodeExpired(wm.InviteCodeCreationDate, wm.InviteCodeExpiry) { //TODO: make failure count comparison with wm.InviteCheckFailureCount configurable?
				// invalidate the invite code after attempting to verify an expired code, or a wrong code three or more times
				// so that a new invite code can be created for this user
				wm.EmptyInviteCode()
				wm.CodeReturned = false
			}
		case *user.HumanEmailVerifiedEvent:
			wm.EmailVerified = true
			wm.EmptyInviteCode()
		case *user.UserLockedEvent:
			wm.UserState = domain.UserStateLocked
		case *user.UserUnlockedEvent:
			wm.UserState = domain.UserStateActive
		case *user.UserDeactivatedEvent:
			wm.UserState = domain.UserStateInactive
		case *user.UserReactivatedEvent:
			wm.UserState = domain.UserStateActive
		case *user.UserRemovedEvent:
			wm.EmptyInviteCode()
			wm.ApplicationName = ""
			wm.AuthRequestID = ""
			wm.URLTemplate = ""
			wm.CodeReturned = false
			wm.EmailVerified = false
			wm.authMethods = nil
			wm.UserState = domain.UserStateDeleted
		case *user.HumanPasswordChangedEvent:
			wm.addAuthMethod(passwordAuthMethodKey())
		case *user.UserIDPLinkAddedEvent:
			wm.addAuthMethod(idpAuthMethodKey(e.IDPConfigID, e.ExternalUserID))
		case *user.UserIDPLinkRemovedEvent:
			delete(wm.authMethods, idpAuthMethodKey(e.IDPConfigID, e.ExternalUserID))
		case *user.UserIDPLinkCascadeRemovedEvent:
			delete(wm.authMethods, idpAuthMethodKey(e.IDPConfigID, e.ExternalUserID))
		case *user.HumanPasswordlessVerifiedEvent:
			wm.addAuthMethod(passkeyAuthMethodKey(e.WebAuthNTokenID))
		case *user.HumanPasswordlessRemovedEvent:
			delete(wm.authMethods, passkeyAuthMethodKey(e.WebAuthNTokenID))
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UserV2InviteWriteModel) SetInviteCode(code *crypto.CryptoValue, expiry time.Duration, creationDate time.Time) {
	wm.InviteCode = code
	wm.InviteCodeExpiry = expiry
	wm.InviteCodeCreationDate = creationDate
	wm.InviteCheckFailureCount = 0
}

func (wm *UserV2InviteWriteModel) EmptyInviteCode() {
	wm.InviteCode = nil
	wm.InviteCodeExpiry = 0
	wm.InviteCodeCreationDate = time.Time{}
	wm.InviteCheckFailureCount = 0
}
func (wm *UserV2InviteWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			user.UserV1AddedType,
			user.HumanAddedType,
			user.UserV1RegisteredType,
			user.HumanRegisteredType,
			user.MachineAddedEventType,
			user.HumanInviteCodeAddedType,
			user.HumanInviteCheckSucceededType,
			user.HumanInviteCheckFailedType,
			user.UserV1EmailVerifiedType,
			user.HumanEmailVerifiedType,
			user.UserLockedType,
			user.UserUnlockedType,
			user.UserDeactivatedType,
			user.UserReactivatedType,
			user.UserRemovedType,
			user.HumanPasswordChangedType,
			user.UserV1PasswordChangedType,
			user.UserIDPLinkAddedType,
			user.UserIDPLinkRemovedType,
			user.UserIDPLinkCascadeRemovedType,
			user.HumanPasswordlessTokenVerifiedType,
			user.HumanPasswordlessTokenRemovedType,
		).Builder()
	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func (wm *UserV2InviteWriteModel) Aggregate() *user.Aggregate {
	return user.NewAggregate(wm.AggregateID, wm.ResourceOwner)
}
