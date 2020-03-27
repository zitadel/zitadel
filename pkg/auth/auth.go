package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/auth"
	app "github.com/caos/zitadel/internal/auth"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/pkg/auth/api"
)

type Config struct {
	App *app.Config
	API *api.Config
}

func Start(ctx context.Context, config Config, authZ auth.Config) error {
	return errors.ThrowUnimplemented(nil, "AUTH-l7Hdx", "not implemented yet") //TODO: implement
}
