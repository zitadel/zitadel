package setup

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 48_fix_in_tx_order_type/01_get_current_type.sql
	getInTxOrderTypeQuery string
	//go:embed 48_fix_in_tx_order_type/02_alter_column.sql
	alterInTxOrderTypeQuery string
)

type FixInTxOrderType struct {
	dbClient *database.DB
}

func (mig *FixInTxOrderType) Execute(ctx context.Context, _ eventstore.Event) error {
	// The INT type in CockroachDB is actually a BIGINT.
	// https://www.cockroachlabs.com/docs/v24.3/int
	// That means altering the type from  BIGINT to INT doesn't make any sense.
	if mig.dbClient.Database.Type() == "cockroach" {
		return nil
	}

	var currentType string
	err := mig.dbClient.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&currentType)
	}, getInTxOrderTypeQuery)
	if err != nil {
		return err
	}
	logging.WithFields("migration", mig.String(), "current_type", currentType).Info("execute statement")
	if strings.EqualFold(currentType, "integer") {
		return nil
	}

	logging.WithFields("migration", mig.String()).Info("executing ALTER TABLE")
	if _, err := mig.dbClient.ExecContext(ctx, alterInTxOrderTypeQuery); err != nil {
		return fmt.Errorf("%s %s: %w", mig.String(), alterInTxOrderTypeQuery, err)
	}
	return nil
}

func (mig *FixInTxOrderType) String() string {
	return "48_fix_in_tx_order_type"
}
