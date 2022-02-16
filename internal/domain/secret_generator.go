package domain

const (
	InitCodeGeneratorType          string = "InitCode"
	VerifyEmailCodeGeneratorType   string = "VerifyEmailCode"
	VerifyPhoneCodeGeneratorType   string = "VerifyPhoneCode"
	PasswordResetCodeGeneratorType string = "PasswordResetCode"
	PasswordlessCodeGeneratorType  string = "PasswordlessInitCode"
	AppSecretGeneratorType         string = "ApplicationSecret"
)

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
