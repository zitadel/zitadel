package setup

import (
	"bytes"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/encryption"
	"github.com/zitadel/zitadel/cmd/hooks"
	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/hook"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	"github.com/zitadel/zitadel/internal/query/projection"
	static_config "github.com/zitadel/zitadel/internal/static/config"
)

type Config struct {
	Database        database.Config
	SystemDefaults  systemdefaults.SystemDefaults
	InternalAuthZ   authz.Config
	ExternalDomain  string
	ExternalPort    uint16
	ExternalSecure  bool
	Log             *logging.Config
	EncryptionKeys  *encryption.EncryptionKeyConfig
	DefaultInstance command.InstanceSetup
	Machine         *id.Config
	Projections     projection.Config
	Eventstore      *eventstore.Config

	InitProjections InitProjections
	AssetStorage    static_config.AssetStorageConfig
	OIDC            oidc.Config
	Login           login.Config
	WebAuthNName    string
	Telemetry       *handlers.TelemetryPusherConfig
	SystemAPIUsers  map[string]*authz.SystemAPIUser
}

type InitProjections struct {
	Enabled          bool
	RetryFailedAfter time.Duration
	MaxFailureCount  uint8
	BulkLimit        uint64
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := v.Unmarshal(config,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hook.Base64ToBytesHookFunc(),
			hook.TagToLanguageHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
			database.DecodeHook,
			hook.EnumHookFunc(authz.MemberTypeString),
			actions.HTTPConfigDecodeHook,
			hooks.MapTypeStringDecode[string, *authz.SystemAPIUser],
			hooks.SliceTypeStringDecode[authz.RoleMapping],
		)),
	)
	logging.OnError(err).Fatal("unable to read default config")

	err = config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	id.Configure(config.Machine)

	return config
}

type Steps struct {
	s1ProjectionTable                      *ProjectionTable
	s2AssetsTable                          *AssetTable
	FirstInstance                          *FirstInstance
	s5LastFailed                           *LastFailed
	s6OwnerRemoveColumns                   *OwnerRemoveColumns
	s7LogstoreTables                       *LogstoreTables
	s8AuthTokens                           *AuthTokenIndexes
	CorrectCreationDate                    *CorrectCreationDate
	s12AddOTPColumns                       *AddOTPColumns
	s13FixQuotaProjection                  *FixQuotaConstraints
	s14NewEventsTable                      *NewEventsTable
	s15CurrentStates                       *CurrentProjectionState
	s16UniqueConstraintsLower              *UniqueConstraintToLower
	s17AddOffsetToUniqueConstraints        *AddOffsetToCurrentStates
	s18AddLowerFieldsToLoginNames          *AddLowerFieldsToLoginNames
	s19AddCurrentStatesIndex               *AddCurrentSequencesIndex
	s20AddByUserSessionIndex               *AddByUserIndexToSession
	s21AddBlockFieldToLimits               *AddBlockFieldToLimits
	s22ActiveInstancesIndex                *ActiveInstanceEvents
	s23CorrectGlobalUniqueConstraints      *CorrectGlobalUniqueConstraints
	s24AddActorToAuthTokens                *AddActorToAuthTokens
	s25User11AddLowerFieldsToVerifiedEmail *User11AddLowerFieldsToVerifiedEmail
}

func MustNewSteps(v *viper.Viper) *Steps {
	v.AutomaticEnv()
	v.SetEnvPrefix("ZITADEL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(defaultSteps))
	logging.OnError(err).Fatal("unable to read setup steps")

	for _, file := range stepFiles {
		v.SetConfigFile(file)
		err := v.MergeInConfig()
		logging.WithFields("file", file).OnError(err).Warn("unable to read setup file")
	}

	steps := new(Steps)
	err = v.Unmarshal(steps,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hook.Base64ToBytesHookFunc(),
			hook.TagToLanguageHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
		)),
	)
	logging.OnError(err).Fatal("unable to read steps")
	return steps
}
