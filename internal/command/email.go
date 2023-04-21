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

	// ReturnCode is used if the Verified field is false
	ReturnCode bool

	// VerificationCode is set by the command
	VerificationCode string
}

func (e *Email) Validate() error {
	return e.Address.Validate()
}

func newEmailCode(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.EncryptionAlgorithm) (*CryptoCodeWithExpiry, error) {
	return newCryptoCodeWithExpiry(ctx, filter, domain.SecretGeneratorTypeVerifyEmailCode, alg)
}
