package projection

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler/v3"
	"github.com/caos/zitadel/internal/query/projection/org/owner"
)

func Start(ctx context.Context, es *eventstore.Eventstore, config Config) error {
	sqlClient, err := config.CRDB.Start()
	if err != nil {
		return err
	}

	handlerConfig := handler.HandlerConfig{
		IteratorConfig: handler.IteratorConfig{
			Client:     sqlClient,
			Eventstore: es,
			Interval:   config.RequeueEvery.Duration,
			BulkLimit:  config.BulkLimit,
			Pool:       iteratorPool(config.MaxIterators),
		},
		PusherConfig: handler.PusherConfig{
			Client:          sqlClient,
			Interval:        config.RetryFailedAfter.Duration,
			MaxFailureCount: config.MaxFailureCount,
		},
	}

	NewOrgProjection(ctx, applyCustomConfig(handlerConfig, config.Customizations["orgs"]))
	NewProjectProjection(ctx, applyCustomConfig(handlerConfig, config.Customizations["projects"]))
	owner.NewOrgOwnerProjection(ctx, applyCustomConfig(handlerConfig, config.Customizations["org_owners"]))
	return nil
}

func applyCustomConfig(config handler.HandlerConfig, customConfig CustomConfig) handler.HandlerConfig {
	if customConfig.BulkLimit != nil {
		config.IteratorConfig.BulkLimit = *customConfig.BulkLimit
	}
	if customConfig.MaxFailureCount != nil {
		config.PusherConfig.MaxFailureCount = *customConfig.MaxFailureCount
	}
	if customConfig.RequeueEvery != nil {
		config.IteratorConfig.Interval = customConfig.RequeueEvery.Duration
	}
	if customConfig.RetryFailedAfter != nil {
		config.PusherConfig.Interval = customConfig.RetryFailedAfter.Duration
	}

	return config
}

func iteratorPool(workerCount int) chan func() {
	if workerCount <= 0 {
		return nil
	}

	queue := make(chan func())
	for i := 0; i < workerCount; i++ {
		go func() {
			for iteration := range queue {
				iteration()
				time.Sleep(2 * time.Second)
			}
		}()
	}
	return queue
}
