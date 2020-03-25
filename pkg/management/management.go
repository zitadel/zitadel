package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/errors"
	app "github.com/caos/zitadel/internal/management"
	"github.com/caos/zitadel/pkg/management/api"
)

type Config struct {
	App *app.Config
	API *api.Config
}

func Start(ctx context.Context, config *Config, authZ *auth.Config) error {
	return errors.ThrowUnimplemented(nil, "MANAG-h3k3x", "not implemented yet") //TODO: implement
}
