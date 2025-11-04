package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func TestFieldChange(t *testing.T) {
	for _, test := range []struct {
		name   string
		change database.Change
		output string
	}{
		{
			name: "single json update",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change := NewFieldChange([]string{"path", "to", "key"}, "value")
				changes := NewJsonChanges(col, change)

				return changes
			}(),
			output: `column = jsonb_set_lax(table.column, '{pathtokey}', $1, true, 'delete_key')`,
		},
		{
			name: "two json update",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change1 := NewFieldChange([]string{"path1"}, "value1")
				change2 := NewFieldChange([]string{"path2"}, "value2")
				changes := NewJsonChanges(col, change1, change2)

				return changes
			}(),
			output: `column = jsonb_set_lax(jsonb_set_lax(table.column, '{path1}', $1, true, 'delete_key'), '{path2}', $2, true, 'delete_key')`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			builder := database.StatementBuilder{}
			err := test.change.Write(&builder)
			require.NoError(t, err)

			assert.Equal(t, test.output, builder.String())
		})
	}
}

func TestArrayChange(t *testing.T) {
	for _, test := range []struct {
		name   string
		change database.Change
		output string
	}{
		{
			name: "one json array add",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change := NewArrayChange([]string{"path"}, "value", false)
				changes := NewJsonChanges(col, change)

				return changes
			}(),
			output: `column = zitadel.jsonb_array_append(zitadel.jsonb_array_remove(table.column, $1, $2::TEXT), $3, $4::TEXT)`,
		},
		{
			name: "two json array add",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change1 := NewArrayChange([]string{"path1)"}, "value1", false)
				change2 := NewArrayChange([]string{"path2"}, "value2", false)
				changes := NewJsonChanges(col, change1, change2)

				return changes
			}(),
			output: `column = zitadel.jsonb_array_append(zitadel.jsonb_array_remove(zitadel.jsonb_array_append(zitadel.jsonb_array_remove(table.column, $1, $2::TEXT), $3, $4::TEXT), $5, $6::TEXT), $7, $8::TEXT)`,
		},
		{
			name: "one json array remove",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change := NewArrayChange([]string{"path"}, "value", true)
				changes := NewJsonChanges(col, change)

				return changes
			}(),
			output: `column = zitadel.jsonb_array_remove(table.column, $1, $2::TEXT)`,
		},
		{
			name: "two json array remove",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change1 := NewArrayChange([]string{"path1)"}, "value1", true)
				change2 := NewArrayChange([]string{"path2"}, "value2", true)
				changes := NewJsonChanges(col, change1, change2)

				return changes
			}(),
			output: `column = zitadel.jsonb_array_remove(zitadel.jsonb_array_remove(table.column, $1, $2::TEXT), $3, $4::TEXT)`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			builder := database.StatementBuilder{}
			err := test.change.Write(&builder)
			require.NoError(t, err)

			assert.Equal(t, test.output, builder.String())
		})
	}
}

func TestArrayMixedChange(t *testing.T) {
	for _, test := range []struct {
		name   string
		change database.Change
		output string
	}{
		{
			name: "one json array add, one json remove",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change1 := NewArrayChange([]string{"path1)"}, "value1", false)
				change2 := NewArrayChange([]string{"path2"}, "value2", true)
				changes := NewJsonChanges(col, change1, change2)

				return changes
			}(),
			output: `column = zitadel.jsonb_array_remove(zitadel.jsonb_array_append(zitadel.jsonb_array_remove(table.column, $1, $2::TEXT), $3, $4::TEXT), $5, $6::TEXT)`,
		},
		{
			name: "one json array remove, one json array add",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change1 := NewArrayChange([]string{"path1)"}, "value1", true)
				change2 := NewArrayChange([]string{"path2"}, "value2", false)
				changes := NewJsonChanges(col, change1, change2)

				return changes
			}(),
			output: `column = zitadel.jsonb_array_append(zitadel.jsonb_array_remove(zitadel.jsonb_array_remove(table.column, $1, $2::TEXT), $3, $4::TEXT), $5, $6::TEXT)`,
		},
		{
			name: "one json array add, one json array remove, one array add",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change1 := NewArrayChange([]string{"path1)"}, "value1", false)
				change2 := NewArrayChange([]string{"path2"}, "value2", true)
				change3 := NewArrayChange([]string{"path3)"}, "value3", false)
				changes := NewJsonChanges(col, change1, change2, change3)

				return changes
			}(),
			output: `column = zitadel.jsonb_array_append(zitadel.jsonb_array_remove(zitadel.jsonb_array_remove(zitadel.jsonb_array_append(zitadel.jsonb_array_remove(table.column, $1, $2::TEXT), $3, $4::TEXT), $5, $6::TEXT), $7, $8::TEXT), $9, $10::TEXT)`,
		},
		{
			name: "one json array remove, one json array add, one array remove",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change1 := NewArrayChange([]string{"path1)"}, "value1", true)
				change2 := NewArrayChange([]string{"path2"}, "value2", false)
				change3 := NewArrayChange([]string{"path3)"}, "value3", true)
				changes := NewJsonChanges(col, change1, change2, change3)

				return changes
			}(),
			output: `column = zitadel.jsonb_array_remove(zitadel.jsonb_array_append(zitadel.jsonb_array_remove(zitadel.jsonb_array_remove(table.column, $1, $2::TEXT), $3, $4::TEXT), $5, $6::TEXT), $7, $8::TEXT)`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			builder := database.StatementBuilder{}
			err := test.change.Write(&builder)
			require.NoError(t, err)

			assert.Equal(t, test.output, builder.String())
		})
	}
}

