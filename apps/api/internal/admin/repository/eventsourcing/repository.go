package eventsourcing

import (
	"context"

	admin_handler "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/handler"
	admin_view "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/static"
)

type Config struct {
	Spooler admin_handler.Config
}

func Start(ctx context.Context, conf Config, static static.Storage, dbClient *database.DB, queries *query.Queries) error {
	view, err := admin_view.StartView(dbClient)
	if err != nil {
		return err
	}

	admin_handler.Register(ctx, conf.Spooler, view, static)
	admin_handler.Start(ctx)

	return nil
}
