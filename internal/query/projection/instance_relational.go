package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const InstanceRelationalTable = "zitadel.instances"

type instanceRelationalProjection struct{}

func newInstanceRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(instanceRelationalProjection))
}

func (*instanceRelationalProjection) Name() string {
	return InstanceRelationalTable
}

func (*instanceRelationalProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(InstanceColumnID, handler.ColumnTypeText),
			handler.NewColumn(InstanceColumnName, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(InstanceColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(InstanceColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(InstanceColumnDefaultOrgID, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(InstanceColumnProjectID, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(InstanceColumnConsoleID, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(InstanceColumnConsoleAppID, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(InstanceColumnDefaultLanguage, handler.ColumnTypeText, handler.Default("")),
		},
			handler.NewPrimaryKey(InstanceColumnID),
		),
	)
}

func (p *instanceRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceAddedEventType,
					Reduce: p.reduceInstanceRelationalAdded,
				},
				// {
				// 	Event:  instance.InstanceChangedEventType,
				// 	Reduce: p.reduceInstanceChanged,
				// },
				// {
				// 	Event:  instance.InstanceRemovedEventType,
				// 	Reduce: reduceInstanceRemovedHelper(InstanceColumnID),
				// },
				// {
				// 	Event:  instance.DefaultOrgSetEventType,
				// 	Reduce: p.reduceDefaultOrgSet,
				// },
				// {
				// 	Event:  instance.ProjectSetEventType,
				// 	Reduce: p.reduceIAMProjectSet,
				// },
				// {
				// 	Event:  instance.ConsoleSetEventType,
				// 	Reduce: p.reduceConsoleSet,
				// },
				// {
				// 	Event:  instance.DefaultLanguageSetEventType,
				// 	Reduce: p.reduceDefaultLanguageSet,
				// },
			},
		},
	}
}

func (p *instanceRelationalProjection) reduceInstanceRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.InstanceAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-29nRr", "reduce.wrong.event.type %s", instance.InstanceAddedEventType)
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(InstanceColumnID, e.Aggregate().InstanceID),
			handler.NewCol(InstanceColumnCreationDate, e.CreationDate()),
			handler.NewCol(InstanceColumnChangeDate, e.CreationDate()),
			handler.NewCol(InstanceColumnSequence, e.Sequence()),
			handler.NewCol(InstanceColumnName, e.Name),
		},
	), nil
}

// func reduceInstanceRemovedHelper(instanceIDCol string) func(event eventstore.Event) (*handler.Statement, error) {
// 	return func(event eventstore.Event) (*handler.Statement, error) {
// 		e, ok := event.(*instance.InstanceRemovedEvent)
// 		if !ok {
// 			return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-22nlS", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
// 		}
// 		return handler.NewDeleteStatement(
// 			e,
// 			[]handler.Condition{
// 				handler.NewCond(instanceIDCol, e.Aggregate().ID),
// 			},
// 		), nil
// 	}
// }

// func (p *instanceRelationalProjection) reduceInstanceChanged(event eventstore.Event) (*handler.Statement, error) {
// 	e, ok := event.(*instance.InstanceChangedEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-so2am1", "reduce.wrong.event.type %s", instance.InstanceChangedEventType)
// 	}
// 	return handler.NewUpdateStatement(
// 		e,
// 		[]handler.Column{
// 			handler.NewCol(InstanceColumnName, e.Name),
// 			handler.NewCol(InstanceColumnChangeDate, e.CreationDate()),
// 			handler.NewCol(InstanceColumnSequence, e.Sequence()),
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(InstanceColumnID, e.Aggregate().InstanceID),
// 		},
// 	), nil
// }

// func (p *instanceRelationalProjection) reduceDefaultOrgSet(event eventstore.Event) (*handler.Statement, error) {
// 	e, ok := event.(*instance.DefaultOrgSetEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-2n9f2", "reduce.wrong.event.type %s", instance.DefaultOrgSetEventType)
// 	}
// 	return handler.NewUpdateStatement(
// 		e,
// 		[]handler.Column{
// 			handler.NewCol(InstanceColumnChangeDate, e.CreationDate()),
// 			handler.NewCol(InstanceColumnSequence, e.Sequence()),
// 			handler.NewCol(InstanceColumnDefaultOrgID, e.OrgID),
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(InstanceColumnID, e.Aggregate().InstanceID),
// 		},
// 	), nil
// }

// func (p *instanceRelationalProjection) reduceIAMProjectSet(event eventstore.Event) (*handler.Statement, error) {
// 	e, ok := event.(*instance.ProjectSetEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-30o0e", "reduce.wrong.event.type %s", instance.ProjectSetEventType)
// 	}
// 	return handler.NewUpdateStatement(
// 		e,
// 		[]handler.Column{
// 			handler.NewCol(InstanceColumnChangeDate, e.CreationDate()),
// 			handler.NewCol(InstanceColumnSequence, e.Sequence()),
// 			handler.NewCol(InstanceColumnProjectID, e.ProjectID),
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(InstanceColumnID, e.Aggregate().InstanceID),
// 		},
// 	), nil
// }

// func (p *instanceRelationalProjection) reduceConsoleSet(event eventstore.Event) (*handler.Statement, error) {
// 	e, ok := event.(*instance.ConsoleSetEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Dgf11", "reduce.wrong.event.type %s", instance.ConsoleSetEventType)
// 	}
// 	return handler.NewUpdateStatement(
// 		e,
// 		[]handler.Column{
// 			handler.NewCol(InstanceColumnChangeDate, e.CreationDate()),
// 			handler.NewCol(InstanceColumnSequence, e.Sequence()),
// 			handler.NewCol(InstanceColumnConsoleID, e.ClientID),
// 			handler.NewCol(InstanceColumnConsoleAppID, e.AppID),
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(InstanceColumnID, e.Aggregate().InstanceID),
// 		},
// 	), nil
// }

// func (p *instanceRelationalProjection) reduceDefaultLanguageSet(event eventstore.Event) (*handler.Statement, error) {
// 	e, ok := event.(*instance.DefaultLanguageSetEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-30o0e", "reduce.wrong.event.type %s", instance.DefaultLanguageSetEventType)
// 	}
// 	return handler.NewUpdateStatement(
// 		e,
// 		[]handler.Column{
// 			handler.NewCol(InstanceColumnChangeDate, e.CreationDate()),
// 			handler.NewCol(InstanceColumnSequence, e.Sequence()),
// 			handler.NewCol(InstanceColumnDefaultLanguage, e.Language.String()),
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(InstanceColumnID, e.Aggregate().InstanceID),
// 		},
// 	), nil
// }
