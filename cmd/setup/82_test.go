package setup

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
)

func TestCurrentStatesInTxOrderSQL(t *testing.T) {
	assert.Contains(t, currentStatesInTxOrderAddColumn, "ADD COLUMN IF NOT EXISTS in_tx_order %s")
	assert.NotContains(t, currentStatesInTxOrderAddColumn, "UPDATE")
	assert.NotContains(t, currentStatesInTxOrderAddColumn, "FROM eventstore")

	assert.Contains(t, currentStatesInTxOrderTriggers, "zitadel.keep_in_tx_order")
	assert.Contains(t, currentStatesInTxOrderTriggers, "BEFORE UPDATE OF in_tx_order")
	assert.Contains(t, currentStatesInTxOrderTriggers, "BEFORE UPDATE ON projections.current_states")
	assert.Contains(t, currentStatesInTxOrderTriggers, "current_states_keep_in_tx_order_opt_in")
	assert.Contains(t, currentStatesInTxOrderTriggers, "NEW.in_tx_order := NULL")
	assert.NotContains(t, currentStatesInTxOrderTriggers, "BEFORE INSERT OR UPDATE")
	assert.NotContains(t, currentStatesInTxOrderTriggers, "ADD COLUMN")
	assert.NotContains(t, currentStatesInTxOrderTriggers, "UPDATE projections.current_states")
	assert.NotContains(t, currentStatesInTxOrderTriggers, "FROM eventstore")

	assert.Contains(t, currentStatesInTxOrderBackfill, "SET in_tx_order = e.in_tx_order")
	assert.Contains(t, currentStatesInTxOrderBackfill, "FROM eventstore.events2")
	assert.NotContains(t, currentStatesInTxOrderBackfill, "ADD COLUMN")
	assert.NotContains(t, currentStatesInTxOrderBackfill, "e.in_tx_order <> 0")
	assert.NotContains(t, currentStatesInTxOrderBackfill, "filter_offset = e.in_tx_order")
}

func TestCurrentStatesInTxOrder_Execute(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT data_type FROM information_schema.columns WHERE table_schema = 'eventstore' AND table_name = 'events2' AND column_name = 'in_tx_order'`)).
		WillReturnRows(sqlmock.NewRows([]string{"data_type"}).AddRow("integer"))
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(currentStatesInTxOrderAddColumn, "integer"))).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta(currentStatesInTxOrderTriggers)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta(currentStatesInTxOrderBackfill)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mig := &CurrentStatesInTxOrder{dbClient: &database.DB{DB: db}}
	assert.NoError(t, mig.Execute(context.Background(), nil))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCurrentStatesInTxOrder_String(t *testing.T) {
	assert.Equal(t, "82_current_states_in_tx_order", new(CurrentStatesInTxOrder).String())
}
