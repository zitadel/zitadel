package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/settings"
)

const (
	OrganizationSettingsTable                          = "projections.organization_settings"
	OrganizationSettingsIDCol                          = "id"
	OrganizationSettingsCreationDateCol                = "creation_date"
	OrganizationSettingsChangeDateCol                  = "change_date"
	OrganizationSettingsResourceOwnerCol               = "resource_owner"
	OrganizationSettingsInstanceIDCol                  = "instance_id"
	OrganizationSettingsSequenceCol                    = "sequence"
	OrganizationSettingsOrganizationScopedUsernamesCol = "organization_scoped_usernames"
)

type organizationSettingsProjection struct{}

func newOrganizationSettingsProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(organizationSettingsProjection))
}

func (*organizationSettingsProjection) Name() string {
	return OrganizationSettingsTable
}

func (*organizationSettingsProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(OrganizationSettingsIDCol, handler.ColumnTypeText),
			handler.NewColumn(OrganizationSettingsCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(OrganizationSettingsChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(OrganizationSettingsResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(OrganizationSettingsInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(OrganizationSettingsSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(OrganizationSettingsOrganizationScopedUsernamesCol, handler.ColumnTypeBool),
		},
			handler.NewPrimaryKey(OrganizationSettingsInstanceIDCol, OrganizationSettingsResourceOwnerCol, OrganizationSettingsIDCol),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{OrganizationSettingsResourceOwnerCol})),
		),
	)
}

func (p *organizationSettingsProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: settings.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  settings.OrganizationSettingsSetEventType,
					Reduce: p.reduceOrganizationSettingsSet,
				},
				{
					Event:  settings.OrganizationSettingsRemovedEventType,
					Reduce: p.reduceOrganizationSettingsRemoved,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOrgRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(OrganizationSettingsInstanceIDCol),
				},
			},
		},
	}
}

func (p *organizationSettingsProjection) reduceOrganizationSettingsSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*settings.OrganizationSettingsSetEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewUpsertStatement(e,
		[]handler.Column{
			handler.NewCol(OrganizationSettingsInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(OrganizationSettingsResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(OrganizationSettingsIDCol, e.Aggregate().ID),
		},
		[]handler.Column{
			handler.NewCol(OrganizationSettingsInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(OrganizationSettingsResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(OrganizationSettingsIDCol, e.Aggregate().ID),
			handler.NewCol(OrganizationSettingsCreationDateCol, handler.OnlySetValueOnInsert(OrganizationSettingsTable, e.CreationDate())),
			handler.NewCol(OrganizationSettingsChangeDateCol, e.CreationDate()),
			handler.NewCol(OrganizationSettingsSequenceCol, e.Sequence()),
			handler.NewCol(OrganizationSettingsOrganizationScopedUsernamesCol, e.OrganizationScopedUsernames),
		},
	), nil
}

func (p *organizationSettingsProjection) reduceOrganizationSettingsRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*settings.OrganizationSettingsRemovedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(OrganizationSettingsInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(OrganizationSettingsResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCond(OrganizationSettingsIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *organizationSettingsProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.OrgRemovedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(OrganizationSettingsInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(OrganizationSettingsResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCond(OrganizationSettingsIDCol, e.Aggregate().ID),
		},
	), nil
}
