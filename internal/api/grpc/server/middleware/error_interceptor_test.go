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
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"no error",
			args{
				ctx:     context.Background(),
				req:     &mockReq{},
				handler: emptyMockHandler,
			},
			&mockReq{},
			false,
		},
		{
			"error",
			args{
				ctx:     context.Background(),
				req:     &mockReq{},
				handler: errorMockHandler,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toGRPCError(tt.args.ctx, tt.args.req, tt.args.handler, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("toGRPCError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toGRPCError() got = %v, want %v", got, tt.want)
			}
		})
	}
}
