package admin

import (
	idp_grpc "github.com/caos/zitadel/internal/api/grpc/idp"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	user_model "github.com/caos/zitadel/internal/user/model"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func addOIDCIDPRequestToDomain(req *admin_pb.AddOIDCIDPRequest) *domain.IDPConfig {
	return &domain.IDPConfig{
		Name:        req.Name,
		OIDCConfig:  addOIDCIDPRequestToDomainOIDCIDPConfig(req),
		StylingType: idp_grpc.IDPStylingTypeToDomain(req.StylingType),
		Type:        domain.IDPConfigTypeOIDC,
	}
}

func addOIDCIDPRequestToDomainOIDCIDPConfig(req *admin_pb.AddOIDCIDPRequest) *domain.OIDCIDPConfig {
	return &domain.OIDCIDPConfig{
		ClientID:              req.ClientId,
		ClientSecretString:    req.ClientSecret,
		Issuer:                req.Issuer,
		Scopes:                req.Scopes,
		IDPDisplayNameMapping: idp_grpc.MappingFieldToDomain(req.DisplayNameMapping),
		UsernameMapping:       idp_grpc.MappingFieldToDomain(req.UsernameMapping),
	}
}

func updateIDPToDomain(req *admin_pb.UpdateIDPRequest) *domain.IDPConfig {
	return &domain.IDPConfig{
		IDPConfigID: req.IdpId,
		Name:        req.Name,
		StylingType: idp_grpc.IDPStylingTypeToDomain(req.StylingType),
	}
}

func updateOIDCConfigToDomain(req *admin_pb.UpdateIDPOIDCConfigRequest) *domain.OIDCIDPConfig {
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

func listIDPsToModel(req *admin_pb.ListIDPsRequest) *iam_model.IDPConfigSearchRequest {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &iam_model.IDPConfigSearchRequest{
		Offset:        offset,
		Limit:         limit,
		Asc:           asc,
		SortingColumn: idp_grpc.FieldNameToModel(req.SortingColumn),
		Queries:       idpQueriesToModel(req.Queries),
	}
}

func idpQueriesToModel(queries []*admin_pb.IDPQuery) []*iam_model.IDPConfigSearchQuery {
	q := make([]*iam_model.IDPConfigSearchQuery, len(queries))
	for i, query := range queries {
		q[i] = idpQueryToModel(query)
	}

	return q
}

func idpQueryToModel(query *admin_pb.IDPQuery) *iam_model.IDPConfigSearchQuery {
	switch q := query.Query.(type) {
	case *admin_pb.IDPQuery_IdpNameQuery:
		return idp_grpc.IDPNameQueryToModel(q.IdpNameQuery)
	case *admin_pb.IDPQuery_IdpIdQuery:
		return idp_grpc.IDPIDQueryToModel(q.IdpIdQuery)
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
