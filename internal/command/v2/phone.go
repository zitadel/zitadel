package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/ttacon/libphonenumber"
)

type Phone struct {
	Number   string
	Verified bool
}

func FormatPhoneNumber(number string) (string, error) {
	if number == "" {
		return "", nil
	}
	phoneNr, err := libphonenumber.Parse(number, libphonenumber.UNKNOWN_REGION)
	if err != nil {
		return "", errors.ThrowInvalidArgument(nil, "EVENT-so0wa", "Errors.User.Phone.Invalid")
	}
	number = libphonenumber.Format(phoneNr, libphonenumber.E164)
	return number, nil
}

func newPhoneCode(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.EncryptionAlgorithm) (value *crypto.CryptoValue, expiry time.Duration, err error) {
	return newCryptoCodeWithExpiry(ctx, filter, domain.SecretGeneratorTypeVerifyPhoneCode, alg)
}
