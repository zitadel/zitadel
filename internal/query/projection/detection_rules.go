package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

const (
	DetectionRulesProjectionTable       = "projections.detection_rules"
	DetectionRulesColumnID              = "id"
	DetectionRulesColumnCreationDate    = "creation_date"
	DetectionRulesColumnChangeDate      = "change_date"
	DetectionRulesColumnInstanceID      = "instance_id"
	DetectionRulesColumnSequence        = "sequence"
	DetectionRulesColumnDescription     = "description"
	DetectionRulesColumnExpr            = "expr"
	DetectionRulesColumnEngine          = "engine"
	DetectionRulesColumnFindingName     = "finding_name"
	DetectionRulesColumnFindingMessage  = "finding_message"
	DetectionRulesColumnFindingBlock    = "finding_block"
	DetectionRulesColumnContextTemplate = "context_template"
	DetectionRulesColumnRateLimitKey    = "rate_limit_key"
	DetectionRulesColumnRateLimitWindow = "rate_limit_window"
	DetectionRulesColumnRateLimitMax    = "rate_limit_max"
	DetectionRulesColumnPriority        = "priority"
	DetectionRulesColumnStopOnMatch     = "stop_on_match"
)

type detectionRulesProjection struct{}

func newDetectionRulesProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(detectionRulesProjection))
}

func (*detectionRulesProjection) Name() string {
	return DetectionRulesProjectionTable
}

