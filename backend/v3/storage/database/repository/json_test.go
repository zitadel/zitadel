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

//
//func TestJSON_MarshalJSON(t *testing.T) {
//	type testCase[T any] struct {
//		name    string
//		j       JSON[T]
//		want    []byte
//		wantErr assert.ErrorAssertionFunc
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := tt.j.MarshalJSON()
//			if !tt.wantErr(t, err, fmt.Sprintf("MarshalJSON()")) {
//				return
//			}
//			assert.Equalf(t, tt.want, got, "MarshalJSON()")
//		})
//	}
//}
//
//func TestJSON_UnmarshalJSON(t *testing.T) {
//	type args struct {
//		data []byte
//	}
//	type testCase[T any] struct {
//		name    string
//		j       JSON[T]
//		args    args
//		wantErr assert.ErrorAssertionFunc
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.wantErr(t, tt.j.UnmarshalJSON(tt.args.data), fmt.Sprintf("UnmarshalJSON(%v)", tt.args.data))
//		})
//	}
//}
