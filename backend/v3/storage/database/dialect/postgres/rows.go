package postgres

import (
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var (
	_ database.Rows            = (*Rows)(nil)
	_ database.CollectableRows = (*Rows)(nil)
	_ database.Row             = (*Row)(nil)
)

type Row struct{ pgx.Row }

// Err implements [database.Row].
// [pgx.Row] does not provide an error so we always return nil.
func (r *Row) Err() error {
	return nil
}

// Scan implements [database.Row].
// Subtle: this method shadows the method ([pgx.Row]).Scan of Row.Row.
func (r *Row) Scan(dest ...any) error {
	return wrapError(r.Row.Scan(dest...))
}

type Rows struct{ pgx.Rows }

// Err implements [database.Rows].
// Subtle: this method shadows the method ([pgx.Rows]).Err of Rows.Rows.
func (r *Rows) Err() error {
	return wrapError(r.Rows.Err())
}

func (r *Rows) Scan(dest ...any) error {
	return wrapError(r.Rows.Scan(dest...))
}

// Collect implements [database.CollectableRows].
// See [this page](https://github.com/georgysavva/scany/blob/master/dbscan/doc.go#L8) for additional details.
func (r *Rows) Collect(dest any) (err error) {
	defer func() {
		closeErr := r.Close()
		if err == nil {
			err = closeErr
		}
	}()
	return wrapError(pgxscan.ScanAll(dest, r.Rows))
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
	return wrapError(pgxscan.ScanRow(dest, r.Rows))
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
	return wrapError(pgxscan.ScanOne(dest, r.Rows))
}

// Close implements [database.Rows].
// Subtle: this method shadows the method (Rows).Close of Rows.Rows.
func (r *Rows) Close() error {
	r.Rows.Close()
	return nil
}
