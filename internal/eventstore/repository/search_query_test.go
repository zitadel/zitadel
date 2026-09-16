package repository

import (
	"reflect"
	"testing"

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
			if got := NewFilter(tt.args.field, tt.args.value, tt.args.operation); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFilter() = %v, want %v", got, tt.want)
			}
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
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Filter.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
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
			if err := tt.fields.columns.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Columns.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueryFromBuilder_eventTypeScans(t *testing.T) {
	base := func() *eventstore.SearchQueryBuilder {
		return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).ScanEventTypesSeparately()
	}
	tests := []struct {
		name    string
		builder *eventstore.SearchQueryBuilder
		want    []EventTypeScan
	}{
		{
			name: "not requested",
			builder: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
				AddQuery().AggregateTypes("user").EventTypes("user.locked").Builder(),
			want: nil,
		},
		{
			name: "every combination of each sub query, sorted and distinct",
			builder: base().
				AddQuery().AggregateTypes("user", "org").EventTypes("user.locked", "org.removed").Builder().
				AddQuery().AggregateTypes("user").EventTypes("user.locked").Builder(),
			want: []EventTypeScan{
				{AggregateType: "org", EventType: "org.removed"},
				{AggregateType: "org", EventType: "user.locked"},
				{AggregateType: "user", EventType: "org.removed"},
				{AggregateType: "user", EventType: "user.locked"},
			},
		},
		{
			name:    "sub query without event types",
			builder: base().AddQuery().AggregateTypes("user").Builder(),
			want:    nil,
		},
		{
			name:    "sub query with aggregate ids",
			builder: base().AddQuery().AggregateTypes("user").AggregateIDs("id").EventTypes("user.locked").Builder(),
			want:    nil,
		},
		{
			name:    "sub query with event data",
			builder: base().AddQuery().AggregateTypes("user").EventTypes("user.locked").EventData(map[string]any{"key": "value"}).Builder(),
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := QueryFromBuilder(tt.builder)
			if err != nil {
				t.Fatalf("QueryFromBuilder() error = %v", err)
			}
			if !reflect.DeepEqual(query.EventTypeScans, tt.want) {
				t.Errorf("EventTypeScans = %v, want %v", query.EventTypeScans, tt.want)
			}
		})
	}
}
