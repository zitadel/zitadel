package management

import (
	"context"

	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetIAM(ctx context.Context, _ *mgmt_pb.GetIAMRequest) (*mgmt_pb.GetIAMResponse, error) {
	instance, err := s.query.Instance(ctx, true)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetIAMResponse{
		GlobalOrgId:  instance.DefaultOrgID,
		IamProjectId: instance.IAMProjectID,
		DefaultOrgId: instance.DefaultOrgID,
	}, nil
}
