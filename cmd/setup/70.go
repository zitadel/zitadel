package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 70.sql
	assetsHashSha256 string
)

type AssetsHashSha256 struct {
	dbClient *database.DB
}

func (mig *AssetsHashSha256) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, assetsHashSha256)
	return err
}

func (mig *AssetsHashSha256) String() string {
	return "70_assets_hash_sha256"
}
