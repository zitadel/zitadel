package server

import (
	"context"
	"net"

	"github.com/caos/logging"
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/api/http"
)

const (
	defaultGrpcPort = "80"
)

type Server interface {
	GRPCPort() string
	GRPCServer() (*grpc.Server, error)
}

func StartServer(ctx context.Context, s Server) {
	port := grpcPort(s.GRPCPort())
	listener := http.CreateListener(port)
	server := createGrpcServer(s)
	serveServer(ctx, server, listener, port)
}

func createGrpcServer(s Server) *grpc.Server {
	grpcServer, err := s.GRPCServer()
	logging.Log("SERVE-k280HZ").OnError(err).Panic("failed to create grpc server")
	return grpcServer
}

func serveServer(ctx context.Context, server *grpc.Server, listener net.Listener, port string) {
	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	go func() {
		err := server.Serve(listener)
		logging.Log("SERVE-Ga3e94").OnError(err).Panic("grpc server serve failed")
	}()
	logging.LogWithFields("SERVE-bZ44QM", "port", port).Info("grpc server is listening")
}

func grpcPort(port string) string {
	if port == "" {
		return defaultGrpcPort
	}
	return port
}
