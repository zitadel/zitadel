package setup

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 80/*.sql
	users14InstanceResourceOwnerIndex embed.FS
)

const users14InstanceResourceOwnerIdx = "users14_instance_resource_owner_idx"

type Users14InstanceResourceOwnerIndex struct {
	dbClient *database.DB
}

func (mig *Users14InstanceResourceOwnerIndex) Execute(ctx context.Context, _ eventstore.Event) error {
	var exists bool
	err := mig.dbClient.QueryRowContext(ctx, func(r *sql.Row) error {
		return r.Scan(&exists)
	}, users14TableExistsQuery)
	if err != nil || !exists {
		return err
	}

	if err := mig.dropInvalidIndex(ctx); err != nil {
		return err
	}

	statements, err := readStatements(users14InstanceResourceOwnerIndex, "80")
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.Info(ctx, "execute statement", "file", stmt.file, "migration", mig.String())
		if _, err := mig.dbClient.ExecContext(ctx, stmt.query); err != nil {
			return fmt.Errorf("%s %s: %w", mig.String(), stmt.file, err)
		}
	}
	return nil
}

func (mig *Users14InstanceResourceOwnerIndex) dropInvalidIndex(ctx context.Context) error {
	var invalid bool
	err := mig.dbClient.QueryRowContext(ctx, func(r *sql.Row) error {
		return r.Scan(&invalid)
	}, users14InvalidIndexQuery, users14InstanceResourceOwnerIdx)
	if err != nil {
		return fmt.Errorf("%s check invalid index %s: %w", mig.String(), users14InstanceResourceOwnerIdx, err)
	}
	if !invalid {
		return nil
	}
	drop := dropInvalidIndexConcurrently + users14InstanceResourceOwnerIdx
	logging.Info(ctx, "drop invalid leftover index", "index", users14InstanceResourceOwnerIdx, "migration", mig.String())
	if _, err := mig.dbClient.ExecContext(ctx, drop); err != nil {
		return fmt.Errorf("%s drop invalid index %s: %w", mig.String(), users14InstanceResourceOwnerIdx, err)
	}
	return nil
}

func (mig *Users14InstanceResourceOwnerIndex) String() string {
	return "80_users14_instance_resource_owner_idx"
}
