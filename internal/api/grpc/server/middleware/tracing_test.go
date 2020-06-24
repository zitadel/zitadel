package middleware

import (
	"context"
	"testing"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/stats"
)

func Test_tracingServerHandler_TagRPC(t *testing.T) {
	type fields struct {
		IgnoredMethods []GRPCMethod
		ServerHandler  ocgrpc.ServerHandler
	}
	type args struct {
		ctx     context.Context
		tagInfo *stats.RPCTagInfo
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantSpan bool
	}{
		{
			"ignored method",
			fields{
				IgnoredMethods: []GRPCMethod{"ignore"},
				ServerHandler:  ocgrpc.ServerHandler{},
			},
			args{
				ctx: context.Background(),
				tagInfo: &stats.RPCTagInfo{
					FullMethodName: "ignore",
				},
			},
			false,
		},
		{
			"tag",
			fields{
				IgnoredMethods: []GRPCMethod{"ignore"},
				ServerHandler:  ocgrpc.ServerHandler{},
			},
			args{
				ctx: context.Background(),
				tagInfo: &stats.RPCTagInfo{
					FullMethodName: "tag",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &tracingServerHandler{
				IgnoredMethods: tt.fields.IgnoredMethods,
				ServerHandler:  tt.fields.ServerHandler,
			}
			got := s.TagRPC(tt.args.ctx, tt.args.tagInfo)
			if (trace.FromContext(got) != nil) != tt.wantSpan {
				t.Errorf("TagRPC() = %v, want %v", got, tt.wantSpan)
			}
		})
	}
}
