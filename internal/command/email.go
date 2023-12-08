package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

type Email struct {
	Address  domain.EmailAddress
	Verified bool

	// ReturnCode is used if the Verified field is false
	ReturnCode bool

	// URLTemplate can be used to specify a custom link to be sent in the mail verification
	URLTemplate string
}

func (e *Email) Validate() error {
	return e.Address.Validate()
}

func (c *Commands) newEmailCode(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.EncryptionAlgorithm) (*CryptoCode, error) {
	return c.newCode(ctx, filter, domain.SecretGeneratorTypeVerifyEmailCode, alg)
}

func (c *Commands) newEmailCodeFunc(alg crypto.EncryptionAlgorithm) getCryptoCodeFunc {
	return func(ctx context.Context) (*CryptoCode, error) {
		return c.newEmailCode(ctx, c.eventstore.Filter, alg)
	}
}

func (c *Commands) verifyEmailCodeFunc(alg crypto.EncryptionAlgorithm) verifyCryptoCodeFunc {
	return func(ctx context.Context, creation time.Time, expiry time.Duration, crypted *crypto.CryptoValue, plain string) error {
		return verifyCryptoCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyEmailCode, alg, creation, expiry, crypted, plain)
	}
}
