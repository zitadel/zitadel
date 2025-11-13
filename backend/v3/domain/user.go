package domain

import (
	"time"

	"golang.org/x/text/language"
)

type User struct {
	InstanceID     string    `db:"instance_id"`
	OrganizationID string    `db:"organization_id"`
	ID             string    `db:"id"`
	Username       string    `db:"username"`
	State          UserState `db:"state"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`

	Machine  *MachineUser    `db:"machine"`
	Human    *HumanUser      `db:"human"`
	Metadata []*UserMetadata `db:"metadata"`
}

type MachineUser struct {
	Name            string                  `db:"name"`
	Description     string                  `db:"description"`
	Secret          string                  `db:"secret"`
	AccessTokenType PersonalAccessTokenType `db:"access_token_type"`

	PATs []*PersonalAccessToken `db:"pats"`
	Keys []*MachineKey          `db:"keys"`
}

type HumanUser struct {
	FirstName         string       `db:"first_name"`
	LastName          string       `db:"last_name"`
	Nickname          string       `db:"nickname"`
	DisplayName       string       `db:"display_name"`
	PreferredLanguage language.Tag `db:"preferred_language"`
	Gender            HumanGender  `db:"gender"`
	AvatarKey         string       `db:"avatar_key"`

	MultifactorInitializationSkippedAt time.Time `db:"multifactor_initialization_skipped_at"`

	Email    HumanEmail    `db:"email"`
	Phone    *HumanPhone   `db:"phone"`
	Passkeys []*Passkey    `db:"passkeys"`
	Password HumanPassword `db:"password"`

	IdentityProviderLinks []*IdentityProviderLink `db:"identity_provider_links"`

	// Verifications are used if the Login does something on behalf of the user that requires verification
	// e.g. starting passkey registration
	Verifications []*Verification `db:"verifications"`
}

type UserMetadata struct {
	Metadata
	UserID string `db:"user_id"`
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
	Password string `db:"password"`
	// IsChangeRequired indicates if the user must change their password
	IsChangeRequired bool `db:"is_change_required"`
	// VerifiedAt is the time when the current password was verified
	VerifiedAt time.Time `db:"verified_at"`
	// Unverified is the verification data for setting a new password
	// If nil, no password change is in progress
	Unverified *Verification
}

type HumanEmail struct {
	Address    string    `db:"email"`
	VerifiedAt time.Time `db:"verified_at"`
	OTP        OTP
	// Unverified is the verification data for setting a new email
	// If nil, no email change is in progress
	Unverified *Verification
}

type HumanPhone struct {
	Number     string    `db:"phone"`
	VerifiedAt time.Time `db:"verified_at"`
	OTP        OTP
	// Unverified is the verification data for setting a new phone number
	// If nil, no phone change is in progress
	Unverified *Verification
}

type HumanTOTP struct {
	VerifiedAt time.Time `db:"verified_at"`

	LastSuccessfullyCheckedAt time.Time `db:"last_successfully_checked_at"`
	// Check is the verified secret
	Check *Check
	// Unverified is the unverified secret
	Unverified *Verification
}

type OTP struct {
	EnabledAt                 time.Time `db:"enabled_at"`
	LastSuccessfullyCheckedAt time.Time `db:"last_successfully_checked_at"`
	// Check is the currently active OTP check
	// If nil, no OTP check is active
	Check *Check
}

type PersonalAccessToken struct {
	ID        string                  `db:"id"`
	CreatedAt time.Time               `db:"created_at"`
	ExpiresAt time.Time               `db:"expires_at"`
	Type      PersonalAccessTokenType `db:"type"`
	PublicKey string                  `db:"public_key"`
	Scopes    []string                `db:"scopes"`
}

//go:generate enumer -type PersonalAccessTokenType -transform lower -trimprefix PersonalAccessTokenType

type PersonalAccessTokenType uint8

const (
	PersonalAccessTokenTypeUnspecified PersonalAccessTokenType = iota
	PersonalAccessTokenTypeBearer
	PersonalAccessTokenTypeJWT
)

type MachineKey struct {
	ID        string         `db:"id"`
	PublicKey []byte         `db:"public_key"`
	CreatedAt time.Time      `db:"created_at"`
	ExpiresAt time.Time      `db:"expires_at"`
	Type      MachineKeyType `db:"type"`
}

//go:generate enumer -type MachineKeyType -transform lower -trimprefix MachineKeyType

type MachineKeyType uint8

const (
	MachineKeyTypeUnspecified MachineKeyType = iota
	MachineKeyTypeJSON
	MachineKeyTypeNone
)

type Passkey struct {
	ID                           string      `db:"id"`
	KeyID                        []byte      `db:"key_id"`
	Name                         string      `db:"name"`
	SignCount                    uint32      `db:"sign_count"`
	PublicKey                    []byte      `db:"public_key"`
	AttestationType              string      `db:"attestation_type"`
	AuthenticatorAttestationGUID []byte      `db:"authenticator_attestation_guid"`
	Type                         PasskeyType `db:"type"`

	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	VerifiedAt time.Time `db:"verified_at"`

	// Challenge is used during device registration
	Challenge []byte `db:"challenge"`
	// RelayingPartyID is used during device registration
	RelayingPartyID string `db:"relaying_party_id"`

	// Initialization is used during user registration
	Initialization *Verification `db:"initialization"`
}

//go:generate enumer -type PasskeyType -transform lower -trimprefix PasskeyType

type PasskeyType uint8

const (
	PasskeyTypeUnspecified PasskeyType = iota
	PasskeyTypePasswordless
	PasskeyTypeU2F
)

type IdentityProviderLink struct {
	ProviderID       string `db:"provider_id"`
	ProvidedUserID   string `db:"provided_user_id"`
	ProvidedUsername string `db:"provided_username"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

//go:generate enumer -type UserType -transform lower -trimprefix UserType

type UserType uint8

const (
	UserTypeUnspecified UserType = iota
	UserTypeHuman
	UserTypeMachine
)
