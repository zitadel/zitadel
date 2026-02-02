package ready

import (
	"context"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	old_logging "github.com/zitadel/logging" //nolint:staticcheck

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/config/hook"
	"github.com/zitadel/zitadel/internal/config/network"
)

type Config struct {
	Instrumentation instrumentation.Config
	Log             *old_logging.Config
	Port            uint16
	TLS             network.TLS
}

func newConfig(ctx context.Context, v *viper.Viper) (*Config, instrumentation.ShutdownFunc, error) {
	config := new(Config)
	err := v.Unmarshal(config,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hook.Base64ToBytesHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
			hook.EnumHookFunc(internal_authz.MemberTypeString),
			mapstructure.TextUnmarshallerHookFunc(),
		)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read default config: %w", err)
	}
	// Force-disable metrics and tracing for ready command
	config.Instrumentation.Metric.Exporter.Type = instrumentation.ExporterTypeNone
	config.Instrumentation.Trace.Exporter.Type = instrumentation.ExporterTypeNone
	// Legacy logger
	err = config.Log.SetLogger()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to set logger: %w", err)
	}

	shutdown, err := instrumentation.Start(ctx, config.Instrumentation)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start instrumentation: %w", err)
	}

	return config, shutdown, nil
}
