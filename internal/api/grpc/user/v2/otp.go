package user

import (
	"context"

	user "github.com/zitadel/zitadel/pkg/grpc/user/v2"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
)

func (s *Server) AddOTPSMS(ctx context.Context, req *user.AddOTPSMSRequest) (*user.AddOTPSMSResponse, error) {
	details, err := s.command.AddHumanOTPSMS(ctx, req.GetUserId(), "")
	if err != nil {
		return nil, err
	}
	return &user.AddOTPSMSResponse{Details: object.DomainToDetailsPb(details)}, nil

}

func (s *Server) RemoveOTPSMS(ctx context.Context, req *user.RemoveOTPSMSRequest) (*user.RemoveOTPSMSResponse, error) {
	objectDetails, err := s.command.RemoveHumanOTPSMS(ctx, req.GetUserId(), "")
	if err != nil {
		return nil, err
	}
	return &user.RemoveOTPSMSResponse{Details: object.DomainToDetailsPb(objectDetails)}, nil
}

func (s *Server) AddOTPEmail(ctx context.Context, req *user.AddOTPEmailRequest) (*user.AddOTPEmailResponse, error) {
	details, err := s.command.AddHumanOTPEmail(ctx, req.GetUserId(), "")
	if err != nil {
		return nil, err
	}
	return &user.AddOTPEmailResponse{Details: object.DomainToDetailsPb(details)}, nil

}

func (s *Server) RemoveOTPEmail(ctx context.Context, req *user.RemoveOTPEmailRequest) (*user.RemoveOTPEmailResponse, error) {
	objectDetails, err := s.command.RemoveHumanOTPEmail(ctx, req.GetUserId(), "")
	if err != nil {
		return nil, err
	}
	return &user.RemoveOTPEmailResponse{Details: object.DomainToDetailsPb(objectDetails)}, nil
}
