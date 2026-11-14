package setup

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

//go:embed 78.sql
var backfillGroupLifecycleFields string

type GroupLifecycleFields struct{ dbClient *database.DB }

func (mig *GroupLifecycleFields) Execute(ctx context.Context, _ eventstore.Event) error {
	// A canceled concurrent build leaves an invalid index that IF NOT EXISTS would skip.
	var valid bool
	err := mig.dbClient.QueryRowContext(ctx, func(row *sql.Row) error { return row.Scan(&valid) },
		`SELECT COALESCE((SELECT indisvalid FROM pg_index WHERE indexrelid = to_regclass('eventstore.f_group_lifecycle_idx')), true)`)
	if err != nil {
		return err
	}
	if !valid {
		if _, err = mig.dbClient.ExecContext(ctx, "DROP INDEX CONCURRENTLY eventstore.f_group_lifecycle_idx"); err != nil {
			return err
		}
	}
	// User lookups use f_text_value_idx; organization cascades need a scoped index too.
	_, err = mig.dbClient.ExecContext(ctx, `CREATE INDEX CONCURRENTLY IF NOT EXISTS f_group_lifecycle_idx ON eventstore.fields (instance_id, resource_owner, aggregate_id) WHERE aggregate_type = 'group'`)
	if err != nil {
		return err
	}
	var instances []string
	err = mig.dbClient.QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				return err
			}
			instances = append(instances, id)
		}
		return rows.Err()
	}, "SELECT DISTINCT instance_id FROM eventstore.events2 WHERE aggregate_type = 'group' ORDER BY instance_id")
	if err != nil {
		return err
	}
	for _, id := range instances {
		if err := mig.backfillInstance(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (mig *GroupLifecycleFields) backfillInstance(ctx context.Context, id string) (err error) {
	tx, err := mig.dbClient.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { err = database.CloseTransaction(tx, err) }()
	_, err = tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock(hashtext('zitadel.group-lifecycle'), hashtext($1))", id)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, backfillGroupLifecycleFields, id)
	return err
}

func (*GroupLifecycleFields) String() string { return "78_group_lifecycle_fields" }
