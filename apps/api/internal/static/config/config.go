package config

import (
	"database/sql"

	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/static"
	"github.com/zitadel/zitadel/internal/static/database"
	"github.com/zitadel/zitadel/internal/static/s3"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AssetStorageConfig struct {
	Type   string
	Cache  middleware.CacheConfig
	Config map[string]interface{} `mapstructure:",remain"`
}

func (a *AssetStorageConfig) NewStorage(client *sql.DB) (static.Storage, error) {
	t, ok := storage[a.Type]
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "STATIC-dsbjh", "config type %s not supported", a.Type)
	}

	return t(client, a.Config)
}

var storage = map[string]static.CreateStorage{
	"db": database.NewStorage,
	"":   database.NewStorage,
	"s3": s3.NewStorage,
}
