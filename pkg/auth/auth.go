package auth

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
)

type Config struct {
}

func Start(ctx context.Context, config Config) error {
	return errors.ThrowUnimplemented(nil, "", "") //TODO: implement
}
