package initialise

import (
	"context"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	old_config "github.com/zitadel/logging"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/id"
)

type Config struct {
	Instrumentation instrumentation.Config
	Database        database.Config
	Machine         *id.Config
	Log             *old_config.Config
}

func NewConfig(ctx context.Context, v *viper.Viper) (*Config, instrumentation.ShutdownFunc, error) {
	config := new(Config)
	err := v.Unmarshal(config,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			database.DecodeHook(false),
			mapstructure.TextUnmarshallerHookFunc(),
		)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read config: %w", err)
	}
	shutdown, err := instrumentation.Start(ctx, config.Instrumentation)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start instrumentation: %w", err)
	}

	err = config.Log.SetLogger()
	if err != nil {
		err = errors.Join(err, shutdown(ctx))
		return nil, nil, fmt.Errorf("unable to set logger: %w", err)
	}

	id.Configure(config.Machine)

	return config, shutdown, nil
}
