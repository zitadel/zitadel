package org

import (
	"context"

	"connectrpc.com/connect"

	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

// import (
// 	"context"

// 	"github.com/zitadel/zitadel/backend/v3/domain"
// 	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
// )

// func CreateOrg(ctx context.Context, req *org.AddOrganizationRequest) (resp *org.AddOrganizationResponse, err error) {
// 	cmd := domain.NewAddOrgCommand(
// 		req.GetName(),
// 		addOrgAdminToCommand(req.GetAdmins()...)...,
// 	)
// 	err = domain.Invoke(ctx, cmd)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &org.AddOrganizationResponse{
// 		OrganizationId: cmd.ID,
// 	}, nil
// }

// func addOrgAdminToCommand(admins ...*org.AddOrganizationRequest_Admin) []*domain.AddMemberCommand {
// 	cmds := make([]*domain.AddMemberCommand, len(admins))
// 	for i, admin := range admins {
// 		cmds[i] = &domain.AddMemberCommand{
// 			UserID: admin.GetUserId(),
// 			Roles:  admin.GetRoles(),
// 		}
// 	}
// 	return cmds
// }

// ActivateOrganization implements [orgconnect.OrganizationServiceHandler].
func (s *Server) ActivateOrganization(ctx context.Context, req *connect.Request[org.ActivateOrganizationRequest]) (*connect.Response[org.ActivateOrganizationResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.ActivateOrganization(ctx, req)
}

// CreateOrganization implements [orgconnect.OrganizationServiceHandler].
func (s *Server) CreateOrganization(ctx context.Context, req *connect.Request[org.CreateOrganizationRequest]) (*connect.Response[org.CreateOrganizationResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.CreateOrganization(ctx, req)
}

// DeactivateOrganization implements [orgconnect.OrganizationServiceHandler].
func (s *Server) DeactivateOrganization(ctx context.Context, req *connect.Request[org.DeactivateOrganizationRequest]) (*connect.Response[org.DeactivateOrganizationResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.DeactivateOrganization(ctx, req)
}

// DeleteOrganization implements [orgconnect.OrganizationServiceHandler].
func (s *Server) DeleteOrganization(ctx context.Context, req *connect.Request[org.DeleteOrganizationRequest]) (*connect.Response[org.DeleteOrganizationResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.DeleteOrganization(ctx, req)
}

// ListOrganizations implements [orgconnect.OrganizationServiceHandler].
func (s *Server) ListOrganizations(ctx context.Context, req *connect.Request[org.ListOrganizationsRequest]) (*connect.Response[org.ListOrganizationsResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.ListOrganizations(ctx, req)
}

// UpdateOrganization implements [orgconnect.OrganizationServiceHandler].
func (s *Server) UpdateOrganization(ctx context.Context, req *connect.Request[org.UpdateOrganizationRequest]) (*connect.Response[org.UpdateOrganizationResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.UpdateOrganization(ctx, req)
}
