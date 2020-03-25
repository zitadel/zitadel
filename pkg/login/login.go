package login

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	app "github.com/caos/zitadel/internal/login"
	"github.com/caos/zitadel/pkg/login/api"
)

type Config struct {
	App *app.Config
	API *api.Config
}

func Start(ctx context.Context, config *Config) error {
	return errors.ThrowUnimplemented(nil, "LOGIN-3fwvD", "not implemented yet") //TODO: implement
}
