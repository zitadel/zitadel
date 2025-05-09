package mock

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/storage/database"
)

type Row struct {
	t *testing.T

	res []any
}

func NewRow(t *testing.T, res ...any) *Row {
	return &Row{t: t, res: res}
}

// Scan implements [database.Row].
func (r *Row) Scan(dest ...any) error {
	require.Len(r.t, dest, len(r.res))
	for i := range dest {
		reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(r.res[i]))
	}
	return nil
}

var _ database.Row = (*Row)(nil)
