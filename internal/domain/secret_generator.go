package domain

type SecretGeneratorType int32

const (
	SecretGeneratorTypeUnspecified SecretGeneratorType = iota
	SecretGeneratorTypeInitCode
	SecretGeneratorTypeVerifyEmailCode
	SecretGeneratorTypeVerifyPhoneCode
	SecretGeneratorTypePasswordResetCode
	SecretGeneratorTypePasswordlessInitCode
	SecretGeneratorTypeAppSecret

	secretGeneratorTypeCount
)

type SecretGeneratorState int32

const (
	SecretGeneratorStateUnspecified SecretGeneratorState = iota
	SecretGeneratorStateActive
	SecretGeneratorStateRemoved
)
