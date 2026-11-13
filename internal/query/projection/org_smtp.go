package projection

import (
	"context"

	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

const (
	OrgSMTPConfigProjectionTable = "projections.org_smtp_configs"
	OrgSMTPConfigTable           = OrgSMTPConfigProjectionTable + "_" + smtpConfigSMTPTableSuffix
	OrgSMTPConfigHTTPTable       = OrgSMTPConfigProjectionTable + "_" + smtpConfigHTTPTableSuffix

	OrgSMTPConfigColumnInstanceID    = "instance_id"
	OrgSMTPConfigColumnResourceOwner = "resource_owner"
	OrgSMTPConfigColumnAggregateID   = "aggregate_id"
	OrgSMTPConfigColumnID            = "id"
	OrgSMTPConfigColumnCreationDate  = "creation_date"
	OrgSMTPConfigColumnChangeDate    = "change_date"
	OrgSMTPConfigColumnSequence      = "sequence"
	OrgSMTPConfigColumnState         = "state"
	OrgSMTPConfigColumnDescription   = "description"

	OrgSMTPConfigSMTPColumnInstanceID                               = "instance_id"
	OrgSMTPConfigSMTPColumnID                                       = "id"
	OrgSMTPConfigSMTPColumnTLS                                      = "tls"
	OrgSMTPConfigSMTPColumnSenderAddress                            = "sender_address"
	OrgSMTPConfigSMTPColumnSenderName                               = "sender_name"
	OrgSMTPConfigSMTPColumnReplyToAddress                           = "reply_to_address"
	OrgSMTPConfigSMTPColumnHost                                     = "host"
	OrgSMTPConfigSMTPColumnUser                                     = "username"
	OrgSMTPConfigSMTPColumnPlainAuthPassword                        = "password"
	OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientId     = "xoauth2auth_client_credentials_client_id"
	OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientSecret = "xoauth2auth_client_credentials_client_secret"
	OrgSMTPConfigSMTPColumnXOAuth2AuthTokenEndpoint                 = "xoauth2auth_token_endpoint"
	OrgSMTPConfigSMTPColumnXOAuth2AuthScope                         = "xoauth2auth_scope"

	OrgSMTPConfigHTTPColumnInstanceID = "instance_id"
	OrgSMTPConfigHTTPColumnID         = "id"
	OrgSMTPConfigHTTPColumnEndpoint   = "endpoint"
	OrgSMTPConfigHTTPColumnSigningKey = "signing_key"
)

type orgSMTPConfigProjection struct{}

func newOrgSMTPConfigProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(orgSMTPConfigProjection))
}

func (*orgSMTPConfigProjection) Name() string {
	return OrgSMTPConfigProjectionTable
}

func (*orgSMTPConfigProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(OrgSMTPConfigColumnID, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigColumnAggregateID, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(OrgSMTPConfigColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(OrgSMTPConfigColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(OrgSMTPConfigColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigColumnDescription, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigColumnState, handler.ColumnTypeEnum),
		},
			handler.NewPrimaryKey(OrgSMTPConfigColumnInstanceID, OrgSMTPConfigColumnID),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(OrgSMTPConfigSMTPColumnID, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigSMTPColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigSMTPColumnTLS, handler.ColumnTypeBool),
			handler.NewColumn(OrgSMTPConfigSMTPColumnSenderAddress, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigSMTPColumnSenderName, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigSMTPColumnReplyToAddress, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigSMTPColumnHost, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigSMTPColumnUser, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigSMTPColumnPlainAuthPassword, handler.ColumnTypeJSONB, handler.Nullable()),
			handler.NewColumn(OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientId, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientSecret, handler.ColumnTypeJSONB, handler.Nullable()),
			handler.NewColumn(OrgSMTPConfigSMTPColumnXOAuth2AuthTokenEndpoint, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(OrgSMTPConfigSMTPColumnXOAuth2AuthScope, handler.ColumnTypeTextArray, handler.Nullable()),
		},
			handler.NewPrimaryKey(OrgSMTPConfigSMTPColumnInstanceID, OrgSMTPConfigSMTPColumnID),
			smtpConfigSMTPTableSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(OrgSMTPConfigHTTPColumnID, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigHTTPColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigHTTPColumnEndpoint, handler.ColumnTypeText),
			handler.NewColumn(OrgSMTPConfigHTTPColumnSigningKey, handler.ColumnTypeJSONB, handler.Nullable()),
		},
			handler.NewPrimaryKey(OrgSMTPConfigHTTPColumnInstanceID, OrgSMTPConfigHTTPColumnID),
			smtpConfigHTTPTableSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
	)
}

func (p *orgSMTPConfigProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgSMTPConfigAddedEventType,
					Reduce: p.reduceOrgSMTPConfigAdded,
				},
				{
					Event:  org.OrgSMTPConfigChangedEventType,
					Reduce: p.reduceOrgSMTPConfigChanged,
				},
				{
					Event:  org.OrgSMTPConfigPasswordChangedEventType,
					Reduce: p.reduceOrgSMTPConfigPasswordChanged,
				},
				{
					Event:  org.OrgSMTPConfigHTTPAddedEventType,
					Reduce: p.reduceOrgSMTPConfigHTTPAdded,
				},
				{
					Event:  org.OrgSMTPConfigHTTPChangedEventType,
					Reduce: p.reduceOrgSMTPConfigHTTPChanged,
				},
				{
					Event:  org.OrgSMTPConfigActivatedEventType,
					Reduce: p.reduceOrgSMTPConfigActivated,
				},
				{
					Event:  org.OrgSMTPConfigDeactivatedEventType,
					Reduce: p.reduceOrgSMTPConfigDeactivated,
				},
				{
					Event:  org.OrgSMTPConfigRemovedEventType,
					Reduce: p.reduceOrgSMTPConfigRemoved,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOrgSMTPOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(OrgSMTPConfigColumnInstanceID),
				},
			},
		},
	}
}
