package grpc

import (
	"path"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

const (
	Healthz    = "/Healthz"
	Readiness  = "/Ready"
	Validation = "/Validate"
)

var (
	Probes = []string{Healthz, Readiness, Validation}
)

func init() {
	Probes = append(Probes, AllPaths(grpc_reflection_v1alpha.ServerReflection_ServiceDesc)...)
}

func AllPaths(sd grpc.ServiceDesc) []string {
	paths := make([]string, 0, len(sd.Methods)+len(sd.Streams))
	for _, method := range sd.Methods {
		paths = append(paths, path.Join("/", sd.ServiceName, method.MethodName))
	}
	for _, stream := range sd.Streams {
		paths = append(paths, path.Join("/", sd.ServiceName, stream.StreamName))
	}
	return paths
}
