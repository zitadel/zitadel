package domain

import "time"

type Check struct {
	Code           []byte     `db:"code"`
	ExpiresAt      *time.Time `db:"expires_at"`
	FailedAttempts uint8      `db:"failed_attempts"`
}

type CheckType interface {
	isCheckType()
}

// CheckTypeInit initializes a new check.
type CheckTypeInit struct {
	// CreatedAt is the time when the check was created.
	// If zero, the current time will be used.
	CreatedAt time.Time
	// Expiry is the duration after which the check expires.
	Expiry *time.Duration
	// Code to be used for the check.
	Code []byte
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

var _ CheckType = (*CheckTypeFailed)(nil)

// CheckTypeSucceeded indicates that the check succeeded.
// It is used to reset the failed attempts counter.
type CheckTypeSucceeded struct {
	SucceededAt time.Time
}

// isCheckType implements [CheckType].
func (v *CheckTypeSucceeded) isCheckType() {}

var _ CheckType = (*CheckTypeSucceeded)(nil)
