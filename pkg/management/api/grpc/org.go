package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
)

func (s *Server) GetOrgByID(ctx context.Context, in *OrgID) (*Org, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-sdo5g", "Not implemented")
}

func (s *Server) GetOrgByDomainGlobal(ctx context.Context, in *OrgDomain) (*Org, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-mop4s", "Not implemented")
}

func (s *Server) DeactivateOrg(ctx context.Context, in *OrgID) (*Org, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-vel3X", "Not implemented")
}

func (s *Server) ReactivateOrg(ctx context.Context, in *OrgID) (*Org, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-Scmk3", "Not implemented")
}

func (s *Server) OrgChanges(ctx context.Context, changesRequest *ChangeRequest) (*Changes, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-mfiF4", "Not implemented")
}
