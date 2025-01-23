package setup

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/execution"
	"github.com/zitadel/zitadel/internal/query/projection"
)

type projectionTables struct {
	es *eventstore.Eventstore

	Version string `json:"version"`
}

func (mig *projectionTables) Check(lastRun map[string]interface{}) bool {
	currentVersion, _ := lastRun["version"].(string)
	return currentVersion != mig.Version
}

func (mig *projectionTables) Execute(ctx context.Context, _ eventstore.Event) error {
	if err := projection.Init(ctx); err != nil {
		return err
	}
	return execution.Init(ctx)
}

func (mig *projectionTables) String() string {
	return "projection_tables"
}
