package setup

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	createAssets = `
CREATE TABLE system.assets (
    instance_id TEXT,
    asset_type TEXT,
    resource_owner TEXT,
    name TEXT,
    content_type TEXT,
    hash TEXT GENERATED ALWAYS AS (md5(data)) STORED,
    data BYTEA,
    updated_at TIMESTAMPTZ,

    PRIMARY KEY (instance_id, resource_owner, name)
);
`
)

type AssetTable struct {
	dbClient database.Pool
}

func (mig *AssetTable) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, createAssets)
	return err
}

func (mig *AssetTable) String() string {
	return "02_assets"
}
