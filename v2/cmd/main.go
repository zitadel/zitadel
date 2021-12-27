package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/caos/zitadel/v2/api/admin"
	"github.com/caos/zitadel/v2/api/auth"
	"github.com/caos/zitadel/v2/api/mgmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func main() {
	l, err := net.Listen("tcp", ":50002")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.TODO()

	m := cmux.New(l)
	grpcL := m.Match(cmux.HTTP2())
	httpL := m.Match(cmux.HTTP1Fast())
	mux := runtime.NewServeMux()

	// services
	mgmtHandler := mgmt.New(ctx)
	adminHandler := admin.New(ctx)
	authHandler := auth.New(ctx)

	//grpc
	grpcServer := grpc.NewServer()
	mgmtHandler.RegisterGRPC(grpcServer)
	adminHandler.RegisterGRPC(grpcServer)
	authHandler.RegisterGRPC(grpcServer)

	//REST
	mgmtHandler.RegisterRESTGateway(ctx, mux)
	adminHandler.RegisterRESTGateway(ctx, mux)
	authHandler.RegisterRESTGateway(ctx, mux)
	httpS := &http.Server{Handler: h2c.NewHandler(mux, &http2.Server{})}

	errs := make(chan error, 2)
	go func() { errs <- httpS.Serve(httpL) }()
	go func() { errs <- grpcServer.Serve(grpcL) }()

	if err := m.Serve(); err != nil {
		log.Panicf("serve failed: %v\n", err)
	}

	for err := range errs {
		log.Fatal(err)
	}
}
