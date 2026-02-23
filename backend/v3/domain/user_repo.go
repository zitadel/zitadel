package domain

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate mockgen -typed -package domainmock -destination ./mock/user.mock.go . UserRepository

type UserRepository interface {
	Repository
	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*User, error)
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*User, error)
	Create(ctx context.Context, client database.QueryExecutor, user *User) error
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)

	Human() HumanUserRepository
	Machine() MachineUserRepository

	LoadMetadata() UserRepository
	LoadPasskeys() UserRepository
	LoadIdentityProviderLinks() UserRepository
	LoadVerifications() UserRepository
	LoadPATs() UserRepository
	LoadKeys() UserRepository

	userConditions
	userChanges
	userColumns
}

type userConditions interface {
	PrimaryKeyCondition(instanceID, userID string) database.Condition
	InstanceIDCondition(instanceID string) database.Condition
	OrgIDCondition(orgID string) database.Condition
	IDCondition(userID string) database.Condition
	UsernameCondition(op database.TextOperation, username string) database.Condition
	StateCondition(state UserState) database.Condition
	TypeCondition(userType UserType) database.Condition
	LoginNameCondition(op database.TextOperation, loginName string) database.Condition

	userMetadataConditions
}

type userMetadataConditions interface {
	ExistsMetadata(condition database.Condition) database.Condition
	MetadataKeyCondition(op database.TextOperation, key string) database.Condition
	MetadataValueCondition(op database.BytesOperation, value []byte) database.Condition
}

type userChanges interface {
	// SetUsername sets the username field
	SetUsername(username string) database.Change
	// SetState sets the state field
	// [UserStateUnspecified] is not allowed and will result in no change
	SetState(state UserState) database.Change
	// SetUpdatedAt sets the updated at field
	// This is used to replay events
	SetUpdatedAt(updatedAt time.Time) database.Change
	// AddMetadata adds metadata for the user
	AddMetadata(metadata ...*Metadata) database.Change
	// RemoveMetadata removes the metadata for the user
	// Use one of the following conditions to specify which metadata to remove:
	//  - [userMetadataConditions.MetadataKeyCondition]
	//  - [userMetadataConditions.MetadataValueCondition]
	//  - nil to remove all metadata for the user
	RemoveMetadata(condition database.Condition) database.Change
}

type userColumns interface {
	PrimaryKeyColumns() []database.Column
	UsernameColumn() database.Column
	StateColumn() database.Column
	TypeColumn() database.Column
	CreatedAtColumn() database.Column
}

//go:generate mockgen -typed -package domainmock -destination ./mock/human_user.mock.go . HumanUserRepository

type HumanUserRepository interface {
	// Update updates the user
	// It ensures that updates are only applied to user
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	humanConditions
	humanChanges
	humanColumns
}

type humanConditions interface {
	userConditions
	FirstNameCondition(op database.TextOperation, firstName string) database.Condition
	LastNameCondition(op database.TextOperation, lastName string) database.Condition
	NicknameCondition(op database.TextOperation, nickname string) database.Condition
	DisplayNameCondition(op database.TextOperation, displayName string) database.Condition
	PreferredLanguageCondition(lang language.Tag) database.Condition
	EmailCondition(op database.TextOperation, email string) database.Condition
	PhoneCondition(op database.TextOperation, phone string) database.Condition
	PasskeyIDCondition(passkeyID string) database.Condition
	IdentityProviderIDCondition(idpID string) database.Condition
	ProvidedUserIDCondition(providedUserID string) database.Condition
	ProvidedUsernameCondition(username string) database.Condition

	humanPasskeyConditions
}

