package domain_test

// import (
// 	"context"
// 	"log/slog"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
// 	sdktrace "go.opentelemetry.io/otel/sdk/trace"
// 	"go.uber.org/mock/gomock"

// 	. "github.com/zitadel/zitadel/backend/v3/domain"
// 	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
// 	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
// 	"github.com/zitadel/zitadel/backend/v3/telemetry/logging"
// 	"github.com/zitadel/zitadel/backend/v3/telemetry/tracing"
// )

// These tests give an overview of how to use the domain package.
// func TestExample(t *testing.T) {
// 	t.Skip("skip example test because it is not a real test")
// 	ctx := context.Background()

// 	ctrl := gomock.NewController(t)
// 	pool := dbmock.NewMockPool(ctrl)
// 	tx := dbmock.NewMockTransaction(ctrl)

// 	pool.EXPECT().Begin(gomock.Any(), gomock.Any()).Return(tx, nil)
// 	tx.EXPECT().End(gomock.Any(), gomock.Any()).Return(nil)
// 	SetPool(pool)

// 	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
// 	require.NoError(t, err)
// 	tracerProvider := sdktrace.NewTracerProvider(
// 		sdktrace.WithSyncer(exporter),
// 	)
// 	otel.SetTracerProvider(tracerProvider)
// 	SetTracer(tracing.Tracer{Tracer: tracerProvider.Tracer("test")})
// 	defer func() { assert.NoError(t, tracerProvider.Shutdown(ctx)) }()

// 	SetLogger(logging.Logger{Logger: slog.Default()})

// 	SetUserRepository(repository.UserRepository)
// 	SetOrgRepository(repository.OrgRepository)
// 	// SetInstanceRepository(repository.Instance)
// 	// SetCryptoRepository(repository.Crypto)

// 	t.Run("create org", func(t *testing.T) {
// 		org := NewAddOrgCommand("testorg", NewAddMemberCommand("testuser", "ORG_OWNER"))
// 		user := NewCreateHumanCommand("testuser")
// 		err := Invoke(ctx, BatchCommands(org, user))
// 		assert.NoError(t, err)
// 	})

// 	t.Run("verified email", func(t *testing.T) {
// 		err := Invoke(ctx, NewSetEmailCommand("u1", "test@example.com", NewEmailVerifiedCommand("u1", true)))
// 		assert.NoError(t, err)
// 	})

// 	t.Run("unverified email", func(t *testing.T) {
// 		err := Invoke(ctx, NewSetEmailCommand("u2", "test2@example.com", NewEmailVerifiedCommand("u2", false)))
// 		assert.NoError(t, err)
// 	})
// }
