package org

import (
	"context"

	"connectrpc.com/connect"

	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

// DeleteOrganizationMetadata implements [orgconnect.OrganizationServiceHandler].
func (s *Server) DeleteOrganizationMetadata(ctx context.Context, req *connect.Request[org.DeleteOrganizationMetadataRequest]) (*connect.Response[org.DeleteOrganizationMetadataResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.DeleteOrganizationMetadata(ctx, req)
}

// ListOrganizationMetadata implements [orgconnect.OrganizationServiceHandler].
func (s *Server) ListOrganizationMetadata(ctx context.Context, req *connect.Request[org.ListOrganizationMetadataRequest]) (*connect.Response[org.ListOrganizationMetadataResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.ListOrganizationMetadata(ctx, req)
}

// SetOrganizationMetadata implements [orgconnect.OrganizationServiceHandler].
func (s *Server) SetOrganizationMetadata(ctx context.Context, req *connect.Request[org.SetOrganizationMetadataRequest]) (*connect.Response[org.SetOrganizationMetadataResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.SetOrganizationMetadata(ctx, req)
}
