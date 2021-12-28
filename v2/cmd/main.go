package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/caos/zitadel/v2/api/admin"
	"github.com/caos/zitadel/v2/api/auth"
	"github.com/caos/zitadel/v2/api/mgmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func main() {
	grpcServ := grpc.NewServer()
	wrappedGrpc := grpcweb.WrapServer(grpcServ)
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", home)
	mux := runtime.NewServeMux()

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// services
	mgmtHandler := mgmt.New()
	adminHandler := admin.New()
	authHandler := auth.New()

	//grpc
	mgmtHandler.RegisterGRPC(grpcServ)
	adminHandler.RegisterGRPC(grpcServ)
	authHandler.RegisterGRPC(grpcServ)

	//REST
	mgmtHandler.RegisterRESTGateway(ctx, mux)
	adminHandler.RegisterRESTGateway(ctx, mux)
	authHandler.RegisterRESTGateway(ctx, mux)

	mixedHandler := newHTTPandGRPCMux(mux, grpcServ, wrappedGrpc)
	http2Server := &http2.Server{}
	http1Server := &http.Server{Handler: h2c.NewHandler(mixedHandler, http2Server)}
	lis, err := net.Listen("tcp", ":50002")
	if err != nil {
		panic(err)
	}

	err = http1Server.Serve(lis)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		panic(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from http handler!\n")
}

func newHTTPandGRPCMux(httpHand, grpcHandler http.Handler, wrappedGrpc *grpcweb.WrappedGrpcServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(r) {
			wrappedGrpc.ServeHTTP(w, r)
			return
		}
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("content-type"), "application/grpc") {
			grpcHandler.ServeHTTP(w, r)
			return
		}
		httpHand.ServeHTTP(w, r)
	})
}

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net"
// 	"net/http"

// 	"github.com/caos/zitadel/v2/api/admin"
// 	"github.com/caos/zitadel/v2/api/auth"
// 	"github.com/caos/zitadel/v2/api/mgmt"
// 	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
// 	"github.com/soheilhy/cmux"
// 	"golang.org/x/net/http2"
// 	"golang.org/x/net/http2/h2c"
// 	"google.golang.org/grpc"
// )

// func main() {
// 	httpMux := http.NewServeMux()
// 	httpMux.HandleFunc("/", home)

// 	l, err := net.Listen("tcp", ":50002")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	ctx := context.TODO()

// 	m := cmux.New(l)
// 	grpcL := m.Match(cmux.HTTP2())
// 	httpL := m.Match(cmux.HTTP1Fast())

// 	// services
// 	mgmtHandler := mgmt.New()
// 	adminHandler := admin.New()
// 	authHandler := auth.New()

// 	//grpc
// 	grpcServer := grpc.NewServer()
// 	mgmtHandler.RegisterGRPC(grpcServer)
// 	adminHandler.RegisterGRPC(grpcServer)
// 	authHandler.RegisterGRPC(grpcServer)

// 	//REST
// 	mux := runtime.NewServeMux()
// 	mgmtHandler.RegisterRESTGateway(ctx, mux)
// 	adminHandler.RegisterRESTGateway(ctx, mux)
// 	authHandler.RegisterRESTGateway(ctx, mux)
// 	httpS := &http.Server{Handler: h2c.NewHandler(mux, &http2.Server{})}

// 	errs := make(chan error, 2)
// 	go func() { errs <- httpS.Serve(httpL) }()
// 	go func() { errs <- grpcServer.Serve(grpcL) }()

// 	if err := m.Serve(); err != nil {
// 		log.Panicf("serve failed: %v\n", err)
// 	}

// 	for err := range errs {
// 		log.Fatal(err)
// 	}
// }

// func home(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "hello from http handler!\n")
// }
