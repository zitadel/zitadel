package mirror

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/cmd/hooks"
	"github.com/zitadel/zitadel/internal/actions"
	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/hook"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/id"
)

type Migration struct {
	Source      database.Config
	Destination database.Config

	EventBulkSize     uint32
	MaxAuthRequestAge time.Duration

	Log             *logging.Config
	Machine         *id.Config
	Instrumentation instrumentation.Config
}

var (
	//go:embed defaults.yaml
	defaultConfig []byte
)

func mustNewMigrationConfig(ctx context.Context, v *viper.Viper) (*Migration, instrumentation.ShutdownFunc, error) {
	config := new(Migration)
	mustNewConfig(v, config)

	shutdown, err := instrumentation.Start(ctx, config.Instrumentation)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start instrumentation: %w", err)
	}

	// Legacy logger
	err = config.Log.SetLogger()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to set logger: %w", err)
	}

	id.Configure(config.Machine)

	return config, shutdown, nil
}

func mustNewProjectionsConfig(v *viper.Viper) *ProjectionsConfig {
	config := new(ProjectionsConfig)
	mustNewConfig(v, config)

	err := config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	id.Configure(config.Machine)

	return config
}

func mustNewConfig(v *viper.Viper, config any) {
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
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
			database.DecodeHook(true),
			actions.HTTPConfigDecodeHook,
			hook.EnumHookFunc(internal_authz.MemberTypeString),
			mapstructure.TextUnmarshallerHookFunc(),
		)),
	)
	logging.OnError(err).Fatal("unable to read default config")
}
