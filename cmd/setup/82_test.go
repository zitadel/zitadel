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
	assert.Contains(t, currentStatesInTxOrder, "BEFORE INSERT OR UPDATE")
	assert.Contains(t, currentStatesInTxOrder, "NEW.in_tx_order := NULL")
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
