package admin

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/logging"

	object_pb "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) ActivateFeatureLoginDefaultOrg(ctx context.Context, _ *admin_pb.ActivateFeatureLoginDefaultOrgRequest) (*admin_pb.ActivateFeatureLoginDefaultOrgResponse, error) {
	details, err := s.command.SetInstanceFeatures(ctx, &command.InstanceFeatures{
		LoginDefaultOrg: gu.Ptr(true),
	})
	if err != nil {
		return nil, err
	}
	_, err = projection.InstanceFeatureProjection.Trigger(ctx, handler.WithAwaitRunning())
	logging.OnError(err).Warn("trigger instance feature projection")

	return &admin_pb.ActivateFeatureLoginDefaultOrgResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil

}
