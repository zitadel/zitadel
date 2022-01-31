package domain

const (
	InitCodeGeneratorType          string = "InitCode"
	VerifyEmailCodeGeneratorType   string = "VerifyEmailCode"
	VerifyPhoneCodeGeneratorType   string = "VerifyPhoneCode"
	PasswordResetCodeGeneratorType string = "PasswordResetCode"
	PasswordlessCodeGeneratorType  string = "PasswordlessInitCode"
	AppSecretGeneratorType         string = "ApplicationSecret"
)

type SecretGeneratorState int32

const (
	SecretGeneratorStateUnspecified SecretGeneratorState = iota
	SecretGeneratorStateActive
	SecretGeneratorStateRemoved
)
