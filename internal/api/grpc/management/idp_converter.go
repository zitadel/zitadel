package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	idp_grpc "github.com/caos/zitadel/internal/api/grpc/idp"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/query"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func addOIDCIDPRequestToDomain(req *mgmt_pb.AddOrgOIDCIDPRequest) *domain.IDPConfig {
	return &domain.IDPConfig{
		Name:         req.Name,
		OIDCConfig:   addOIDCIDPRequestToDomainOIDCIDPConfig(req),
		StylingType:  idp_grpc.IDPStylingTypeToDomain(req.StylingType),
		Type:         domain.IDPConfigTypeOIDC,
		AutoRegister: req.AutoRegister,
	}
}

func addOIDCIDPRequestToDomainOIDCIDPConfig(req *mgmt_pb.AddOrgOIDCIDPRequest) *domain.OIDCIDPConfig {
	return &domain.OIDCIDPConfig{
		ClientID:              req.ClientId,
		ClientSecretString:    req.ClientSecret,
		Issuer:                req.Issuer,
		Scopes:                req.Scopes,
		IDPDisplayNameMapping: idp_grpc.MappingFieldToDomain(req.DisplayNameMapping),
		UsernameMapping:       idp_grpc.MappingFieldToDomain(req.UsernameMapping),
	}
}

func addJWTIDPRequestToDomain(req *mgmt_pb.AddOrgJWTIDPRequest) *domain.IDPConfig {
	return &domain.IDPConfig{
		Name:         req.Name,
		JWTConfig:    addJWTIDPRequestToDomainJWTIDPConfig(req),
		StylingType:  idp_grpc.IDPStylingTypeToDomain(req.StylingType),
		Type:         domain.IDPConfigTypeJWT,
		AutoRegister: req.AutoRegister,
	}
}

func addJWTIDPRequestToDomainJWTIDPConfig(req *mgmt_pb.AddOrgJWTIDPRequest) *domain.JWTIDPConfig {
	return &domain.JWTIDPConfig{
		JWTEndpoint:  req.JwtEndpoint,
		Issuer:       req.Issuer,
		KeysEndpoint: req.KeysEndpoint,
		HeaderName:   req.HeaderName,
	}
}

func updateIDPToDomain(req *mgmt_pb.UpdateOrgIDPRequest) *domain.IDPConfig {
	return &domain.IDPConfig{
		IDPConfigID:  req.IdpId,
		Name:         req.Name,
		StylingType:  idp_grpc.IDPStylingTypeToDomain(req.StylingType),
		AutoRegister: req.AutoRegister,
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

func updateJWTConfigToDomain(req *mgmt_pb.UpdateOrgIDPJWTConfigRequest) *domain.JWTIDPConfig {
	return &domain.JWTIDPConfig{
		IDPConfigID:  req.IdpId,
		JWTEndpoint:  req.JwtEndpoint,
		Issuer:       req.Issuer,
		KeysEndpoint: req.KeysEndpoint,
		HeaderName:   req.HeaderName,
	}
}

func listIDPsToModel(ctx context.Context, req *mgmt_pb.ListOrgIDPsRequest) (queries *query.IDPSearchQueries, err error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	q, err := idpQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	resourceOwnerQuery, err := query.NewIDPResourceOwnerListSearchQuery(authz.GetInstance(ctx).ID, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	q = append(q, resourceOwnerQuery)
	return &query.IDPSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: idp_grpc.FieldNameToModel(req.SortingColumn),
		},
		Queries: q,
	}, nil
}

func idpQueriesToModel(queries []*mgmt_pb.IDPQuery) (q []query.SearchQuery, err error) {
	q = make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = idpQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}

	return q, nil
}

func idpQueryToModel(idpQuery *mgmt_pb.IDPQuery) (query.SearchQuery, error) {
	switch q := idpQuery.Query.(type) {
	case *mgmt_pb.IDPQuery_IdpNameQuery:
		return query.NewIDPNameSearchQuery(object.TextMethodToQuery(q.IdpNameQuery.Method), q.IdpNameQuery.Name)
	case *mgmt_pb.IDPQuery_IdpIdQuery:
		return query.NewIDPIDSearchQuery(q.IdpIdQuery.Id)
	case *mgmt_pb.IDPQuery_OwnerTypeQuery:
		return query.NewIDPOwnerTypeSearchQuery(idp_grpc.IDPProviderTypeFromPb(q.OwnerTypeQuery.OwnerType))
	default:
		return nil, errors.ThrowInvalidArgument(nil, "MANAG-WtLPV", "List.Query.Invalid")
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

func userLinksToDomain(idps []*query.IDPUserLink) []*domain.UserIDPLink {
	links := make([]*domain.UserIDPLink, len(idps))
	for i, idp := range idps {
		links[i] = &domain.UserIDPLink{
			ObjectRoot: models.ObjectRoot{
				AggregateID:   idp.UserID,
				ResourceOwner: idp.ResourceOwner,
			},
			IDPConfigID:    idp.IDPID,
			ExternalUserID: idp.ProvidedUserID,
			DisplayName:    idp.ProvidedUsername,
		}
	}
	return links
}
