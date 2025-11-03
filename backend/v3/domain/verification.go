package domain

import "time"

type Verification struct {
	Value          *[]byte    `db:"value"`
	Code           *[]byte    `db:"code"`
	ExpiresAt      *time.Time `db:"expires_at"`
	FailedAttempts uint8      `db:"failed_attempts"`
}

type VerificationType interface {
	isVerificationType()
}

// VerificationTypeSuccessful indicates that the verification was successful.
// The Value of the verification is used as the new value to be set.
// If VerifiedAt is zero, the current time will be used.
type VerificationTypeSuccessful struct {
	VerifiedAt time.Time
}

func (v *VerificationTypeSuccessful) isVerificationType() {}

// VerificationTypeSkipVerification indicates that the verification is skipped.
// The Value is used as the new value to be set.
// If VerifiedAt is zero, the current time will be used.
type VerificationTypeSkipVerification struct {
	VerifiedAt *time.Time
	Value      string
}

func (v *VerificationTypeSkipVerification) isVerificationType() {}

var _ VerificationType = (*VerificationTypeSkipVerification)(nil)

var _ VerificationType = (*VerificationTypeSuccessful)(nil)

// VerificationTypeInitCode indicates that a verification is started
type VerificationTypeInitCode struct {
	Code  []byte
	Value *string
	// CreatedAt is the time when the verification was created.
	// If zero, the current time will be used.
	CreatedAt time.Time
	Expiry    *time.Duration
}

func (v *VerificationTypeInitCode) isVerificationType() {}

var _ VerificationType = (*VerificationTypeInitCode)(nil)

// VerificationTypeUpdate indicates that a verification is updated.
// The fields that are non-nil will be updated.
type VerificationTypeUpdate struct {
	Code   []byte
	Value  *string
	Expiry *time.Duration
}

func (v *VerificationTypeUpdate) isVerificationType() {}

var _ VerificationType = (*VerificationTypeUpdate)(nil)
