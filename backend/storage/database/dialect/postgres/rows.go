package postgres

import (
	"github.com/jackc/pgx/v5"

	"github.com/zitadel/zitadel/backend/storage/database"
)

var _ database.Rows = (*Rows)(nil)

type Rows struct{ pgx.Rows }

// Close implements [database.Rows].
// Subtle: this method shadows the method (Rows).Close of Rows.Rows.
func (r *Rows) Close() error {
	r.Rows.Close()
	return nil
}
