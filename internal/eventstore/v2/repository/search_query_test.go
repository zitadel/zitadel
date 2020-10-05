package repository

import (
	"reflect"
	"testing"
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
				field:     Field_AggregateID,
				value:     "hodor",
				operation: Operation_Equals,
			},
			want: &Filter{Field: Field_AggregateID, Operation: Operation_Equals, Value: "hodor"},
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
				field:     Field_LatestSequence,
				operation: Operation_Greater,
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
				operation: Operation_Greater,
				value:     uint64(235),
			},
			wantErr: true,
		},
		{
			name: "no value error",
			fields: fields{
				field:     Field_LatestSequence,
				operation: Operation_Greater,
			},
			wantErr: true,
		},
		{
			name: "no operation error",
			fields: fields{
				field: Field_LatestSequence,
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
