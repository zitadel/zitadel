package domain

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
)

type Verification struct {
	ID             string              `json:"id" db:"id"`
	Value          *string             `json:"value" db:"value"`
	Code           *crypto.CryptoValue `json:"code" db:"code"`
	CreatedAt      time.Time           `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time           `json:"updatedAt" db:"updated_at"`
	ExpiresAt      *time.Time          `json:"expiresAt" db:"expires_at"`
	FailedAttempts uint8               `json:"failedAttempts" db:"failed_attempts"`
}

type VerificationType interface {
	isVerificationType()
}

// VerificationTypeInit indicates that an object which needs verification is created.
type VerificationTypeInit struct {
	// ID is the ID of the verification.
	// The id must be set if the the object can have multiple verifications.
	ID *string
	// CreatedAt is the time when the verification was created.
	// If zero, the current time will be used.
	CreatedAt time.Time
	// Expiry is the duration after which the verification expires.
	Expiry *time.Duration
	// Code is the code to be used for verification.
	Code *crypto.CryptoValue
	// Value is the value to be set after successful verification.
	Value *string
}

// isVerificationType implements [VerificationType].
func (v *VerificationTypeInit) isVerificationType() {}

var _ VerificationType = (*VerificationTypeInit)(nil)

// VerificationTypeVerified indicates that the verification was successful.
// If VerifiedAt is zero, the current time will be used.
// If the Value of the verification is present, it is used as the new value to be set.
type VerificationTypeVerified struct {
	// ID is the ID of the verification.
	// The id must be set if the the object can have multiple verifications.
	ID         *string
	VerifiedAt time.Time
}

// isVerificationType implements [VerificationType].
func (v *VerificationTypeVerified) isVerificationType() {}

var _ VerificationType = (*VerificationTypeVerified)(nil)

// VerificationTypeUpdate updates an existing verification.
// Non-nil fields get updated.
type VerificationTypeUpdate struct {
	// ID is the ID of the verification.
	// The id must be set if the the object can have multiple verifications.
	ID     *string
	Code   *crypto.CryptoValue
	Value  *string
	Expiry *time.Duration
}

// isVerificationType implements [VerificationType].
func (v *VerificationTypeUpdate) isVerificationType() {}

var _ VerificationType = (*VerificationTypeUpdate)(nil)

// VerificationTypeSkipped indicates that the verification was skipped.
// If Value is present, it is used as the new value to be set.
type VerificationTypeSkipped struct {
	// ID is the ID of the verification.
	// The id must be set if the the object can have multiple verifications.
	ID        *string
	SkippedAt time.Time
	Value     *string
}

// isVerificationType implements [VerificationType].
func (v *VerificationTypeSkipped) isVerificationType() {}

var _ VerificationType = (*VerificationTypeSkipped)(nil)

type VerificationTypeFailed struct {
	// ID is the ID of the verification.
	// The id must be set if the the object can have multiple verifications.
	ID       *string
	FailedAt time.Time
}

// isVerificationType implements [VerificationType].
func (v *VerificationTypeFailed) isVerificationType() {}

var _ VerificationType = (*VerificationTypeFailed)(nil)
