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
	Name            string          `json:"name,omitempty" db:"-"`
	Description     string          `json:"description,omitempty" db:"-"`
	Secret          string          `json:"secret,omitempty" db:"-"`
	AccessTokenType AccessTokenType `json:"accessTokenType,omitempty" db:"-"`

	PATs []*PersonalAccessToken `json:"pats,omitempty" db:"-"`
	Keys []*MachineKey          `json:"keys,omitempty" db:"-"`
}

type HumanUser struct {
	FirstName         string       `json:"firstName,omitempty" db:"-"`
	LastName          string       `json:"lastName,omitempty" db:"-"`
	Nickname          string       `json:"nickname,omitempty" db:"-"`
	DisplayName       string       `json:"displayName,omitempty" db:"-"`
	PreferredLanguage language.Tag `json:"preferredLanguage,omitzero" db:"-"`
	Gender            HumanGender  `json:"gender,omitempty" db:"-"`
	AvatarKey         string       `json:"avatarKey,omitempty" db:"-"`

	MultifactorInitializationSkippedAt time.Time `json:"multifactorInitializationSkippedAt,omitzero" db:"-"`

	Email    HumanEmail    `json:"email,omitzero" db:"-"`
	Phone    *HumanPhone   `json:"phone,omitempty" db:"-"`
	Passkeys []*Passkey    `json:"passkeys,omitempty" db:"-"`
	Password HumanPassword `json:"password,omitzero" db:"-"`
	TOTP     *HumanTOTP    `json:"totp,omitempty" db:"-"`

	IdentityProviderLinks []*IdentityProviderLink `json:"identityProviderLinks,omitempty" db:"-"`

	Invite *HumanInvite `json:"invite,omitempty" db:"-"`

	// Verifications are used if the Login does something on behalf of the user that requires verification
	// e.g. starting passkey registration
	Verifications []*Verification `json:"verifications,omitempty" db:"-"`
}

type UserMetadata struct {
	Metadata
	UserID string `json:"userId,omitempty" db:"-"`
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
	// Hash is the hashed password
	Hash string `json:"hash" db:"-"`
	// IsChangeRequired indicates if the user must change their password
	IsChangeRequired bool `json:"isChangeRequired,omitempty" db:"-"`
	// ChangedAt is the time when the current password was last updated
	ChangedAt time.Time `json:"changedAt,omitzero" db:"-"`
	// PendingVerification is the verification data for setting a new password
	// If nil, no password change is in progress
	PendingVerification *Verification `json:"pendingVerification,omitempty" db:"-"`
	// LastSuccessfullyCheckedAt is the time when the password was last successfully checked
	LastSuccessfullyCheckedAt *time.Time `json:"lastSuccessfullyCheckedAt,omitzero" db:"-"`
	// FailedAttempts is the number of consecutive failed password attempts
	// It is reset to 0 on successful verification
	FailedAttempts uint8 `json:"failedAttempts,omitempty" db:"-"`
}

type HumanEmail struct {
	Address           string    `json:"address" db:"-"`
	UnverifiedAddress string    `json:"unverifiedAddress,omitempty" db:"-"`
	VerifiedAt        time.Time `json:"verifiedAt,omitzero" db:"-"`
	OTP               OTP       `json:"otp" db:"-"`
	// PendingVerification is the verification data for UnverifiedAddress
	// If nil, no email change is in progress
	PendingVerification *Verification `json:"pendingVerification,omitempty" db:"-"`
}

type HumanPhone struct {
	Number           string    `json:"number" db:"-"`
	UnverifiedNumber string    `json:"unverifiedNumber,omitempty" db:"-"`
	VerifiedAt       time.Time `json:"verifiedAt,omitzero" db:"-"`
	OTP              OTP       `json:"otp,omitzero" db:"-"`
	// PendingVerification is the verification data for setting a new phone number
	// If nil, no phone change is in progress
	PendingVerification *Verification `json:"pendingVerification,omitempty" db:"-"`
}

type HumanTOTP struct {
	VerifiedAt time.Time           `json:"verifiedAt,omitzero" db:"-"`
	Secret     *crypto.CryptoValue `json:"secret,omitempty" db:"-"`
	// LastSuccessfullyCheckedAt is the time when the TOTP was last successfully checked
	LastSuccessfullyCheckedAt *time.Time `json:"lastSuccessfullyCheckedAt,omitzero" db:"-"`

	// FailedAttempts is the number of consecutive failed password attempts
	// It is reset to 0 on successful verification
	FailedAttempts uint8 `json:"failedAttempts,omitempty" db:"-"`
}

