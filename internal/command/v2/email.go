package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
)

type Email struct {
	Address  string
	Verified bool
}

func (e *Email) Valid() bool {
	return e.Address != "" && domain.EmailRegex.MatchString(e.Address)
}

func newEmailCode(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.EncryptionAlgorithm) (value *crypto.CryptoValue, expiry time.Duration, err error) {
	return newCryptoCode(ctx, filter, domain.SecretGeneratorTypeVerifyEmailCode, alg)
}
