package domain

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
)

type generateCodeCommand struct {
	code  string
	value *crypto.CryptoValue
}

type CryptoRepository interface {
	GetEncryptionConfig(ctx context.Context) (*crypto.GeneratorConfig, error)
}

func (cmd *generateCodeCommand) Execute(ctx context.Context, opts *CommandOpts) error {
	config, err := cryptoRepo(opts.DB).GetEncryptionConfig(ctx)
	if err != nil {
		return err
	}
	generator := crypto.NewEncryptionGenerator(*config, userCodeAlgorithm)
	cmd.value, cmd.code, err = crypto.NewCode(generator)
	return err
}
