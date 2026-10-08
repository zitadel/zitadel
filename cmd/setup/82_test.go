package setup

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
)

func TestCurrentStatesInTxOrderSQL(t *testing.T) {
	assert.Contains(t, currentStatesInTxOrder, "ADD COLUMN IF NOT EXISTS in_tx_order INTEGER")
	assert.Contains(t, currentStatesInTxOrder, "zitadel.keep_in_tx_order")
	assert.Contains(t, currentStatesInTxOrder, "BEFORE UPDATE OF in_tx_order")
	assert.Contains(t, currentStatesInTxOrder, "BEFORE UPDATE ON projections.current_states")
	assert.Contains(t, currentStatesInTxOrder, "current_states_keep_in_tx_order_opt_in")
	assert.Contains(t, currentStatesInTxOrder, "NEW.in_tx_order := NULL")
	assert.Contains(t, currentStatesInTxOrder, "SET in_tx_order = e.in_tx_order")
	assert.Contains(t, currentStatesInTxOrder, "FROM eventstore.events2")
	assert.NotContains(t, currentStatesInTxOrder, "BEFORE INSERT OR UPDATE")
	assert.NotContains(t, currentStatesInTxOrder, "filter_offset = e.in_tx_order")
}

func TestCurrentStatesInTxOrder_Execute(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(currentStatesInTxOrder)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mig := &CurrentStatesInTxOrder{dbClient: &database.DB{DB: db}}
	assert.NoError(t, mig.Execute(context.Background(), nil))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCurrentStatesInTxOrder_String(t *testing.T) {
	assert.Equal(t, "82_current_states_in_tx_order", new(CurrentStatesInTxOrder).String())
}
