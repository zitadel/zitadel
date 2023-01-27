package domain

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
)

func NewMachineClientSecret(generator crypto.Generator) (*crypto.CryptoValue, string, error) {
	cryptoValue, stringSecret, err := crypto.NewCode(generator)
	if err != nil {
		return nil, "", errors.ThrowInternal(err, "MODEL-gH2Wl", "Errors.Project.CouldNotGenerateClientSecret")
	}
	return cryptoValue, stringSecret, nil
}
