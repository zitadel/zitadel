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

	// URLTemplate can be used to specify a custom link to be sent in the mail verification
	URLTemplate string
}

func (e *Email) Validate() error {
	return e.Address.Validate()
}

func (c *Commands) newEmailCode(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.EncryptionAlgorithm) (*CryptoCode, error) {
	return c.newCode(ctx, filter, domain.SecretGeneratorTypeVerifyEmailCode, alg)
}
