package mirror

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	old_logging "github.com/zitadel/logging" //nolint:staticcheck

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/cmd/hooks"
	"github.com/zitadel/zitadel/internal/actions"
	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/hook"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/denylist"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/id"
)

type Migration struct {
	Source      database.Config
	Destination database.Config

	EventBulkSize     uint32
	MaxAuthRequestAge time.Duration

	Log             *old_logging.Config
	Machine         *id.Config
	Instrumentation instrumentation.Config
	Metrics         instrumentation.LegacyMetricConfig
}

var (
	//go:embed defaults.yaml
	defaultConfig []byte
)

func newMigrationConfig(cmd *cobra.Command, v *viper.Viper) (*Migration, instrumentation.ShutdownFunc, error) {
	config := new(Migration)
	err := newConfig(v, config)
	if err != nil {
		return nil, nil, err
	}
	shutdown, err := startInstrumentation(cmd, config.Instrumentation, config.Log)
	if err != nil {
		return nil, nil, err
	}
	id.Configure(config.Machine)
	return config, shutdown, nil
}

func newProjectionsConfig(cmd *cobra.Command, v *viper.Viper) (*ProjectionsConfig, instrumentation.ShutdownFunc, error) {
	config := new(ProjectionsConfig)
	err := newConfig(v, config)
	if err != nil {
		return nil, nil, err
	}
	shutdown, err := startInstrumentation(cmd, config.Instrumentation, config.Log)
	if err != nil {
		return nil, nil, err
	}
	id.Configure(config.Machine)
	return config, shutdown, nil
}

func startInstrumentation(cmd *cobra.Command, cfg instrumentation.Config, logConfig *old_logging.Config) (instrumentation.ShutdownFunc, error) {
	cfg.Log.SetLegacyConfig(logConfig)
	shutdown, err := instrumentation.Start(cmd.Context(), cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to start instrumentation: %w", err)
	}
	// Legacy logger
	err = logConfig.SetLogger()
	if err != nil {
		return nil, fmt.Errorf("unable to set logger: %w", err)
	}
	cmd.SetContext(logging.NewCtx(cmd.Context(), logging.StreamRuntime))
	return shutdown, nil
}

func newConfig(v *viper.Viper, config any) error {
	err := v.Unmarshal(config,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hooks.SliceTypeStringDecode[*domain.CustomMessageText],
			hooks.SliceTypeStringDecode[*command.SetQuota],
			hooks.SliceTypeStringDecode[internal_authz.RoleMapping],
			hooks.MapTypeStringDecode[string, *internal_authz.SystemAPIUser],
			hooks.MapTypeStringDecode[domain.Feature, any],
			hooks.MapHTTPHeaderStringDecode,
			hook.Base64ToBytesHookFunc(),
			hook.TagToLanguageHookFunc(),
			hook.StringToURLHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
			database.DecodeHook(true),
			actions.HTTPConfigDecodeHook,
			denylist.DenyListDecodeHook,
			hook.EnumHookFunc(internal_authz.MemberTypeString),
			mapstructure.TextUnmarshallerHookFunc(),
		)),
	)
	if err != nil {
		return fmt.Errorf("unable to decode config: %w", err)
	}
	return nil
}
