package projection

import (
	"context"

	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
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

func newsystemFeatureProjection(ctx context.Context, config handler.Config) *handler.Handler {
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
	return []handler.AggregateReducer{}
}
