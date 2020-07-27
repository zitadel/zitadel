package models

import (
	"reflect"
	"testing"
)

func TestSearchQuery_setFilter(t *testing.T) {
	type fields struct {
		query *SearchQuery
	}
	type args struct {
		filters []*Filter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *SearchQuery
	}{
		{
			name:   "set idFilter",
			fields: fields{query: NewSearchQuery()},
			args: args{filters: []*Filter{
				{field: Field_AggregateID, operation: Operation_Equals, value: "hodor"},
			}},
			want: &SearchQuery{Filters: []*Filter{
				{field: Field_AggregateID, operation: Operation_Equals, value: "hodor"},
			}},
		},
		{
			name:   "overwrite idFilter",
			fields: fields{query: NewSearchQuery()},
			args: args{filters: []*Filter{
				{field: Field_AggregateID, operation: Operation_Equals, value: "hodor"},
				{field: Field_AggregateID, operation: Operation_Equals, value: "ursli"},
			}},
			want: &SearchQuery{Filters: []*Filter{
				{field: Field_AggregateID, operation: Operation_Equals, value: "ursli"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.query
			for _, filter := range tt.args.filters {
				got = got.setFilter(filter)
			}
			for _, wantFilter := range tt.want.Filters {
				found := false
				for _, gotFilter := range got.Filters {
					if gotFilter.field == wantFilter.field {
						found = true
						if !reflect.DeepEqual(wantFilter, gotFilter) {
							t.Errorf("filter not as expected: want: %v got %v", wantFilter, gotFilter)
						}
					}
				}
				if !found {
					t.Errorf("filter field %v not found", wantFilter.field)
				}
			}
		})
	}
}
