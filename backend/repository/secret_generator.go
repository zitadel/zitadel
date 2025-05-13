package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

type SecretGeneratorOptions struct{}

type SecretGenerator struct {
	options[SecretGeneratorOptions]
}

func NewSecretGenerator(opts ...Option[SecretGeneratorOptions]) *SecretGenerator {
	i := new(SecretGenerator)

	for _, opt := range opts {
		opt.apply(&i.options)
	}
	return i
}

type SecretGeneratorType = domain.SecretGeneratorType

func (sg *SecretGenerator) GeneratorConfigByType(ctx context.Context, client database.Querier, typ SecretGeneratorType) (*crypto.GeneratorConfig, error) {
	return tracing.Wrap(sg.tracer, "secretGenerator.GeneratorConfigByType",
		query(client).SecretGeneratorConfigByType,
	)(ctx, typ)
}
