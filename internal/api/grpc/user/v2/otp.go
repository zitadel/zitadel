package user

import (
	"context"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddOTPSMS(ctx context.Context, req *connect.Request[user.AddOTPSMSRequest]) (*connect.Response[user.AddOTPSMSResponse], error) {
	details, err := s.command.AddHumanOTPSMS(ctx, req.Msg.GetUserId(), "")
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.AddOTPSMSResponse{Details: object.DomainToDetailsPb(details)}), nil
}

func (s *Server) RemoveOTPSMS(ctx context.Context, req *connect.Request[user.RemoveOTPSMSRequest]) (*connect.Response[user.RemoveOTPSMSResponse], error) {
	objectDetails, err := s.command.RemoveHumanOTPSMS(ctx, req.Msg.GetUserId(), "")
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.RemoveOTPSMSResponse{Details: object.DomainToDetailsPb(objectDetails)}), nil
}

func (s *Server) AddOTPEmail(ctx context.Context, req *connect.Request[user.AddOTPEmailRequest]) (*connect.Response[user.AddOTPEmailResponse], error) {
	details, err := s.command.AddHumanOTPEmail(ctx, req.Msg.GetUserId(), "")
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.AddOTPEmailResponse{Details: object.DomainToDetailsPb(details)}), nil

}

func (s *Server) RemoveOTPEmail(ctx context.Context, req *connect.Request[user.RemoveOTPEmailRequest]) (*connect.Response[user.RemoveOTPEmailResponse], error) {
	objectDetails, err := s.command.RemoveHumanOTPEmail(ctx, req.Msg.GetUserId(), "")
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.RemoveOTPEmailResponse{Details: object.DomainToDetailsPb(objectDetails)}), nil
}
