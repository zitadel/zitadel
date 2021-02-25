package admin

import (
	"github.com/caos/zitadel/internal/api/grpc/idp"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/v2/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func addOIDCIDPRequestToDomain(req *admin_pb.AddOIDCIDPRequest) *domain.IDPConfig {
	return &domain.IDPConfig{
		Name:        req.Name,
		OIDCConfig:  addOIDCIDPRequestToDomainOIDCIDPConfig(req),
		StylingType: idp.IDPStylingTypeToDomain(req.StylingType),
	}
}

func addOIDCIDPRequestToDomainOIDCIDPConfig(req *admin_pb.AddOIDCIDPRequest) *domain.OIDCIDPConfig {
	var clientSecret *crypto.CryptoValue
	return &domain.OIDCIDPConfig{
		ClientID:              req.ClientId,
		ClientSecret:          clientSecret,
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
		IDPConfigID: req.IdpId,
		// ClientID: req.,
		// ClientSecretString: req.,
		Issuer:                req.Issuer,
		Scopes:                req.Scopes,
		IDPDisplayNameMapping: idp.MappingFieldToDomain(req.DisplayNameMapping),
		UsernameMapping:       idp.MappingFieldToDomain(req.UsernameMapping),
	}
}
