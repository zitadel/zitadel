package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func TestAllPaths(t *testing.T) {
	type args struct {
		sd grpc.ServiceDesc
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "server reflection",
			args: args{grpc_reflection_v1alpha.ServerReflection_ServiceDesc},
			want: []string{"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AllPaths(tt.args.sd)
			assert.Equal(t, tt.want, got)
		})
	}
}
