package user

const (
	UserV1AddedType                     = userEventTypePrefix + "added"
	UserV1RegisteredType                = userEventTypePrefix + "selfregistered"
	UserV1InitialCodeAddedType          = userEventTypePrefix + "initialization.code.added"
	UserV1InitialCodeSentType           = userEventTypePrefix + "initialization.code.sent"
	UserV1InitializedCheckSucceededType = userEventTypePrefix + "initialization.check.succeeded"
	UserV1InitializedCheckFailedType    = userEventTypePrefix + "initialization.check.failed"
	UserV1SignedOutType                 = userEventTypePrefix + "signed.out"

	userV1PasswordEventTypePrefix    = userEventTypePrefix + "password."
	UserV1PasswordChangedType        = userV1PasswordEventTypePrefix + "changed"
	UserV1PasswordCodeAddedType      = userV1PasswordEventTypePrefix + "code.added"
	UserV1PasswordCodeSentType       = userV1PasswordEventTypePrefix + "code.sent"
	UserV1PasswordCheckSucceededType = userV1PasswordEventTypePrefix + "check.succeeded"
	UserV1PasswordCheckFailedType    = userV1PasswordEventTypePrefix + "check.failed"

	userV1EmailEventTypePrefix        = userEventTypePrefix + "email."
	UserV1EmailChangedType            = userV1EmailEventTypePrefix + "changed"
	UserV1EmailVerifiedType           = userV1EmailEventTypePrefix + "verified"
	UserV1EmailVerificationFailedType = userV1EmailEventTypePrefix + "verification.failed"
	UserV1EmailCodeAddedType          = userV1EmailEventTypePrefix + "code.added"
	UserV1EmailCodeSentType           = userV1EmailEventTypePrefix + "code.sent"

	userV1PhoneEventTypePrefix        = userEventTypePrefix + "phone."
	UserV1PhoneChangedType            = userV1PhoneEventTypePrefix + "changed"
	UserV1PhoneRemovedType            = userV1PhoneEventTypePrefix + "removed"
	UserV1PhoneVerifiedType           = userV1PhoneEventTypePrefix + "verified"
	UserV1PhoneVerificationFailedType = userV1PhoneEventTypePrefix + "verification.failed"
	UserV1PhoneCodeAddedType          = userV1PhoneEventTypePrefix + "code.added"
	UserV1PhoneCodeSentType           = userV1PhoneEventTypePrefix + "code.sent"

	userV1ProfileEventTypePrefix = userEventTypePrefix + "profile."
	UserV1ProfileChangedType     = userV1ProfileEventTypePrefix + "changed"

	userV1AddressEventTypePrefix = userEventTypePrefix + "address."
	UserV1AddressChangedType     = userV1AddressEventTypePrefix + "changed"

	userV1MFAEventTypePrefix = userEventTypePrefix + "mfa."
	UserV1MFAInitSkippedType = userV1MFAOTPEventTypePrefix + "init.skipped"

	userV1MFAOTPEventTypePrefix    = userV1MFAEventTypePrefix + "otp."
	UserV1MFAOTPAddedType          = userV1MFAOTPEventTypePrefix + "added"
	UserV1MFAOTPRemovedType        = userV1MFAOTPEventTypePrefix + "removed"
	UserV1MFAOTPVerifiedType       = userV1MFAOTPEventTypePrefix + "verified"
	UserV1MFAOTPCheckSucceededType = userV1MFAOTPEventTypePrefix + "check.succeeded"
	UserV1MFAOTPCheckFailedType    = userV1MFAOTPEventTypePrefix + "check.failed"
)
