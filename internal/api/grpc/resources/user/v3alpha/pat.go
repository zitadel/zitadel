package user

import (
	"context"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"

	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	z_oidc "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/command"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

func (s *Server) AddPersonalAccessToken(ctx context.Context, req *user.AddPersonalAccessTokenRequest) (_ *user.AddPersonalAccessTokenResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	pat := addPersonalAccessTokenRequestToAddPAT(req)
	details, err := s.command.AddPAT(ctx, pat)
	if err != nil {
		return nil, err
	}
	return &user.AddPersonalAccessTokenResponse{
		Details:               resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
		PersonalAccessTokenId: details.ID,
		PersonalAccessToken:   pat.PAT.Token,
	}, nil
}

func addPersonalAccessTokenRequestToAddPAT(req *user.AddPersonalAccessTokenRequest) *command.AddPAT {
	if req == nil {
		return nil
	}
	return &command.AddPAT{
		ResourceOwner: organizationToUpdateResourceOwner(req.Organization),
		UserID:        req.GetId(),
		PAT:           setPersonalAccessTokenToAddPAT(req.GetPersonalAccessToken()),
	}
}

func setPersonalAccessTokenToAddPAT(set *user.SetPersonalAccessToken) *command.PAT {
	if set == nil {
		return nil
	}
	expDate := time.Time{}
	if set.GetExpirationDate() != nil {
		expDate = set.GetExpirationDate().AsTime()
	}
	return &command.PAT{
		ExpirationDate: expDate,
		Scopes:         []string{oidc.ScopeOpenID, oidc.ScopeProfile, z_oidc.ScopeUserMetaData, z_oidc.ScopeResourceOwner},
	}
}

func (s *Server) RemovePersonalAccessToken(ctx context.Context, req *user.RemovePersonalAccessTokenRequest) (_ *user.RemovePersonalAccessTokenResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeletePAT(ctx, organizationToUpdateResourceOwner(req.Organization), req.GetId(), req.GetPersonalAccessTokenId())
	if err != nil {
		return nil, err
	}
	return &user.RemovePersonalAccessTokenResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}
