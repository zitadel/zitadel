package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/restrictions"
)

const (
	RestrictionsProjectionTable = "projections.restrictions"

	RestrictionsColumnAggregateID   = "aggregate_id"
	RestrictionsColumnCreationDate  = "creation_date"
	RestrictionsColumnChangeDate    = "change_date"
	RestrictionsColumnResourceOwner = "resource_owner"
	RestrictionsColumnInstanceID    = "instance_id"
	RestrictionsColumnSequence      = "sequence"

	RestrictionsColumnDisallowPublicOrgRegistration = "disallow_public_org_registration"
)

type restrictionsProjection struct{}

func newRestrictionsProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, &restrictionsProjection{})
}

func (*restrictionsProjection) Name() string {
	return RestrictionsProjectionTable
}

func (*restrictionsProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(RestrictionsColumnAggregateID, handler.ColumnTypeText),
			handler.NewColumn(RestrictionsColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(RestrictionsColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(RestrictionsColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(RestrictionsColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(RestrictionsColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(RestrictionsColumnDisallowPublicOrgRegistration, handler.ColumnTypeBool),
		},
			handler.NewPrimaryKey(RestrictionsColumnInstanceID, RestrictionsColumnResourceOwner),
		),
	)
}

func (p *restrictionsProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: restrictions.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  restrictions.SetEventType,
					Reduce: p.reduceRestrictionsSet,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(RestrictionsColumnInstanceID),
				},
			},
		},
	}
}

func (p *restrictionsProjection) reduceRestrictionsSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*restrictions.SetEvent](event)
	if err != nil {
		return nil, err
	}
	conflictCols := []handler.Column{
		handler.NewCol(RestrictionsColumnInstanceID, e.Aggregate().InstanceID),
		handler.NewCol(RestrictionsColumnResourceOwner, e.Aggregate().ResourceOwner),
	}
	updateCols := []handler.Column{
		handler.NewCol(RestrictionsColumnInstanceID, e.Aggregate().InstanceID),
		handler.NewCol(RestrictionsColumnResourceOwner, e.Aggregate().ResourceOwner),
		handler.NewCol(RestrictionsColumnCreationDate, handler.OnlySetValueOnInsert(RestrictionsProjectionTable, e.CreationDate())),
		handler.NewCol(RestrictionsColumnChangeDate, e.CreationDate()),
		handler.NewCol(RestrictionsColumnSequence, e.Sequence()),
		handler.NewCol(RestrictionsColumnAggregateID, e.Aggregate().ID),
	}
	if e.DisallowPublicOrgRegistrations != nil {
		updateCols = append(updateCols, handler.NewCol(RestrictionsColumnDisallowPublicOrgRegistration, *e.DisallowPublicOrgRegistrations))
	}
	return handler.NewUpsertStatement(e, conflictCols, updateCols), nil
}
