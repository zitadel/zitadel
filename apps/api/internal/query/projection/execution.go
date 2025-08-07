package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	exec "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/target"
)

const (
	ExecutionTable           = "projections.executions1"
	ExecutionIDCol           = "id"
	ExecutionCreationDateCol = "creation_date"
	ExecutionChangeDateCol   = "change_date"
	ExecutionInstanceIDCol   = "instance_id"
	ExecutionSequenceCol     = "sequence"

	ExecutionTargetSuffix         = "targets"
	ExecutionTargetExecutionIDCol = "execution_id"
	ExecutionTargetInstanceIDCol  = "instance_id"
	ExecutionTargetPositionCol    = "position"
	ExecutionTargetTargetIDCol    = "target_id"
	ExecutionTargetIncludeCol     = "include"
)

type executionProjection struct{}

func newExecutionProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(executionProjection))
}

func (*executionProjection) Name() string {
	return ExecutionTable
}

func (*executionProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(ExecutionIDCol, handler.ColumnTypeText),
			handler.NewColumn(ExecutionCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(ExecutionChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(ExecutionSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(ExecutionInstanceIDCol, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(ExecutionInstanceIDCol, ExecutionIDCol),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(ExecutionTargetInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(ExecutionTargetExecutionIDCol, handler.ColumnTypeText),
			handler.NewColumn(ExecutionTargetPositionCol, handler.ColumnTypeInt64),
			handler.NewColumn(ExecutionTargetIncludeCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(ExecutionTargetTargetIDCol, handler.ColumnTypeText, handler.Nullable()),
		},
			handler.NewPrimaryKey(ExecutionTargetInstanceIDCol, ExecutionTargetExecutionIDCol, ExecutionTargetPositionCol),
			ExecutionTargetSuffix,
			handler.WithForeignKey(handler.NewForeignKey("execution", []string{ExecutionTargetInstanceIDCol, ExecutionTargetExecutionIDCol}, []string{ExecutionInstanceIDCol, ExecutionIDCol})),
			handler.WithIndex(handler.NewIndex("execution", []string{ExecutionTargetInstanceIDCol, ExecutionTargetExecutionIDCol})),
		),
	)
}

func (p *executionProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: exec.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  exec.SetEventV2Type,
					Reduce: p.reduceExecutionSet,
				},
				{
					Event:  exec.RemovedEventType,
					Reduce: p.reduceExecutionRemoved,
				},
			},
		},
		{
			Aggregate: target.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  target.RemovedEventType,
					Reduce: p.reduceTargetRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(ExecutionInstanceIDCol),
				},
			},
		},
	}
}

func (p *executionProjection) reduceExecutionSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*exec.SetEventV2](event)
	if err != nil {
		return nil, err
	}

	stmts := []func(eventstore.Event) handler.Exec{
		handler.AddUpsertStatement(
			[]handler.Column{
				handler.NewCol(ExecutionInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(ExecutionIDCol, e.Aggregate().ID),
			},
			[]handler.Column{
				handler.NewCol(ExecutionInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(ExecutionIDCol, e.Aggregate().ID),
				handler.NewCol(ExecutionCreationDateCol, handler.OnlySetValueOnInsert(ExecutionTable, e.CreationDate())),
				handler.NewCol(ExecutionChangeDateCol, e.CreationDate()),
				handler.NewCol(ExecutionSequenceCol, e.Sequence()),
			},
		),
		// cleanup execution targets to re-insert them
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(ExecutionTargetInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCond(ExecutionTargetExecutionIDCol, e.Aggregate().ID),
			},
			handler.WithTableSuffix(ExecutionTargetSuffix),
		),
	}

	if len(e.Targets) > 0 {
		for i, target := range e.Targets {
			var targetStr, includeStr string
			switch target.Type {
			case domain.ExecutionTargetTypeTarget:
				targetStr = target.Target
			case domain.ExecutionTargetTypeInclude:
				includeStr = target.Target
			case domain.ExecutionTargetTypeUnspecified:
				continue
			default:
				continue
			}

			stmts = append(stmts,
				handler.AddCreateStatement(
					[]handler.Column{
						handler.NewCol(ExecutionTargetInstanceIDCol, e.Aggregate().InstanceID),
						handler.NewCol(ExecutionTargetExecutionIDCol, e.Aggregate().ID),
						handler.NewCol(ExecutionTargetPositionCol, i+1),
						handler.NewCol(ExecutionTargetIncludeCol, includeStr),
						handler.NewCol(ExecutionTargetTargetIDCol, targetStr),
					},
					handler.WithTableSuffix(ExecutionTargetSuffix),
				),
			)
		}
	}

	return handler.NewMultiStatement(e, stmts...), nil
}

func (p *executionProjection) reduceTargetRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*target.RemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ExecutionTargetInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(ExecutionTargetTargetIDCol, e.Aggregate().ID),
		},
		handler.WithTableSuffix(ExecutionTargetSuffix),
	), nil
}

func (p *executionProjection) reduceExecutionRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*exec.RemovedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(ExecutionInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(ExecutionIDCol, e.Aggregate().ID),
		},
	), nil
}
