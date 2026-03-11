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
DetectionSettingsProjectionTable             = "projections.detection_settings"
DetectionSettingsColumnCreationDate          = "creation_date"
DetectionSettingsColumnChangeDate            = "change_date"
DetectionSettingsColumnInstanceID            = "instance_id"
DetectionSettingsColumnSequence              = "sequence"
DetectionSettingsColumnEnabled               = "enabled"
DetectionSettingsColumnFailOpen              = "fail_open"
DetectionSettingsColumnFailureBurstThreshold = "failure_burst_threshold"
DetectionSettingsColumnHistoryWindow         = "history_window"
DetectionSettingsColumnContextChangeWindow   = "context_change_window"
DetectionSettingsColumnMaxSignalsPerUser     = "max_signals_per_user"
DetectionSettingsColumnMaxSignalsPerSession  = "max_signals_per_session"
DetectionSettingsColumnRulesManaged          = "rules_managed"
)

type detectionSettingsProjection struct{}

func newDetectionSettingsProjection(ctx context.Context, config handler.Config) *handler.Handler {
return handler.NewHandler(ctx, &config, new(detectionSettingsProjection))
}

func (*detectionSettingsProjection) Name() string {
return DetectionSettingsProjectionTable
}

func (*detectionSettingsProjection) Init() *old_handler.Check {
return handler.NewTableCheck(
handler.NewTable([]*handler.InitColumn{
handler.NewColumn(DetectionSettingsColumnCreationDate, handler.ColumnTypeTimestamp),
handler.NewColumn(DetectionSettingsColumnChangeDate, handler.ColumnTypeTimestamp),
handler.NewColumn(DetectionSettingsColumnInstanceID, handler.ColumnTypeText),
handler.NewColumn(DetectionSettingsColumnSequence, handler.ColumnTypeInt64),
handler.NewColumn(DetectionSettingsColumnEnabled, handler.ColumnTypeBool, handler.Default(false)),
handler.NewColumn(DetectionSettingsColumnFailOpen, handler.ColumnTypeBool, handler.Default(false)),
handler.NewColumn(DetectionSettingsColumnFailureBurstThreshold, handler.ColumnTypeInt64, handler.Default(0)),
handler.NewColumn(DetectionSettingsColumnHistoryWindow, handler.ColumnTypeInterval, handler.Default("0 seconds")),
handler.NewColumn(DetectionSettingsColumnContextChangeWindow, handler.ColumnTypeInterval, handler.Default("0 seconds")),
handler.NewColumn(DetectionSettingsColumnMaxSignalsPerUser, handler.ColumnTypeInt64, handler.Default(0)),
handler.NewColumn(DetectionSettingsColumnMaxSignalsPerSession, handler.ColumnTypeInt64, handler.Default(0)),
handler.NewColumn(DetectionSettingsColumnRulesManaged, handler.ColumnTypeBool, handler.Default(false)),
}, handler.NewPrimaryKey(DetectionSettingsColumnInstanceID)),
)
}

func (p *detectionSettingsProjection) Reducers() []handler.AggregateReducer {
return []handler.AggregateReducer{{
Aggregate: instance.AggregateType,
EventReducers: []handler.EventReducer{
{Event: instance.DetectionSettingsSetEventType, Reduce: p.reduceDetectionSettingsSet},
{Event: instance.InstanceRemovedEventType, Reduce: reduceInstanceRemovedHelper(DetectionSettingsColumnInstanceID)},
},
}}
}

func (p *detectionSettingsProjection) reduceDetectionSettingsSet(event eventstore.Event) (*handler.Statement, error) {
e, ok := event.(*instance.DetectionSettingsSetEvent)
if !ok {
return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-eK8sd", "reduce.wrong.event.type %s", instance.DetectionSettingsSetEventType)
}
changes := []handler.Column{
handler.NewCol(DetectionSettingsColumnCreationDate, handler.OnlySetValueOnInsert(DetectionSettingsProjectionTable, e.CreationDate())),
handler.NewCol(DetectionSettingsColumnChangeDate, e.CreationDate()),
handler.NewCol(DetectionSettingsColumnInstanceID, e.Aggregate().InstanceID),
handler.NewCol(DetectionSettingsColumnSequence, e.Sequence()),
}
if e.Enabled != nil {
changes = append(changes, handler.NewCol(DetectionSettingsColumnEnabled, *e.Enabled))
}
if e.FailOpen != nil {
changes = append(changes, handler.NewCol(DetectionSettingsColumnFailOpen, *e.FailOpen))
}
if e.FailureBurstThreshold != nil {
changes = append(changes, handler.NewCol(DetectionSettingsColumnFailureBurstThreshold, *e.FailureBurstThreshold))
}
if e.HistoryWindow != nil {
changes = append(changes, handler.NewCol(DetectionSettingsColumnHistoryWindow, *e.HistoryWindow))
}
if e.ContextChangeWindow != nil {
changes = append(changes, handler.NewCol(DetectionSettingsColumnContextChangeWindow, *e.ContextChangeWindow))
}
if e.MaxSignalsPerUser != nil {
changes = append(changes, handler.NewCol(DetectionSettingsColumnMaxSignalsPerUser, *e.MaxSignalsPerUser))
}
if e.MaxSignalsPerSession != nil {
changes = append(changes, handler.NewCol(DetectionSettingsColumnMaxSignalsPerSession, *e.MaxSignalsPerSession))
}
if e.RulesManaged != nil {
changes = append(changes, handler.NewCol(DetectionSettingsColumnRulesManaged, *e.RulesManaged))
}
return handler.NewUpsertStatement(e, []handler.Column{handler.NewCol(DetectionSettingsColumnInstanceID, "")}, changes), nil
}
