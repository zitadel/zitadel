package json

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func TestFieldChange(t *testing.T) {
	for _, test := range []struct {
		name   string
		change database.Change
		want   string
	}{
		{
			name: "single josn update",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change := NewFieldChange("path", "value")
				changes := NewJsonChanges(col, change)

				return changes
			}(),
			want: "missing condition for column",
		},
		{
			name: "two josn update",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change1 := NewFieldChange("path", "value")
				change2 := NewFieldChange("path", "value")
				changes := NewJsonChanges(col, change1, change2)

				return changes
			}(),
			want: "missing condition for column",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			builder := database.StatementBuilder{}
			err := test.change.Write(&builder)
			require.NoError(t, err)
			fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> builder.String() = %+v\n", builder.String())

			// assert.Equal(t, test.want, test.err.Error())
		})
	}
}

func TestArrayChange(t *testing.T) {
	for _, test := range []struct {
		name   string
		change database.Change
		want   string
	}{
		{
			name: "single josn array",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change := NewArrayChange("path", "value", false)
				changes := NewJsonChanges(col, change)

				return changes
			}(),
			want: "missing condition for column",
		},
		{
			name: "two josn array",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change1 := NewArrayChange("path1", "value", false)
				change2 := NewArrayChange("path2", 33, false)
				changes := NewJsonChanges(col, change1, change2)

				return changes
			}(),
			want: "missing condition for column",
		},
		{
			name: "remove josn array",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change1 := NewArrayChange("path1", "value", true)
				// change2 := NewArrayChange("path2", 33, false)
				changes := NewJsonChanges(col, change1)

				return changes
			}(),
			want: "missing condition for column",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			builder := database.StatementBuilder{}
			err := test.change.Write(&builder)
			require.NoError(t, err)
			fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> builder.String() = %+v\n", builder.String())

			// assert.Equal(t, test.want, test.err.Error())
		})
	}
}
