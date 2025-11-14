package connect_middleware

import (
	"context"
	"reflect"
	"testing"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/authz"
)

func Test_toGRPCError(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     connect.AnyRequest
		handler func(t *testing.T) connect.UnaryFunc
	}
	type res struct {
		want    interface{}
		wantErr bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"no error",
			args{
				ctx:     context.Background(),
				req:     &mockReq[struct{}]{},
				handler: emptyMockHandler(&connect.Response[struct{}]{}, authz.CtxData{}),
			},
			res{
				&connect.Response[struct{}]{},
				false,
			},
		},
		{
			"error",
			args{
				ctx:     context.Background(),
				req:     &mockReq[struct{}]{},
				handler: errorMockHandler(),
			},
			res{
				nil,
				true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toConnectError(tt.args.ctx, tt.args.req, tt.args.handler(t))
			if (err != nil) != tt.res.wantErr {
				t.Errorf("toGRPCError() error = %v, wantErr %v", err, tt.res.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.res.want) {
				t.Errorf("toGRPCError() got = %v, want %v", got, tt.res.want)
			}
		})
	}
}
