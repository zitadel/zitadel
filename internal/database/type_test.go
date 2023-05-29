package database

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMap_Scan(t *testing.T) {
	type args struct {
		src any
	}
	type res[V any] struct {
		want Map[V]
		err  bool
	}
	type testCase[V any] struct {
		name string
		m    Map[V]
		args args
		res[V]
	}
	tests := []testCase[string]{
		{
			"null",
			Map[string]{},
			args{src: "invalid"},
			res[string]{
				want: Map[string]{},
				err:  true,
			},
		},
		{
			"null",
			Map[string]{},
			args{src: nil},
			res[string]{
				want: Map[string]{},
			},
		},
		{
			"empty",
			Map[string]{},
			args{src: []byte(`{}`)},
			res[string]{
				want: Map[string]{},
			},
		},
		{
			"set",
			Map[string]{},
			args{src: []byte(`{"key": "value"}`)},
			res[string]{
				want: Map[string]{
					"key": "value",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Scan(tt.args.src); (err != nil) != tt.res.err {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.res.err)
			}
			assert.Equal(t, tt.res.want, tt.m)
		})
	}
}

func TestMap_Value(t *testing.T) {
	type res struct {
		want driver.Value
		err  bool
	}
	type testCase[V any] struct {
		name string
		m    Map[V]
		res  res
	}
	tests := []testCase[string]{
		{
			"nil",
			nil,
			res{
				want: nil,
			},
		},
		{
			"empty",
			Map[string]{},
			res{
				want: nil,
			},
		},
		{
			"set",
			Map[string]{
				"key": "value",
			},
			res{
				want: driver.Value([]byte(`{"key":"value"}`)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Value()
			if tt.res.err {
				assert.Error(t, err)
			}
			if !tt.res.err {
				require.NoError(t, err)
				assert.Equalf(t, tt.res.want, got, "Value()")
			}
		})
	}
}

func TestArray_ScanInt32(t *testing.T) {
	type args struct {
		src any
	}
	type res[V arrayField] struct {
		want Array[V]
		err  bool
	}
	type testCase[V arrayField] struct {
		name string
		m    Array[V]
		args args
		res[V]
	}
	tests := []testCase[int32]{
		{
			"number",
			Array[int32]{},
			args{src: "{1,2}"},
			res[int32]{
				want: []int32{1, 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Scan(tt.args.src); (err != nil) != tt.res.err {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.res.err)
			}

			assert.Equal(t, tt.res.want, tt.m)
		})
	}
}

func TestArray_Value(t *testing.T) {
	type res struct {
		want driver.Value
		err  bool
	}
	type testCase[V arrayField] struct {
		name string
		a    Array[V]
		res  res
	}
	tests := []testCase[int32]{
		{
			"nil",
			nil,
			res{
				want: nil,
			},
		},
		{
			"empty",
			Array[int32]{},
			res{
				want: nil,
			},
		},
		{
			"set",
			Array[int32]([]int32{1, 2}),
			res{
				want: driver.Value(string([]byte(`{1,2}`))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Value()
			if tt.res.err {
				assert.Error(t, err)
			}
			if !tt.res.err {
				require.NoError(t, err)
				assert.Equalf(t, tt.res.want, got, "Value()")
			}
		})
	}
}
