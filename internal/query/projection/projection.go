package projection

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
)

const (
	currentSeqTable   = "projections.current_sequences"
	locksTable        = "projections.locks"
	failedEventsTable = "projections.failed_events"
)

func Start(ctx context.Context, es *eventstore.Eventstore, config Config) error {
	sqlClient, err := config.CRDB.Start()
	if err != nil {
		return err
	}

	projectionConfig := crdb.StatementHandlerConfig{
		ProjectionHandlerConfig: handler.ProjectionHandlerConfig{
			HandlerConfig: handler.HandlerConfig{
				Eventstore: es,
			},
			RequeueEvery: config.RequeueEvery.Duration,
		},
		Client:            sqlClient,
		SequenceTable:     currentSeqTable,
		LockTable:         locksTable,
		FailedEventsTable: failedEventsTable,
		MaxFailureCount:   config.MaxFailureCount,
		BulkLimit:         config.BulkLimit,
	}

	NewOrgProjection(ctx, projectionConfig)
	return nil
}
