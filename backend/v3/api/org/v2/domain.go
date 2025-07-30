package org

import (
	"context"

	"connectrpc.com/connect"

	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

// AddOrganizationDomain implements [orgconnect.OrganizationServiceHandler].
func (s *Server) AddOrganizationDomain(ctx context.Context, req *connect.Request[org.AddOrganizationDomainRequest]) (*connect.Response[org.AddOrganizationDomainResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.AddOrganizationDomain(ctx, req)
}

// DeleteOrganizationDomain implements [orgconnect.OrganizationServiceHandler].
func (s *Server) DeleteOrganizationDomain(ctx context.Context, req *connect.Request[org.DeleteOrganizationDomainRequest]) (*connect.Response[org.DeleteOrganizationDomainResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.DeleteOrganizationDomain(ctx, req)
}

// GenerateOrganizationDomainValidation implements [orgconnect.OrganizationServiceHandler].
func (s *Server) GenerateOrganizationDomainValidation(ctx context.Context, req *connect.Request[org.GenerateOrganizationDomainValidationRequest]) (*connect.Response[org.GenerateOrganizationDomainValidationResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.GenerateOrganizationDomainValidation(ctx, req)
}

// ListOrganizationDomains implements [orgconnect.OrganizationServiceHandler].
func (s *Server) ListOrganizationDomains(ctx context.Context, req *connect.Request[org.ListOrganizationDomainsRequest]) (*connect.Response[org.ListOrganizationDomainsResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.ListOrganizationDomains(ctx, req)
}

// VerifyOrganizationDomain implements [orgconnect.OrganizationServiceHandler].
func (s *Server) VerifyOrganizationDomain(ctx context.Context, req *connect.Request[org.VerifyOrganizationDomainRequest]) (*connect.Response[org.VerifyOrganizationDomainResponse], error) {
	return s.UnimplementedOrganizationServiceHandler.VerifyOrganizationDomain(ctx, req)
}
