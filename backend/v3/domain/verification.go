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

type VerificationTypeIsVerified struct {
	VerifiedAt time.Time
}

func (v *VerificationTypeIsVerified) isVerificationType() {}

var _ VerificationType = (*VerificationTypeIsVerified)(nil)

// VerificationTypeSet is used to update an existing verification
// It sets the fields which are not nil
type VerificationTypeSet struct {
	CreatedAt *time.Time
	Expiry    *time.Duration
	Code      *[]byte
	Value     *string
}

func (v *VerificationTypeSet) isVerificationType() {}

var _ VerificationType = (*VerificationTypeSet)(nil)

// VerificationTypeVerifiedValue is used to skip the verification process by directly setting the verified value
// It includes updating the updatedAt field of the object based on VerifiedAt.
type VerificationTypeVerifiedValue struct {
	Value string
	// VerifiedAt is the time when the verification was completed
	// If nil or [time.Time.IsZero], NOW() will be used
	VerifiedAt *time.Time
}

func (v *VerificationTypeVerifiedValue) isVerificationType() {}

var _ VerificationType = (*VerificationTypeVerifiedValue)(nil)
