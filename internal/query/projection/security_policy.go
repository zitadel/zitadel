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
	SecurityPolicyProjectionTable             = "projections.security_policies2"
	SecurityPolicyColumnInstanceID            = "instance_id"
	SecurityPolicyColumnCreationDate          = "creation_date"
	SecurityPolicyColumnChangeDate            = "change_date"
	SecurityPolicyColumnSequence              = "sequence"
	SecurityPolicyColumnEnableIframeEmbedding = "enable_iframe_embedding"
	SecurityPolicyColumnAllowedOrigins        = "origins"
	SecurityPolicyColumnEnableImpersonation   = "enable_impersonation"
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
	return handler.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(SecurityPolicyColumnInstanceID, ""),
		},
		changes,
	), nil
}
