package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) SetRestrictions(ctx context.Context, req *admin.SetRestrictionsRequest) (*admin.SetRestrictionsResponse, error) {
	details, err := s.command.SetInstanceRestrictions(ctx, &command.SetRestrictions{DisallowPublicOrgRegistration: req.DisallowPublicOrgRegistration})
	if err != nil {
		return nil, err
	}
	return &admin.SetRestrictionsResponse{
		Details: object.ChangeToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) GetRestrictions(ctx context.Context, _ *admin.GetRestrictionsRequest) (*admin.GetRestrictionsResponse, error) {
	restrictions, err := s.query.GetInstanceRestrictions(ctx)
	if err != nil {
		return nil, err
	}
	return &admin.GetRestrictionsResponse{
		Details:                       object.ToViewDetailsPb(restrictions.Sequence, restrictions.CreationDate, restrictions.ChangeDate, restrictions.ResourceOwner),
		DisallowPublicOrgRegistration: restrictions.DisallowPublicOrgRegistration,
	}, nil
}
