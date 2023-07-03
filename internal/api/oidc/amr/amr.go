package amr

const (
	// DEPRECATED: use `PWD` instead
	Password     = "password"
	PWD          = "pwd"
	MFA          = "mfa"
	OTP          = "otp"
	UserPresence = "user"
)

type AuthenticationMethodReference interface {
	IsPasswordChecked() bool
	IsPasskeyChecked() bool
	IsU2FChecked() bool
	IsOTPChecked() bool
}

func List(model AuthenticationMethodReference) []string {
	amr := make([]string, 0)
	if model.IsPasswordChecked() {
		amr = append(amr, PWD)
	}
	if model.IsPasskeyChecked() || model.IsU2FChecked() {
		amr = append(amr, UserPresence)
	}
	if model.IsOTPChecked() {
		amr = append(amr, OTP)
	}
	if model.IsPasskeyChecked() || len(amr) >= 2 {
		amr = append(amr, MFA)
	}
	return amr
}
