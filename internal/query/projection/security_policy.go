package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

const (
	SecurityPolicyProjectionTable      = "projections.security_policies"
	SecurityPolicyColumnInstanceID     = "instance_id"
	SecurityPolicyColumnCreationDate   = "creation_date"
	SecurityPolicyColumnChangeDate     = "change_date"
	SecurityPolicyColumnSequence       = "sequence"
	SecurityPolicyColumnEnabled        = "enabled"
	SecurityPolicyColumnAllowedOrigins = "origins"
)

type securityPolicyProjection struct {
	crdb.StatementHandler
}

func newSecurityPolicyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *securityPolicyProjection {
	p := new(securityPolicyProjection)
	config.ProjectionName = SecurityPolicyProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(SecurityPolicyColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(SecurityPolicyColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(SecurityPolicyColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(SecurityPolicyColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(SecurityPolicyColumnEnabled, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(SecurityPolicyColumnAllowedOrigins, crdb.ColumnTypeTextArray, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(SecurityPolicyColumnInstanceID),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *securityPolicyProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-D3g87", "reduce.wrong.event.type %s", instance.SecurityPolicySetEventType)
	}
	changes := []handler.Column{
		handler.NewCol(SecurityPolicyColumnCreationDate, e.CreationDate()),
		handler.NewCol(SecurityPolicyColumnChangeDate, e.CreationDate()),
		handler.NewCol(SecurityPolicyColumnInstanceID, e.Aggregate().InstanceID),
		handler.NewCol(SecurityPolicyColumnSequence, e.Sequence()),
	}
	if e.Enabled != nil {
		changes = append(changes, handler.NewCol(SecurityPolicyColumnEnabled, *e.Enabled))
	}
	if e.AllowedOrigins != nil {
		changes = append(changes, handler.NewCol(SecurityPolicyColumnAllowedOrigins, e.AllowedOrigins))
	}
	return crdb.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(SecurityPolicyColumnInstanceID, ""),
		},
		changes,
	), nil
}
