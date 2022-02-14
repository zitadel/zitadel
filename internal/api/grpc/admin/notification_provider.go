package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/notification/channels/fs"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetFileSystemNotificationProvider(ctx context.Context, req *admin_pb.GetFileSystemNotificationProviderRequest) (*admin_pb.GetFileSystemNotificationProviderResponse, error) {
	//result, err := s.query.SMSProviderConfigByID(ctx, req.Id)
	//if err != nil {
	//	return nil, err
	//
	//}
	//return &admin_pb.GetFileSystemNotificationProviderResponse{
	//	Config: SMTPConfigToPb(result),
	//}, nil
	return nil, nil
}

func (s *Server) AddFileSystemNotificationProvider(ctx context.Context, req *admin_pb.AddFileSystemNotificationProviderRequest) (*admin_pb.AddFileSystemNotificationProviderResponse, error) {
	result, err := s.command.AddDebugNotificationProviderFile(ctx, &fs.FSConfig{Compact: req.Compact, Enabled: req.Enabled})
	if err != nil {
		return nil, err

	}
	return &admin_pb.AddFileSystemNotificationProviderResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}

func (s *Server) UpdateFileSystemNotificationProvider(ctx context.Context, req *admin_pb.UpdateFileSystemNotificationProviderRequest) (*admin_pb.UpdateFileSystemNotificationProviderResponse, error) {
	result, err := s.command.ChangeDefaultNotificationFile(ctx, &fs.FSConfig{Compact: req.Compact, Enabled: req.Enabled})
	if err != nil {
		return nil, err

	}
	return &admin_pb.UpdateFileSystemNotificationProviderResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}

func (s *Server) RemoveFileSystemNotificationProvider(ctx context.Context, _ *admin_pb.RemoveFileSystemNotificationProviderRequest) (*admin_pb.RemoveFileSystemNotificationProviderResponse, error) {
	result, err := s.command.RemoveDefaultNotificationFile(ctx)
	if err != nil {
		return nil, err

	}
	return &admin_pb.RemoveFileSystemNotificationProviderResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}

func (s *Server) GetLogNotificationProvider(ctx context.Context, req *admin_pb.GetLogNotificationProviderRequest) (*admin_pb.GetLogNotificationProviderResponse, error) {
	//result, err := s.query.SMSProviderConfigByID(ctx, req.Id)
	//if err != nil {
	//	return nil, err
	//
	//}
	//return &admin_pb.GetLogNotificationProviderResponse{
	//	Config: SMTPConfigToPb(result),
	//}, nil
	return nil, nil
}

func (s *Server) AddLogNotificationProvider(ctx context.Context, req *admin_pb.AddLogNotificationProviderRequest) (*admin_pb.AddLogNotificationProviderResponse, error) {
	result, err := s.command.AddDebugNotificationProviderFile(ctx, &fs.FSConfig{Compact: req.Compact, Enabled: req.Enabled})
	if err != nil {
		return nil, err

	}
	return &admin_pb.AddLogNotificationProviderResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}

func (s *Server) UpdateLogNotificationProvider(ctx context.Context, req *admin_pb.UpdateLogNotificationProviderRequest) (*admin_pb.UpdateLogNotificationProviderResponse, error) {
	result, err := s.command.ChangeDefaultNotificationFile(ctx, &fs.FSConfig{Compact: req.Compact, Enabled: req.Enabled})
	if err != nil {
		return nil, err

	}
	return &admin_pb.UpdateLogNotificationProviderResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}

func (s *Server) RemoveLogNotificationProvider(ctx context.Context, _ *admin_pb.RemoveLogNotificationProviderRequest) (*admin_pb.RemoveLogNotificationProviderResponse, error) {
	result, err := s.command.RemoveDefaultNotificationFile(ctx)
	if err != nil {
		return nil, err

	}
	return &admin_pb.RemoveLogNotificationProviderResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}
