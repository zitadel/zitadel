package domain

import "time"

type Check struct {
	Code           *[]byte    `db:"code"`
	ExpiresAt      *time.Time `db:"expires_at"`
	FailedAttempts uint8      `db:"failed_attempts"`
}

type CheckType interface {
	isCheckType()
}

// CheckTypeInit indicates that an object which needs verification is created.
type CheckTypeInit struct {
	// CreatedAt is the time when the verification was created.
	// If zero, the current time will be used.
	CreatedAt time.Time
	// Expiry is the duration after which the verification expires.
	Expiry *time.Duration
	// Code is the code to be used for verification.
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

// CheckTypeSucceeded indicates that the verification was Succeeded.
// It is used to reset the failed attempts counter.
type CheckTypeSucceeded struct {
	SucceededAt time.Time
}

// isCheckType implements [CheckType].
func (v *CheckTypeSucceeded) isCheckType() {}

var _ CheckType = (*CheckTypeSucceeded)(nil)
