package management

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetIAM(ctx context.Context, req *mgmt_pb.GetIAMRequest) (*mgmt_pb.GetIAMResponse, error) {
	iam, err := s.query.IAMByID(ctx, domain.IAMID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetIAMResponse{
		GlobalOrgId:  iam.GlobalOrgID,
		IamProjectId: iam.IAMProjectID,
	}, nil
}
