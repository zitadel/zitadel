package projection

import (
	"context"

	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
)

const (
	InstanceFeatureTable = "projections.instance_features"

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
	return []handler.AggregateReducer{}
}
