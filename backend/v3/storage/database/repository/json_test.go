package repository_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestJSON_Scan(t *testing.T) {
	type want struct {
		err error
		res *testObject
	}
	tests := []struct {
		name string
		src  any
		want want
	}{
		{
			name: "nil",
			src:  nil,
			want: want{
				err: nil,
				res: nil,
			},
		},
		{
			name: "[]byte with data",
			src:  []byte(`{"id":"1","number":1,"active":true}`),
			want: want{
				err: nil,
				res: &testObject{ID: "1", Number: 1, Active: true},
			},
		},
		{
			name: "[]byte without data",
			src:  []byte(`{}`),
			want: want{
				err: nil,
				res: &testObject{},
			},
		},
		{
			name: "nil []byte",
			src:  []byte(nil),
			want: want{
				err: nil,
				res: nil,
			},
		},
		{
			name: "string with data",
			src:  `{"id":"2","number":2,"active":false}`,
			want: want{
				err: nil,
				res: &testObject{ID: "2", Number: 2, Active: false},
			},
		},
		{
			name: "string without data",
			src:  `{}`,
			want: want{
				err: nil,
				res: &testObject{},
			},
		},
		{
			name: "empty string",
			src:  ``,
			want: want{
				err: nil,
				res: nil,
			},
		},
		{
			name: "unsupported type",
			src:  12345,
			want: want{
				err: repository.ErrScanSource,
				res: nil,
			},
		},
		{
			name: "invalid json",
			src:  []byte(`this is not json`),
			want: want{
				err: new(database.ScanError),
				res: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a repository.JSON[testObject]
			gotErr := a.Scan(tt.src)
			require.ErrorIs(t, gotErr, tt.want.err)
		})
	}
}
