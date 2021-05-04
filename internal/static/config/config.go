package config

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/static"
	"github.com/caos/zitadel/internal/static/s3"
)

type AssetStorageConfig struct {
	Type   string
	Config static.Config
}

var storage = map[string]func() static.Config{
	"s3": func() static.Config { return &s3.Config{} },
}

func (c *AssetStorageConfig) UnmarshalJSON(data []byte) error {
	var rc struct {
		Type   string
		Config json.RawMessage
	}

	if err := json.Unmarshal(data, &rc); err != nil {
		return errors.ThrowInternal(err, "STATIC-Bfn5r", "error parsing config")
	}

	c.Type = rc.Type

	var err error
	c.Config, err = newStorageConfig(c.Type, rc.Config)
	if err != nil {
		return err
	}

	return nil
}

func newStorageConfig(storageType string, configData []byte) (static.Config, error) {
	t, ok := storage[storageType]
	if !ok {
		return nil, errors.ThrowInternalf(nil, "STATIC-dsbjh", "config type %s not supported", storageType)
	}

	staticConfig := t()
	if len(configData) == 0 {
		return staticConfig, nil
	}

	if err := json.Unmarshal(configData, staticConfig); err != nil {
		return nil, errors.ThrowInternal(err, "STATIC-GB4nw", "Could not read config: %v")
	}

	return staticConfig, nil
}
