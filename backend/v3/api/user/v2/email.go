package userv2

// import (
// 	"context"

// 	"github.com/zitadel/zitadel/backend/v3/domain"
// 	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
// )

// func SetEmail(ctx context.Context, req *user.SetEmailRequest) (resp *user.SetEmailResponse, err error) {
// 	var (
// 		verification domain.SetEmailOpt
// 		returnCode   *domain.ReturnCodeCommand
// 	)

// 	switch req.GetVerification().(type) {
// 	case *user.SetEmailRequest_IsVerified:
// 		verification = domain.NewEmailVerifiedCommand(req.GetUserId(), req.GetIsVerified())
// 	case *user.SetEmailRequest_SendCode:
// 		verification = domain.NewSendCodeCommand(req.GetUserId(), req.GetSendCode().UrlTemplate)
// 	case *user.SetEmailRequest_ReturnCode:
// 		returnCode = domain.NewReturnCodeCommand(req.GetUserId())
// 		verification = returnCode
// 	default:
// 		verification = domain.NewSendCodeCommand(req.GetUserId(), nil)
// 	}

// 	err = domain.Invoke(ctx, domain.NewSetEmailCommand(req.GetUserId(), req.GetEmail(), verification))
// 	if err != nil {
// 		return nil, err
// 	}

// 	var code *string
// 	if returnCode != nil && returnCode.Code != "" {
// 		code = &returnCode.Code
// 	}

// 	return &user.SetEmailResponse{
// 		VerificationCode: code,
// 	}, nil
// }

// func SendEmailCode(ctx context.Context, req *user.SendEmailCodeRequest) (resp *user.SendEmailCodeResponse, err error) {
// 	var (
// 		returnCode *domain.ReturnCodeCommand
// 		cmd        domain.Commander
// 	)

// 	switch req.GetVerification().(type) {
// 	case *user.SendEmailCodeRequest_SendCode:
// 		cmd = domain.NewSendCodeCommand(req.GetUserId(), req.GetSendCode().UrlTemplate)
// 	case *user.SendEmailCodeRequest_ReturnCode:
// 		returnCode = domain.NewReturnCodeCommand(req.GetUserId())
// 		cmd = returnCode
// 	default:
// 		cmd = domain.NewSendCodeCommand(req.GetUserId(), req.GetSendCode().UrlTemplate)
// 	}
// 	err = domain.Invoke(ctx, cmd)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp = new(user.SendEmailCodeResponse)
// 	if returnCode != nil {
// 		resp.VerificationCode = &returnCode.Code
// 	}
// 	return resp, nil
// }

// func ResendEmailCode(ctx context.Context, req *user.ResendEmailCodeRequest) (resp *user.SendEmailCodeResponse, err error) {
// 	var (
// 		returnCode *domain.ReturnCodeCommand
// 		cmd        domain.Commander
// 	)

// 	switch req.GetVerification().(type) {
// 	case *user.ResendEmailCodeRequest_SendCode:
// 		cmd = domain.NewSendCodeCommand(req.GetUserId(), req.GetSendCode().UrlTemplate)
// 	case *user.ResendEmailCodeRequest_ReturnCode:
// 		returnCode = domain.NewReturnCodeCommand(req.GetUserId())
// 		cmd = returnCode
// 	default:
// 		cmd = domain.NewSendCodeCommand(req.GetUserId(), req.GetSendCode().UrlTemplate)
// 	}
// 	err = domain.Invoke(ctx, cmd)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp = new(user.SendEmailCodeResponse)
// 	if returnCode != nil {
// 		resp.VerificationCode = &returnCode.Code
// 	}
// 	return resp, nil
// }
