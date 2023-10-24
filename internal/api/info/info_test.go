package info

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HTTPPath(t *testing.T) {
	type args struct {
		path string
		ok   bool
	}
	type want struct {
		path string
		ok   bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"not set",
			args{
				ok: false,
			},
			want{
				ok: false,
			},
		},
		{
			"set empty",
			args{
				ok:   true,
				path: "",
			},
			want{
				ok:   true,
				path: "",
			},
		},
		{
			"set",
			args{
				ok:   true,
				path: "set",
			},
			want{
				ok:   true,
				path: "set",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.args.ok {
				ctx = HTTPPathIntoContext(tt.args.path)(ctx)
			}
			path, ok := HTTPPathFromContext()(ctx)
			assert.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.path, path)
		})
	}
}

func Test_RPCMethod(t *testing.T) {
	type args struct {
		method string
		ok     bool
	}
	type want struct {
		method string
		ok     bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"not set",
			args{
				ok: false,
			},
			want{
				ok: false,
			},
		},
		{
			"set empty",
			args{
				ok:     true,
				method: "",
			},
			want{
				ok:     true,
				method: "",
			},
		},
		{
			"set",
			args{
				ok:     true,
				method: "set",
			},
			want{
				ok:     true,
				method: "set",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.args.ok {
				ctx = RPCMethodIntoContext(tt.args.method)(ctx)
			}
			method, ok := RPCMethodFromContext()(ctx)
			assert.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.method, method)
		})
	}
}

func Test_RequestMethod(t *testing.T) {
	type args struct {
		method string
		ok     bool
	}
	type want struct {
		method string
		ok     bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"not set",
			args{
				ok: false,
			},
			want{
				ok: false,
			},
		},
		{
			"set empty",
			args{
				ok:     true,
				method: "",
			},
			want{
				ok:     true,
				method: "",
			},
		},
		{
			"set",
			args{
				ok:     true,
				method: "set",
			},
			want{
				ok:     true,
				method: "set",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.args.ok {
				ctx = RequestMethodIntoContext(tt.args.method)(ctx)
			}
			method, ok := RequestMethodFromContext()(ctx)
			assert.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.method, method)
		})
	}
}
