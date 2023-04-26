package middleware

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/grpc"
)

func Test_toGRPCError(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     interface{}
		handler grpc.UnaryHandler
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
				req:     &mockReq{},
				handler: emptyMockHandler,
			},
			res{
				&mockReq{},
				false,
			},
		},
		{
			"error",
			args{
				ctx:     context.Background(),
				req:     &mockReq{},
				handler: errorMockHandler,
			},
			res{
				nil,
				true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toGRPCError(tt.args.ctx, tt.args.req, tt.args.handler)
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
