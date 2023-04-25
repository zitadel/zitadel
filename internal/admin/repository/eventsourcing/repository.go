package eventsourcing

import (
	"context"

	auth_handler "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/handler"
	admin_view "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/static"
)

type Config struct {
	Spooler auth_handler.Config
}

func Start(ctx context.Context, conf Config, static static.Storage, dbClient *database.DB) error {
	view, err := admin_view.StartView(dbClient)
	if err != nil {
		return err
	}

	auth_handler.Register(ctx, conf.Spooler, view, static)

	return nil
}
