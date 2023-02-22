package projection

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

const (
	IDPTemplateTable                 = "projections.idp_templates"
	IDPTemplateOAuthTable            = IDPTemplateTable + "_" + IDPTemplateOAuthSuffix
	IDPTemplateGitHubTable           = IDPTemplateTable + "_" + IDPTemplateGitHubSuffix
	IDPTemplateGitHubEnterpriseTable = IDPTemplateTable + "_" + IDPTemplateGitHubEnterpriseSuffix
	IDPTemplateGoogleTable           = IDPTemplateTable + "_" + IDPTemplateGoogleSuffix
	IDPTemplateLDAPTable             = IDPTemplateTable + "_" + IDPTemplateLDAPSuffix

	IDPTemplateOAuthSuffix            = "oauth"
	IDPTemplateGitHubSuffix           = "github"
	IDPTemplateGitHubEnterpriseSuffix = "github_enterprise"
	IDPTemplateGoogleSuffix           = "google"
	IDPTemplateLDAPSuffix             = "ldap"

	IDPTemplateIDCol                = "id"
	IDPTemplateCreationDateCol      = "creation_date"
	IDPTemplateChangeDateCol        = "change_date"
	IDPTemplateSequenceCol          = "sequence"
	IDPTemplateResourceOwnerCol     = "resource_owner"
	IDPTemplateInstanceIDCol        = "instance_id"
	IDPTemplateStateCol             = "state"
	IDPTemplateNameCol              = "name"
	IDPTemplateOwnerTypeCol         = "owner_type"
	IDPTemplateTypeCol              = "type"
	IDPTemplateOwnerRemovedCol      = "owner_removed"
	IDPTemplateIsCreationAllowedCol = "is_creation_allowed"
	IDPTemplateIsLinkingAllowedCol  = "is_linking_allowed"
	IDPTemplateIsAutoCreationCol    = "is_auto_creation"
	IDPTemplateIsAutoUpdateCol      = "is_auto_update"

	OAuthIDCol                    = "idp_id"
	OAuthInstanceIDCol            = "instance_id"
	OAuthClientIDCol              = "client_id"
	OAuthClientSecretCol          = "client_secret"
	OAuthAuthorizationEndpointCol = "authorization_endpoint"
	OAuthTokenEndpointCol         = "token_endpoint"
	OAuthUserEndpointCol          = "user_endpoint"
	OAuthScopesCol                = "scopes"

	GitHubIDCol           = "idp_id"
	GitHubInstanceIDCol   = "instance_id"
	GitHubClientIDCol     = "client_id"
	GitHubClientSecretCol = "client_secret"
	GitHubScopesCol       = "scopes"

	GitHubEnterpriseIDCol                    = "idp_id"
	GitHubEnterpriseInstanceIDCol            = "instance_id"
	GitHubEnterpriseClientIDCol              = "client_id"
	GitHubEnterpriseClientSecretCol          = "client_secret"
	GitHubEnterpriseAuthorizationEndpointCol = "authorization_endpoint"
	GitHubEnterpriseTokenEndpointCol         = "token_endpoint"
	GitHubEnterpriseUserEndpointCol          = "user_endpoint"
	GitHubEnterpriseScopesCol                = "scopes"

	GoogleIDCol           = "idp_id"
	GoogleInstanceIDCol   = "instance_id"
	GoogleClientIDCol     = "client_id"
	GoogleClientSecretCol = "client_secret"
	GoogleScopesCol       = "scopes"

	LDAPIDCol                         = "idp_id"
	LDAPInstanceIDCol                 = "instance_id"
	LDAPHostCol                       = "host"
	LDAPPortCol                       = "port"
	LDAPTlsCol                        = "tls"
	LDAPBaseDNCol                     = "base_dn"
	LDAPUserObjectClassCol            = "user_object_class"
	LDAPUserUniqueAttributeCol        = "user_unique_attribute"
	LDAPAdminCol                      = "admin"
	LDAPPasswordCol                   = "password"
	LDAPIDAttributeCol                = "id_attribute"
	LDAPFirstNameAttributeCol         = "first_name_attribute"
	LDAPLastNameAttributeCol          = "last_name_attribute"
	LDAPDisplayNameAttributeCol       = "display_name_attribute"
	LDAPNickNameAttributeCol          = "nick_name_attribute"
	LDAPPreferredUsernameAttributeCol = "preferred_username_attribute"
	LDAPEmailAttributeCol             = "email_attribute"
	LDAPEmailVerifiedAttributeCol     = "email_verified"
	LDAPPhoneAttributeCol             = "phone_attribute"
	LDAPPhoneVerifiedAttributeCol     = "phone_verified_attribute"
	LDAPPreferredLanguageAttributeCol = "preferred_language_attribute"
	LDAPAvatarURLAttributeCol         = "avatar_url_attribute"
	LDAPProfileAttributeCol           = "profile_attribute"
)

