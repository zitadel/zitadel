package domain

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
)

type Verification struct {
	ID             string              `json:"id" db:"id"`
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
}

// isVerificationType implements [VerificationType].
func (v *VerificationTypeInit) isVerificationType() {}

var _ VerificationType = (*VerificationTypeInit)(nil)

// VerificationTypeSucceeded indicates that the verification was successful.
// If VerifiedAt is zero, the current time will be used.
type VerificationTypeSucceeded struct {
	// ID is the ID of the verification.
	// The id must be set if the the object can have multiple verifications.
	ID *string
	// If zero, the current time will be used.
	VerifiedAt time.Time
}

// isVerificationType implements [VerificationType].
func (v *VerificationTypeSucceeded) isVerificationType() {}

var _ VerificationType = (*VerificationTypeSucceeded)(nil)

// VerificationTypeUpdate updates an existing verification.
// Non-nil fields get updated.
type VerificationTypeUpdate struct {
	// ID is the ID of the verification.
	// The id must be set if the the object can have multiple verifications.
	ID     *string
	Code   *crypto.CryptoValue
	Expiry *time.Duration
	// If zero, the current time will be used.
	UpdatedAt time.Time
}

// isVerificationType implements [VerificationType].
func (v *VerificationTypeUpdate) isVerificationType() {}

var _ VerificationType = (*VerificationTypeUpdate)(nil)

// VerificationTypeSkipped indicates that the verification was skipped.
type VerificationTypeSkipped struct {
	// ID is the ID of the verification.
	// The id must be set if the the object can have multiple verifications.
	ID *string
	// If zero, the current time will be used.
	SkippedAt time.Time
}

// isVerificationType implements [VerificationType].
func (v *VerificationTypeSkipped) isVerificationType() {}

var _ VerificationType = (*VerificationTypeSkipped)(nil)

type VerificationTypeFailed struct {
	// ID is the ID of the verification.
	// The id must be set if the the object can have multiple verifications.
	ID *string
	// If zero, the current time will be used.
	FailedAt time.Time
}

// isVerificationType implements [VerificationType].
func (v *VerificationTypeFailed) isVerificationType() {}

var _ VerificationType = (*VerificationTypeFailed)(nil)