// TODO
func TestFieldArrayMixedChange(t *testing.T) {
	for _, test := range []struct {
		name   string
		change database.Change
		output string
	}{
		{
			name: "one field change, one json array add, one json remove",
			change: func() database.Change {
				col := database.NewColumn("table", "column")
				change1 := NewFieldChange([]string{"path", "to", "key"}, "value")
				change2 := NewArrayChange([]string{"path1)"}, "value1", false)
				change3 := NewArrayChange([]string{"path3"}, "value3", true)
				changes := NewJsonChanges(col, change1, change2, change3)

				return changes
			}(),
			output: `column = zitadel.jsonb_array_remove(zitadel.jsonb_array_append(zitadel.jsonb_array_remove(table.column, $1, $2::TEXT), $3, $4::TEXT), $5, $6::TEXT)`,
		},
		// {
		// 	name: "one json array remove, one json array add",
		// 	change: func() database.Change {
		// 		col := database.NewColumn("table", "column")
		// 		change1 := NewArrayChange([]string{"path1)"}, "value1", true)
		// 		change2 := NewArrayChange([]string{"path2"}, "value2", false)
		// 		changes := NewJsonChanges(col, change1, change2)

		// 		return changes
		// 	}(),
		// 	output: `column = zitadel.jsonb_array_append(zitadel.jsonb_array_remove(zitadel.jsonb_array_remove(table.column, $1, $2::TEXT), $3, $4::TEXT), $5, $6::TEXT)`,
		// },
		// {
		// 	name: "one json array add, one json array remove, one array add",
		// 	change: func() database.Change {
		// 		col := database.NewColumn("table", "column")
		// 		change1 := NewArrayChange([]string{"path1)"}, "value1", false)
		// 		change2 := NewArrayChange([]string{"path2"}, "value2", true)
		// 		change3 := NewArrayChange([]string{"path3)"}, "value3", false)
		// 		changes := NewJsonChanges(col, change1, change2, change3)

		// 		return changes
		// 	}(),
		// 	output: `column = zitadel.jsonb_array_append(zitadel.jsonb_array_remove(zitadel.jsonb_array_remove(zitadel.jsonb_array_append(zitadel.jsonb_array_remove(table.column, $1, $2::TEXT), $3, $4::TEXT), $5, $6::TEXT), $7, $8::TEXT), $9, $10::TEXT)`,
		// },
		// {
		// 	name: "one json array remove, one json array add, one array remove",
		// 	change: func() database.Change {
		// 		col := database.NewColumn("table", "column")
		// 		change1 := NewArrayChange([]string{"path1)"}, "value1", true)
		// 		change2 := NewArrayChange([]string{"path2"}, "value2", false)
		// 		change3 := NewArrayChange([]string{"path3)"}, "value3", true)
		// 		changes := NewJsonChanges(col, change1, change2, change3)

		// 		return changes
		// 	}(),
		// 	output: `column = zitadel.jsonb_array_remove(zitadel.jsonb_array_append(zitadel.jsonb_array_remove(zitadel.jsonb_array_remove(table.column, $1, $2::TEXT), $3, $4::TEXT), $5, $6::TEXT), $7, $8::TEXT)`,
		// },
	} {
		t.Run(test.name, func(t *testing.T) {
			builder := database.StatementBuilder{}
			err := test.change.Write(&builder)
			require.NoError(t, err)

			assert.Equal(t, test.output, builder.String())
		})
	}
}
