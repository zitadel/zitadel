package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func TestNewFilter(t *testing.T) {
	type args struct {
		field     Field
		value     interface{}
		operation Operation
	}
	tests := []struct {
		name string
		args args
		want *Filter
	}{
		{
			name: "aggregateID equals",
			args: args{
				field:     FieldAggregateID,
				value:     "hodor",
				operation: OperationEquals,
			},
			want: &Filter{Field: FieldAggregateID, Operation: OperationEquals, Value: "hodor"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFilter(tt.args.field, tt.args.value, tt.args.operation)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFilter_Validate(t *testing.T) {
	type fields struct {
		field     Field
		value     interface{}
		operation Operation
		isNil     bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "correct filter",
			fields: fields{
				field:     FieldSequence,
				operation: OperationGreater,
				value:     uint64(235),
			},
			wantErr: false,
		},
		{
			name:    "filter is nil",
			fields:  fields{isNil: true},
			wantErr: true,
		},
		{
			name: "no field error",
			fields: fields{
				operation: OperationGreater,
				value:     uint64(235),
			},
			wantErr: true,
		},
		{
			name: "no value error",
			fields: fields{
				field:     FieldSequence,
				operation: OperationGreater,
			},
			wantErr: true,
		},
		{
			name: "no operation error",
			fields: fields{
				field: FieldSequence,
				value: uint64(235),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f *Filter
			if !tt.fields.isNil {
				f = &Filter{
					Field:     tt.fields.field,
					Value:     tt.fields.value,
					Operation: tt.fields.operation,
				}
			}
			err := f.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestColumns_Validate(t *testing.T) {
	type fields struct {
		columns eventstore.Columns
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "correct filter",
			fields: fields{
				columns: eventstore.ColumnsEvent,
			},
			wantErr: false,
		},
		{
			name: "columns too low",
			fields: fields{
				columns: 0,
			},
			wantErr: true,
		},
		{
			name: "columns too high",
			fields: fields{
				columns: 100,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.columns.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestQueryFromBuilder_ExcludeRelationalEvents(t *testing.T) {
	t.Run("set exclude relational events", func(t *testing.T) {
		builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			ExcludeRelationalEvents()

		query, err := QueryFromBuilder(builder)
		require.NoError(t, err)
		assert.True(t, query.ExcludeRelationalEvents)
	})

	t.Run("exclude relational events defaults to false", func(t *testing.T) {
		builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent)

		query, err := QueryFromBuilder(builder)
		require.NoError(t, err)
		assert.False(t, query.ExcludeRelationalEvents)
	})
}
