package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteUnqualified(t *testing.T) {
	for _, tests := range []struct {
		name     string
		column   Column
		expected string
	}{
		{
			name:     "column",
			column:   NewColumn("table", "column"),
			expected: "column",
		},
		{
			name: "columns",
			column: Columns{
				NewColumn("table", "column1"),
				NewColumn("table", "column2"),
			},
			expected: "column1, column2",
		},
		{
			name:     "function column",
			column:   SHA256Column(NewColumn("table", "column")),
			expected: "SHA256(column)",
		},
	} {
		t.Run(tests.name, func(t *testing.T) {
			var builder StatementBuilder
			tests.column.WriteUnqualified(&builder)
			assert.Equal(t, tests.expected, builder.String())
		})
	}
}

func TestWriteQualified(t *testing.T) {
	for _, tests := range []struct {
		name     string
		column   Column
		expected string
	}{
		{
			name:     "column",
			column:   NewColumn("table", "column"),
			expected: "table.column",
		},
		{
			name: "columns",
			column: Columns{
				NewColumn("table", "column1"),
				NewColumn("table", "column2"),
			},
			expected: "table.column1, table.column2",
		},
		{
			name:     "function column",
			column:   SHA256Column(NewColumn("table", "column")),
			expected: "SHA256(table.column)",
		},
	} {
		t.Run(tests.name, func(t *testing.T) {
			var builder StatementBuilder
			tests.column.WriteQualified(&builder)
			assert.Equal(t, tests.expected, builder.String())
		})
	}
}

func TestEquals(t *testing.T) {
	for _, tests := range []struct {
		name     string
		column   Column
		toCheck  Column
		expected bool
	}{
		{
			name:     "column equal",
			column:   NewColumn("table", "column"),
			toCheck:  NewColumn("table", "column"),
			expected: true,
		},
		{
			name:     "column nil check",
			column:   NewColumn("table", "column"),
			toCheck:  nil,
			expected: false,
		},
		{
			name:     "column both nil",
			column:   (*column)(nil),
			toCheck:  nil,
			expected: true,
		},
		{
			name:     "column not equal (different name)",
			column:   NewColumn("table", "column"),
			toCheck:  NewColumn("table", "column2"),
			expected: false,
		},
		{
			name:     "column not equal (different type)",
			column:   NewColumn("table", "column"),
			toCheck:  SHA256Column(NewColumn("table", "column")),
			expected: false,
		},
		{
			name: "columns equal",
			column: Columns{
				NewColumn("table", "column1"),
				NewColumn("table", "column2"),
			},
			toCheck: Columns{
				NewColumn("table", "column1"),
				NewColumn("table", "column2"),
			},
			expected: true,
		},
		{
			name: "columns nil check",
			column: Columns{
				NewColumn("table", "column1"),
				NewColumn("table", "column2"),
			},
			toCheck:  nil,
			expected: false,
		},
		{
			name:     "columns both nil",
			column:   Columns(nil),
			toCheck:  nil,
			expected: true,
		},
		{
			name: "columns not equal (different type)",
			column: Columns{
				NewColumn("table", "column1"),
				NewColumn("table", "column2"),
			},
			toCheck:  NewColumn("table", "column1"),
			expected: false,
		},
		{
			name: "columns not equal (different length)",
			column: Columns{
				NewColumn("table", "column1"),
				NewColumn("table", "column2"),
			},
			toCheck: Columns{
				NewColumn("table", "column1"),
			},
			expected: false,
		},
		{
			name: "columns not equal (different order)",
			column: Columns{
				NewColumn("table", "column1"),
				NewColumn("table", "column2"),
			},
			toCheck: Columns{
				NewColumn("table", "column2"),
				NewColumn("table", "column1"),
			},
			expected: false,
		},
		{
			name:     "function column equal",
			column:   SHA256Column(NewColumn("table", "column")),
			toCheck:  SHA256Column(NewColumn("table", "column")),
			expected: true,
		},
		{
			name:     "function nil check",
			column:   SHA256Column(NewColumn("table", "column")),
			toCheck:  nil,
			expected: false,
		},
		{
			name:     "function both nil",
			column:   (*functionColumn)(nil),
			toCheck:  nil,
			expected: true,
		},
		{
			name:     "function column not equal (different function)",
			column:   SHA256Column(NewColumn("table", "column")),
			toCheck:  LowerColumn(NewColumn("table", "column")),
			expected: false,
		},
		{
			name:     "function column not equal (different inner column)",
			column:   SHA256Column(NewColumn("table", "column")),
			toCheck:  SHA256Column(NewColumn("table", "column2")),
			expected: false,
		},
		{
			name:     "function column not equal (different type)",
			column:   SHA256Column(NewColumn("table", "column")),
			toCheck:  NewColumn("table", "column2"),
			expected: false,
		},
	} {
		t.Run(tests.name, func(t *testing.T) {
			assert.Equal(t, tests.expected, tests.column.Equals(tests.toCheck))
		})
	}
}
