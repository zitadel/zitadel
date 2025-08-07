package domain

//go:generate enumer -type SecretGeneratorType -transform snake -trimprefix SecretGeneratorType
type SecretGeneratorType int32

const (
	SecretGeneratorTypeUnspecified SecretGeneratorType = iota
	SecretGeneratorTypeInitCode
	SecretGeneratorTypeVerifyEmailCode
	SecretGeneratorTypeVerifyPhoneCode
	SecretGeneratorTypeVerifyDomain
	SecretGeneratorTypePasswordResetCode
	SecretGeneratorTypePasswordlessInitCode
	SecretGeneratorTypeAppSecret
	SecretGeneratorTypeOTPSMS
	SecretGeneratorTypeOTPEmail
	SecretGeneratorTypeInviteCode
	SecretGeneratorTypeSigningKey

	secretGeneratorTypeCount
)

func (t SecretGeneratorType) Valid() bool {
	return t > SecretGeneratorTypeUnspecified && t < secretGeneratorTypeCount
}

type SecretGeneratorState int32

const (
	SecretGeneratorStateUnspecified SecretGeneratorState = iota
	SecretGeneratorStateActive
	SecretGeneratorStateRemoved
)
