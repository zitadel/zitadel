package domain

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
