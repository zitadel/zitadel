package admin

import (
	"context"

	object_pb "github.com/zitadel/zitadel/v2/internal/api/grpc/object"
	"github.com/zitadel/zitadel/v2/internal/domain"
	admin_pb "github.com/zitadel/zitadel/v2/pkg/grpc/admin"
)

func (s *Server) ActivateFeatureLoginDefaultOrg(ctx context.Context, _ *admin_pb.ActivateFeatureLoginDefaultOrgRequest) (*admin_pb.ActivateFeatureLoginDefaultOrgResponse, error) {
	details, err := s.command.SetBooleanInstanceFeature(ctx, domain.FeatureLoginDefaultOrg, true)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ActivateFeatureLoginDefaultOrgResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil

}
