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

func TestUniqueConstraintOwners_Execute(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(uniqueConstraintsOwnersColumn)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta(uniqueConstraintsOwnersGIN)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mig := &UniqueConstraintOwners{dbClient: &database.DB{DB: db}}
	assert.NoError(t, mig.Execute(context.Background(), nil))
	assert.NoError(t, mock.ExpectationsWereMet())
}
