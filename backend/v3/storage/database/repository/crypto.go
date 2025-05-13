package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
)

type cryptoRepo struct {
	database.QueryExecutor
}

func Crypto(db database.QueryExecutor) domain.CryptoRepository {
	return &cryptoRepo{
		QueryExecutor: db,
	}
}

const getEncryptionConfigQuery = "SELECT" +
	" length" +
	", expiry" +
	", should_include_lower_letters" +
	", should_include_upper_letters" +
	", should_include_digits" +
	", should_include_symbols" +
	" FROM encryption_config"

func (repo *cryptoRepo) GetEncryptionConfig(ctx context.Context) (*crypto.GeneratorConfig, error) {
	var config crypto.GeneratorConfig
	row := repo.QueryRow(ctx, getEncryptionConfigQuery)
	err := row.Scan(
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
	return &config, nil
}
