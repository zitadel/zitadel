package mirror

import (
	_ "embed"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/hooks"
	"github.com/zitadel/zitadel/internal/actions"
	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/hook"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/id_generator/sonyflake"
)

type Migration struct {
	Source      database.Config
	Destination database.Config

	EventBulkSize uint32

	Log         *logging.Config
	IDGenerator id_generator.GeneratorType
	Machine     *sonyflake.Config
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

	id_generator.SetGeneratorWithConfig(config.IDGenerator, config.Machine)

	return config
}

func mustNewProjectionsConfig(v *viper.Viper) *ProjectionsConfig {
	config := new(ProjectionsConfig)
	mustNewConfig(v, config)

	err := config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	id_generator.SetGeneratorWithConfig(config.IDGenerator, config.Machine)

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
		)),
	)
	logging.OnError(err).Fatal("unable to read default config")
}
