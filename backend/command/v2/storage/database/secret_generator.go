package database

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
)

type query struct{ Querier }

func Query(querier Querier) *query {
	return &query{Querier: querier}
}

const getEncryptionConfigQuery = "SELECT" +
	" length" +
	", expiry" +
	", should_include_lower_letters" +
	", should_include_upper_letters" +
	", should_include_digits" +
	", should_include_symbols" +
	" FROM encryption_config"

func (q query) GetEncryptionConfig(ctx context.Context) (*crypto.GeneratorConfig, error) {
	var config crypto.GeneratorConfig
	row := q.QueryRow(ctx, getEncryptionConfigQuery)
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
