package grpc

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestGetHeader(t *testing.T) {
	type args struct {
		ctx        context.Context
		headername string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"empty context",
			args{
				ctx:        context.Background(),
				headername: "header",
			},
			"",
		},
		{
			"context without header",
			args{
				ctx:        metadata.NewIncomingContext(context.Background(), metadata.Pairs("header2", "value")),
				headername: "header",
			},
			"",
		},
		{
			"context with header",
			args{
				ctx:        metadata.NewIncomingContext(context.Background(), metadata.Pairs("header", "value")),
				headername: "header",
			},
			"value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHeader(tt.args.ctx, tt.args.headername); got != tt.want {
				t.Errorf("GetHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAuthorizationHeader(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"empty context",
			args{
				ctx: context.Background(),
			},
			"",
		},
		{
			"context without header",
			args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("header", "value")),
			},
			"",
		},
		{
			"context with header",
			args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "value")),
			},
			"value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAuthorizationHeader(tt.args.ctx); got != tt.want {
				t.Errorf("GetAuthorizationHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
