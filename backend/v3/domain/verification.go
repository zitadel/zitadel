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

type OTPVerificationType interface {
	isOTPVerificationType()
}

// VerificationTypeSuccessful indicates that the verification was successful.
// If the Value of the verification is present, it is used as the new value to be set.
// If VerifiedAt is zero, the current time will be used.
type VerificationTypeSuccessful struct {
	VerifiedAt time.Time
}

// isOTPVerificationType implements [OTPVerificationType].
func (v *VerificationTypeSuccessful) isOTPVerificationType() {}

// isVerificationType implements [VerificationType].
func (v *VerificationTypeSuccessful) isVerificationType() {}

var (
	_ VerificationType    = (*VerificationTypeSuccessful)(nil)
	_ OTPVerificationType = (*VerificationTypeSuccessful)(nil)
)

// VerificationTypeFailed indicates that the verification was failed.
// This can be used to increment the failed attempts counter.
type VerificationTypeFailed struct{}

// isOTPVerificationType implements [OTPVerificationType].
func (v *VerificationTypeFailed) isOTPVerificationType() {}

// isVerificationType implements [VerificationType].
func (v *VerificationTypeFailed) isVerificationType() {}

var (
	_ VerificationType    = (*VerificationTypeFailed)(nil)
	_ OTPVerificationType = (*VerificationTypeFailed)(nil)
)

// VerificationTypeSkipVerification indicates that the verification is skipped.
// The Value is used as the new value to be set.
// If VerifiedAt is zero, the current time will be used.
type VerificationTypeSkipVerification struct {
	VerifiedAt *time.Time
	Value      string
}

// isVerificationType implements [VerificationType].
func (v *VerificationTypeSkipVerification) isVerificationType() {}

var _ VerificationType = (*VerificationTypeSkipVerification)(nil)

// VerificationTypeInitCode indicates that a verification is started
type VerificationTypeInitCode struct {
	Code  []byte
	Value *string
	// CreatedAt is the time when the verification was created.
	// If zero, the current time will be used.
	CreatedAt time.Time
	Expiry    *time.Duration
}

// isOTPVerificationType implements [OTPVerificationType.]
func (v *VerificationTypeInitCode) isOTPVerificationType() {}

// isVerificationType implements [VerificationType.]
func (v *VerificationTypeInitCode) isVerificationType() {}

var (
	_ VerificationType    = (*VerificationTypeInitCode)(nil)
	_ OTPVerificationType = (*VerificationTypeInitCode)(nil)
)

// VerificationTypeUpdate indicates that a verification is updated.
// The fields that are non-nil will be updated.
type VerificationTypeUpdate struct {
	Code   []byte
	Value  *string
	Expiry *time.Duration
}

func (v *VerificationTypeUpdate) isVerificationType() {}

var _ VerificationType = (*VerificationTypeUpdate)(nil)
