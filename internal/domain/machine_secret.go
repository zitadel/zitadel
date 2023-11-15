package domain

import (
	"github.com/zitadel/zitadel/v2/internal/crypto"
	"github.com/zitadel/zitadel/v2/internal/errors"
)

func NewMachineClientSecret(generator crypto.Generator) (*crypto.CryptoValue, string, error) {
	cryptoValue, stringSecret, err := crypto.NewCode(generator)
	if err != nil {
		return nil, "", errors.ThrowInternal(err, "MODEL-57cjsiw", "Errors.User.Machine.Secret.CouldNotGenerate")
	}
	return cryptoValue, stringSecret, nil
}
