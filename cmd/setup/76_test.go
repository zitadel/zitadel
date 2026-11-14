package setup

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/cmd/initialise"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/postgres"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_es "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	new_es "github.com/zitadel/zitadel/internal/eventstore/v3"
	"github.com/zitadel/zitadel/internal/migration"
)

type noopMigration struct{ ran bool }

func (m *noopMigration) Execute(context.Context, eventstore.Event) error {
	m.ran = true
	return nil
}

func (*noopMigration) String() string { return "test_noop_migration" }

// A database created by an older release has no claimant columns, but every migration pushes a started event with a unique constraint.
func TestUniqueConstraintClaimants_BootstrapBeforeMigrations(t *testing.T) {
	config, stop := postgres.StartEmbedded()
	defer stop()
	db, err := sql.Open("pgx", config.GetConnectionURL())
	require.NoError(t, err)
	defer db.Close()
	client := &database.DB{DB: db, Database: &postgres.Config{Database: "postgres"}}
	ctx := context.Background()
	dbConfig := database.Config{}
	dbConfig.SetConnector(&postgres.Config{User: postgres.User{Username: "postgres"}, Database: "postgres"})
	require.NoError(t, initialise.VerifyZitadel(ctx, client, dbConfig))
	_, err = db.ExecContext(ctx, "ALTER TABLE eventstore.unique_constraints DROP COLUMN claimant_type, DROP COLUMN claimant_id")
	require.NoError(t, err)

	esV3 := new_es.NewEventstore(client)
	es := eventstore.NewEventstore(&eventstore.Config{Querier: old_es.NewPostgres(client), Pusher: esV3, Searcher: esV3})
	step := new(noopMigration)
	require.Error(t, migration.Migrate(ctx, es, step), "the started event must not be pushable before the bootstrap")
	require.False(t, step.ran)

	require.NoError(t, (&UniqueConstraintClaimants{dbClient: client}).Execute(ctx, nil))
	require.NoError(t, migration.Migrate(ctx, es, step))
	require.True(t, step.ran)
}
