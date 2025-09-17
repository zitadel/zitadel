package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsColumn(t *testing.T) {
	col := NewColumn("orgs", "instance_id")
	tests := []struct {
		name      string
		condition Condition
		want      bool
	}{
		{
			name:      "value condition with instance",
			condition: NewTextCondition(col, TextOperationEqual, "id"),
			want:      true,
		},
		{
			name:      "value condition without instance",
			condition: NewTextCondition(NewColumn("orgs", "name"), TextOperationEqual, "name"),
			want:      false,
		},
		{
			name:      "column condition with instance",
			condition: NewColumnCondition(col, NewColumn("orgs", "other")),
			want:      true,
		},
		{
			name:      "column condition without instance",
			condition: NewColumnCondition(NewColumn("orgs", "id"), NewColumn("orgs", "other")),
			want:      false,
		},
		{
			name: "and with instance",
			condition: And(
				NewTextCondition(col, TextOperationEqual, "id"),
				NewTextCondition(NewColumn("orgs", "name"), TextOperationEqual, "name"),
			),
			want: true,
		},
		{
			name: "and without instance",
			condition: And(
				NewTextCondition(NewColumn("orgs", "id"), TextOperationEqual, "id"),
				NewTextCondition(NewColumn("orgs", "name"), TextOperationEqual, "name"),
			),
			want: false,
		},
		{
			name: "or with partial instance",
			condition: Or(
				NewTextCondition(col, TextOperationEqual, "id"),
				NewTextCondition(NewColumn("orgs", "name"), TextOperationEqual, "name"),
			),
			want: false,
		},
		{
			name: "or with only instance",
			condition: Or(
				NewTextCondition(col, TextOperationEqual, "id"),
				NewTextCondition(col, TextOperationEqual, "id2"),
			),
			want: true,
		},
		{
			name: "or without instance",
			condition: Or(
				NewTextCondition(NewColumn("orgs", "id"), TextOperationEqual, "id"),
				NewTextCondition(NewColumn("orgs", "name"), TextOperationEqual, "name"),
			),
			want: false,
		},
		{
			name: "nested and/or with instance",
			condition: And(
				Or(
					NewTextCondition(NewColumn("orgs", "id"), TextOperationEqual, "id"),
					NewTextCondition(NewColumn("orgs", "name"), TextOperationEqual, "name"),
				),
				NewTextCondition(col, TextOperationEqual, "id"),
			),
			want: true,
		},
		{
			name: "nested and/or with in or instance",
			condition: And(
				Or(
					NewTextCondition(col, TextOperationEqual, "id"),
					NewTextCondition(col, TextOperationEqual, "id2"),
				),
				NewTextCondition(NewColumn("orgs", "name"), TextOperationEqual, "name"),
			),
			want: true,
		},
		{
			name: "nested and/or without instance",
			condition: And(
				Or(
					NewTextCondition(NewColumn("orgs", "id"), TextOperationEqual, "id"),
					NewTextCondition(NewColumn("orgs", "name"), TextOperationEqual, "name"),
				),
				NewTextCondition(NewColumn("orgs", "other"), TextOperationEqual, "id"),
			),
			want: false,
		},
		{
			name:      "is null with instance",
			condition: IsNull(col),
			want:      true,
		},
		{
			name:      "is null without instance",
			condition: IsNull(NewColumn("orgs", "name")),
			want:      false,
		},
		{
			name:      "is not null with instance",
			condition: IsNotNull(col),
			want:      true,
		},
		{
			name:      "is not null without instance",
			condition: IsNotNull(NewColumn("orgs", "name")),
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.condition.ContainsColumn(col)
			assert.Equal(t, tt.want, got)
		})
	}
}
