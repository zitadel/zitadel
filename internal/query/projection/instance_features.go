package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/feature"
	feature_v1 "github.com/zitadel/zitadel/internal/repository/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	InstanceFeatureTable = "projections.instance_features2"

	InstanceFeatureInstanceIDCol   = "instance_id"
	InstanceFeatureKeyCol          = "key"
	InstanceFeatureCreationDateCol = "creation_date"
	InstanceFeatureChangeDateCol   = "change_date"
	InstanceFeatureSequenceCol     = "sequence"
	InstanceFeatureValueCol        = "value"
)

type instanceFeatureProjection struct{}

func newInstanceFeatureProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(instanceFeatureProjection))
}

func (*instanceFeatureProjection) Name() string {
	return InstanceFeatureTable
}

func (*instanceFeatureProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(handler.NewTable(
		[]*handler.InitColumn{
			handler.NewColumn(InstanceFeatureInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(InstanceFeatureKeyCol, handler.ColumnTypeText),
			handler.NewColumn(InstanceFeatureCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(InstanceFeatureChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(InstanceFeatureSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(InstanceFeatureValueCol, handler.ColumnTypeJSONB),
		},
		handler.NewPrimaryKey(InstanceFeatureInstanceIDCol, InstanceFeatureKeyCol),
	))
}

func (*instanceFeatureProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{{
		Aggregate: feature_v2.AggregateType,
		EventReducers: []handler.EventReducer{
			{
				Event:  feature_v1.DefaultLoginInstanceEventType,
				Reduce: reduceSetDefaultLoginInstance_v1,
			},
			{
				Event:  feature_v2.InstanceResetEventType,
				Reduce: reduceInstanceResetFeatures,
			},
			{
				Event:  feature_v2.InstanceLoginDefaultOrgEventType,
				Reduce: reduceInstanceSetFeature[bool],
			},
			{
				Event:  feature_v2.InstanceTriggerIntrospectionProjectionsEventType,
				Reduce: reduceInstanceSetFeature[bool],
			},
			{
				Event:  feature_v2.InstanceLegacyIntrospectionEventType,
				Reduce: reduceInstanceSetFeature[bool],
			},
			{
				Event:  feature_v2.InstanceUserSchemaEventType,
				Reduce: reduceInstanceSetFeature[bool],
			},
			{
				Event:  feature_v2.InstanceTokenExchangeEventType,
				Reduce: reduceInstanceSetFeature[bool],
			},
			{
				Event:  feature_v2.InstanceActionsEventType,
				Reduce: reduceInstanceSetFeature[bool],
			},
			{
				Event:  feature_v2.InstanceImprovedPerformanceEventType,
				Reduce: reduceInstanceSetFeature[[]feature.ImprovedPerformanceType],
			},
			{
				Event:  feature_v2.InstanceWebKeyEventType,
				Reduce: reduceInstanceSetFeature[bool],
			},
			{
				Event:  feature_v2.InstanceDebugOIDCParentErrorEventType,
				Reduce: reduceInstanceSetFeature[bool],
			},
			{
				Event:  feature_v2.InstanceOIDCSingleV1SessionTerminationEventType,
				Reduce: reduceInstanceSetFeature[bool],
			},
			{
				Event:  feature_v2.InstanceDisableUserTokenEvent,
				Reduce: reduceInstanceSetFeature[bool],
			},
			{
				Event:  feature_v2.InstanceEnableBackChannelLogout,
				Reduce: reduceInstanceSetFeature[bool],
			},
			{
				Event:  instance.InstanceRemovedEventType,
				Reduce: reduceInstanceRemovedHelper(InstanceDomainInstanceIDCol),
			},
		},
	}}
}

func reduceSetDefaultLoginInstance_v1(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*feature_v1.SetEvent[feature_v1.Boolean])
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-in2Xo", "reduce.wrong.event.type %T", event)
	}
	return reduceInstanceSetFeature[bool](
		feature_v1.DefaultLoginInstanceEventToV2(e),
	)
}

func reduceInstanceSetFeature[T any](event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*feature_v2.SetEvent[T])
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-uPh8O", "reduce.wrong.event.type %T", event)
	}
	f, err := e.FeatureJSON()
	if err != nil {
		return nil, err
	}
	columns := []handler.Column{
		handler.NewCol(InstanceFeatureInstanceIDCol, e.Aggregate().ID),
		handler.NewCol(InstanceFeatureKeyCol, f.Key.String()),
		handler.NewCol(InstanceFeatureCreationDateCol, handler.OnlySetValueOnInsert(InstanceFeatureTable, e.CreationDate())),
		handler.NewCol(InstanceFeatureChangeDateCol, e.CreationDate()),
		handler.NewCol(InstanceFeatureSequenceCol, e.Sequence()),
		handler.NewCol(InstanceFeatureValueCol, f.Value),
	}
	return handler.NewUpsertStatement(e, columns[0:2], columns), nil
}

func reduceInstanceResetFeatures(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*feature_v2.ResetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-roo6A", "reduce.wrong.event.type %T", event)
	}
	return handler.NewDeleteStatement(e, []handler.Condition{
		handler.NewCond(InstanceFeatureInstanceIDCol, e.Aggregate().ID),
	}), nil
}
