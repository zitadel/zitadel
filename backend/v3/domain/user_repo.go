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
	OrganizationIDCondition(orgID string) database.Condition
	IDCondition(userID string) database.Condition
	UsernameCondition(op database.TextOperation, username string) database.Condition
	StateCondition(state UserState) database.Condition
	TypeCondition(userType UserType) database.Condition
	LoginNameCondition(op database.TextOperation, loginName string) database.Condition

	ExistsMetadata(condition database.Condition) database.Condition
	MetadataConditions() UserMetadataConditions
}

type UserMetadataConditions interface {
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
	InstanceIDColumn() database.Column
	IDColumn() database.Column
	UsernameColumn() database.Column
	StateColumn() database.Column
	TypeColumn() database.Column
	CreatedAtColumn() database.Column
}

//go:generate mockgen -typed -package domainmock -destination ./mock/human_user.mock.go . HumanUserRepository

type HumanUserRepository interface {
	// Update updates the human user
	// It ensures that updates are only applied to human users
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
	LinkedIdentityProviderIDCondition(idpID string) database.Condition
	ProvidedUserIDCondition(providedUserID string) database.Condition
	ProvidedUsernameCondition(username string) database.Condition

	ExistsPasskey(condition database.Condition) database.Condition
	PasskeyConditions() HumanPasskeyConditions
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
	// SkipMultifactorInitializationAt sets the multifactor initialization skipped at field
	SkipMultifactorInitializationAt(skippedAt time.Time) database.Change
	// SkipMultifactorInitializationAt sets the multifactor initialization skipped at field
	SkipMultifactorInitialization() database.Change

	// SetVerification sets the verification based on the verification type
	SetVerification(verification VerificationType) database.Change

	humanPasswordChanges
	humanEmailChanges
	humanPhoneChanges
	humanTOTPChanges
	identityProviderLinkChanges
	humanPasskeyChanges
}

type humanPasswordChanges interface {
	// SetPassword sets the password based on the verification
	SetPassword(verification VerificationType) database.Change
	// SetPasswordChangeRequired sets whether a password change is required
	SetPasswordChangeRequired(required bool) database.Change
	// SetLastSuccessfulPasswordCheck sets the last successful password check time
	// If checkedAt is zero, it will be set to NOW()
	SetLastSuccessfulPasswordCheck(checkedAt time.Time) database.Change
	// IncrementPasswordFailedAttempts increments the password failed attempts
	IncrementPasswordFailedAttempts() database.Change
	// ResetPasswordFailedAttempts resets the password failed attempts
	ResetPasswordFailedAttempts() database.Change
}

type humanEmailChanges interface {
	// SetEmail sets the email based on the verification
	// 	* [VerificationTypeInit] to initialize email verification, previously verified email remains verified
	// 	* [VerificationTypeVerified] to mark email as verified, a verification must exist
	// 	* [VerificationTypeUpdate] to update email verification, a verification must exist (e.g. resend code)
	// 	* [VerificationTypeSkipped] to skip email verification, existing verification is removed (e.g. admin set email)
	SetEmail(verification VerificationType) database.Change
	// EnableEmailOTPAt enables the email OTP
	// If enabledAt is zero, it will be set to NOW()
	EnableEmailOTPAt(enabledAt time.Time) database.Change
	// EnableEmailOTP sets the enabled at time to NOW()
	EnableEmailOTP() database.Change
	// DisableEmailOTP clears the enabled at time
	DisableEmailOTP() database.Change
	// SetLastSuccessfulEmailOTPCheck sets the last successful email OTP check time
	// If checkedAt is zero, it will be set to NOW()
	SetLastSuccessfulEmailOTPCheck(checkedAt time.Time) database.Change
	// IncrementEmailOTPFailedAttempts increments the email OTP failed attempts
	IncrementEmailOTPFailedAttempts() database.Change
	// ResetEmailOTPFailedAttempts resets the email OTP failed attempts
	ResetEmailOTPFailedAttempts() database.Change
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
	// EnableSMSOTPAt enables the SMS OTP
	// If enabledAt is zero, it will be set to NOW()
	EnableSMSOTPAt(enabledAt time.Time) database.Change
	// EnableSMSOTP sets the enabled at time to NOW()
	EnableSMSOTP() database.Change
	// DisableSMSOTP clears the enabled at time
	DisableSMSOTP() database.Change
	// SetLastSuccessfulSMSOTPCheck sets the last successful SMS OTP check time
	// If checkedAt is zero, it will be set to NOW()
	SetLastSuccessfulSMSOTPCheck(checkedAt time.Time) database.Change
	// IncrementPhoneOTPFailedAttempts increments the phone OTP failed attempts
	IncrementPhoneOTPFailedAttempts() database.Change
	// ResetPhoneOTPFailedAttempts resets the phone OTP failed attempts
	ResetPhoneOTPFailedAttempts() database.Change
}

type humanTOTPChanges interface {
	// SetTOTP changes the TOTP secret and verification timestamp
	SetTOTPSecret(secret []byte) database.Change
	// SetTOTPVerifiedAt sets the TOTP verified at time
	// If verifiedAt is zero, it will be set to NOW()
	SetTOTPVerifiedAt(verifiedAt time.Time) database.Change
	// RemoveTOTP removes the TOTP
	RemoveTOTP() database.Change
	// SetLastSuccessfulTOTPCheck sets the last successful TOTP check time
	// If checkedAt is zero, it will be set to NOW()
	SetLastSuccessfulTOTPCheck(checkedAt time.Time) database.Change
}

type identityProviderLinkChanges interface {
	// AddIdentityProviderLink adds or updates an identity provider link
	AddIdentityProviderLink(link *IdentityProviderLink) database.Change
	UpdateIdentityProviderLink(condition database.Condition, changes ...database.Change) database.Change
	// RemoveIdentityProviderLink removes an identity provider link based on the condition
	RemoveIdentityProviderLink(providerID, providedUserID string) database.Change

	// SetIdentityProviderLinkUsername sets the username for an identity provider link
	SetIdentityProviderLinkUsername(username string) database.Change
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
	// Update updates the machine user
	// It ensures that updates are only applied to machine users
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
	SetAccessTokenType(tokenType AccessTokenType) database.Change

	// AddKey adds a key for the machine user
	AddKey(key *MachineKey) database.Change
	// RemoveKey removes a key for the machine user
	RemoveKey(id string) database.Change

	// AddPersonalAccessToken adds a personal access token for the machine user
	AddPersonalAccessToken(pat *PersonalAccessToken) database.Change
	// RemovePersonalAccessToken removes a personal access token for the machine user
	RemovePersonalAccessToken(id string) database.Change
}

type machineColumns interface {
	userColumns
}

type HumanPasskeyConditions interface {
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
