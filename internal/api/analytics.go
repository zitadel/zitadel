package api

// import (
// 	"context"

// 	"github.com/zitadel/zitadel/internal/repository/analytics"
// 	"github.com/zitadel/zitadel/internal/telemetry/tracing"
// )

// const (
// 	analyticsEventsTable = "analytics.events"
// 	insertEventStmt      = "INSERT INTO " + analyticsEventsTable + " (type, event_data) VALUES ($1, $2)"
// )

// func (q *API) InsertAnalyticsEvent(ctx context.Context, eventType string, data []byte) (err error) {
// 	ctx, span := tracing.NewSpan(ctx)
// 	defer func() { span.EndWithError(err) }()

// 	_, err = q.client.ExecContext(ctx, insertEventStmt, eventType, data)
// 	return analytics.With(err, "QUERY-aGhe4", "unable to insert analytics event")
// }
