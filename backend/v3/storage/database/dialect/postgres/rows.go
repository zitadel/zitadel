package postgres

import (
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var (
	_ database.Rows            = (*Rows)(nil)
	_ database.CollectableRows = (*Rows)(nil)
)

type Rows struct{ pgx.Rows }

// Collect implements [database.CollectableRows].
// See [this page](https://github.com/georgysavva/scany/blob/master/dbscan/doc.go#L8) for additional details.
func (r *Rows) Collect(dest any) (err error) {
	defer func() {
		closeErr := r.Close()
		if err == nil {
			err = closeErr
		}
	}()
	return pgxscan.ScanAll(dest, r.Rows)
}

// CollectFirst implements [database.CollectableRows].
// See [this page](https://github.com/georgysavva/scany/blob/master/dbscan/doc.go#L8) for additional details.
func (r *Rows) CollectFirst(dest any) (err error) {
	defer func() {
		closeErr := r.Close()
		if err == nil {
			err = closeErr
		}
	}()
	return pgxscan.ScanRow(dest, r.Rows)
}

// CollectExactlyOneRow implements [database.CollectableRows].
// See [this page](https://github.com/georgysavva/scany/blob/master/dbscan/doc.go#L8) for additional details.
func (r *Rows) CollectExactlyOneRow(dest any) (err error) {
	defer func() {
		closeErr := r.Close()
		if err == nil {
			err = closeErr
		}
	}()
	return pgxscan.ScanOne(dest, r.Rows)
}

// Close implements [database.Rows].
// Subtle: this method shadows the method (Rows).Close of Rows.Rows.
func (r *Rows) Close() error {
	r.Rows.Close()
	return nil
}
