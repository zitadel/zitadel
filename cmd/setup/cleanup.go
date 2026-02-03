package setup

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_es "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	new_es "github.com/zitadel/zitadel/internal/eventstore/v3"
	"github.com/zitadel/zitadel/internal/migration"
)

func NewCleanup() *cobra.Command {
	return &cobra.Command{
		Use:   "cleanup",
		Short: "cleans up migration if they got stuck",
		Long:  `cleans up migration if they got stuck`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				logging.OnError(cmd.Context(), err).Error("zitadel setup cleanup command failed")
			}()
			config, shutdown, err := NewConfig(cmd, viper.GetViper())
			if err != nil {
				return err
			}
			defer func() {
				err = errors.Join(err, shutdown(cmd.Context()))
			}()
			return Cleanup(cmd.Context(), config)
		},
	}
}

func Cleanup(ctx context.Context, config *Config) error {
	logging.Info(ctx, "cleanup started")

	dbClient, err := database.Connect(config.Database, false)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	config.Eventstore.Pusher = new_es.NewEventstore(dbClient)
	config.Eventstore.Querier = old_es.NewPostgres(dbClient)
	es := eventstore.NewEventstore(config.Eventstore)

	step, err := migration.LastStuckStep(ctx, es)
	if err != nil {
		return fmt.Errorf("unable to query latest migration: %w", err)
	}

	if step == nil {
		logging.Info(ctx, "there is no stuck migration please run `zitadel setup`")
		return nil
	}

	logging.Info(ctx, "cleanup migration", "name", step.Name)

	err = migration.CancelStep(ctx, es, step)
	if err != nil {
		return fmt.Errorf("cleanup migration failed please retry: %w", err)
	}
	return nil
}
