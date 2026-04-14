package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_toGRPCError(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     any
		handler grpc.UnaryHandler
	}
	tests := []struct {
		name     string
		args     args
		want     any
		wantCode codes.Code
	}{
		{
			name: "no error",
			args: args{
				ctx:     context.Background(),
				req:     &mockReq{},
				handler: emptyMockHandler,
			},
			want: &mockReq{},
		},
		{
			name: "error",
			args: args{
				ctx:     context.Background(),
				req:     &mockReq{},
				handler: errorMockHandler,
			},
			want:     nil,
			wantCode: codes.FailedPrecondition,
		},
		{
			name: "grpc status error",
			args: args{
				ctx: context.Background(),
				req: &mockReq{},
				handler: func(context.Context, any) (any, error) {
					return nil, status.Error(codes.Unauthenticated, "auth header missing")
				},
			},
			want:     nil,
			wantCode: codes.Unauthenticated,
		},
		{
			name: "panic with string",
			args: args{
				ctx:     context.Background(),
				req:     &mockReq{},
				handler: panicMockHandler("test panic"),
			},
			want:     nil,
			wantCode: codes.Internal,
		},
		{
			name: "panic with error",
			args: args{
				ctx:     context.Background(),
				req:     &mockReq{},
				handler: panicMockHandler(errors.New("oops")),
			},
			want:     nil,
			wantCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toGRPCError(tt.args.ctx, tt.args.req, tt.args.handler)
			if tt.wantCode != 0 {
				status, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.wantCode, status.Code())
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
