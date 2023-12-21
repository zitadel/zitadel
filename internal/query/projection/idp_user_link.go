package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
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

type idpUserLinkProjection struct{}

func newIDPUserLinkProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(idpUserLinkProjection))
}

func (*idpUserLinkProjection) Name() string {
	return IDPUserLinkTable
}

func (*idpUserLinkProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(IDPUserLinkIDPIDCol, handler.ColumnTypeText),
			handler.NewColumn(IDPUserLinkUserIDCol, handler.ColumnTypeText),
			handler.NewColumn(IDPUserLinkExternalUserIDCol, handler.ColumnTypeText),
			handler.NewColumn(IDPUserLinkCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(IDPUserLinkChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(IDPUserLinkSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(IDPUserLinkResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(IDPUserLinkInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(IDPUserLinkDisplayNameCol, handler.ColumnTypeText),
			handler.NewColumn(IDPUserLinkOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(IDPUserLinkInstanceIDCol, IDPUserLinkIDPIDCol, IDPUserLinkExternalUserIDCol),
			handler.WithIndex(handler.NewIndex("user_id", []string{IDPUserLinkUserIDCol})),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{IDPUserLinkOwnerRemovedCol})),
		),
	)
}

func (p *idpUserLinkProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
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
				{
					Event:  user.UserIDPExternalUsernameChangedType,
					Reduce: p.reduceExternalUsernameChanged,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-DpmXq", "reduce.wrong.event.type %s", user.UserIDPLinkAddedType)
	}

	return handler.NewCreateStatement(e,
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-AZmfJ", "reduce.wrong.event.type %s", user.UserIDPLinkRemovedType)
	}

	return handler.NewDeleteStatement(e,
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-jQpv9", "reduce.wrong.event.type %s", user.UserIDPLinkCascadeRemovedType)
	}

	return handler.NewDeleteStatement(e,
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-PGiAY", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(IDPUserLinkResourceOwnerCol, e.Aggregate().ID),
			handler.NewCond(IDPUserLinkInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpUserLinkProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-uwlWE", "reduce.wrong.event.type %s", user.UserRemovedType)
	}

	return handler.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(IDPUserLinkUserIDCol, e.Aggregate().ID),
			handler.NewCond(IDPUserLinkInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpUserLinkProjection) reduceExternalIDMigrated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.UserIDPExternalIDMigratedEvent](event)
	if err != nil {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-AS3th", "reduce.wrong.event.type %s", user.UserIDPExternalIDMigratedType)
	}

	return handler.NewUpdateStatement(e,
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

func (p *idpUserLinkProjection) reduceExternalUsernameChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.UserIDPExternalUsernameEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(e,
		[]handler.Column{
			handler.NewCol(IDPUserLinkChangeDateCol, e.CreationDate()),
			handler.NewCol(IDPUserLinkSequenceCol, e.Sequence()),
			handler.NewCol(IDPUserLinkDisplayNameCol, e.ExternalUsername),
		},
		[]handler.Condition{
			handler.NewCond(IDPUserLinkIDPIDCol, e.IDPConfigID),
			handler.NewCond(IDPUserLinkUserIDCol, e.Aggregate().ID),
			handler.NewCond(IDPUserLinkExternalUserIDCol, e.ExternalUserID),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-iCKSj", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigRemovedEventType, instance.IDPConfigRemovedEventType})
	}

	return handler.NewDeleteStatement(event,
		[]handler.Condition{
			handler.NewCond(IDPUserLinkIDPIDCol, idpID),
			handler.NewCond(IDPUserLinkResourceOwnerCol, event.Aggregate().ResourceOwner),
			handler.NewCond(IDPUserLinkInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}
