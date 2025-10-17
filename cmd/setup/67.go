package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 67.sql
	analyticsEvents string
)

type AnalyticsEvents struct {
	dbClient *database.DB
}

func (mig *AnalyticsEvents) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, analyticsEvents)
	return err
}

func (mig *AnalyticsEvents) String() string {
	return "67_analytics_events"
}
