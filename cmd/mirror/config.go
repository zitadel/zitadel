package mirror

import (
	_ "embed"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/v2/cmd/hooks"
	"github.com/zitadel/zitadel/v2/internal/actions"
	internal_authz "github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/command"
	"github.com/zitadel/zitadel/v2/internal/config/hook"
	"github.com/zitadel/zitadel/v2/internal/database"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/id"
)

type Migration struct {
	Source      database.Config
	Destination database.Config

	EventBulkSize uint32

	Log     *logging.Config
	Machine *id.Config
}

var (
	//go:embed defaults.yaml
	defaultConfig []byte
)

func mustNewMigrationConfig(v *viper.Viper) *Migration {
	config := new(Migration)
	mustNewConfig(v, config)

	err := config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	id.Configure(config.Machine)

	return config
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
			database.DecodeHook,
			actions.HTTPConfigDecodeHook,
			hook.EnumHookFunc(internal_authz.MemberTypeString),
			mapstructure.TextUnmarshallerHookFunc(),
		)),
	)
	logging.OnError(err).Fatal("unable to read default config")
}
