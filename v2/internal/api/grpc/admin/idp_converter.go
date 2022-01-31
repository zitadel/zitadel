package admin

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/query"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
	idp_grpc "github.com/caos/zitadel/v2/internal/api/grpc/idp"
	"github.com/caos/zitadel/v2/internal/api/grpc/object"
)

func addOIDCIDPRequestToDomain(req *admin_pb.AddOIDCIDPRequest) *domain.IDPConfig {
	return &domain.IDPConfig{
		Name:         req.Name,
		OIDCConfig:   addOIDCIDPRequestToDomainOIDCIDPConfig(req),
		StylingType:  idp_grpc.IDPStylingTypeToDomain(req.StylingType),
		Type:         domain.IDPConfigTypeOIDC,
		AutoRegister: req.AutoRegister,
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

func addJWTIDPRequestToDomain(req *admin_pb.AddJWTIDPRequest) *domain.IDPConfig {
	return &domain.IDPConfig{
		Name:         req.Name,
		JWTConfig:    addJWTIDPRequestToDomainJWTIDPConfig(req),
		StylingType:  idp_grpc.IDPStylingTypeToDomain(req.StylingType),
		Type:         domain.IDPConfigTypeJWT,
		AutoRegister: req.AutoRegister,
	}
}

func addJWTIDPRequestToDomainJWTIDPConfig(req *admin_pb.AddJWTIDPRequest) *domain.JWTIDPConfig {
	return &domain.JWTIDPConfig{
		JWTEndpoint:  req.JwtEndpoint,
		Issuer:       req.Issuer,
		KeysEndpoint: req.KeysEndpoint,
		HeaderName:   req.HeaderName,
	}
}

func updateIDPToDomain(req *admin_pb.UpdateIDPRequest) *domain.IDPConfig {
	return &domain.IDPConfig{
		IDPConfigID:  req.IdpId,
		Name:         req.Name,
		StylingType:  idp_grpc.IDPStylingTypeToDomain(req.StylingType),
		AutoRegister: req.AutoRegister,
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

func updateJWTConfigToDomain(req *admin_pb.UpdateIDPJWTConfigRequest) *domain.JWTIDPConfig {
	return &domain.JWTIDPConfig{
		IDPConfigID:  req.IdpId,
		JWTEndpoint:  req.JwtEndpoint,
		Issuer:       req.Issuer,
		KeysEndpoint: req.KeysEndpoint,
		HeaderName:   req.HeaderName,
	}
}

func listIDPsToModel(req *admin_pb.ListIDPsRequest) (*query.IDPSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := idpQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	iamQuery, err := query.NewIDPResourceOwnerSearchQuery(domain.IAMID)
	if err != nil {
		return nil, err
	}
	queries = append(queries, iamQuery)
	return &query.IDPSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: idp_grpc.FieldNameToModel(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func idpQueriesToModel(queries []*admin_pb.IDPQuery) (q []query.SearchQuery, err error) {
	q = make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = idpQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}

	return q, nil
}

func idpQueryToModel(idpQuery *admin_pb.IDPQuery) (query.SearchQuery, error) {
	switch q := idpQuery.Query.(type) {
	case *admin_pb.IDPQuery_IdpNameQuery:
		return query.NewIDPNameSearchQuery(object.TextMethodToQuery(q.IdpNameQuery.Method), q.IdpNameQuery.Name)
	case *admin_pb.IDPQuery_IdpIdQuery:
		return query.NewIDPIDSearchQuery(q.IdpIdQuery.Id)
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ADMIN-VmqQu", "List.Query.Invalid")
	}
}

func idpsToDomain(idps []*query.IDP) []*domain.IDPProvider {
	idpProvider := make([]*domain.IDPProvider, len(idps))
	for i, idp := range idps {
		idpProvider[i] = &domain.IDPProvider{
			ObjectRoot: models.ObjectRoot{
				AggregateID: idp.ResourceOwner,
			},
			IDPConfigID: idp.ID,
			Type:        idp.OwnerType,
		}
	}
	return idpProvider
}

func idpUserLinksToDomain(idps []*query.IDPUserLink) []*domain.UserIDPLink {
	externalIDPs := make([]*domain.UserIDPLink, len(idps))
	for i, idp := range idps {
		externalIDPs[i] = &domain.UserIDPLink{
			ObjectRoot: models.ObjectRoot{
				AggregateID:   idp.UserID,
				ResourceOwner: idp.ResourceOwner,
			},
			IDPConfigID:    idp.IDPID,
			ExternalUserID: idp.ProvidedUserID,
			DisplayName:    idp.ProvidedUsername,
		}
	}
	return externalIDPs
}
