package mock

//go:generate mockgen -package mock -destination queue.mock.go github.com/zitadel/zitadel/internal/serviceping Queue
//go:generate mockgen -package mock -destination queries.mock.go github.com/zitadel/zitadel/internal/serviceping Queries
//go:generate mockgen -package mock -destination telemetry.mock.go github.com/zitadel/zitadel/pkg/grpc/analytics/v2beta TelemetryServiceClient