type humanChanges interface {
	userChanges
	// SetFirstName sets the first name field
	SetFirstName(firstName string) database.Change
	// SetLastName sets the last name field
	SetLastName(lastName string) database.Change
	// SetNickname sets the nickname field
	SetNickname(nickname string) database.Change
	// SetDisplayName sets the display name field
	SetDisplayName(displayName string) database.Change
	// SetPreferredLanguage sets the preferred language field
	// If preferredLanguage is [language.Und], it will be set to NULL in the database
	SetPreferredLanguage(preferredLanguage language.Tag) database.Change
	// SetGender sets the gender field
	// If gender is [HumanGenderUnspecified], it will be set to NULL in the database
	SetGender(gender HumanGender) database.Change
	// SetAvatarKey sets the avatar key field
	// If avatarKey is nil, it will be set to NULL in the database
	SetAvatarKey(avatarKey *string) database.Change
	// SetMultifactorInitializationSkippedAt sets the multifactor initialization skipped at field
	SetMultifactorInitializationSkippedAt(skippedAt time.Time) database.Change

	// SetPassword sets the password based on the verification
	SetPassword(verification VerificationType) database.Change
	// SetPasswordChangeRequired sets whether a password change is required
	SetPasswordChangeRequired(required bool) database.Change
	// CheckPassword sets the password check based on the check type
	//  * [CheckTypeFailed] increments failed attempts
	//  * [CheckTypeSucceeded] to mark the check as succeeded, removes the check and updates verified at time
	CheckPassword(check PasswordCheckType) database.Change

	// SetVerification sets the verification based on the verification type
	SetVerification(id string, verification VerificationType) database.Change

	humanEmailChanges
	humanPhoneChanges
	humanTOTPChanges
	identityProviderLinkChanges
	humanPasskeyChanges
}

type humanEmailChanges interface {
	// SetEmail sets the email based on the verification
	// 	* [VerificationTypeInit] to initialize email verification, previously verified email remains verified
	// 	* [VerificationTypeVerified] to mark email as verified, a verification must exist
	// 	* [VerificationTypeUpdate] to update email verification, a verification must exist (e.g. resend code)
	// 	* [VerificationTypeSkipped] to skip email verification, existing verification is removed (e.g. admin set email)
	SetEmail(verification VerificationType) database.Change
	// CheckEmailOTP sets the OTP Email based on the check
	//  * [CheckTypeInit] to initialize a new check, previous check is overwritten
	//  * [CheckTypeFailed] increments failed attempts
	//  * [CheckTypeSucceeded] to mark the check as succeeded, removes the check and updates verified at time
	CheckEmailOTP(check CheckType) database.Change
	// EnableEmailOTPAt enables the OTP Email
	// If enabledAt is zero, it will be set to NOW()
	EnableEmailOTPAt(enabledAt time.Time) database.Change
	// EnableEmailOTP sets the enabled at time to NOW()
	EnableEmailOTP() database.Change
	// DisableEmailOTP clears the enabled at time
	DisableEmailOTP() database.Change
}

type humanPhoneChanges interface {
	// SetPhone sets the phone based on the verification
	// 	* [VerificationTypeInit] to initialize phone verification, previously verified phone remains verified
	// 	* [VerificationTypeVerified] to mark phone as verified, a verification must exist
	// 	* [VerificationTypeUpdate] to update phone verification, a verification must exist (e.g. resend code)
	// 	* [VerificationTypeSkipped] to skip phone verification, existing verification is removed (e.g. admin set phone)
	SetPhone(verification VerificationType) database.Change
	// RemovePhone removes the phone number
	RemovePhone() database.Change
	// CheckSMSOTP sets the OTP SMS based on the check
	//  * [CheckTypeInit] to initialize a new check, previous check is overwritten
	//  * [CheckTypeFailed] increments failed attempts
	//  * [CheckTypeSucceeded] to mark the check as succeeded, removes the check and updates verified at time
	CheckSMSOTP(check CheckType) database.Change
	// EnableSMSOTPAt enables the OTP SMS
	// If enabledAt is zero, it will be set to NOW()
	EnableSMSOTPAt(enabledAt time.Time) database.Change
	// EnableSMSOTP sets the enabled at time to NOW()
	EnableSMSOTP() database.Change
	// DisableSMSOTP clears the enabled at time
	DisableSMSOTP() database.Change
}

type humanTOTPChanges interface {
	// SetTOTP changes the TOTP secret based on the verification
	// 	* [VerificationTypeInit] to initialize totp verification, previously verified totp remains verified
	// 	* [VerificationTypeVerified] to mark totp as verified, a verification must exist
	// 	* [VerificationTypeUpdate] to update totp verification, a verification must exist (e.g. resend code)
	// 	* [VerificationTypeSkipped] to skip totp verification, existing verification is removed (e.g. admin set totp)
	SetTOTP(verification VerificationType) database.Change
	// SetTOTPCheck sets the TOTP check based on the check type
	//  * [CheckTypeFailed] increments failed attempts
	//  * [CheckTypeSucceeded] to mark the check as succeeded, removes the check and updates verified at time
	//  * [CheckTypeInit] is not allowed for TOTP
	CheckTOTP(check CheckType) database.Change
	// RemoveTOTP removes the TOTP
	RemoveTOTP() database.Change
}

