package grpc

import (
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"strings"
)

type GatewayConfig struct {
	Port          string
	GRPCEndpoint  string
	CustomHeaders []string
}

type Gateway struct {
	grpcEndpoint string
	port         string
	cutomHeaders []string
}

func StartGateway(conf GatewayConfig) *Gateway {
	return &Gateway{
		grpcEndpoint: conf.GRPCEndpoint,
		port:         conf.Port,
		cutomHeaders: conf.CustomHeaders,
	}
}

func (gw *Gateway) Gateway() server.GatewayFunc {
	return RegisterAdminServiceHandlerFromEndpoint
}

func (gw *Gateway) GRPCEndpoint() string {
	return ":" + gw.grpcEndpoint
}

func (gw *Gateway) GatewayPort() string {
	return gw.port
}

func (gw *Gateway) GatewayServeMuxOptions() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithIncomingHeaderMatcher(func(header string) (string, bool) {
			for _, customHeader := range gw.cutomHeaders {
				if strings.HasPrefix(strings.ToLower(header), customHeader) {
					return header, true
				}
			}
			return header, false
		}),
	}
}
