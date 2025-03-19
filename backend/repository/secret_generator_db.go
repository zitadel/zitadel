package repository

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
)

const secretGeneratorByTypeStmt = `SELECT * FROM secret_generators WHERE instance_id = $1 AND type = $2`

func (q querier) SecretGeneratorConfigByType(ctx context.Context, typ SecretGeneratorType) (config *crypto.GeneratorConfig, err error) {
	err = q.client.QueryRow(ctx, secretGeneratorByTypeStmt, authz.GetInstance(ctx).InstanceID, typ).Scan(
		&config.Length,
		&config.Expiry,
		&config.IncludeLowerLetters,
		&config.IncludeUpperLetters,
		&config.IncludeDigits,
		&config.IncludeSymbols,
	)
	if err != nil {
		return nil, err
	}
	return config, nil
}
