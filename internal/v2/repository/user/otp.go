package user

const (
	otpEventPrefix                = mfaEventPrefix + "otp."
	HumanMFAOTPAddedType          = otpEventPrefix + "added"
	HumanMFAOTPVerifiedType       = otpEventPrefix + "verified"
	HumanMFAOTPRemovedType        = otpEventPrefix + "removed"
	HumanMFAOTPCheckSucceededType = otpEventPrefix + "check.succeeded"
	HumanMFAOTPCheckFailedType    = otpEventPrefix + "check.failed"
)
