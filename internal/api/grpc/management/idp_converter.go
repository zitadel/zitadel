package management

import (
	idp_grpc "github.com/caos/zitadel/internal/api/grpc/idp"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	user_model "github.com/caos/zitadel/internal/user/model"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func addOIDCIDPRequestToDomain(req *mgmt_pb.AddOrgOIDCIDPRequest) *domain.OIDCIDPConfig {
	return &domain.OIDCIDPConfig{
		ClientID:              req.ClientId,
		ClientSecretString:    req.ClientSecret,
		Issuer:                req.Issuer,
		Scopes:                req.Scopes,
		IDPDisplayNameMapping: idp_grpc.MappingFieldToDomain(req.DisplayNameMapping),
		UsernameMapping:       idp_grpc.MappingFieldToDomain(req.UsernameMapping),
		CommonIDPConfig: domain.CommonIDPConfig{
			Name:        req.Name,
			StylingType: idp_grpc.IDPStylingTypeToDomain(req.StylingType),
			Type:        domain.IDPConfigTypeOIDC,
		},
	}
}

func addAuthConnectorIDPRequestToDomain(req *mgmt_pb.AddOrgAuthConnectorIDPRequest) *domain.AuthConnectorIDPConfig {
	return &domain.AuthConnectorIDPConfig{
		BaseURL:    req.BaseUrl,
		ProviderID: req.ProviderId,
		MachineID:  req.MachineId,
		CommonIDPConfig: domain.CommonIDPConfig{
			Name:        req.Name,
			StylingType: idp_grpc.IDPStylingTypeToDomain(req.StylingType),
			Type:        domain.IDPConfigTypeAuthConnector,
		},
	}
}

func updateIDPToDomain(req *mgmt_pb.UpdateOrgIDPRequest) *domain.CommonIDPConfig {
	return &domain.CommonIDPConfig{
		IDPConfigID: req.IdpId,
		Name:        req.Name,
		StylingType: idp_grpc.IDPStylingTypeToDomain(req.StylingType),
	}
}

func updateOIDCConfigToDomain(req *mgmt_pb.UpdateOrgIDPOIDCConfigRequest) *domain.OIDCIDPConfig {
	return &domain.OIDCIDPConfig{
		IDPConfigID:           req.IdpId,
		ClientID:              req.ClientId,
		ClientSecretString:    req.ClientSecret,
		Issuer:                req.Issuer,
		Scopes:                req.Scopes,
		IDPDisplayNameMapping: idp_grpc.MappingFieldToDomain(req.DisplayNameMapping),
		UsernameMapping:       idp_grpc.MappingFieldToDomain(req.UsernameMapping),
	}
}

func updateAuthConnectorConfigToDomain(req *mgmt_pb.UpdateOrgIDPAuthConnectorConfigRequest) *domain.AuthConnectorIDPConfig {
	return &domain.AuthConnectorIDPConfig{
		CommonIDPConfig: domain.CommonIDPConfig{
			IDPConfigID: req.IdpId,
		},
		BaseURL:    req.BaseUrl,
		ProviderID: req.ProviderId,
		MachineID:  req.MachineId,
	}
}

func listIDPsToModel(req *mgmt_pb.ListOrgIDPsRequest) *iam_model.IDPConfigSearchRequest {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &iam_model.IDPConfigSearchRequest{
		Offset:        offset,
		Limit:         limit,
		Asc:           asc,
		SortingColumn: idp_grpc.FieldNameToModel(req.SortingColumn),
		Queries:       idpQueriesToModel(req.Queries),
	}
}

func idpQueriesToModel(queries []*mgmt_pb.IDPQuery) []*iam_model.IDPConfigSearchQuery {
	q := make([]*iam_model.IDPConfigSearchQuery, len(queries))
	for i, query := range queries {
		q[i] = idpQueryToModel(query)
	}

	return q
}

func idpQueryToModel(query *mgmt_pb.IDPQuery) *iam_model.IDPConfigSearchQuery {
	switch q := query.Query.(type) {
	case *mgmt_pb.IDPQuery_IdpNameQuery:
		return idp_grpc.IDPNameQueryToModel(q.IdpNameQuery)
	case *mgmt_pb.IDPQuery_IdpIdQuery:
		return idp_grpc.IDPIDQueryToModel(q.IdpIdQuery)
	case *mgmt_pb.IDPQuery_OwnerTypeQuery:
		return idp_grpc.IDPOwnerTypeQueryToModel(q.OwnerTypeQuery)
	default:
		return nil
	}
}

func idpProviderViewsToDomain(idps []*iam_model.IDPProviderView) []*domain.IDPProvider {
	idpProvider := make([]*domain.IDPProvider, len(idps))
	for i, idp := range idps {
		idpProvider[i] = &domain.IDPProvider{
			ObjectRoot: models.ObjectRoot{
				AggregateID: idp.AggregateID,
			},
			IDPConfigID: idp.IDPConfigID,
			Type:        idpConfigTypeToDomain(idp.IDPProviderType),
		}
	}
	return idpProvider
}

func idpConfigTypeToDomain(idpType iam_model.IDPProviderType) domain.IdentityProviderType {
	switch idpType {
	case iam_model.IDPProviderTypeOrg:
		return domain.IdentityProviderTypeOrg
	default:
		return domain.IdentityProviderTypeSystem
	}
}

func externalIDPViewsToDomain(idps []*user_model.ExternalIDPView) []*domain.ExternalIDP {
	externalIDPs := make([]*domain.ExternalIDP, len(idps))
	for i, idp := range idps {
		externalIDPs[i] = &domain.ExternalIDP{
			ObjectRoot: models.ObjectRoot{
				AggregateID:   idp.UserID,
				ResourceOwner: idp.ResourceOwner,
			},
			IDPConfigID:    idp.IDPConfigID,
			ExternalUserID: idp.ExternalUserID,
			DisplayName:    idp.UserDisplayName,
		}
	}
	return externalIDPs
}
