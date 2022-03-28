package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
)

func newCryptoCode(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm) (value *crypto.CryptoValue, expiry time.Duration, err error) {
	wm := command.NewInstanceSecretGeneratorConfigWriteModel(typ)
	events, err := filter(ctx, wm.Query())
	if err != nil {
		return nil, -1, err
	}
	wm.AppendEvents(events...)
	if err := wm.Reduce(); err != nil {
		return nil, -1, err
	}
	config := crypto.GeneratorConfig{
		Length:              wm.Length,
		Expiry:              wm.Expiry,
		IncludeLowerLetters: wm.IncludeLowerLetters,
		IncludeUpperLetters: wm.IncludeUpperLetters,
		IncludeDigits:       wm.IncludeDigits,
		IncludeSymbols:      wm.IncludeSymbols,
	}
	value, _, err = crypto.NewCode(crypto.NewEncryptionGenerator(config, alg))
	if err != nil {
		return nil, -1, err
	}
	return value, config.Expiry, nil
}
