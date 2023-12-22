package system

import (
	"context"

	object_pb "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func (s *Server) SetInstanceFeature(ctx context.Context, req *system_pb.SetInstanceFeatureRequest) (*system_pb.SetInstanceFeatureResponse, error) {
	details, err := s.setInstanceFeature(ctx, req)
	if err != nil {
		return nil, err
	}
	return &system_pb.SetInstanceFeatureResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil

}

func (s *Server) setInstanceFeature(ctx context.Context, req *system_pb.SetInstanceFeatureRequest) (*domain.ObjectDetails, error) {
	feat := domain.Feature(req.FeatureId)
	if !feat.IsAFeature() {
		return nil, zerrors.ThrowInvalidArgument(nil, "SYST-SGV45", "Errors.Feature.NotExisting")
	}
	switch t := req.Value.(type) {
	case *system_pb.SetInstanceFeatureRequest_Bool:
		return s.command.SetBooleanInstanceFeature(ctx, feat, t.Bool)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "SYST-dag5g", "Errors.Feature.TypeNotSupported")
	}
}
