package domain

import (
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
)

type User struct {
	InstanceID     string    `json:"instanceId,omitempty" db:"instance_id"`
	OrganizationID string    `json:"organizationId,omitempty" db:"organization_id"`
	ID             string    `json:"id,omitempty" db:"id"`
	Username       string    `json:"username,omitempty" db:"username"`
	State          UserState `json:"state,omitempty" db:"state"`
	CreatedAt      time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt,omitzero" db:"updated_at"`

	Machine  *MachineUser `json:"machine,omitempty" db:"-"`
	Human    *HumanUser   `json:"human,omitempty" db:"-"`
	Metadata []*Metadata  `json:"metadata,omitempty" db:"metadata"`
}

type MachineUser struct {
	Name            string          `json:"name,omitempty" db:"name"`
	Description     string          `json:"description,omitempty" db:"description"`
	Secret          []byte          `json:"secret,omitempty" db:"secret"`
	AccessTokenType AccessTokenType `json:"accessTokenType,omitempty" db:"access_token_type"`

	PATs []*PersonalAccessToken `json:"pats,omitempty" db:"pats"`
	Keys []*MachineKey          `json:"keys,omitempty" db:"keys"`
}

type HumanUser struct {
	FirstName         string       `json:"firstName,omitempty" db:"first_name"`
	LastName          string       `json:"lastName,omitempty" db:"last_name"`
	Nickname          string       `json:"nickname,omitempty" db:"nickname"`
	DisplayName       string       `json:"displayName,omitempty" db:"display_name"`
	PreferredLanguage language.Tag `json:"preferredLanguage,omitzero" db:"preferred_language"`
	Gender            HumanGender  `json:"gender,omitempty" db:"gender"`
	AvatarKey         string       `json:"avatarKey,omitempty" db:"avatar_key"`

	MultifactorInitializationSkippedAt time.Time `json:"multifactorInitializationSkippedAt,omitzero" db:"multifactor_initialization_skipped_at"`

	Email    HumanEmail    `json:"email,omitzero" db:"email"`
	Phone    *HumanPhone   `json:"phone,omitempty" db:"phone"`
	Passkeys []*Passkey    `json:"passkeys,omitempty" db:"passkeys"`
	Password HumanPassword `json:"password,omitzero" db:"password"`
	TOTP     *HumanTOTP    `json:"totp,omitempty" db:"totp"`

	IdentityProviderLinks []*IdentityProviderLink `json:"identityProviderLinks,omitempty" db:"identity_provider_links"`

	// Verifications are used if the Login does something on behalf of the user that requires verification
	// e.g. starting passkey registration
	Verifications []*Verification `db:"verifications" json:"verifications,omitempty"`
}

type UserMetadata struct {
	Metadata
	UserID string `json:"userId,omitempty" db:"user_id"`
}

//go:generate enumer -type UserState -transform lower -trimprefix UserState -sql

type UserState uint8

const (
	UserStateUnspecified UserState = iota
	UserStateActive
	UserStateInactive
	UserStateLocked
	UserStateInitial
)

//go:generate enumer -type HumanGender -transform lower -trimprefix HumanGender

type HumanGender uint8

const (
	HumanGenderUnspecified HumanGender = iota
	HumanGenderFemale
	HumanGenderMale
	HumanGenderDiverse
)

type HumanPassword struct {
	// Password is the hashed password
	Password *crypto.CryptoValue `json:"password" db:"password"` //TODO: make sure password is not marshalled
	// IsChangeRequired indicates if the user must change their password
	IsChangeRequired bool `json:"isChangeRequired,omitempty" db:"is_change_required"`
	// ChangedAt is the time when the current password was last updated
	ChangedAt time.Time `json:"changedAt,omitzero" db:"changed_at"`
	// Unverified is the verification data for setting a new password
	// If nil, no password change is in progress
	Unverified *Verification `json:"pendingVerification,omitempty" db:"-"`
	// LastSuccessfullyCheckedAt is the time when the password was last successfully checked
	LastSuccessfullyCheckedAt *time.Time `json:"lastSuccessfullyCheckedAt,omitzero" db:"last_successfully_checked_at"`
	// FailedAttempts is the number of consecutive failed password attempts
	// It is reset to 0 on successful verification
	FailedAttempts uint8 `json:"failedAttempts,omitempty" db:"failed_attempts"`
}

type HumanEmail struct {
	Address    string    `json:"address" db:"address"`
	VerifiedAt time.Time `json:"verifiedAt,omitzero" db:"verified_at"`
	OTP        OTP       `json:"otp" db:"otp"`
	// Unverified is the verification data for setting a new email
	// If nil, no email change is in progress
	Unverified *Verification `json:"pendingVerification,omitempty" db:"-"`
}

