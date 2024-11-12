package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	SystemFeatureTable = "projections.system_features"

	SystemFeatureKeyCol          = "key"
	SystemFeatureCreationDateCol = "creation_date"
	SystemFeatureChangeDateCol   = "change_date"
	SystemFeatureSequenceCol     = "sequence"
	SystemFeatureValueCol        = "value"
)

type systemFeatureProjection struct{}

func newSystemFeatureProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(systemFeatureProjection))
}

func (*systemFeatureProjection) Name() string {
	return SystemFeatureTable
}

func (*systemFeatureProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(handler.NewTable(
		[]*handler.InitColumn{
			handler.NewColumn(SystemFeatureKeyCol, handler.ColumnTypeText),
			handler.NewColumn(SystemFeatureCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(SystemFeatureChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(SystemFeatureSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(SystemFeatureValueCol, handler.ColumnTypeJSONB),
		},
		handler.NewPrimaryKey(SystemFeatureKeyCol),
	))
}

func (*systemFeatureProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{{
		Aggregate: feature_v2.AggregateType,
		EventReducers: []handler.EventReducer{
			{
				Event:  feature_v2.SystemResetEventType,
				Reduce: reduceSystemResetFeatures,
			},
			{
				Event:  feature_v2.SystemLoginDefaultOrgEventType,
				Reduce: reduceSystemSetFeature[bool],
			},
			{
				Event:  feature_v2.SystemTriggerIntrospectionProjectionsEventType,
				Reduce: reduceSystemSetFeature[bool],
			},
			{
				Event:  feature_v2.SystemLegacyIntrospectionEventType,
				Reduce: reduceSystemSetFeature[bool],
			},
			{
				Event:  feature_v2.SystemUserSchemaEventType,
				Reduce: reduceSystemSetFeature[bool],
			},
			{
				Event:  feature_v2.SystemTokenExchangeEventType,
				Reduce: reduceSystemSetFeature[bool],
			},
			{
				Event:  feature_v2.SystemActionsEventType,
				Reduce: reduceSystemSetFeature[bool],
			},
			{
				Event:  feature_v2.SystemImprovedPerformanceEventType,
				Reduce: reduceSystemSetFeature[[]feature.ImprovedPerformanceType],
			},
			{
				Event:  feature_v2.SystemDisableUserTokenEvent,
				Reduce: reduceSystemSetFeature[bool],
			},
			{
				Event:  feature_v2.SystemEnableBackChannelLogout,
				Reduce: reduceSystemSetFeature[bool],
			},
		},
	}}
}

func reduceSystemSetFeature[T any](event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*feature_v2.SetEvent[T])
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-uPh8O", "reduce.wrong.event.type %T", event)
	}
	f, err := e.FeatureJSON()
	if err != nil {
		return nil, err
	}
	columns := []handler.Column{
		handler.NewCol(SystemFeatureKeyCol, f.Key.String()),
		handler.NewCol(SystemFeatureCreationDateCol, handler.OnlySetValueOnInsert(SystemFeatureTable, e.CreationDate())),
		handler.NewCol(SystemFeatureChangeDateCol, e.CreationDate()),
		handler.NewCol(SystemFeatureSequenceCol, e.Sequence()),
		handler.NewCol(SystemFeatureValueCol, f.Value),
	}
	return handler.NewUpsertStatement(e, columns[0:1], columns), nil
}

func reduceSystemResetFeatures(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*feature_v2.ResetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-roo6A", "reduce.wrong.event.type %T", event)
	}
	return handler.NewDeleteStatement(e, []handler.Condition{
		// Hack: need at least one condition or the query builder will throw us an error
		handler.NewIsNotNullCond(SystemFeatureKeyCol),
	}), nil
}
