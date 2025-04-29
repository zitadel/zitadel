package command

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/v2/pattern"
	"github.com/zitadel/zitadel/internal/crypto"
)

type generateCode struct {
	set       func(code string)
	generator pattern.Query[crypto.Generator]
}

func GenerateCode(set func(code string), generator pattern.Query[crypto.Generator]) *generateCode {
	return &generateCode{
		set:       set,
		generator: generator,
	}
}

var _ pattern.Command = (*generateCode)(nil)

// Execute implements [pattern.Command].
func (cmd *generateCode) Execute(ctx context.Context) error {
	if err := cmd.generator.Execute(ctx); err != nil {
		return err
	}
	value, code, err := crypto.NewCode(cmd.generator.Result())
	_ = value
	if err != nil {
		return err
	}
	cmd.set(code)
	return nil
}

// Name implements [pattern.Command].
func (*generateCode) Name() string {
	return "command.generate_code"
}
