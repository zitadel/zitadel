package migration_test

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/embedded"
)

func TestMigrate(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		stmt string
		args []any
		res  []any
	}{
		{
			name: "schema",
			stmt: "SELECT EXISTS(SELECT 1 FROM information_schema.schemata where schema_name = 'zitadel') ;",
			res:  []any{true},
		},
		{
			name: "001",
			stmt: "SELECT EXISTS(SELECT 1 FROM pg_catalog.pg_tables WHERE schemaname = 'zitadel' and tablename=$1)",
			args: []any{"instances"},
			res:  []any{true},
		},
	}

	ctx := context.Background()

	connector, stop, err := embedded.StartEmbedded()
	require.NoError(t, err, "failed to start embedded postgres")
	defer stop()

	client, err := connector.Connect(ctx)
	require.NoError(t, err, "failed to connect to embedded postgres")

	err = client.(database.Migrator).Migrate(ctx)
	require.NoError(t, err, "failed to execute migration steps")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]any, len(tt.res))
			for i := range got {
				got[i] = new(any)
				tt.res[i] = gu.Ptr(tt.res[i])
			}

			require.NoError(t, client.QueryRow(ctx, tt.stmt, tt.args...).Scan(got...), "failed to execute check query")

			assert.Equal(t, tt.res, got, "query result does not match")
		})
	}
}
