package database

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArray_Scan_string(t *testing.T) {
	const src = "{foo,bar}"
	want := Array[string]{"foo", "bar"}

	var got Array[string]
	err := got.Scan(src)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestArray_Value_string(t *testing.T) {
	src := Array[string]{"foo", "bar"}
	const want = "{foo,bar}"

	got, err := src.Value()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

type myEnum int32

const (
	myEnumZero myEnum = iota
	myEnumOne
	myEnumTwo
)

func TestEnumArray_Scan(t *testing.T) {
	const src = "{0,1,2}"
	want := EnumArray[myEnum]{myEnumZero, myEnumOne, myEnumTwo}

	var got EnumArray[myEnum]
	err := got.Scan(src)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestEnumArray_Value(t *testing.T) {
	src := EnumArray[myEnum]{myEnumZero, myEnumOne, myEnumTwo}
	const want = "{0,1,2}"

	got, err := src.Value()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

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
