package domain

import (
	"time"

	"golang.org/x/text/language"
)

type User struct {
	InstanceID     string    `json:"instanceId,omitempty" db:"instance_id"`
	OrganizationID string    `json:"organizationId,omitempty" db:"organization_id"`
	ID             string    `json:"id,omitempty" db:"id"`
	Username       string    `json:"username,omitempty" db:"username"`
	State          UserState `json:"state,omitempty" db:"state"`
	CreatedAt      time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt,omitzero" db:"updated_at"`

	Machine  *MachineUser    `json:"machine,omitempty" db:"machine"`
	Human    *HumanUser      `json:"human,omitempty" db:"human"`
	Metadata []*UserMetadata `json:"metadata,omitempty" db:"metadata"`
}

type MachineUser struct {
	Name            string                  `json:"name,omitempty" db:"name"`
	Description     string                  `json:"description,omitempty" db:"description"`
	Secret          string                  `json:"secret,omitempty" db:"secret"`
	AccessTokenType PersonalAccessTokenType `json:"access_token_type,omitempty" db:"access_token_type"`

	PATs []*PersonalAccessToken `json:"pats,omitempty" db:"pats"`
	Keys []*MachineKey          `json:"keys,omitempty" db:"keys"`
}

type HumanUser struct {
	FirstName         string       `json:"firstName,omitempty" db:"first_name"`
	LastName          string       `json:"lastName,omitempty" db:"last_name"`
	Nickname          string       `json:"nickname,omitempty" db:"nickname"`
	DisplayName       string       `json:"displayName,omitempty" db:"display_name"`
	PreferredLanguage language.Tag `json:"preferredLanguage,omitempty" db:"preferred_language"`
	Gender            HumanGender  `json:"gender,omitempty" db:"gender"`
	AvatarKey         string       `json:"avatarKey,omitempty" db:"avatar_key"`

	MultifactorInitializationSkippedAt time.Time `json:"multifactorInitializationSkippedAt,omitempty" db:"multifactor_initialization_skipped_at"`

	Email    HumanEmail    `json:"email,omitempty" db:"email"`
	Phone    *HumanPhone   `json:"phone,omitempty" db:"phone"`
	Passkeys []*Passkey    `json:"passkeys,omitempty" db:"passkeys"`
	Password HumanPassword `json:"password,omitempty" db:"password"`
	TOTP     HumanTOTP     `json:"totp,omitempty" db:"totp"`

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
	Password string `json:"-" db:"password"`
	// IsChangeRequired indicates if the user must change their password
	IsChangeRequired bool `json:"isChangeRequired,omitempty" db:"is_change_required"`
	// VerifiedAt is the time when the current password was verified
	VerifiedAt time.Time `json:"verifiedAt,omitzero" db:"verified_at"`
	// Unverified is the verification data for setting a new password
	// If nil, no password change is in progress
	Unverified *Verification `json:"-"`
}

type HumanEmail struct {
	Address    string    `json:"address" db:"email"`
	VerifiedAt time.Time `json:"verifiedAt,omitzero" db:"verified_at"`
	OTP        OTP       `json:"-"`
	// Unverified is the verification data for setting a new email
	// If nil, no email change is in progress
	Unverified *Verification `json:"-"`
}

type HumanPhone struct {
	Number     string    `json:"number" db:"phone"`
	VerifiedAt time.Time `json:"verifiedAt,omitzero" db:"verified_at"`
	OTP        OTP       `json:"-"`
	// Unverified is the verification data for setting a new phone number
	// If nil, no phone change is in progress
	Unverified *Verification `json:"-"`
}

type HumanTOTP struct {
	VerifiedAt time.Time `json:"verifiedAt,omitzero" db:"verified_at"`

	LastSuccessfullyCheckedAt time.Time `json:"lastSuccessfullyCheckedAt,omitzero" db:"last_successfully_checked_at"`
	// Check is the verified secret
	Check *Check `json:"-"`
	// Unverified is the unverified secret
	Unverified *Verification `json:"-"`
}

type OTP struct {
	EnabledAt                 time.Time `json:"enabledAt,omitzero" db:"enabled_at"`
	LastSuccessfullyCheckedAt time.Time `json:"lastSuccessfullyCheckedAt,omitzero" db:"last_successfully_checked_at"`
	// Check is the currently active OTP check
	// If nil, no OTP check is active
	Check *Check `json:"-"`
}

type PersonalAccessToken struct {
	ID        string                  `json:"id" db:"id"`
	CreatedAt time.Time               `json:"createdAt,omitzero" db:"created_at"`
	ExpiresAt time.Time               `json:"expiresAt,omitzero" db:"expires_at"`
	Type      PersonalAccessTokenType `json:"type" db:"type"`
	PublicKey string                  `json:"-" db:"public_key"`
	Scopes    []string                `json:"scopes" db:"scopes"`
}

//go:generate enumer -type PersonalAccessTokenType -transform lower -trimprefix PersonalAccessTokenType

type PersonalAccessTokenType uint8

const (
	PersonalAccessTokenTypeUnspecified PersonalAccessTokenType = iota
	PersonalAccessTokenTypeBearer
	PersonalAccessTokenTypeJWT
)

type MachineKey struct {
	ID        string         `json:"id" db:"id"`
	PublicKey []byte         `json:"-" db:"public_key"`
	CreatedAt time.Time      `json:"createdAt,omitzero" db:"created_at"`
	ExpiresAt time.Time      `json:"expiresAt,omitzero" db:"expires_at"`
	Type      MachineKeyType `json:"type" db:"type"`
}

//go:generate enumer -type MachineKeyType -transform lower -trimprefix MachineKeyType

type MachineKeyType uint8

const (
	MachineKeyTypeUnspecified MachineKeyType = iota
	MachineKeyTypeJSON
	MachineKeyTypeNone
)

type Passkey struct {
	ID                           string      `json:"id" db:"token_id"`
	KeyID                        []byte      `json:"-" db:"key_id"`
	Name                         string      `json:"name" db:"name"`
	SignCount                    uint32      `json:"-" db:"sign_count"`
	PublicKey                    []byte      `json:"-" db:"public_key"`
	AttestationType              string      `json:"attestationType" db:"attestation_type"`
	AuthenticatorAttestationGUID []byte      `json:"aaGuid" db:"authenticator_attestation_guid"`
	Type                         PasskeyType `json:"type" db:"type"`

	CreatedAt  time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt,omitzero" db:"updated_at"`
	VerifiedAt time.Time `json:"verifiedAt,omitzero" db:"verified_at"`

	// Challenge is used during device registration
	Challenge []byte `json:"-" db:"challenge"`
	// RelyingPartyID is used during device registration
	RelyingPartyID string `json:"-" db:"relying_party_id"`

	// Initialization is used during user registration
	Initialization *Verification `json:"-" db:"initialization"`
}

//go:generate enumer -type PasskeyType -transform lower -trimprefix PasskeyType

type PasskeyType uint8

const (
	PasskeyTypeUnspecified PasskeyType = iota
	PasskeyTypePasswordless
	PasskeyTypeU2F
)

type IdentityProviderLink struct {
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
