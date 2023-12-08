package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

type Phone struct {
	Number   domain.PhoneNumber
	Verified bool

	// ReturnCode is used if the Verified field is false
	ReturnCode bool
}

func (c *Commands) newPhoneCode(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.EncryptionAlgorithm) (*CryptoCode, error) {
	return c.newCode(ctx, filter, domain.SecretGeneratorTypeVerifyPhoneCode, alg)
}

func (c *Commands) newPhoneCodeFunc(alg crypto.EncryptionAlgorithm) getCryptoCodeFunc {
	return func(ctx context.Context) (*CryptoCode, error) {
		return c.newPhoneCode(ctx, c.eventstore.Filter, alg)
	}
}

func (c *Commands) verifyPhoneCodeFunc(alg crypto.EncryptionAlgorithm) verifyCryptoCodeFunc {
	return func(ctx context.Context, creation time.Time, expiry time.Duration, crypted *crypto.CryptoValue, plain string) error {
		return verifyCryptoCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyPhoneCode, alg, creation, expiry, crypted, plain)
	}
}
