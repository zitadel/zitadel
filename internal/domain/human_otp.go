package domain

import (
	"slices"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// These constants are not exported by the totp library,
// but are important for our reuse prevention logic.
// They are taken from the [totp.Validate] implementation,
// and are compatible with Google-Authenticator and most clients.
// See [totp.ValidateOpts] for the meaning of these values.
const (
	TOTPPeriod = 30
	TOTPSkew   = 1
)

type TOTP struct {
	*ObjectDetails

	Secret string
	URI    string
}

func NewTOTPKey(issuer, accountName string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
		Period:      TOTPPeriod,
	})
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "TOTP-ieY3o", "Errors.Internal")
	}
	return key, nil
}

func VerifyTOTP(code string, secret *crypto.CryptoValue, cryptoAlg crypto.EncryptionAlgorithm) error {
	decrypt, err := crypto.DecryptString(secret, cryptoAlg)
	if err != nil {
		return err
	}

	// Use ValidateCustom so we have compile-time alignment of constants,
	// as [totp.Validate] does not export the default values.
	valid, err := totp.ValidateCustom(
		code,
		decrypt,
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    TOTPPeriod,
			Skew:      TOTPSkew,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		})
	if err != nil || !valid {
		return zerrors.ThrowInvalidArgument(err, "EVENT-8isk2", "Errors.User.MFA.OTP.InvalidCode")
	}
	return nil
}

const (
	checkPeriods  = 2 + 2*TOTPSkew // 1 extra period to account for DB clock skew.
	checkDuration = TOTPPeriod * time.Second * checkPeriods
)

type TOTPHistory struct {
	start        time.Time
	recentValues []*crypto.HMACValue
}

// AddRecent appends the HMAC value of a used code if its timestamp is
// within the current period window and value is not nil.
// A nil value is ignored, as it originates from a legacy event
// which did not record the used code yet.
func (h *TOTPHistory) AddRecent(ts time.Time, value *crypto.HMACValue) {
	if value == nil {
		return
	}

	// If the start is zero, we assume this is the first code being added,
	// and set the start to the beginning of the current window.
	// This ensures that we only store codes that are within the valid window.
	if h.start.IsZero() {
		h.start = time.Now().Add(-checkDuration)
	}
	if ts.After(h.start) {
		h.recentValues = append(h.recentValues, value)
	}
}

// CheckReuse returns an error if the provided code was recently used.
func (h *TOTPHistory) CheckReuse(code string) error {
	if slices.ContainsFunc(h.recentValues, func(v *crypto.HMACValue) bool {
		return v.Equal(code)
	}) {
		return zerrors.ThrowInvalidArgument(nil, "TOTP-Auw0a", "Errors.User.MFA.OTP.Reused")
	}
	return nil
}
