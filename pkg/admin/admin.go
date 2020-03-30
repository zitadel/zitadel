package admin

import (
	"context"

	app "github.com/caos/zitadel/internal/admin"
	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/pkg/admin/api"
)

type Config struct {
	App app.Config
	API api.Config
}

func Start(ctx context.Context, config Config, authZ auth.Config) error {
	return errors.ThrowUnimplemented(nil, "ADMIN-n8vw5", "not implemented yet") //TODO: implement
}
