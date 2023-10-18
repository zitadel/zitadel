package database

import (
	"database/sql/driver"
	"testing"
	"time"

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

func TestNullDuration_Scan(t *testing.T) {
	type args struct {
		src any
	}
	type res struct {
		want NullDuration
		err  bool
	}
	type testCase struct {
		name string
		args args
		res  res
	}
	tests := []testCase{
		{
			"invalid",
			args{src: "invalid"},
			res{
				want: NullDuration{
					Valid: false,
				},
				err: true,
			},
		},
		{
			"null",
			args{src: nil},
			res{
				want: NullDuration{
					Valid: false,
				},
				err: false,
			},
		},
		{
			"valid",
			args{src: "1:0:0"},
			res{
				want: NullDuration{
					Valid:    true,
					Duration: time.Hour,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := new(NullDuration)
			if err := d.Scan(tt.args.src); (err != nil) != tt.res.err {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.res.err)
			}
			assert.Equal(t, tt.res.want, *d)
		})
	}
}
