package domain

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
)

type Check struct {
	Code           *crypto.CryptoValue `json:"code" db:"code"`
	ExpiresAt      *time.Time          `json:"expiresAt" db:"expires_at"`
	FailedAttempts uint8               `json:"failedAttempts" db:"failed_attempts"`
}

type CheckType interface {
	isCheckType()
}

type PasswordCheckType interface {
	isPasswordCheckType()
}

// CheckTypeInit initializes a new check.
type CheckTypeInit struct {
	// CreatedAt is the time when the check was created.
	// If zero, the current time will be used.
	CreatedAt time.Time
	// Expiry is the duration after which the check expires.
	Expiry *time.Duration
	// Code to be used for the check.
	Code *crypto.CryptoValue
}

// isCheckType implements [CheckType].
func (v *CheckTypeInit) isCheckType() {}

var _ CheckType = (*CheckTypeInit)(nil)

// CheckTypeFailed indicates that the check failed.
// It is used to increment the failed attempts counter.
type CheckTypeFailed struct {
	FailedAt time.Time
}

// isCheckType implements [CheckType].
func (v *CheckTypeFailed) isCheckType() {}

// isPasswordCheckType implements [PasswordCheckType].
func (v *CheckTypeFailed) isPasswordCheckType() {}

var (
	_ CheckType         = (*CheckTypeFailed)(nil)
	_ PasswordCheckType = (*CheckTypeFailed)(nil)
)

// CheckTypeSucceeded indicates that the check succeeded.
// It is used to reset the failed attempts counter.
type CheckTypeSucceeded struct {
	SucceededAt time.Time
}

// isCheckType implements [CheckType].
func (v *CheckTypeSucceeded) isCheckType() {}

// isPasswordCheckType implements [PasswordCheckType].
func (v *CheckTypeSucceeded) isPasswordCheckType() {}

var (
	_ CheckType         = (*CheckTypeSucceeded)(nil)
	_ PasswordCheckType = (*CheckTypeSucceeded)(nil)
)