type OTP struct {
	EnabledAt time.Time `json:"enabledAt,omitzero" db:"-"`
	// LastSuccessfullyCheckedAt is the time when the OTP was last successfully checked
	LastSuccessfullyCheckedAt *time.Time `json:"lastSuccessfullyCheckedAt,omitzero" db:"-"`
	// FailedAttempts is the number of consecutive failed password attempts
	// It is reset to 0 on successful verification
	FailedAttempts uint8 `json:"failedAttempts,omitempty" db:"-"`
}

type PersonalAccessToken struct {
	ID        string    `json:"id" db:"-"`
	CreatedAt time.Time `json:"createdAt,omitzero" db:"-"`
	ExpiresAt time.Time `json:"expiresAt,omitzero" db:"-"`
	Scopes    []string  `json:"scopes" db:"-"`
}

//go:generate enumer -type AccessTokenType -transform lower -trimprefix AccessTokenType
type AccessTokenType uint8

const (
	AccessTokenTypeUnspecified AccessTokenType = iota
	AccessTokenTypeBearer
	AccessTokenTypeJWT
)

type MachineKey struct {
	ID        string         `json:"id" db:"-"`
	PublicKey []byte         `json:"publicKey" db:"-"`
	CreatedAt time.Time      `json:"createdAt,omitzero" db:"-"`
	ExpiresAt time.Time      `json:"expiresAt,omitzero" db:"-"`
	Type      MachineKeyType `json:"type" db:"-"`
}

//go:generate enumer -type MachineKeyType -transform lower -trimprefix MachineKeyType

type MachineKeyType uint8

const (
	MachineKeyTypeNone MachineKeyType = iota
	MachineKeyTypeJSON
)

type Passkey struct {
	ID                           string      `json:"id" db:"-"`
	KeyID                        []byte      `json:"keyId" db:"-"`
	Name                         string      `json:"name" db:"-"`
	SignCount                    uint32      `json:"signCount" db:"-"`
	PublicKey                    []byte      `json:"publicKey" db:"-"`
	AttestationType              string      `json:"attestationType" db:"-"`
	AuthenticatorAttestationGUID []byte      `json:"aaGuid" db:"-"`
	Type                         PasskeyType `json:"type" db:"-"`

	CreatedAt  time.Time `json:"createdAt,omitzero" db:"-"`
	UpdatedAt  time.Time `json:"updatedAt,omitzero" db:"-"`
	VerifiedAt time.Time `json:"verifiedAt,omitzero" db:"-"`

	// Challenge is used during device registration
	Challenge []byte `json:"challenge" db:"-"`
	// RelyingPartyID is used during device registration
	RelyingPartyID string `json:"rpId" db:"-"`

	// Initialization is used during user registration
	Initialization *Verification `json:"-" db:"-"`
}

//go:generate enumer -type PasskeyType -transform lower -trimprefix PasskeyType -json -sql

type PasskeyType uint8

const (
	PasskeyTypeUnspecified PasskeyType = iota
	PasskeyTypePasswordless
	PasskeyTypeU2F
)

type IdentityProviderLink struct {
	ProviderID       string `json:"providerId" db:"-"`
	ProvidedUserID   string `json:"providedUserId" db:"-"`
	ProvidedUsername string `json:"providedUsername" db:"-"`

	CreatedAt time.Time `json:"createdAt,omitzero" db:"-"`
	UpdatedAt time.Time `json:"updatedAt,omitzero" db:"-"`
}

//go:generate enumer -type UserType -transform lower -trimprefix UserType

type UserType uint8

const (
	UserTypeUnspecified UserType = iota
	UserTypeHuman
	UserTypeMachine
)

type HumanInvite struct {
	AcceptedAt time.Time `json:"acceptedAt,omitzero" db:"-"`
	// PendingVerification is the verification data for UnverifiedAddress
	// If nil, no email change is in progress
	PendingVerification *Verification `json:"pendingVerification,omitempty" db:"-"`
	// FailedAttempts is the number of consecutive failed invite attempts
	FailedAttempts uint8 `json:"failedAttempts,omitempty" db:"-"`
}
