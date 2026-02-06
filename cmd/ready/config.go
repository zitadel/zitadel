package ready

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/config/hook"
	"github.com/zitadel/zitadel/internal/config/network"
)

type Config struct {
	Instrumentation instrumentation.Config
	Port            uint16
	TLS             network.TLS
}

func newConfig(cmd *cobra.Command, v *viper.Viper) (*Config, instrumentation.ShutdownFunc, error) {
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

	shutdown, err := instrumentation.Start(cmd.Context(), config.Instrumentation)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start instrumentation: %w", err)
	}
	cmd.SetContext(logging.NewCtx(cmd.Context(), logging.StreamRuntime))

	return config, shutdown, nil
}