type HumanPhone struct {
	Number     string    `json:"number" db:"number"`
	VerifiedAt time.Time `json:"verifiedAt,omitzero" db:"verified_at"`
	OTP        OTP       `json:"otp,omitzero" db:"otp"`
	// Unverified is the verification data for setting a new phone number
	// If nil, no phone change is in progress
	Unverified *Verification `json:"pendingVerification,omitempty" db:"-"`
}

type HumanTOTP struct {
	VerifiedAt time.Time `json:"verifiedAt,omitzero" db:"verified_at"`
	Secret     []byte    `json:"secret,omitempty" db:"secret"`
	// LastSuccessfullyCheckedAt is the time when the TOTP was last successfully checked
	LastSuccessfullyCheckedAt *time.Time `json:"lastSuccessfullyCheckedAt,omitzero" db:"last_successfully_checked_at"`
	// FailedAttempts is the number of consecutive failed password attempts
	// It is reset to 0 on successful verification
	FailedAttempts uint8 `json:"failedAttempts,omitempty" db:"failed_attempts"`
}

type OTP struct {
	EnabledAt time.Time `json:"enabledAt,omitzero" db:"enabled_at"`
	// LastSuccessfullyCheckedAt is the time when the OTP was last successfully checked
	LastSuccessfullyCheckedAt *time.Time `json:"lastSuccessfullyCheckedAt,omitzero" db:"last_successfully_checked_at"`
	// FailedAttempts is the number of consecutive failed password attempts
	// It is reset to 0 on successful verification
	FailedAttempts uint8 `json:"failedAttempts,omitempty" db:"failed_attempts"`
}

type PersonalAccessToken struct {
	ID        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"createdAt,omitzero" db:"created_at"`
	ExpiresAt time.Time `json:"expiresAt,omitzero" db:"expires_at"`
	Scopes    []string  `json:"scopes" db:"scopes"`
}

//go:generate enumer -type AccessTokenType -transform lower -trimprefix AccessTokenType
type AccessTokenType uint8

const (
	AccessTokenTypeUnspecified AccessTokenType = iota
	AccessTokenTypeBearer
	AccessTokenTypeJWT
)

type MachineKey struct {
	ID        string         `json:"id" db:"id"`
	PublicKey []byte         `json:"publicKey" db:"public_key"`
	CreatedAt time.Time      `json:"createdAt,omitzero" db:"created_at"`
	ExpiresAt time.Time      `json:"expiresAt,omitzero" db:"expires_at"`
	Type      MachineKeyType `json:"type" db:"type"`
}

//go:generate enumer -type MachineKeyType -transform lower -trimprefix MachineKeyType

type MachineKeyType uint8

const (
	MachineKeyTypeNone MachineKeyType = iota
	MachineKeyTypeJSON
)

type Passkey struct {
	ID                           string      `json:"id" db:"id"`
	KeyID                        []byte      `json:"keyId" db:"key_id"`
	Name                         string      `json:"name" db:"name"`
	SignCount                    uint32      `json:"signCount" db:"sign_count"`
	PublicKey                    []byte      `json:"publicKey" db:"public_key"`
	AttestationType              string      `json:"attestationType" db:"attestation_type"`
	AuthenticatorAttestationGUID []byte      `json:"aaGuid" db:"authenticator_attestation_guid"`
	Type                         PasskeyType `json:"type" db:"type"`

	CreatedAt  time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt,omitzero" db:"updated_at"`
	VerifiedAt time.Time `json:"verifiedAt,omitzero" db:"verified_at"`

	// Challenge is used during device registration
	Challenge []byte `json:"challenge" db:"challenge"`
	// RelyingPartyID is used during device registration
	RelyingPartyID string `json:"rpId" db:"relying_party_id"`

	// Initialization is used during user registration
	Initialization *Verification `json:"-" db:"initialization"`
}

//go:generate enumer -type PasskeyType -transform lower -trimprefix PasskeyType -json -sql

type PasskeyType uint8

const (
	PasskeyTypeUnspecified PasskeyType = iota
	PasskeyTypePasswordless
	PasskeyTypeU2F
)

type IdentityProviderLink struct {
	// TODO(adlerhurst): double check with marcos pr
	ProviderID       string `json:"providerId" db:"provider_id"`
	ProvidedUserID   string `json:"providedUserId" db:"provided_user_id"`
	ProvidedUsername string `json:"providedUsername" db:"provided_username"`

	CreatedAt time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt,omitzero" db:"updated_at"`
}

//go:generate enumer -type UserType -transform lower -trimprefix UserType

type UserType uint8

const (
	UserTypeUnspecified UserType = iota
	UserTypeHuman
	UserTypeMachine
)
