package domain

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

// HashTOTPCode creates a deterministic hash of the user-provided TOTP code.
// This is mainly for compliance, so that we don't store clear text codes.
// As codes are only 6 digits, they have a very low entropy:
// instanceID and userID are combined as salt to counter rainbow table usage.
func HashTOTPCode(instanceID, userID, code string) string {
	h := hmac.New(sha256.New, []byte(instanceID+userID))
	h.Write([]byte(code))
	return hex.EncodeToString(h.Sum(nil))
}

const (
	validPeriods  = 1 + 2*TOTPSkew
	validDuration = TOTPPeriod * time.Second * validPeriods
)

type TOTPHistory struct {
	start       time.Time
	recentCodes []string
}

// AddRecent appends a used codeHash if its timestamp is
// within the current period window.
func (h *TOTPHistory) AddRecent(ts time.Time, codeHash string) {
	if h.start.IsZero() {
		h.start = time.Now().Add(-validDuration)
	}
	if ts.After(h.start) {
		h.recentCodes = append(h.recentCodes, codeHash)
	}
}

// CheckReuse returns an error if the provided codeHash was recently used.
func (h *TOTPHistory) CheckReuse(codeHash string) error {
	if slices.Contains(h.recentCodes, codeHash) {
		return zerrors.ThrowInvalidArgument(nil, "TOTP-Auw0a", "Errors.User.MFA.OTP.Reused")
	}
	return nil
}
