package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user/schema"
)

const (
	UserSchemaTable = "projections.user_schemas"

	UserSchemaIDCol                     = "id"
	UserSchemaChangeDateCol             = "change_date"
	UserSchemaSequenceCol               = "sequence"
	UserSchemaInstanceIDCol             = "instance_id"
	UserSchemaStateCol                  = "state"
	UserSchemaTypeCol                   = "type"
	UserSchemaRevisionCol               = "revision"
	UserSchemaSchemaCol                 = "schema"
	UserSchemaPossibleAuthenticatorsCol = "possible_authenticators"
)

type userSchemaProjection struct{}

func newUserSchemaProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(userSchemaProjection))
}

func (*userSchemaProjection) Name() string {
	return UserSchemaTable
}

func (*userSchemaProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(UserSchemaIDCol, handler.ColumnTypeText),
			handler.NewColumn(UserSchemaChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(UserSchemaSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(UserSchemaStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(UserSchemaInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(UserSchemaTypeCol, handler.ColumnTypeText),
			handler.NewColumn(UserSchemaRevisionCol, handler.ColumnTypeInt64),
			handler.NewColumn(UserSchemaSchemaCol, handler.ColumnTypeJSONB, handler.Nullable()),
			handler.NewColumn(UserSchemaPossibleAuthenticatorsCol, handler.ColumnTypeEnumArray, handler.Nullable()),
		},
			handler.NewPrimaryKey(UserSchemaInstanceIDCol, UserSchemaIDCol),
		),
	)
}

func (p *userSchemaProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: schema.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  schema.CreatedType,
					Reduce: p.reduceCreated,
				},
				{
					Event:  schema.UpdatedType,
					Reduce: p.reduceUpdated,
				},
				{
					Event:  schema.DeactivatedType,
					Reduce: p.reduceDeactivated,
				},
				{
					Event:  schema.ReactivatedType,
					Reduce: p.reduceReactivated,
				},
				{
					Event:  schema.DeletedType,
					Reduce: p.reduceDeleted,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(UserSchemaInstanceIDCol),
				},
			},
		},
	}
}

func (p *userSchemaProjection) reduceCreated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*schema.CreatedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserSchemaIDCol, event.Aggregate().ID),
			handler.NewCol(UserSchemaChangeDateCol, event.CreatedAt()),
			handler.NewCol(UserSchemaSequenceCol, event.Sequence()),
			handler.NewCol(UserSchemaInstanceIDCol, event.Aggregate().InstanceID),
			handler.NewCol(UserSchemaStateCol, domain.UserSchemaStateActive),
			handler.NewCol(UserSchemaTypeCol, e.SchemaType),
			handler.NewCol(UserSchemaRevisionCol, 1),
			handler.NewCol(UserSchemaSchemaCol, e.Schema),
			handler.NewCol(UserSchemaPossibleAuthenticatorsCol, e.PossibleAuthenticators),
		},
	), nil
}

func (p *userSchemaProjection) reduceUpdated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*schema.UpdatedEvent](event)
	if err != nil {
		return nil, err
	}

	cols := []handler.Column{
		handler.NewCol(UserSchemaChangeDateCol, event.CreatedAt()),
		handler.NewCol(UserSchemaSequenceCol, event.Sequence()),
	}
	if e.SchemaType != nil {
		cols = append(cols, handler.NewCol(UserSchemaTypeCol, *e.SchemaType))
	}

	if len(e.Schema) > 0 {
		cols = append(cols, handler.NewCol(UserSchemaSchemaCol, e.Schema))
		cols = append(cols, handler.NewIncrementCol(UserSchemaRevisionCol, 1))
	}

	if len(e.PossibleAuthenticators) > 0 {
		cols = append(cols, handler.NewCol(UserSchemaPossibleAuthenticatorsCol, e.PossibleAuthenticators))
	}

	return handler.NewUpdateStatement(
		event,
		cols,
		[]handler.Condition{
			handler.NewCond(UserSchemaIDCol, event.Aggregate().ID),
			handler.NewCond(UserSchemaInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *userSchemaProjection) reduceDeactivated(event eventstore.Event) (*handler.Statement, error) {
	_, err := assertEvent[*schema.DeactivatedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserSchemaChangeDateCol, event.CreatedAt()),
			handler.NewCol(UserSchemaSequenceCol, event.Sequence()),
			handler.NewCol(UserSchemaStateCol, domain.UserSchemaStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(UserSchemaIDCol, event.Aggregate().ID),
			handler.NewCond(UserSchemaInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *userSchemaProjection) reduceReactivated(event eventstore.Event) (*handler.Statement, error) {
	_, err := assertEvent[*schema.ReactivatedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserSchemaChangeDateCol, event.CreatedAt()),
			handler.NewCol(UserSchemaSequenceCol, event.Sequence()),
			handler.NewCol(UserSchemaStateCol, domain.UserSchemaStateActive),
		},
		[]handler.Condition{
			handler.NewCond(UserSchemaIDCol, event.Aggregate().ID),
			handler.NewCond(UserSchemaInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *userSchemaProjection) reduceDeleted(event eventstore.Event) (*handler.Statement, error) {
	_, err := assertEvent[*schema.DeletedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(UserSchemaIDCol, event.Aggregate().ID),
			handler.NewCond(UserSchemaInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}