type idpTemplateProjection struct {
	crdb.StatementHandler
}

func newIDPTemplateProjection(ctx context.Context, config crdb.StatementHandlerConfig) *idpTemplateProjection {
	p := new(idpTemplateProjection)
	config.ProjectionName = IDPTemplateTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(IDPTemplateIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPTemplateCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(IDPTemplateChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(IDPTemplateSequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(IDPTemplateResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPTemplateInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPTemplateStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(IDPTemplateNameCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(IDPTemplateOwnerTypeCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(IDPTemplateTypeCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(IDPTemplateOwnerRemovedCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(IDPTemplateIsCreationAllowedCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(IDPTemplateIsLinkingAllowedCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(IDPTemplateIsAutoCreationCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(IDPTemplateIsAutoUpdateCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(IDPTemplateInstanceIDCol, IDPTemplateIDCol),
			crdb.WithIndex(crdb.NewIndex("resource_owner", []string{IDPTemplateResourceOwnerCol})),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{IDPTemplateOwnerRemovedCol})),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(OAuthIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(OAuthInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(OAuthClientIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(OAuthClientSecretCol, crdb.ColumnTypeJSONB),
			crdb.NewColumn(OAuthAuthorizationEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(OAuthTokenEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(OAuthUserEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(OAuthScopesCol, crdb.ColumnTypeTextArray, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(OAuthInstanceIDCol, OAuthIDCol),
			IDPTemplateOAuthSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(GitHubIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubClientIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubClientSecretCol, crdb.ColumnTypeJSONB),
			crdb.NewColumn(GitHubScopesCol, crdb.ColumnTypeTextArray, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(GitHubInstanceIDCol, GitHubIDCol),
			IDPTemplateGitHubSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(GitHubEnterpriseIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubEnterpriseInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubEnterpriseClientIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubEnterpriseClientSecretCol, crdb.ColumnTypeJSONB),
			crdb.NewColumn(GitHubEnterpriseAuthorizationEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubEnterpriseTokenEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubEnterpriseUserEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubEnterpriseScopesCol, crdb.ColumnTypeTextArray, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(OAuthInstanceIDCol, OAuthIDCol),
			IDPTemplateOAuthSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(GoogleIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GoogleInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GoogleClientIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GoogleClientSecretCol, crdb.ColumnTypeJSONB),
			crdb.NewColumn(GoogleScopesCol, crdb.ColumnTypeTextArray, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(GoogleInstanceIDCol, GoogleIDCol),
			IDPTemplateGoogleSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(LDAPIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(LDAPInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(LDAPHostCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPPortCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPTlsCol, crdb.ColumnTypeBool, crdb.Nullable()),
			crdb.NewColumn(LDAPBaseDNCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPUserObjectClassCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPUserUniqueAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPAdminCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPPasswordCol, crdb.ColumnTypeJSONB, crdb.Nullable()),
			crdb.NewColumn(LDAPIDAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPFirstNameAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPLastNameAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPDisplayNameAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPNickNameAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPPreferredUsernameAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPEmailAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPEmailVerifiedAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPPhoneAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPPhoneVerifiedAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPPreferredLanguageAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPAvatarURLAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPProfileAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(LDAPInstanceIDCol, LDAPIDCol),
			IDPTemplateLDAPSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *idpTemplateProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.OAuthIDPAddedEventType,
					Reduce: p.reduceOAuthIDPAdded,
				},
				{
					Event:  instance.OAuthIDPChangedEventType,
					Reduce: p.reduceOAuthIDPChanged,
				},
				{
					Event:  instance.GitHubIDPAddedEventType,
					Reduce: p.reduceGitHubIDPAdded,
				},
				{
					Event:  instance.GitHubIDPChangedEventType,
					Reduce: p.reduceGitHubIDPChanged,
				},
				{
					Event:  instance.GitHubEnterpriseIDPAddedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPAdded,
				},
				{
					Event:  instance.GitHubEnterpriseIDPChangedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPChanged,
				},
				{
					Event:  instance.GoogleIDPAddedEventType,
					Reduce: p.reduceGoogleIDPAdded,
				},
				{
					Event:  instance.GoogleIDPChangedEventType,
					Reduce: p.reduceGoogleIDPChanged,
				},
				{
					Event:  instance.LDAPIDPAddedEventType,
					Reduce: p.reduceLDAPIDPAdded,
				},
				{
					Event:  instance.LDAPIDPChangedEventType,
					Reduce: p.reduceLDAPIDPChanged,
				},
				{
					Event:  instance.IDPRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(IDPTemplateInstanceIDCol),
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OAuthIDPAddedEventType,
					Reduce: p.reduceOAuthIDPAdded,
				},
				{
					Event:  org.OAuthIDPChangedEventType,
					Reduce: p.reduceOAuthIDPChanged,
				},
				{
					Event:  org.GitHubEnterpriseIDPAddedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPAdded,
				},
				{
					Event:  org.GitHubEnterpriseIDPChangedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPChanged,
				},
				{
					Event:  org.GoogleIDPAddedEventType,
					Reduce: p.reduceGoogleIDPAdded,
				},
				{
					Event:  org.GoogleIDPChangedEventType,
					Reduce: p.reduceGoogleIDPChanged,
				},
				{
					Event:  org.LDAPIDPAddedEventType,
					Reduce: p.reduceLDAPIDPAdded,
				},
				{
					Event:  org.LDAPIDPChangedEventType,
					Reduce: p.reduceLDAPIDPChanged,
				},
				{
					Event:  org.IDPRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
	}
}

func (p *idpTemplateProjection) reduceOAuthIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OAuthIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.OAuthIDPAddedEvent:
		idpEvent = e.OAuthIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.OAuthIDPAddedEvent:
		idpEvent = e.OAuthIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ap9ihb", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPAddedEventType, instance.OAuthIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeOAuth),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(OAuthIDCol, idpEvent.ID),
				handler.NewCol(OAuthInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(OAuthClientIDCol, idpEvent.ClientID),
				handler.NewCol(OAuthClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(OAuthAuthorizationEndpointCol, idpEvent.AuthorizationEndpoint),
				handler.NewCol(OAuthTokenEndpointCol, idpEvent.TokenEndpoint),
				handler.NewCol(OAuthUserEndpointCol, idpEvent.UserEndpoint),
				handler.NewCol(OAuthScopesCol, database.StringArray(idpEvent.Scopes)),
			},
			crdb.WithTableSuffix(IDPTemplateOAuthSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceOAuthIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OAuthIDPChangedEvent
	switch e := event.(type) {
	case *org.OAuthIDPChangedEvent:
		idpEvent = e.OAuthIDPChangedEvent
	case *instance.OAuthIDPChangedEvent:
		idpEvent = e.OAuthIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPChangedEventType, instance.OAuthIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	oauthCols := reduceOAuthIDPChangedColumns(idpEvent)
	if len(oauthCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				oauthCols,
				[]handler.Condition{
					handler.NewCond(OAuthIDCol, idpEvent.ID),
					handler.NewCond(OAuthInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateOAuthSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceGitHubIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitHubIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.GitHubIDPAddedEvent:
		idpEvent = e.GitHubIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.GitHubIDPAddedEvent:
		idpEvent = e.GitHubIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPAddedEventType, instance.GitHubIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGitHub),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GitHubIDCol, idpEvent.ID),
				handler.NewCol(GitHubInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GitHubClientIDCol, idpEvent.ClientID),
				handler.NewCol(GitHubClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GitHubScopesCol, database.StringArray(idpEvent.Scopes)),
			},
			crdb.WithTableSuffix(IDPTemplateGitHubSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceGitHubEnterpriseIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitHubEnterpriseIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.GitHubEnterpriseIDPAddedEvent:
		idpEvent = e.GitHubEnterpriseIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.GitHubEnterpriseIDPAddedEvent:
		idpEvent = e.GitHubEnterpriseIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Sf3g2a", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubEnterpriseIDPAddedEventType, instance.GitHubEnterpriseIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGitHubEnterprise),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GitHubEnterpriseIDCol, idpEvent.ID),
				handler.NewCol(GitHubEnterpriseInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GitHubEnterpriseClientIDCol, idpEvent.ClientID),
				handler.NewCol(GitHubEnterpriseClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GitHubEnterpriseAuthorizationEndpointCol, idpEvent.AuthorizationEndpoint),
				handler.NewCol(GitHubEnterpriseTokenEndpointCol, idpEvent.TokenEndpoint),
				handler.NewCol(GitHubEnterpriseUserEndpointCol, idpEvent.UserEndpoint),
				handler.NewCol(GitHubEnterpriseScopesCol, database.StringArray(idpEvent.Scopes)),
			},
			crdb.WithTableSuffix(IDPTemplateGitHubEnterpriseSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceGitHubIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitHubIDPChangedEvent
	switch e := event.(type) {
	case *org.GitHubIDPChangedEvent:
		idpEvent = e.GitHubIDPChangedEvent
	case *instance.GitHubIDPChangedEvent:
		idpEvent = e.GitHubIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPChangedEventType, instance.GitHubIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	githubCols := reduceGitHubIDPChangedColumns(idpEvent)
	if len(githubCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				githubCols,
				[]handler.Condition{
					handler.NewCond(GitHubIDCol, idpEvent.ID),
					handler.NewCond(GitHubInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateGitHubSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceGitHubEnterpriseIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitHubEnterpriseIDPChangedEvent
	switch e := event.(type) {
	case *org.GitHubEnterpriseIDPChangedEvent:
		idpEvent = e.GitHubEnterpriseIDPChangedEvent
	case *instance.GitHubEnterpriseIDPChangedEvent:
		idpEvent = e.GitHubEnterpriseIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SDg3g", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubEnterpriseIDPChangedEventType, instance.GitHubEnterpriseIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	githubCols := reduceGitHubEnterpriseIDPChangedColumns(idpEvent)
	if len(githubCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				githubCols,
				[]handler.Condition{
					handler.NewCond(GitHubEnterpriseIDCol, idpEvent.ID),
					handler.NewCond(GitHubEnterpriseInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateGitHubEnterpriseSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceGoogleIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GoogleIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.GoogleIDPAddedEvent:
		idpEvent = e.GoogleIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.GoogleIDPAddedEvent:
		idpEvent = e.GoogleIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ap9ihb", "reduce.wrong.event.type %v", []eventstore.EventType{org.GoogleIDPAddedEventType, instance.GoogleIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGoogle),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GoogleIDCol, idpEvent.ID),
				handler.NewCol(GoogleInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GoogleClientIDCol, idpEvent.ClientID),
				handler.NewCol(GoogleClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GoogleScopesCol, database.StringArray(idpEvent.Scopes)),
			},
			crdb.WithTableSuffix(IDPTemplateGoogleSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceGoogleIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GoogleIDPChangedEvent
	switch e := event.(type) {
	case *org.GoogleIDPChangedEvent:
		idpEvent = e.GoogleIDPChangedEvent
	case *instance.GoogleIDPChangedEvent:
		idpEvent = e.GoogleIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GoogleIDPChangedEventType, instance.GoogleIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	googleCols := reduceGoogleIDPChangedColumns(idpEvent)
	if len(googleCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				googleCols,
				[]handler.Condition{
					handler.NewCond(GoogleIDCol, idpEvent.ID),
					handler.NewCond(GoogleInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateGoogleSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceLDAPIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.LDAPIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.LDAPIDPAddedEvent:
		idpEvent = e.LDAPIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.LDAPIDPAddedEvent:
		idpEvent = e.LDAPIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-9s02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPAddedEventType, instance.LDAPIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeLDAP),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(LDAPIDCol, idpEvent.ID),
				handler.NewCol(LDAPInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(LDAPHostCol, idpEvent.Host),
				handler.NewCol(LDAPPortCol, idpEvent.Port),
				handler.NewCol(LDAPTlsCol, idpEvent.TLS),
				handler.NewCol(LDAPBaseDNCol, idpEvent.BaseDN),
				handler.NewCol(LDAPUserObjectClassCol, idpEvent.UserObjectClass),
				handler.NewCol(LDAPUserUniqueAttributeCol, idpEvent.UserUniqueAttribute),
				handler.NewCol(LDAPAdminCol, idpEvent.Admin),
				handler.NewCol(LDAPPasswordCol, idpEvent.Password),
				handler.NewCol(LDAPIDAttributeCol, idpEvent.IDAttribute),
				handler.NewCol(LDAPFirstNameAttributeCol, idpEvent.FirstNameAttribute),
				handler.NewCol(LDAPLastNameAttributeCol, idpEvent.LastNameAttribute),
				handler.NewCol(LDAPDisplayNameAttributeCol, idpEvent.DisplayNameAttribute),
				handler.NewCol(LDAPNickNameAttributeCol, idpEvent.NickNameAttribute),
				handler.NewCol(LDAPPreferredUsernameAttributeCol, idpEvent.PreferredUsernameAttribute),
				handler.NewCol(LDAPEmailAttributeCol, idpEvent.EmailAttribute),
				handler.NewCol(LDAPEmailVerifiedAttributeCol, idpEvent.EmailVerifiedAttribute),
				handler.NewCol(LDAPPhoneAttributeCol, idpEvent.PhoneAttribute),
				handler.NewCol(LDAPPhoneVerifiedAttributeCol, idpEvent.PhoneVerifiedAttribute),
				handler.NewCol(LDAPPreferredLanguageAttributeCol, idpEvent.PreferredLanguageAttribute),
				handler.NewCol(LDAPAvatarURLAttributeCol, idpEvent.AvatarURLAttribute),
				handler.NewCol(LDAPProfileAttributeCol, idpEvent.ProfileAttribute),
			},
			crdb.WithTableSuffix(IDPTemplateLDAPSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceLDAPIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.LDAPIDPChangedEvent
	switch e := event.(type) {
	case *org.LDAPIDPChangedEvent:
		idpEvent = e.LDAPIDPChangedEvent
	case *instance.LDAPIDPChangedEvent:
		idpEvent = e.LDAPIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPChangedEventType, instance.LDAPIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)

	ldapCols := reduceLDAPIDPChangedColumns(idpEvent)
	if len(ldapCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				ldapCols,
				[]handler.Condition{
					handler.NewCond(LDAPIDCol, idpEvent.ID),
					handler.NewCond(LDAPInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateLDAPSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceIDPRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.RemovedEvent
	switch e := event.(type) {
	case *org.IDPRemovedEvent:
		idpEvent = e.RemovedEvent
	case *instance.IDPRemovedEvent:
		idpEvent = e.RemovedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-xbcvwin2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPRemovedEventType, instance.IDPRemovedEventType})
	}

	return crdb.NewDeleteStatement(
		&idpEvent,
		[]handler.Condition{
			handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpTemplateProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-Jp0D2K", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(IDPTemplateChangeDateCol, e.CreationDate()),
			handler.NewCol(IDPTemplateSequenceCol, e.Sequence()),
			handler.NewCol(IDPTemplateOwnerRemovedCol, true),
		},
		[]handler.Condition{
			handler.NewCond(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(IDPTemplateResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}

func reduceIDPChangedTemplateColumns(name *string, creationDate time.Time, sequence uint64, optionChanges idp.OptionChanges) []handler.Column {
	cols := make([]handler.Column, 0, 7)
	if name != nil {
		cols = append(cols, handler.NewCol(IDPTemplateNameCol, *name))
	}
	if optionChanges.IsCreationAllowed != nil {
		cols = append(cols, handler.NewCol(IDPTemplateIsCreationAllowedCol, *optionChanges.IsCreationAllowed))
	}
	if optionChanges.IsLinkingAllowed != nil {
		cols = append(cols, handler.NewCol(IDPTemplateIsLinkingAllowedCol, *optionChanges.IsLinkingAllowed))
	}
	if optionChanges.IsAutoCreation != nil {
		cols = append(cols, handler.NewCol(IDPTemplateIsAutoCreationCol, *optionChanges.IsAutoCreation))
	}
	if optionChanges.IsAutoUpdate != nil {
		cols = append(cols, handler.NewCol(IDPTemplateIsAutoUpdateCol, *optionChanges.IsAutoUpdate))
	}
	return append(cols,
		handler.NewCol(IDPTemplateChangeDateCol, creationDate),
		handler.NewCol(IDPTemplateSequenceCol, sequence),
	)
}

func reduceOAuthIDPChangedColumns(idpEvent idp.OAuthIDPChangedEvent) []handler.Column {
	oauthCols := make([]handler.Column, 0, 6)
	if idpEvent.ClientID != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.AuthorizationEndpoint != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthAuthorizationEndpointCol, *idpEvent.AuthorizationEndpoint))
	}
	if idpEvent.TokenEndpoint != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthTokenEndpointCol, *idpEvent.TokenEndpoint))
	}
	if idpEvent.UserEndpoint != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthUserEndpointCol, *idpEvent.UserEndpoint))
	}
	if idpEvent.Scopes != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthScopesCol, database.StringArray(idpEvent.Scopes)))
	}
	return oauthCols
}
func reduceGitHubIDPChangedColumns(idpEvent idp.GitHubIDPChangedEvent) []handler.Column {
	oauthCols := make([]handler.Column, 0, 3)
	if idpEvent.ClientID != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.Scopes != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubScopesCol, database.StringArray(idpEvent.Scopes)))
	}
	return oauthCols
}

func reduceGitHubEnterpriseIDPChangedColumns(idpEvent idp.GitHubEnterpriseIDPChangedEvent) []handler.Column {
	oauthCols := make([]handler.Column, 0, 6)
	if idpEvent.ClientID != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.AuthorizationEndpoint != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseAuthorizationEndpointCol, *idpEvent.AuthorizationEndpoint))
	}
	if idpEvent.TokenEndpoint != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseTokenEndpointCol, *idpEvent.TokenEndpoint))
	}
	if idpEvent.UserEndpoint != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseUserEndpointCol, *idpEvent.UserEndpoint))
	}
	if idpEvent.Scopes != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseScopesCol, database.StringArray(idpEvent.Scopes)))
	}
	return oauthCols
}

func reduceGoogleIDPChangedColumns(idpEvent idp.GoogleIDPChangedEvent) []handler.Column {
	googleCols := make([]handler.Column, 0, 3)
	if idpEvent.ClientID != nil {
		googleCols = append(googleCols, handler.NewCol(GoogleClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		googleCols = append(googleCols, handler.NewCol(GoogleClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.Scopes != nil {
		googleCols = append(googleCols, handler.NewCol(GoogleScopesCol, database.StringArray(idpEvent.Scopes)))
	}
	return googleCols
}

func reduceLDAPIDPChangedColumns(idpEvent idp.LDAPIDPChangedEvent) []handler.Column {
	ldapCols := make([]handler.Column, 0, 4)
	if idpEvent.Host != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPHostCol, *idpEvent.Host))
	}
	if idpEvent.Port != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPPortCol, *idpEvent.Port))
	}
	if idpEvent.TLS != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPTlsCol, *idpEvent.TLS))
	}
	if idpEvent.BaseDN != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPBaseDNCol, *idpEvent.BaseDN))
	}
	if idpEvent.UserObjectClass != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPUserObjectClassCol, *idpEvent.UserObjectClass))
	}
	if idpEvent.UserUniqueAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPUserUniqueAttributeCol, *idpEvent.UserUniqueAttribute))
	}
	if idpEvent.Admin != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPAdminCol, *idpEvent.Admin))
	}
	if idpEvent.Password != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPPasswordCol, *idpEvent.Password))
	}
	if idpEvent.IDAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPIDAttributeCol, *idpEvent.IDAttribute))
	}
	if idpEvent.FirstNameAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPFirstNameAttributeCol, *idpEvent.FirstNameAttribute))
	}
	if idpEvent.LastNameAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPLastNameAttributeCol, *idpEvent.LastNameAttribute))
	}
	if idpEvent.DisplayNameAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPDisplayNameAttributeCol, *idpEvent.DisplayNameAttribute))
	}
	if idpEvent.NickNameAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPNickNameAttributeCol, *idpEvent.NickNameAttribute))
	}
	if idpEvent.PreferredUsernameAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPPreferredUsernameAttributeCol, *idpEvent.PreferredUsernameAttribute))
	}
	if idpEvent.EmailAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPEmailAttributeCol, *idpEvent.EmailAttribute))
	}
	if idpEvent.EmailVerifiedAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPEmailVerifiedAttributeCol, *idpEvent.EmailVerifiedAttribute))
	}
	if idpEvent.PhoneAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPPhoneAttributeCol, *idpEvent.PhoneAttribute))
	}
	if idpEvent.PhoneVerifiedAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPPhoneVerifiedAttributeCol, *idpEvent.PhoneVerifiedAttribute))
	}
	if idpEvent.PreferredLanguageAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPPreferredLanguageAttributeCol, *idpEvent.PreferredLanguageAttribute))
	}
	if idpEvent.AvatarURLAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPAvatarURLAttributeCol, *idpEvent.AvatarURLAttribute))
	}
	if idpEvent.ProfileAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPProfileAttributeCol, *idpEvent.ProfileAttribute))
	}
	return ldapCols
}