func (*detectionRulesProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(DetectionRulesColumnID, handler.ColumnTypeText),
			handler.NewColumn(DetectionRulesColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(DetectionRulesColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(DetectionRulesColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(DetectionRulesColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(DetectionRulesColumnDescription, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(DetectionRulesColumnExpr, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(DetectionRulesColumnEngine, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(DetectionRulesColumnFindingName, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(DetectionRulesColumnFindingMessage, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(DetectionRulesColumnFindingBlock, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(DetectionRulesColumnContextTemplate, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(DetectionRulesColumnRateLimitKey, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(DetectionRulesColumnRateLimitWindow, handler.ColumnTypeInterval, handler.Default("0 seconds")),
			handler.NewColumn(DetectionRulesColumnRateLimitMax, handler.ColumnTypeInt64, handler.Default(0)),
			handler.NewColumn(DetectionRulesColumnPriority, handler.ColumnTypeInt64, handler.Default(0)),
			handler.NewColumn(DetectionRulesColumnStopOnMatch, handler.ColumnTypeBool, handler.Default(false)),
		}, handler.NewPrimaryKey(DetectionRulesColumnInstanceID, DetectionRulesColumnID)),
	)
}

func (p *detectionRulesProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{{
		Aggregate: instance.AggregateType,
		EventReducers: []handler.EventReducer{
			{Event: instance.DetectionRuleAddedEventType, Reduce: p.reduceDetectionRuleAdded},
			{Event: instance.DetectionRuleChangedEventType, Reduce: p.reduceDetectionRuleChanged},
			{Event: instance.DetectionRuleRemovedEventType, Reduce: p.reduceDetectionRuleRemoved},
			{Event: instance.InstanceRemovedEventType, Reduce: reduceInstanceRemovedHelper(DetectionRulesColumnInstanceID)},
		},
	}}
}

func (p *detectionRulesProjection) reduceDetectionRuleAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.DetectionRuleAddedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewUpsertStatement(e,
		[]handler.Column{
			handler.NewCol(DetectionRulesColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(DetectionRulesColumnID, e.RuleID),
		},
		[]handler.Column{
			handler.NewCol(DetectionRulesColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(DetectionRulesColumnID, e.RuleID),
			handler.NewCol(DetectionRulesColumnCreationDate, handler.OnlySetValueOnInsert(DetectionRulesProjectionTable, e.CreationDate())),
			handler.NewCol(DetectionRulesColumnChangeDate, e.CreationDate()),
			handler.NewCol(DetectionRulesColumnSequence, e.Sequence()),
			handler.NewCol(DetectionRulesColumnDescription, e.Description),
			handler.NewCol(DetectionRulesColumnExpr, e.Expr),
			handler.NewCol(DetectionRulesColumnEngine, e.Engine),
			handler.NewCol(DetectionRulesColumnFindingName, e.FindingName),
			handler.NewCol(DetectionRulesColumnFindingMessage, e.FindingMessage),
			handler.NewCol(DetectionRulesColumnFindingBlock, e.FindingBlock),
			handler.NewCol(DetectionRulesColumnContextTemplate, e.ContextTemplate),
			handler.NewCol(DetectionRulesColumnRateLimitKey, e.RateLimitKey),
			handler.NewCol(DetectionRulesColumnRateLimitWindow, e.RateLimitWindow),
			handler.NewCol(DetectionRulesColumnRateLimitMax, e.RateLimitMax),
			handler.NewCol(DetectionRulesColumnPriority, e.Priority),
			handler.NewCol(DetectionRulesColumnStopOnMatch, e.StopOnMatch),
		},
	), nil
}

func (p *detectionRulesProjection) reduceDetectionRuleChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.DetectionRuleChangedEvent](event)
	if err != nil {
		return nil, err
	}
	changes := []handler.Column{
		handler.NewCol(DetectionRulesColumnChangeDate, e.CreationDate()),
		handler.NewCol(DetectionRulesColumnSequence, e.Sequence()),
	}
	if e.Description != nil {
		changes = append(changes, handler.NewCol(DetectionRulesColumnDescription, *e.Description))
	}
	if e.Expr != nil {
		changes = append(changes, handler.NewCol(DetectionRulesColumnExpr, *e.Expr))
	}
	if e.Engine != nil {
		changes = append(changes, handler.NewCol(DetectionRulesColumnEngine, *e.Engine))
	}
	if e.FindingName != nil {
		changes = append(changes, handler.NewCol(DetectionRulesColumnFindingName, *e.FindingName))
	}
	if e.FindingMessage != nil {
		changes = append(changes, handler.NewCol(DetectionRulesColumnFindingMessage, *e.FindingMessage))
	}
	if e.FindingBlock != nil {
		changes = append(changes, handler.NewCol(DetectionRulesColumnFindingBlock, *e.FindingBlock))
	}
	if e.ContextTemplate != nil {
		changes = append(changes, handler.NewCol(DetectionRulesColumnContextTemplate, *e.ContextTemplate))
	}
	if e.RateLimitKey != nil {
		changes = append(changes, handler.NewCol(DetectionRulesColumnRateLimitKey, *e.RateLimitKey))
	}
	if e.RateLimitWindow != nil {
		changes = append(changes, handler.NewCol(DetectionRulesColumnRateLimitWindow, *e.RateLimitWindow))
	}
	if e.RateLimitMax != nil {
		changes = append(changes, handler.NewCol(DetectionRulesColumnRateLimitMax, *e.RateLimitMax))
	}
	if e.Priority != nil {
		changes = append(changes, handler.NewCol(DetectionRulesColumnPriority, *e.Priority))
	}
	if e.StopOnMatch != nil {
		changes = append(changes, handler.NewCol(DetectionRulesColumnStopOnMatch, *e.StopOnMatch))
	}
	return handler.NewUpdateStatement(e, changes, []handler.Condition{
		handler.NewCond(DetectionRulesColumnInstanceID, e.Aggregate().InstanceID),
		handler.NewCond(DetectionRulesColumnID, e.RuleID),
	}), nil
}

func (p *detectionRulesProjection) reduceDetectionRuleRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.DetectionRuleRemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewDeleteStatement(e, []handler.Condition{
		handler.NewCond(DetectionRulesColumnInstanceID, e.Aggregate().InstanceID),
		handler.NewCond(DetectionRulesColumnID, e.RuleID),
	}), nil
}
