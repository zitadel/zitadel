package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

type Email struct {
	Address  domain.EmailAddress
	Verified bool
}

func (e *Email) Validate() error {
	return e.Address.Validate()
}

func newEmailCode(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.EncryptionAlgorithm) (*cryptoCode, error) {
	return newCryptoCode(ctx, filter, domain.SecretGeneratorTypeVerifyEmailCode, alg)
}
