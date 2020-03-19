package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
)

type Config struct {
}

func Start(ctx context.Context, config Config) error {
	return errors.ThrowUnimplemented(nil, "EVENT-1hfiu", "not implemented yet") //TODO: implement
}
