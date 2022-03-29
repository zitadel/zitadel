package setup

import (
	"context"
	"database/sql"
	_ "embed"
)

const (
	createAssets = `
CREATE TABLE system.assets (
    instance_id TEXT,
    asset_type TEXT,
    resource_owner TEXT,
    name TEXT,
    content_type TEXT,
    hash TEXT,
    data BYTES,
    updated_at TIMESTAMPTZ,

    PRIMARY KEY (instance_id, resource_owner, name)
);
`
)

type AssetTable struct {
	dbClient *sql.DB
}

func (mig *AssetTable) Execute(ctx context.Context) error {
	_, err := mig.dbClient.ExecContext(ctx, createAssets)
	return err
}

func (mig *AssetTable) String() string {
	return "02_assets"
}
