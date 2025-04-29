package userv2

import (
	"context"

	"github.com/muhlemmer/gu"
	"github.com/zitadel/zitadel/backend/command/v2/domain"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) SetEmail(ctx context.Context, req *user.SetEmailRequest) (resp *user.SetEmailResponse, err error) {
	request := &domain.SetUserEmail{
		UserID: req.GetUserId(),
		Email:  req.GetEmail(),
	}
	switch req.GetVerification().(type) {
	case *user.SetEmailRequest_IsVerified:
		request.IsVerified = gu.Ptr(req.GetIsVerified())
	case *user.SetEmailRequest_SendCode:
		request.SendCode = &domain.SendCode{
			URLTemplate: req.GetSendCode().UrlTemplate,
		}
	case *user.SetEmailRequest_ReturnCode:
		request.ReturnCode = new(domain.ReturnCode)
	}
	if err := s.domain.SetUserEmail(ctx, request); err != nil {
		return nil, err
	}

	response := new(user.SetEmailResponse)
	if request.ReturnCode != nil {
		response.VerificationCode = &request.ReturnCode.Code
	}
	return response, nil
}
