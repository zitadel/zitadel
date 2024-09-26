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

	// ReturnCode is used if the Verified field is false
	ReturnCode bool
}

// newPhoneCode generates a new code to be sent out to via SMS or
// returns the ID of the external code provider (e.g. when using Twilio verification API)
func (c *Commands) newPhoneCode(ctx context.Context, filter preparation.FilterToQueryReducer, secretGeneratorType domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, defaultConfig *crypto.GeneratorConfig) (*EncryptedCode, string, error) {
	externalID, err := c.activeSMSProvider(ctx)
	if err != nil {
		return nil, "", err
	}
	if externalID != "" {
		return nil, externalID, nil
	}
	code, err := c.newEncryptedCodeWithDefault(ctx, filter, secretGeneratorType, alg, defaultConfig)
	return code, "", err
}
