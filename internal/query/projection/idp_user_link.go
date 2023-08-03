package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	IDPUserLinkTable             = "projections.idp_user_links3"
	IDPUserLinkIDPIDCol          = "idp_id"
	IDPUserLinkUserIDCol         = "user_id"
	IDPUserLinkExternalUserIDCol = "external_user_id"
	IDPUserLinkCreationDateCol   = "creation_date"
	IDPUserLinkChangeDateCol     = "change_date"
	IDPUserLinkSequenceCol       = "sequence"
	IDPUserLinkResourceOwnerCol  = "resource_owner"
	IDPUserLinkInstanceIDCol     = "instance_id"
	IDPUserLinkDisplayNameCol    = "display_name"
	IDPUserLinkOwnerRemovedCol   = "owner_removed"
)

type idpUserLinkProjection struct {
	crdb.StatementHandler
}

func newIDPUserLinkProjection(ctx context.Context, config crdb.StatementHandlerConfig) *idpUserLinkProjection {
	p := new(idpUserLinkProjection)
	config.ProjectionName = IDPUserLinkTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(IDPUserLinkIDPIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPUserLinkUserIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPUserLinkExternalUserIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPUserLinkCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(IDPUserLinkChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(IDPUserLinkSequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(IDPUserLinkResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPUserLinkInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPUserLinkDisplayNameCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPUserLinkOwnerRemovedCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(IDPUserLinkInstanceIDCol, IDPUserLinkIDPIDCol, IDPUserLinkExternalUserIDCol),
			crdb.WithIndex(crdb.NewIndex("user_id", []string{IDPUserLinkUserIDCol})),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{IDPUserLinkOwnerRemovedCol})),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *idpUserLinkProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.UserIDPLinkAddedType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  user.UserIDPLinkCascadeRemovedType,
					Reduce: p.reduceCascadeRemoved,
				},
				{
					Event:  user.UserIDPLinkRemovedType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},
				{
					Event:  user.UserIDPExternalIDMigratedType,
					Reduce: p.reduceExternalIDMigrated,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPConfigRemoved,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPConfigRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(IDPUserLinkInstanceIDCol),
				},
			},
		},
	}
}

func (p *idpUserLinkProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserIDPLinkAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-DpmXq", "reduce.wrong.event.type %s", user.UserIDPLinkAddedType)
	}

	return crdb.NewCreateStatement(e,
		[]handler.Column{
			handler.NewCol(IDPUserLinkIDPIDCol, e.IDPConfigID),
			handler.NewCol(IDPUserLinkUserIDCol, e.Aggregate().ID),
			handler.NewCol(IDPUserLinkExternalUserIDCol, e.ExternalUserID),
			handler.NewCol(IDPUserLinkCreationDateCol, e.CreationDate()),
			handler.NewCol(IDPUserLinkChangeDateCol, e.CreationDate()),
			handler.NewCol(IDPUserLinkSequenceCol, e.Sequence()),
			handler.NewCol(IDPUserLinkResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(IDPUserLinkInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(IDPUserLinkDisplayNameCol, e.DisplayName),
		},
	), nil
}

func (p *idpUserLinkProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserIDPLinkRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-AZmfJ", "reduce.wrong.event.type %s", user.UserIDPLinkRemovedType)
	}

	return crdb.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(IDPUserLinkIDPIDCol, e.IDPConfigID),
			handler.NewCond(IDPUserLinkUserIDCol, e.Aggregate().ID),
			handler.NewCond(IDPUserLinkExternalUserIDCol, e.ExternalUserID),
			handler.NewCond(IDPUserLinkInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpUserLinkProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserIDPLinkCascadeRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-jQpv9", "reduce.wrong.event.type %s", user.UserIDPLinkCascadeRemovedType)
	}

	return crdb.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(IDPUserLinkIDPIDCol, e.IDPConfigID),
			handler.NewCond(IDPUserLinkUserIDCol, e.Aggregate().ID),
			handler.NewCond(IDPUserLinkExternalUserIDCol, e.ExternalUserID),
			handler.NewCond(IDPUserLinkInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpUserLinkProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-PGiAY", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(IDPUserLinkChangeDateCol, e.CreationDate()),
			handler.NewCol(IDPUserLinkSequenceCol, e.Sequence()),
			handler.NewCol(IDPUserLinkOwnerRemovedCol, true),
		},
		[]handler.Condition{
			handler.NewCond(IDPUserLinkResourceOwnerCol, e.Aggregate().ID),
			handler.NewCond(IDPUserLinkInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpUserLinkProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-uwlWE", "reduce.wrong.event.type %s", user.UserRemovedType)
	}

	return crdb.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(IDPUserLinkUserIDCol, e.Aggregate().ID),
			handler.NewCond(IDPUserLinkInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpUserLinkProjection) reduceExternalIDMigrated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.UserIDPExternalIDMigratedEvent](event)
	if err != nil {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-AS3th", "reduce.wrong.event.type %s", user.UserIDPExternalIDMigratedType)
	}

	return crdb.NewUpdateStatement(e,
		[]handler.Column{
			handler.NewCol(IDPUserLinkChangeDateCol, e.CreationDate()),
			handler.NewCol(IDPUserLinkSequenceCol, e.Sequence()),
			handler.NewCol(IDPUserLinkExternalUserIDCol, e.NewID),
		},
		[]handler.Condition{
			handler.NewCond(IDPUserLinkIDPIDCol, e.IDPConfigID),
			handler.NewCond(IDPUserLinkUserIDCol, e.Aggregate().ID),
			handler.NewCond(IDPUserLinkExternalUserIDCol, e.PreviousID),
			handler.NewCond(IDPUserLinkInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpUserLinkProjection) reduceIDPConfigRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idpID string

	switch e := event.(type) {
	case *org.IDPConfigRemovedEvent:
		idpID = e.ConfigID
	case *instance.IDPConfigRemovedEvent:
		idpID = e.ConfigID
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-iCKSj", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigRemovedEventType, instance.IDPConfigRemovedEventType})
	}

	return crdb.NewDeleteStatement(event,
		[]handler.Condition{
			handler.NewCond(IDPUserLinkIDPIDCol, idpID),
			handler.NewCond(IDPUserLinkResourceOwnerCol, event.Aggregate().ResourceOwner),
			handler.NewCond(IDPUserLinkInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}