type identityProviderLinkChanges interface {
	// AddIdentityProviderLink adds or updates an identity provider link
	AddIdentityProviderLink(link *IdentityProviderLink) database.Change
	UpdateIdentityProviderLink(changes ...database.Change) database.Change
	// RemoveIdentityProviderLink removes an identity provider link based on the condition
	RemoveIdentityProviderLink(providerID, providedUserID string) database.Change

	// SetIdentityProviderLinkUsername sets the username for an identity provider link
	SetIdentityProviderLinkUsername(providerID, providedUserID, username string) database.Change
	// SetIdentityProviderLinkProvidedID sets the provided user ID for an identity provider link
	SetIdentityProviderLinkProvidedID(providerID, currentProvidedUserID, newProvidedUserID string) database.Change
}

type humanColumns interface {
	userColumns
	FirstNameColumn() database.Column
	LastNameColumn() database.Column
	NicknameColumn() database.Column
	DisplayNameColumn() database.Column
	EmailColumn() database.Column
}

//go:generate mockgen -typed -package domainmock -destination ./mock/machine_user.mock.go . MachineUserRepository

type MachineUserRepository interface {
	// Update updates the service account
	// It ensures that updates are only applied to service accounts
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	machineConditions
	machineChanges
	machineColumns
}

type machineConditions interface {
	userConditions
}

type machineChanges interface {
	userChanges
	// SetName sets the name field
	SetName(name string) database.Change
	// SetDescription sets the description field
	SetDescription(description string) database.Change
	// SetSecret sets the secret field
	// nil will set the secret to NULL in the database
	SetSecret(secret *string) database.Change
	// SetAccessTokenType sets the personal access token type field
	SetAccessTokenType(tokenType PersonalAccessTokenType) database.Change

	// AddKey adds a key for the service account
	AddKey(key *MachineKey) database.Change
	// RemoveKey removes a key for the service account
	RemoveKey(id string) database.Change

	// AddPersonalAccessToken adds a personal access token for the service account
	AddPersonalAccessToken(pat *PersonalAccessToken) database.Change
	// RemovePersonalAccessToken removes a personal access token for the service account
	RemovePersonalAccessToken(id string) database.Change
}

type machineColumns interface {
	userColumns
}

type humanPasskeyConditions interface {
	PasskeyIDCondition(passkeyID string) database.Condition
	PasskeyKeyIDCondition(keyID string) database.Condition
	PasskeyChallengeCondition(challenge string) database.Condition
	PasskeyTypeCondition(passkeyType PasskeyType) database.Condition
}

type humanPasskeyChanges interface {
	// AddPasskey adds a passkey for the user
	AddPasskey(passkey *Passkey) database.Change
	// UpdatePasskey updates a passkey for the user.
	// The condition must specify which passkey to update, use [humanPasskeyConditions].
	// Use [humanPasskeyFieldChanges] to specify the fields to update.
	UpdatePasskey(condition database.Condition, changes ...database.Change) database.Change
	// RemovePasskey removes a passkey for the user
	// The condition must specify which passkey to update, use [humanPasskeyConditions].
	RemovePasskey(condition database.Condition) database.Change

	passkeyFieldChanges
}

type passkeyFieldChanges interface {
	// SetPasskeyKeyID sets the key ID field
	SetPasskeyKeyID(keyID []byte) database.Change
	// SetPasskeyPublicKey sets the public key field
	SetPasskeyPublicKey(publicKey []byte) database.Change
	// SetPasskeyAttestationType sets the attestation type field
	SetPasskeyAttestationType(attestationType string) database.Change
	// SetPasskeyAuthenticatorAttestationGUID sets the AAGUID field
	SetPasskeyAuthenticatorAttestationGUID(aaguid []byte) database.Change
	// SetPasskeyName sets the name field
	SetPasskeyName(name string) database.Change
	// SetPasskeySignCount sets the sign count field
	SetPasskeySignCount(signCount uint32) database.Change
	// SetPasskeyUpdatedAt sets the updated at field
	SetPasskeyUpdatedAt(updatedAt time.Time) database.Change
	// SetPasskeyVerifiedAt sets the verified at field
	// If verifiedAt is zero, it will be set to NOW()
	SetPasskeyVerifiedAt(verifiedAt time.Time) database.Change
}
