package command

import (
	"context"
	"os"
	"testing"
	"time"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"

	new_db "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/cmd/initialise"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/postgres"
	"github.com/zitadel/zitadel/internal/eventstore"
	es_sql "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	es_v3 "github.com/zitadel/zitadel/internal/eventstore/v3"
	repo "github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type pausedMembershipPusher struct {
	eventstore.Pusher
	ready, release chan struct{}
}

func (p *pausedMembershipPusher) WithInstanceLock(ctx context.Context, change func(context.Context) error) error {
	return p.Pusher.(interface {
		WithInstanceLock(context.Context, func(context.Context) error) error
	}).WithInstanceLock(ctx, change)
}

func (p *pausedMembershipPusher) Push(ctx context.Context, client new_db.QueryExecutor, commands ...eventstore.Command) ([]eventstore.Event, error) {
	for _, command := range commands {
		if command.Type() == repo.GroupUserAddedEventType {
			close(p.ready)
			<-p.release
			break
		}
	}
	return p.Pusher.Push(ctx, client, commands...)
}

func TestGroupLifecycle_database(t *testing.T) {
	config, stop := postgres.StartEmbedded()
	t.Cleanup(stop)
	connection, err := pgxpool.ParseConfig(config.GetConnectionURL())
	require.NoError(t, err)
	connection.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return es_v3.RegisterEventstoreTypes(ctx, conn)
	}
	pool, err := pgxpool.NewWithConfig(t.Context(), connection)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	client := &database.DB{DB: stdlib.OpenDBFromPool(pool), Pool: pool, Database: &postgres.Config{Database: "postgres", User: postgres.User{Username: "postgres"}}}
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	dbConfig := new(database.Config)
	dbConfig.SetConnector(&postgres.Config{Database: "postgres", User: postgres.User{Username: "postgres"}})
	require.NoError(t, initialise.ReadStmts())
	require.NoError(t, initialise.VerifyZitadel(t.Context(), client, *dbConfig))
	fields, err := os.ReadFile("../../cmd/setup/28.sql")
	require.NoError(t, err)
	_, err = client.ExecContext(t.Context(), string(fields))
	require.NoError(t, err)
	pool.Reset()
	for _, deletion := range []string{"user", "org"} {
		t.Run(deletion, func(t *testing.T) {
			ctx := authz.WithInstanceID(t.Context(), t.Name())
			pusher := &pausedMembershipPusher{Pusher: es_v3.NewEventstore(client), ready: make(chan struct{}), release: make(chan struct{})}
			es := eventstore.NewEventstore(&eventstore.Config{Pusher: pusher, Querier: &es_sql.Postgres{DB: client}})
			c := &Commands{eventstore: es, checkPermission: newMockPermissionCheckAllowed()}
			_, err := es.Push(ctx, org.NewOrgAddedEvent(ctx, &org.NewAggregate("org").Aggregate, "org"), org.NewDomainPolicyAddedEvent(ctx, &org.NewAggregate("org").Aggregate, true, true, true), addNewUserEvent("user", "org"))
			require.NoError(t, err)
			_, err = c.CreateGroup(ctx, &CreateGroup{ObjectRoot: models.ObjectRoot{AggregateID: "group", ResourceOwner: "org"}, Name: "group"})
			require.NoError(t, err)
			added, deleted := make(chan error, 1), make(chan error, 1)
			go func() {
				_, err := c.MutateGroup(ctx, "group", "org", []GroupMutation{{AddMembers: []string{"user"}}})
				added <- err
			}()
			<-pusher.ready // member existence has been validated, but nothing has committed yet.
			go func() {
				var err error
				if deletion == "user" {
					_, err = c.RemoveUserV2(ctx, "user", "org", nil, nil, nil)
				} else {
					_, err = c.RemoveOrg(ctx, "org", nil, true)
				}
				deleted <- err
			}()
			select {
			case err := <-deleted:
				t.Fatalf("deletion overtook the membership transaction: %v", err)
			case <-time.After(100 * time.Millisecond):
			}
			close(pusher.release)
			require.NoError(t, <-added)
			require.NoError(t, <-deleted)
			group, err := c.getGroupWriteModelByID(ctx, "group", "org", []string{})
			require.NoError(t, err)
			require.Empty(t, group.existingUserIDs)
			// A fresh attempt after deletion must also reject the stale reference.
			_, err = c.MutateGroup(ctx, "group", "org", []GroupMutation{{AddMembers: []string{"user"}}})
			require.Error(t, err)
			if deletion == "org" {
				require.False(t, group.State.Exists())
			}
		})
	}
}
