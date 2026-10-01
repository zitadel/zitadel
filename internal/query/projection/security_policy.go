package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	SecurityPolicyProjectionTable             = "projections.security_policies3"
	SecurityPolicyColumnInstanceID            = "instance_id"
	SecurityPolicyColumnCreationDate          = "creation_date"
	SecurityPolicyColumnChangeDate            = "change_date"
	SecurityPolicyColumnSequence              = "sequence"
	SecurityPolicyColumnEnableIframeEmbedding = "enable_iframe_embedding"
	SecurityPolicyColumnAllowedOrigins        = "origins"
	SecurityPolicyColumnEnableImpersonation   = "enable_impersonation"

	SecurityPolicyColumnEnableDynamicClientRegistration               = "enable_dynamic_client_registration"
	SecurityPolicyColumnAllowUnauthenticatedDynamicClientRegistration = "allow_unauthenticated_dynamic_client_registration"

	SecurityPolicyColumnEnableClientIDMetadataDocument      = "enable_client_id_metadata_document"
	SecurityPolicyColumnClientIDMetadataDocumentAllowedURLs = "client_id_metadata_document_allowed_urls"
	SecurityPolicyColumnClientIDMetadataDocumentAllowAnyURL = "client_id_metadata_document_allow_any_url"
)

type securityPolicyProjection struct{}

func newSecurityPolicyProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(securityPolicyProjection))
}

func (*securityPolicyProjection) Name() string {
	return SecurityPolicyProjectionTable
}

func (*securityPolicyProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(SecurityPolicyColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SecurityPolicyColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SecurityPolicyColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(SecurityPolicyColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(SecurityPolicyColumnEnableIframeEmbedding, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(SecurityPolicyColumnAllowedOrigins, handler.ColumnTypeTextArray, handler.Nullable()),
			handler.NewColumn(SecurityPolicyColumnEnableImpersonation, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(SecurityPolicyColumnEnableDynamicClientRegistration, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(SecurityPolicyColumnAllowUnauthenticatedDynamicClientRegistration, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(SecurityPolicyColumnEnableClientIDMetadataDocument, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(SecurityPolicyColumnClientIDMetadataDocumentAllowedURLs, handler.ColumnTypeTextArray, handler.Nullable()),
			handler.NewColumn(SecurityPolicyColumnClientIDMetadataDocumentAllowAnyURL, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(SecurityPolicyColumnInstanceID),
		),
	)
}

func (p *securityPolicyProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.SecurityPolicySetEventType,
					Reduce: p.reduceSecurityPolicySet,
				},
				{
					Event:  instance.SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEventType,
					Reduce: p.reduceClientIDMetadataDocumentAllowedURLAdded,
				},
				{
					Event:  instance.SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEventType,
					Reduce: p.reduceClientIDMetadataDocumentAllowedURLRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(SecurityPolicyColumnInstanceID),
				},
			},
		},
	}
}

func (p *securityPolicyProjection) reduceSecurityPolicySet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SecurityPolicySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-D3g87", "reduce.wrong.event.type %s", instance.SecurityPolicySetEventType)
	}
	changes := []handler.Column{
		handler.NewCol(SecurityPolicyColumnCreationDate, handler.OnlySetValueOnInsert(SecurityPolicyProjectionTable, e.CreationDate())),
		handler.NewCol(SecurityPolicyColumnChangeDate, e.CreationDate()),
		handler.NewCol(SecurityPolicyColumnInstanceID, e.Aggregate().InstanceID),
		handler.NewCol(SecurityPolicyColumnSequence, e.Sequence()),
	}
	if e.EnableIframeEmbedding != nil {
		changes = append(changes, handler.NewCol(SecurityPolicyColumnEnableIframeEmbedding, *e.EnableIframeEmbedding))
	} else if e.Enabled != nil {
		changes = append(changes, handler.NewCol(SecurityPolicyColumnEnableIframeEmbedding, *e.Enabled))
	}
	if e.AllowedOrigins != nil {
		changes = append(changes, handler.NewCol(SecurityPolicyColumnAllowedOrigins, e.AllowedOrigins))
	}
	if e.EnableImpersonation != nil {
		changes = append(changes, handler.NewCol(SecurityPolicyColumnEnableImpersonation, e.EnableImpersonation))
	}
	if e.EnableDynamicClientRegistration != nil {
		changes = append(changes, handler.NewCol(SecurityPolicyColumnEnableDynamicClientRegistration, e.EnableDynamicClientRegistration))
	}
	if e.AllowUnauthenticatedDynamicClientRegistration != nil {
		changes = append(changes, handler.NewCol(SecurityPolicyColumnAllowUnauthenticatedDynamicClientRegistration, e.AllowUnauthenticatedDynamicClientRegistration))
	}
	if e.EnableClientIDMetadataDocument != nil {
		changes = append(changes, handler.NewCol(SecurityPolicyColumnEnableClientIDMetadataDocument, e.EnableClientIDMetadataDocument))
	}
	if e.ClientIDMetadataDocumentAllowedURLs != nil {
		changes = append(changes, handler.NewCol(SecurityPolicyColumnClientIDMetadataDocumentAllowedURLs, e.ClientIDMetadataDocumentAllowedURLs))
	}
	if e.ClientIDMetadataDocumentAllowAnyURL != nil {
		changes = append(changes, handler.NewCol(SecurityPolicyColumnClientIDMetadataDocumentAllowAnyURL, e.ClientIDMetadataDocumentAllowAnyURL))
	}
	return handler.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(SecurityPolicyColumnInstanceID, ""),
		},
		changes,
	), nil
}

func (p *securityPolicyProjection) reduceClientIDMetadataDocumentAllowedURLAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Jx4mW", "reduce.wrong.event.type %s", instance.SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEventType)
	}
	return p.reduceClientIDMetadataDocumentAllowedURLs(e, handler.NewArrayAppendCol(SecurityPolicyColumnClientIDMetadataDocumentAllowedURLs, e.URL)), nil
}

func (p *securityPolicyProjection) reduceClientIDMetadataDocumentAllowedURLRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-c8Rq2", "reduce.wrong.event.type %s", instance.SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEventType)
	}
	return p.reduceClientIDMetadataDocumentAllowedURLs(e, handler.NewArrayRemoveCol(SecurityPolicyColumnClientIDMetadataDocumentAllowedURLs, e.URL)), nil
}

// reduceClientIDMetadataDocumentAllowedURLs applies an array change to the allowed client_id
// URLs. The policy row only exists once the policy was set, so it is upserted first and the
// change is applied to it afterwards.
func (p *securityPolicyProjection) reduceClientIDMetadataDocumentAllowedURLs(e eventstore.Event, change handler.Column) *handler.Statement {
	return handler.NewMultiStatement(
		e,
		handler.AddUpsertStatement(
			[]handler.Column{
				handler.NewCol(SecurityPolicyColumnInstanceID, ""),
			},
			[]handler.Column{
				handler.NewCol(SecurityPolicyColumnCreationDate, handler.OnlySetValueOnInsert(SecurityPolicyProjectionTable, e.CreatedAt())),
				handler.NewCol(SecurityPolicyColumnChangeDate, e.CreatedAt()),
				handler.NewCol(SecurityPolicyColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SecurityPolicyColumnSequence, e.Sequence()),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{change},
			[]handler.Condition{
				handler.NewCond(SecurityPolicyColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	)
}
