package setup

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query/projection"
)

type projectionTables struct {
	es             *eventstore.Eventstore
	currentVersion string

	Version string `json:"version"`
}

func (mig *projectionTables) SetLastExecution(lastRun map[string]interface{}) {
	mig.currentVersion, _ = lastRun["version"].(string)
}

func (mig *projectionTables) Check() bool {
	return mig.currentVersion != mig.Version
}

func (mig *projectionTables) Execute(ctx context.Context, _ eventstore.Event) error {
	return projection.Init(ctx)
}

func (mig *projectionTables) String() string {
	return "projection_tables"
}
