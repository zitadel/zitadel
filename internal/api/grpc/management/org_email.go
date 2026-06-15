package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetOrgEmailProvider(ctx context.Context, req *mgmt_pb.GetOrgEmailProviderRequest) (*mgmt_pb.GetOrgEmailProviderResponse, error) {
	smtp, err := s.query.OrgSMTPConfigActive(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetOrgEmailProviderResponse{
		Config: orgEmailProviderToProviderPb(smtp),
	}, nil
}

func (s *Server) GetOrgEmailProviderById(ctx context.Context, req *mgmt_pb.GetOrgEmailProviderByIdRequest) (*mgmt_pb.GetOrgEmailProviderByIdResponse, error) {
	smtp, err := s.query.OrgSMTPConfigByID(ctx, authz.GetCtxData(ctx).OrgID, req.Id)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetOrgEmailProviderByIdResponse{
		Config: orgEmailProviderToProviderPb(smtp),
	}, nil
}

func (s *Server) ListOrgEmailProviders(ctx context.Context, req *mgmt_pb.ListOrgEmailProvidersRequest) (*mgmt_pb.ListOrgEmailProvidersResponse, error) {
	queries, err := listOrgEmailProvidersToModel(req)
	if err != nil {
		return nil, err
	}
	result, err := s.query.SearchOrgSMTPConfigs(ctx, authz.GetCtxData(ctx).OrgID, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListOrgEmailProvidersResponse{
		Details: object.ToListDetails(result.Count, result.Sequence, result.LastRun),
		Result:  orgEmailProvidersToPb(result.Configs),
	}, nil
}

func (s *Server) AddOrgEmailProviderSMTP(ctx context.Context, req *mgmt_pb.AddOrgEmailProviderSMTPRequest) (*mgmt_pb.AddOrgEmailProviderSMTPResponse, error) {
	config := addOrgEmailProviderSMTPToConfig(ctx, req)
	if err := s.command.AddOrgSMTPConfig(ctx, config); err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgEmailProviderSMTPResponse{
		Details: object.DomainToChangeDetailsPb(config.Details),
		Id:      config.ID,
	}, nil
}

func (s *Server) UpdateOrgEmailProviderSMTP(ctx context.Context, req *mgmt_pb.UpdateOrgEmailProviderSMTPRequest) (*mgmt_pb.UpdateOrgEmailProviderSMTPResponse, error) {
	config := updateOrgEmailProviderSMTPToConfig(ctx, req)
	if err := s.command.ChangeOrgSMTPConfig(ctx, config); err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOrgEmailProviderSMTPResponse{
		Details: object.DomainToChangeDetailsPb(config.Details),
	}, nil
}

func (s *Server) UpdateOrgEmailProviderSMTPPassword(ctx context.Context, req *mgmt_pb.UpdateOrgEmailProviderSMTPPasswordRequest) (*mgmt_pb.UpdateOrgEmailProviderSMTPPasswordResponse, error) {
	details, err := s.command.ChangeOrgSMTPConfigPassword(ctx, authz.GetCtxData(ctx).OrgID, req.Id, req.Password)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOrgEmailProviderSMTPPasswordResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddOrgEmailProviderHTTP(ctx context.Context, req *mgmt_pb.AddOrgEmailProviderHTTPRequest) (*mgmt_pb.AddOrgEmailProviderHTTPResponse, error) {
	config := addOrgEmailProviderHTTPToConfig(ctx, req)
	if err := s.command.AddOrgSMTPConfigHTTP(ctx, config); err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgEmailProviderHTTPResponse{
		Details:    object.DomainToChangeDetailsPb(config.Details),
		Id:         config.ID,
		SigningKey: config.SigningKey,
	}, nil
}

func (s *Server) UpdateOrgEmailProviderHTTP(ctx context.Context, req *mgmt_pb.UpdateOrgEmailProviderHTTPRequest) (*mgmt_pb.UpdateOrgEmailProviderHTTPResponse, error) {
	config := updateOrgEmailProviderHTTPToConfig(ctx, req)
	if err := s.command.ChangeOrgSMTPConfigHTTP(ctx, config); err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOrgEmailProviderHTTPResponse{
		Details:    object.DomainToChangeDetailsPb(config.Details),
		SigningKey: config.SigningKey,
	}, nil
}

func (s *Server) ActivateOrgEmailProvider(ctx context.Context, req *mgmt_pb.ActivateOrgEmailProviderRequest) (*mgmt_pb.ActivateOrgEmailProviderResponse, error) {
	result, err := s.command.ActivateOrgSMTPConfig(ctx, authz.GetCtxData(ctx).OrgID, req.Id)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ActivateOrgEmailProviderResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}

func (s *Server) DeactivateOrgEmailProvider(ctx context.Context, req *mgmt_pb.DeactivateOrgEmailProviderRequest) (*mgmt_pb.DeactivateOrgEmailProviderResponse, error) {
	result, err := s.command.DeactivateOrgSMTPConfig(ctx, authz.GetCtxData(ctx).OrgID, req.Id)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateOrgEmailProviderResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}

func (s *Server) RemoveOrgEmailProvider(ctx context.Context, req *mgmt_pb.RemoveOrgEmailProviderRequest) (*mgmt_pb.RemoveOrgEmailProviderResponse, error) {
	details, err := s.command.RemoveOrgSMTPConfig(ctx, authz.GetCtxData(ctx).OrgID, req.Id)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveOrgEmailProviderResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) TestOrgEmailProviderSMTP(ctx context.Context, req *mgmt_pb.TestOrgEmailProviderSMTPRequest) (*mgmt_pb.TestOrgEmailProviderSMTPResponse, error) {
	smtpConfig := testOrgEmailProviderSMTPToConfig(req)
	if err := s.command.TestOrgSMTPConfig(ctx, authz.GetCtxData(ctx).OrgID, "", req.ReceiverAddress, smtpConfig); err != nil {
		return nil, err
	}
	return &mgmt_pb.TestOrgEmailProviderSMTPResponse{}, nil
}

func (s *Server) TestOrgEmailProviderSMTPById(ctx context.Context, req *mgmt_pb.TestOrgEmailProviderSMTPByIdRequest) (*mgmt_pb.TestOrgEmailProviderSMTPByIdResponse, error) {
	if err := s.command.TestOrgSMTPConfigById(ctx, authz.GetCtxData(ctx).OrgID, req.Id, req.ReceiverAddress); err != nil {
		return nil, err
	}
	return &mgmt_pb.TestOrgEmailProviderSMTPByIdResponse{}, nil
}
