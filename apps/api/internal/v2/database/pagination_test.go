package database

import (
	"testing"
)

func TestPagination_Write(t *testing.T) {
	type fields struct {
		Limit  uint32
		Offset uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   wantQuery
	}{
		{
			name: "no values",
			fields: fields{
				Limit:  0,
				Offset: 0,
			},
			want: wantQuery{
				query: "",
				args:  []any{},
			},
		},
		{
			name: "limit",
			fields: fields{
				Limit:  10,
				Offset: 0,
			},
			want: wantQuery{
				query: " LIMIT $1",
				args:  []any{uint32(10)},
			},
		},
		{
			name: "offset",
			fields: fields{
				Limit:  0,
				Offset: 10,
			},
			want: wantQuery{
				query: " OFFSET $1",
				args:  []any{uint32(10)},
			},
		},
		{
			name: "both",
			fields: fields{
				Limit:  10,
				Offset: 10,
			},
			want: wantQuery{
				query: " LIMIT $1 OFFSET $2",
				args:  []any{uint32(10), uint32(10)},
			},
		},
	}
	for _, tt := range tests {
		var stmt Statement
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination{
				Limit:  tt.fields.Limit,
				Offset: tt.fields.Offset,
			}
			p.Write(&stmt)
			assertQuery(t, &stmt, tt.want)
		})
	}
}
