package eventstore

import (
	"context"
	_ "embed"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

func Test_handleSearchDelete(t *testing.T) {
	type args struct {
		clauses map[eventstore.FieldType]any
	}
	type want struct {
		stmt string
		args []any
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "1 condition",
			args: args{
				clauses: map[eventstore.FieldType]any{
					eventstore.FieldTypeInstanceID: "i_id",
				},
			},
			want: want{
				stmt: "DELETE FROM eventstore.fields WHERE instance_id = $1",
				args: []any{"i_id"},
			},
		},
		{
			name: "2 conditions",
			args: args{
				clauses: map[eventstore.FieldType]any{
					eventstore.FieldTypeInstanceID:  "i_id",
					eventstore.FieldTypeAggregateID: "a_id",
				},
			},
			want: want{
				stmt: "DELETE FROM eventstore.fields WHERE aggregate_id = $1 AND instance_id = $2",
				args: []any{"a_id", "i_id"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, args := writeDeleteField(tt.args.clauses)
			if stmt != tt.want.stmt {
				t.Errorf("handleSearchDelete() stmt = %q, want %q", stmt, tt.want.stmt)
			}
			assert.Equal(t, tt.want.args, args)
		})
	}
}

func Test_writeUpsertField(t *testing.T) {
	type args struct {
		fields []eventstore.FieldType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1 field",
			args: args{
				fields: []eventstore.FieldType{
					eventstore.FieldTypeInstanceID,
				},
			},
			want: "WITH upsert AS (UPDATE eventstore.fields SET (instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value, value_must_be_unique, should_index) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) WHERE instance_id = $1 RETURNING * ) INSERT INTO eventstore.fields (instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value, value_must_be_unique, should_index) SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11 WHERE NOT EXISTS (SELECT 1 FROM upsert)",
		},
		{
			name: "2 fields",
			args: args{
				fields: []eventstore.FieldType{
					eventstore.FieldTypeInstanceID,
					eventstore.FieldTypeAggregateType,
				},
			},
			want: "WITH upsert AS (UPDATE eventstore.fields SET (instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value, value_must_be_unique, should_index) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) WHERE instance_id = $1 AND aggregate_type = $3 RETURNING * ) INSERT INTO eventstore.fields (instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value, value_must_be_unique, should_index) SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11 WHERE NOT EXISTS (SELECT 1 FROM upsert)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := writeUpsertField(tt.args.fields); got != tt.want {
				t.Errorf("writeUpsertField() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_buildSearchCondition(t *testing.T) {
	type args struct {
		index      int
		conditions map[eventstore.FieldType]any
	}
	type want struct {
		stmt string
		args []any
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "1 condition",
			args: args{
				index: 1,
				conditions: map[eventstore.FieldType]any{
					eventstore.FieldTypeAggregateID: "a_id",
				},
			},
			want: want{
				stmt: "aggregate_id = $1",
				args: []any{"a_id"},
			},
		},
		{
			name: "3 condition",
			args: args{
				index: 1,
				conditions: map[eventstore.FieldType]any{
					eventstore.FieldTypeAggregateID:   "a_id",
					eventstore.FieldTypeInstanceID:    "i_id",
					eventstore.FieldTypeAggregateType: "a_type",
				},
			},
			want: want{
				stmt: "aggregate_type = $1 AND aggregate_id = $2 AND instance_id = $3",
				args: []any{"a_type", "a_id", "i_id"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder

			if got := buildSearchCondition(&builder, tt.args.index, tt.args.conditions); !reflect.DeepEqual(got, tt.want.args) {
				t.Errorf("buildSearchCondition() = %v, want %v", got, tt.want)
			}
			if tt.want.stmt != builder.String() {
				t.Errorf("buildSearchCondition() stmt = %q, want %q", builder.String(), tt.want.stmt)
			}
		})
	}
}

func Test_buildSearchStatement(t *testing.T) {
	type args struct {
		index      int
		conditions []map[eventstore.FieldType]any
	}
	type want struct {
		stmt string
		args []any
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "1 condition with 1 field",
			args: args{
				index: 1,
				conditions: []map[eventstore.FieldType]any{
					{
						eventstore.FieldTypeAggregateID: "a_id",
					},
				},
			},
			want: want{
				stmt: "SELECT instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value FROM eventstore.fields WHERE instance_id = $1 AND aggregate_id = $2",
				args: []any{"a_id"},
			},
		},
		{
			name: "1 condition with 3 fields",
			args: args{
				index: 1,
				conditions: []map[eventstore.FieldType]any{
					{
						eventstore.FieldTypeAggregateID:   "a_id",
						eventstore.FieldTypeInstanceID:    "i_id",
						eventstore.FieldTypeAggregateType: "a_type",
					},
				},
			},
			want: want{
				stmt: "SELECT instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value FROM eventstore.fields WHERE instance_id = $1 AND (aggregate_type = $2 AND aggregate_id = $3 AND instance_id = $4)",
				args: []any{"a_type", "a_id", "i_id"},
			},
		},
		{
			name: "2 condition with 1 field",
			args: args{
				index: 1,
				conditions: []map[eventstore.FieldType]any{
					{
						eventstore.FieldTypeAggregateID: "a_id",
					},
					{
						eventstore.FieldTypeAggregateType: "a_type",
					},
				},
			},
			want: want{
				stmt: "SELECT instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value FROM eventstore.fields WHERE instance_id = $1 AND (aggregate_id = $2 OR aggregate_type = $3)",
				args: []any{"a_id", "a_type"},
			},
		},
		{
			name: "2 condition with 2 fields",
			args: args{
				index: 1,
				conditions: []map[eventstore.FieldType]any{
					{
						eventstore.FieldTypeAggregateID:   "a_id1",
						eventstore.FieldTypeAggregateType: "a_type1",
					},
					{
						eventstore.FieldTypeAggregateID:   "a_id2",
						eventstore.FieldTypeAggregateType: "a_type2",
					},
				},
			},
			want: want{
				stmt: "SELECT instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, object_revision, field_name, value FROM eventstore.fields WHERE instance_id = $1 AND ((aggregate_type = $2 AND aggregate_id = $3) OR (aggregate_type = $4 AND aggregate_id = $5))",
				args: []any{"a_type1", "a_id1", "a_type2", "a_id2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			tt.want.args = append([]any{"i_id"}, tt.want.args...)
			ctx := authz.WithInstanceID(context.Background(), "i_id")

			if got := buildSearchStatement(ctx, &builder, tt.args.conditions...); !reflect.DeepEqual(got, tt.want.args) {
				t.Errorf("buildSearchStatement() = %v, want %v", got, tt.want)
			}
			if tt.want.stmt != builder.String() {
				t.Errorf("buildSearchStatement() stmt = %q, want %q", builder.String(), tt.want.stmt)
			}
		})
	}
}
