package domain

import (
	"github.com/zitadel/zitadel/backend/command/v2/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"go.opentelemetry.io/otel/trace"
)

type Domain struct {
	pool        database.Pool
	tracer      trace.Tracer
	userCodeAlg crypto.EncryptionAlgorithm
}
