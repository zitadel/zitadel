package noopdb

import "github.com/zitadel/zitadel/backend/v3/storage/database"

type rows struct{}

// Close implements [database.Rows].
func (r *rows) Close() error {
	return nil
}

// Err implements [database.Rows].
func (r *rows) Err() error {
	return nil
}

// Next implements [database.Rows].
func (r *rows) Next() bool {
	return false
}

// Scan implements [database.Rows].
func (r *rows) Scan(dest ...any) error {
	return nil
}

var _ database.Rows = (*rows)(nil)

type row struct{}

// Scan implements database.Row.
func (r *row) Scan(dest ...any) error {
	return nil
}

var _ database.Row = (*row)(nil)
