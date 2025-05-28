package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/org"
)

const (
	HostedLoginTranslationTable = "projections.hosted_login_translations"

	HostedLoginTranslationInstaceIDCol     = "instance_id"
	HostedLoginTranslationCreationDateCol  = "creation_date"
	HostedLoginTranslationChangeDateCol    = "change_date"
	HostedLoginTranslationAggregateIDCol   = "aggregate_id"
	HostedLoginTranslationAggregateTypeCol = "aggregate_type"
	HostedLoginTranslationSequenceCol      = "sequence"
	HostedLoginTranslationLocaleCol        = "locale"
	HostedLoginTranslationFileCol          = "file"
)

type hostedLoginTranslationProjection struct{}

func newHostedLoginTranslationProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(hostedLoginTranslationProjection))
}

func (hltp *hostedLoginTranslationProjection) Name() string {
	return HostedLoginTranslationTable
}

func (hltp *hostedLoginTranslationProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{},
			},
		},
	}
}
