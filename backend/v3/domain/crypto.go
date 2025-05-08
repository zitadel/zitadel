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

// String implements [Commander].
func (cmd *generateCodeCommand) String() string {
	return "generateCodeCommand"
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

var _ Commander = (*generateCodeCommand)(nil)
