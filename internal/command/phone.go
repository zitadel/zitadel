package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

type Phone struct {
	Number   domain.PhoneNumber
	Verified bool
}

func (c *Commands) newPhoneCode(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.EncryptionAlgorithm) (*CryptoCode, error) {
	return c.newCode(ctx, filter, domain.SecretGeneratorTypeVerifyPhoneCode, alg)
}
