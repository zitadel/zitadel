package console

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
)

type Config struct {
	Port      string
	StaticDir string
}

func Start(ctx context.Context, config *Config) error {
	return errors.ThrowUnimplemented(nil, "CONSO-4cT5D", "not implemented yet") //TODO: implement
}
