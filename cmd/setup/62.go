package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 62.sql
	addHTTPProviderSigningKey string
)

type HTTPProviderAddSigningKey struct {
	dbClient *database.DB
}

func (mig *HTTPProviderAddSigningKey) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addHTTPProviderSigningKey)
	return err
}

func (mig *HTTPProviderAddSigningKey) String() string {
	return "62_http_provider_add_signing_key"
}
