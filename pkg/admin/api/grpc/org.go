package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
)

func (s *Server) GetOrgByID(ctx context.Context, orgID *OrgID) (_ *Org, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-ruc8e", "Not implemented")
}

func (s *Server) SearchOrgs(ctx context.Context, request *OrgSearchRequest) (_ *OrgSearchResponse, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-ruc8e", "Not implemented")
}

func (s *Server) IsOrgUnique(ctx context.Context, request *UniqueOrgRequest) (org *UniqueOrgResponse, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-ruc8e", "Not implemented")
}

func (s *Server) SetUpOrg(ctx context.Context, orgSetUp *OrgSetUpRequest) (_ *OrgSetUpResponse, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-ruc8e", "Not implemented")
}
