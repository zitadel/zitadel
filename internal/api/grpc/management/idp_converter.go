package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	idp_grpc "github.com/zitadel/zitadel/internal/api/grpc/idp"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func AddOIDCIDPRequestToDomain(req *mgmt_pb.AddOrgOIDCIDPRequest) *domain.IDPConfig {
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

func AddJWTIDPRequestToDomain(req *mgmt_pb.AddOrgJWTIDPRequest) *domain.IDPConfig {
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
	resourceOwnerQuery, err := query.NewIDPResourceOwnerListSearchQuery(authz.GetInstance(ctx).InstanceID(), authz.GetCtxData(ctx).OrgID)
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

func listProvidersToQuery(ctx context.Context, req *mgmt_pb.ListProvidersRequest) (*query.IDPTemplateSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := providerQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	resourceOwnerQuery, err := query.NewIDPTemplateResourceOwnerListSearchQuery(authz.GetInstance(ctx).InstanceID(), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	queries = append(queries, resourceOwnerQuery)
	return &query.IDPTemplateSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

func providerQueriesToQuery(queries []*mgmt_pb.ProviderQuery) (q []query.SearchQuery, err error) {
	q = make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = providerQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}

	return q, nil
}

func providerQueryToQuery(idpQuery *mgmt_pb.ProviderQuery) (query.SearchQuery, error) {
	switch q := idpQuery.Query.(type) {
	case *mgmt_pb.ProviderQuery_IdpNameQuery:
		return query.NewIDPTemplateNameSearchQuery(object.TextMethodToQuery(q.IdpNameQuery.Method), q.IdpNameQuery.Name)
	case *mgmt_pb.ProviderQuery_IdpIdQuery:
		return query.NewIDPTemplateIDSearchQuery(q.IdpIdQuery.Id)
	case *mgmt_pb.ProviderQuery_OwnerTypeQuery:
		return query.NewIDPTemplateOwnerTypeSearchQuery(idp_grpc.IDPProviderTypeFromPb(q.OwnerTypeQuery.OwnerType))
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ORG-Dr2aa", "List.Query.Invalid")
	}
}

func addGenericOAuthProviderToCommand(req *mgmt_pb.AddGenericOAuthProviderRequest) command.GenericOAuthProvider {
	return command.GenericOAuthProvider{
		Name:                  req.Name,
		ClientID:              req.ClientId,
		ClientSecret:          req.ClientSecret,
		AuthorizationEndpoint: req.AuthorizationEndpoint,
		TokenEndpoint:         req.TokenEndpoint,
		UserEndpoint:          req.UserEndpoint,
		Scopes:                req.Scopes,
		IDAttribute:           req.IdAttribute,
		IDPOptions:            idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateGenericOAuthProviderToCommand(req *mgmt_pb.UpdateGenericOAuthProviderRequest) command.GenericOAuthProvider {
	return command.GenericOAuthProvider{
		Name:                  req.Name,
		ClientID:              req.ClientId,
		ClientSecret:          req.ClientSecret,
		AuthorizationEndpoint: req.AuthorizationEndpoint,
		TokenEndpoint:         req.TokenEndpoint,
		UserEndpoint:          req.UserEndpoint,
		Scopes:                req.Scopes,
		IDAttribute:           req.IdAttribute,
		IDPOptions:            idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func addGenericOIDCProviderToCommand(req *mgmt_pb.AddGenericOIDCProviderRequest) command.GenericOIDCProvider {
	return command.GenericOIDCProvider{
		Name:         req.Name,
		Issuer:       req.Issuer,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateGenericOIDCProviderToCommand(req *mgmt_pb.UpdateGenericOIDCProviderRequest) command.GenericOIDCProvider {
	return command.GenericOIDCProvider{
		Name:         req.Name,
		Issuer:       req.Issuer,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func addJWTProviderToCommand(req *mgmt_pb.AddJWTProviderRequest) command.JWTProvider {
	return command.JWTProvider{
		Name:        req.Name,
		Issuer:      req.Issuer,
		JWTEndpoint: req.JwtEndpoint,
		KeyEndpoint: req.KeysEndpoint,
		HeaderName:  req.HeaderName,
		IDPOptions:  idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateJWTProviderToCommand(req *mgmt_pb.UpdateJWTProviderRequest) command.JWTProvider {
	return command.JWTProvider{
		Name:        req.Name,
		Issuer:      req.Issuer,
		JWTEndpoint: req.JwtEndpoint,
		KeyEndpoint: req.KeysEndpoint,
		HeaderName:  req.HeaderName,
		IDPOptions:  idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func addGitHubProviderToCommand(req *mgmt_pb.AddGitHubProviderRequest) command.GitHubProvider {
	return command.GitHubProvider{
		Name:         req.Name,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateGitHubProviderToCommand(req *mgmt_pb.UpdateGitHubProviderRequest) command.GitHubProvider {
	return command.GitHubProvider{
		Name:         req.Name,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func addGitHubEnterpriseProviderToCommand(req *mgmt_pb.AddGitHubEnterpriseServerProviderRequest) command.GitHubEnterpriseProvider {
	return command.GitHubEnterpriseProvider{
		Name:                  req.Name,
		ClientID:              req.ClientId,
		ClientSecret:          req.ClientSecret,
		AuthorizationEndpoint: req.AuthorizationEndpoint,
		TokenEndpoint:         req.TokenEndpoint,
		UserEndpoint:          req.UserEndpoint,
		Scopes:                req.Scopes,
		IDPOptions:            idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateGitHubEnterpriseProviderToCommand(req *mgmt_pb.UpdateGitHubEnterpriseServerProviderRequest) command.GitHubEnterpriseProvider {
	return command.GitHubEnterpriseProvider{
		Name:                  req.Name,
		ClientID:              req.ClientId,
		ClientSecret:          req.ClientSecret,
		AuthorizationEndpoint: req.AuthorizationEndpoint,
		TokenEndpoint:         req.TokenEndpoint,
		UserEndpoint:          req.UserEndpoint,
		Scopes:                req.Scopes,
		IDPOptions:            idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func addGoogleProviderToCommand(req *mgmt_pb.AddGoogleProviderRequest) command.GoogleProvider {
	return command.GoogleProvider{
		Name:         req.Name,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateGoogleProviderToCommand(req *mgmt_pb.UpdateGoogleProviderRequest) command.GoogleProvider {
	return command.GoogleProvider{
		Name:         req.Name,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func addLDAPProviderToCommand(req *mgmt_pb.AddLDAPProviderRequest) command.LDAPProvider {
	return command.LDAPProvider{
		Name:                req.Name,
		Host:                req.Host,
		Port:                req.Port,
		TLS:                 req.Tls,
		BaseDN:              req.BaseDn,
		UserObjectClass:     req.UserObjectClass,
		UserUniqueAttribute: req.UserUniqueAttribute,
		Admin:               req.Admin,
		Password:            req.Password,
		LDAPAttributes:      idp_grpc.LDAPAttributesToCommand(req.Attributes),
		IDPOptions:          idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateLDAPProviderToCommand(req *mgmt_pb.UpdateLDAPProviderRequest) command.LDAPProvider {
	return command.LDAPProvider{
		Name:                req.Name,
		Host:                req.Host,
		Port:                req.Port,
		TLS:                 req.Tls,
		BaseDN:              req.BaseDn,
		UserObjectClass:     req.UserObjectClass,
		UserUniqueAttribute: req.UserUniqueAttribute,
		Admin:               req.Admin,
		Password:            req.Password,
		LDAPAttributes:      idp_grpc.LDAPAttributesToCommand(req.Attributes),
		IDPOptions:          idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}
