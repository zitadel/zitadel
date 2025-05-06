package domain_test

// import (
// 	"context"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
// 	sdktrace "go.opentelemetry.io/otel/sdk/trace"

// 	. "github.com/zitadel/zitadel/backend/v3/domain"
// 	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
// 	"github.com/zitadel/zitadel/backend/v3/telemetry/tracing"
// )

// func TestExample(t *testing.T) {
// 	ctx := context.Background()

// 	// SetPool(pool)

// 	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
// 	require.NoError(t, err)
// 	tracerProvider := sdktrace.NewTracerProvider(
// 		sdktrace.WithSyncer(exporter),
// 	)
// 	otel.SetTracerProvider(tracerProvider)
// 	SetTracer(tracing.Tracer{Tracer: tracerProvider.Tracer("test")})
// 	defer func() { assert.NoError(t, tracerProvider.Shutdown(ctx)) }()

// 	SetUserRepository(repository.User)
// 	SetInstanceRepository(repository.Instance)
// 	SetCryptoRepository(repository.Crypto)

// 	t.Run("verified email", func(t *testing.T) {
// 		err := Invoke(ctx, NewSetEmailCommand("u1", "test@example.com", NewEmailVerifiedCommand("u1", true)))
// 		assert.NoError(t, err)
// 	})

// 	t.Run("unverified email", func(t *testing.T) {
// 		err := Invoke(ctx, NewSetEmailCommand("u2", "test2@example.com", NewEmailVerifiedCommand("u2", false)))
// 		assert.NoError(t, err)
// 	})
// }
