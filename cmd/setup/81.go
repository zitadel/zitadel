package setup

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

//go:embed 81.sql
var backfillGroupLifecycleFields string

// Skips through the events2 primary key one instance at a time instead of scanning all events.
const groupInstancesQuery = `WITH RECURSIVE instances AS (
    (SELECT instance_id FROM eventstore.events2 ORDER BY instance_id LIMIT 1)
    UNION ALL
    SELECT (SELECT e.instance_id FROM eventstore.events2 e WHERE e.instance_id > i.instance_id ORDER BY e.instance_id LIMIT 1)
    FROM instances i WHERE i.instance_id IS NOT NULL
)
SELECT i.instance_id FROM instances i
-- LIMIT keeps this a per-instance primary key probe; an EXISTS semi-join may be planned as a full scan.
CROSS JOIN LATERAL (SELECT 1 FROM eventstore.events2 e WHERE e.instance_id = i.instance_id AND e.aggregate_type = 'group' LIMIT 1) AS g`

// GroupLifecycleFields is repeatable: raising Rerun runs the backfill again,
// e.g. after older binaries wrote group events without maintaining the lifecycle fields.
type GroupLifecycleFields struct {
	dbClient *database.DB
	Rerun    uint64 `json:"rerun"`
}

func (mig *GroupLifecycleFields) Check(lastRun map[string]any) bool {
	rerun, _ := lastRun["rerun"].(float64)
	return uint64(rerun) != mig.Rerun
}

func (mig *GroupLifecycleFields) Execute(ctx context.Context, _ eventstore.Event) error {
	// User lookups use f_text_value_idx. Organization cascades search by owner; the index is not partial so generic plans binding aggregate_type can use it.
	// Member writes and removals match a single object of a group.
	for _, index := range [][2]string{
		{"f_owner_object_type_idx", "instance_id, resource_owner, aggregate_type, object_type"},
		{"f_aggregate_object_idx", "instance_id, aggregate_type, aggregate_id, object_type, object_id"},
	} {
		if err := mig.createIndex(ctx, index[0], index[1]); err != nil {
			return err
		}
	}
	var instances []string
	err := mig.dbClient.QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				return err
			}
			instances = append(instances, id)
		}
		return rows.Err()
	}, groupInstancesQuery)
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

func (mig *GroupLifecycleFields) createIndex(ctx context.Context, name, columns string) error {
	// A canceled concurrent build leaves an invalid index that IF NOT EXISTS would skip.
	var valid bool
	err := mig.dbClient.QueryRowContext(ctx, func(row *sql.Row) error { return row.Scan(&valid) },
		`SELECT COALESCE((SELECT indisvalid FROM pg_index WHERE indexrelid = to_regclass('eventstore.' || $1)), true)`, name)
	if err != nil {
		return err
	}
	if !valid {
		if _, err = mig.dbClient.ExecContext(ctx, "DROP INDEX CONCURRENTLY eventstore."+name); err != nil {
			return err
		}
	}
	_, err = mig.dbClient.ExecContext(ctx, "CREATE INDEX CONCURRENTLY IF NOT EXISTS "+name+" ON eventstore.fields ("+columns+")")
	return err
}

func (mig *GroupLifecycleFields) backfillInstance(ctx context.Context, id string) (err error) {
	tx, err := mig.dbClient.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { err = database.CloseTransaction(tx, err) }()
	_, err = tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock(hashtextextended('zitadel.group-lifecycle:' || $1, 0))", id)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, backfillGroupLifecycleFields, id)
	return err
}

func (*GroupLifecycleFields) String() string { return "81_group_lifecycle_fields" }
