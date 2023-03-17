package admin

import (
	idp_grpc "github.com/zitadel/zitadel/internal/api/grpc/idp"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
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

func listIDPsToModel(instanceID string, req *admin_pb.ListIDPsRequest) (*query.IDPSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := idpQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	iamQuery, err := query.NewIDPResourceOwnerSearchQuery(instanceID)
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

func listProvidersToQuery(instanceID string, req *admin_pb.ListProvidersRequest) (*query.IDPTemplateSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := providerQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	iamQuery, err := query.NewIDPTemplateResourceOwnerSearchQuery(instanceID)
	if err != nil {
		return nil, err
	}
	queries = append(queries, iamQuery)
	return &query.IDPTemplateSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

func providerQueriesToQuery(queries []*admin_pb.ProviderQuery) (q []query.SearchQuery, err error) {
	q = make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = providerQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}

	return q, nil
}

func providerQueryToQuery(idpQuery *admin_pb.ProviderQuery) (query.SearchQuery, error) {
	switch q := idpQuery.Query.(type) {
	case *admin_pb.ProviderQuery_IdpNameQuery:
		return query.NewIDPTemplateNameSearchQuery(object.TextMethodToQuery(q.IdpNameQuery.Method), q.IdpNameQuery.Name)
	case *admin_pb.ProviderQuery_IdpIdQuery:
		return query.NewIDPTemplateIDSearchQuery(q.IdpIdQuery.Id)
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ADMIN-Dr2aa", "List.Query.Invalid")
	}
}

func addGenericOAuthProviderToCommand(req *admin_pb.AddGenericOAuthProviderRequest) command.GenericOAuthProvider {
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

func updateGenericOAuthProviderToCommand(req *admin_pb.UpdateGenericOAuthProviderRequest) command.GenericOAuthProvider {
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

func addGenericOIDCProviderToCommand(req *admin_pb.AddGenericOIDCProviderRequest) command.GenericOIDCProvider {
	return command.GenericOIDCProvider{
		Name:         req.Name,
		Issuer:       req.Issuer,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateGenericOIDCProviderToCommand(req *admin_pb.UpdateGenericOIDCProviderRequest) command.GenericOIDCProvider {
	return command.GenericOIDCProvider{
		Name:         req.Name,
		Issuer:       req.Issuer,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func addJWTProviderToCommand(req *admin_pb.AddJWTProviderRequest) command.JWTProvider {
	return command.JWTProvider{
		Name:        req.Name,
		Issuer:      req.Issuer,
		JWTEndpoint: req.JwtEndpoint,
		KeyEndpoint: req.KeysEndpoint,
		HeaderName:  req.HeaderName,
		IDPOptions:  idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateJWTProviderToCommand(req *admin_pb.UpdateJWTProviderRequest) command.JWTProvider {
	return command.JWTProvider{
		Name:        req.Name,
		Issuer:      req.Issuer,
		JWTEndpoint: req.JwtEndpoint,
		KeyEndpoint: req.KeysEndpoint,
		HeaderName:  req.HeaderName,
		IDPOptions:  idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func addAzureADProviderToCommand(req *admin_pb.AddAzureADProviderRequest) command.AzureADProvider {
	return command.AzureADProvider{
		Name:          req.Name,
		ClientID:      req.ClientId,
		ClientSecret:  req.ClientSecret,
		Scopes:        req.Scopes,
		Tenant:        idp_grpc.AzureADTenantToCommand(req.Tenant),
		EmailVerified: req.EmailVerified,
		IDPOptions:    idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateAzureADProviderToCommand(req *admin_pb.UpdateAzureADProviderRequest) command.AzureADProvider {
	return command.AzureADProvider{
		Name:          req.Name,
		ClientID:      req.ClientId,
		ClientSecret:  req.ClientSecret,
		Scopes:        req.Scopes,
		Tenant:        idp_grpc.AzureADTenantToCommand(req.Tenant),
		EmailVerified: req.EmailVerified,
		IDPOptions:    idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func addGitHubProviderToCommand(req *admin_pb.AddGitHubProviderRequest) command.GitHubProvider {
	return command.GitHubProvider{
		Name:         req.Name,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateGitHubProviderToCommand(req *admin_pb.UpdateGitHubProviderRequest) command.GitHubProvider {
	return command.GitHubProvider{
		Name:         req.Name,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func addGitHubEnterpriseProviderToCommand(req *admin_pb.AddGitHubEnterpriseServerProviderRequest) command.GitHubEnterpriseProvider {
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

func updateGitHubEnterpriseProviderToCommand(req *admin_pb.UpdateGitHubEnterpriseServerProviderRequest) command.GitHubEnterpriseProvider {
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

func addGitLabProviderToCommand(req *admin_pb.AddGitLabProviderRequest) command.GitLabProvider {
	return command.GitLabProvider{
		Name:         req.Name,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateGitLabProviderToCommand(req *admin_pb.UpdateGitLabProviderRequest) command.GitLabProvider {
	return command.GitLabProvider{
		Name:         req.Name,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func addGitLabSelfHostedProviderToCommand(req *admin_pb.AddGitLabSelfHostedProviderRequest) command.GitLabSelfHostedProvider {
	return command.GitLabSelfHostedProvider{
		Name:         req.Name,
		Issuer:       req.Issuer,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateGitLabSelfHostedProviderToCommand(req *admin_pb.UpdateGitLabSelfHostedProviderRequest) command.GitLabSelfHostedProvider {
	return command.GitLabSelfHostedProvider{
		Name:         req.Name,
		Issuer:       req.Issuer,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func addGoogleProviderToCommand(req *admin_pb.AddGoogleProviderRequest) command.GoogleProvider {
	return command.GoogleProvider{
		Name:         req.Name,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func updateGoogleProviderToCommand(req *admin_pb.UpdateGoogleProviderRequest) command.GoogleProvider {
	return command.GoogleProvider{
		Name:         req.Name,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		Scopes:       req.Scopes,
		IDPOptions:   idp_grpc.OptionsToCommand(req.ProviderOptions),
	}
}

func addLDAPProviderToCommand(req *admin_pb.AddLDAPProviderRequest) command.LDAPProvider {
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

func updateLDAPProviderToCommand(req *admin_pb.UpdateLDAPProviderRequest) command.LDAPProvider {
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
