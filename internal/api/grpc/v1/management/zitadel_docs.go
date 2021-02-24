package management

import (
	"context"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetZitadelDocs(ctx context.Context, _ *empty.Empty) (*management.ZitadelDocs, error) {
	return &management.ZitadelDocs{
		Issuer:            s.systemDefaults.ZitadelDocs.Issuer,
		DiscoveryEndpoint: s.systemDefaults.ZitadelDocs.DiscoveryEndpoint,
	}, nil
}
