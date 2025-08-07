package info

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ActivityInfo(t *testing.T) {
	type args struct {
		ctx           context.Context
		ok            bool
		path          string
		method        string
		requestMethod string
	}
	type want struct {
		ok            bool
		path          string
		method        string
		requestMethod string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"already set",
			args{
				ctx: ctxWithActivityInfo(context.Background(), "set", "set", "set"),
				ok:  false,
			},
			want{
				ok:            true,
				path:          "set",
				method:        "set",
				requestMethod: "set",
			},
		},
		{
			"not set, empty",
			args{
				ctx: context.Background(),
				ok:  false,
			},
			want{
				ok: true,
			},
		},
		{
			"set empty",
			args{
				ctx: context.Background(),
				ok:  true,
			},
			want{
				ok: true,
			},
		},
		{
			"set",
			args{
				ctx:           context.Background(),
				ok:            true,
				path:          "set",
				method:        "set",
				requestMethod: "set",
			},
			want{
				ok:            true,
				path:          "set",
				method:        "set",
				requestMethod: "set",
			},
		},
		{
			"reset",
			args{
				ctx:           ctxWithActivityInfo(context.Background(), "set", "set", "set"),
				ok:            true,
				path:          "set2",
				method:        "set2",
				requestMethod: "set2",
			},
			want{
				ok:            true,
				path:          "set2",
				method:        "set2",
				requestMethod: "set2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ai := &ActivityInfo{}
			ai.SetMethod(tt.args.method).SetPath(tt.args.path).SetRequestMethod(tt.args.requestMethod)
			if tt.args.ok {
				tt.args.ctx = ai.IntoContext(tt.args.ctx)
			}

			res := ActivityInfoFromContext(tt.args.ctx)
			if tt.want.ok {
				assert.NotNil(t, res)
			}
			assert.Equal(t, tt.want.path, res.Path)
			assert.Equal(t, tt.want.method, res.Method)
			assert.Equal(t, tt.want.requestMethod, res.RequestMethod)
		})
	}
}

func ctxWithActivityInfo(ctx context.Context, method, path, requestMethod string) context.Context {
	ai := &ActivityInfo{}
	return ai.SetPath(path).SetRequestMethod(requestMethod).SetMethod(method).IntoContext(ctx)
}
