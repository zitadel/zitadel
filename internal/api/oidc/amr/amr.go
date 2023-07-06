// Package amr maps zitadel session factors to Authentication Method Reference Values
// as defined in [RFC 8176, section 2].
//
// [RFC 8176, section 2]: https://datatracker.ietf.org/doc/html/rfc8176#section-2
package amr

const (
	// Password states that the users password has been verified
	// Deprecated: use `PWD` instead
	Password = "password"
	// PWD states that the users password has been verified
	PWD = "pwd"
	// MFA states that multiple factors have been verified (e.g. pwd and otp or passkey)
	MFA = "mfa"
	// OTP states that a one time password has been verified (e.g. TOTP)
	OTP = "otp"
	// UserPresence states that the end users presence has been verified (e.g. passkey and u2f)
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
