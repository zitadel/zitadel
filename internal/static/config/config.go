package config

import (
	"database/sql"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/static"
	"github.com/caos/zitadel/internal/static/database"
	"github.com/caos/zitadel/internal/static/s3"
)

type AssetStorageConfig struct {
	Type   string
	Config map[string]interface{} `mapstructure:",remain"`
}

func (a *AssetStorageConfig) NewStorage(client *sql.DB) (static.Storage, error) {
	t, ok := storage[a.Type]
	if !ok {
		return nil, errors.ThrowInternalf(nil, "STATIC-dsbjh", "config type %s not supported")
	}

	return t(client, a.Config)
}

var storage = map[string]static.CreateStorage{
	"db": database.NewStorage,
	"":   database.NewStorage,
	"s3": s3.NewStorage,
}
