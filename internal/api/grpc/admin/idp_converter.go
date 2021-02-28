package admin

import (
	"github.com/caos/zitadel/internal/api/grpc/idp"
	"github.com/caos/zitadel/internal/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func addOIDCIDPRequestToDomain(req *admin_pb.AddOIDCIDPRequest) *domain.IDPConfig {
	return &domain.IDPConfig{
		Name:        req.Name,
		OIDCConfig:  addOIDCIDPRequestToDomainOIDCIDPConfig(req),
		StylingType: idp.IDPStylingTypeToDomain(req.StylingType),
		Type:        domain.IDPConfigTypeOIDC,
	}
}

func addOIDCIDPRequestToDomainOIDCIDPConfig(req *admin_pb.AddOIDCIDPRequest) *domain.OIDCIDPConfig {
	return &domain.OIDCIDPConfig{
		ClientID:              req.ClientId,
		ClientSecretString:    req.ClientSecret,
		Issuer:                req.Issuer,
		Scopes:                req.Scopes,
		IDPDisplayNameMapping: idp.MappingFieldToDomain(req.DisplayNameMapping),
		UsernameMapping:       idp.MappingFieldToDomain(req.UsernameMapping),
	}
}

func updateIDPToDomain(req *admin_pb.UpdateIDPRequest) *domain.IDPConfig {
	return &domain.IDPConfig{
		IDPConfigID: req.Id,
		Name:        req.Name,
		StylingType: idp.IDPStylingTypeToDomain(req.StylingType),
	}
}

func updateOIDCConfigToDomain(req *admin_pb.UpdateIDPOIDCConfigRequest) *domain.OIDCIDPConfig {
	return &domain.OIDCIDPConfig{
		IDPConfigID:           req.IdpId,
		ClientID:              req.ClientId,
		ClientSecretString:    req.ClientSecret,
		Issuer:                req.Issuer,
		Scopes:                req.Scopes,
		IDPDisplayNameMapping: idp.MappingFieldToDomain(req.DisplayNameMapping),
		UsernameMapping:       idp.MappingFieldToDomain(req.UsernameMapping),
	}
}
