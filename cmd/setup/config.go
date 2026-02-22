package setup

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	old_logging "github.com/zitadel/logging" //nolint:staticcheck

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/cmd/encryption"
	"github.com/zitadel/zitadel/cmd/hooks"
	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/cache/connector"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/hook"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/execution"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	"github.com/zitadel/zitadel/internal/query/projection"
	static_config "github.com/zitadel/zitadel/internal/static/config"
)

type Config struct {
	ForMirror       bool
	Database        database.Config
	Caches          *connector.CachesConfig
	SystemDefaults  systemdefaults.SystemDefaults
	InternalAuthZ   authz.Config
	SystemAuthZ     authz.Config
	ExternalDomain  string
	ExternalPort    uint16
	ExternalSecure  bool
	Instrumentation instrumentation.Config
	Log             *old_logging.Config
	Metrics         *instrumentation.LegacyMetricConfig
	EncryptionKeys  *encryption.EncryptionKeyConfig
	DefaultInstance command.InstanceSetup
	Machine         *id.Config
	Projections     projection.Config
	Notifications   handlers.WorkerConfig
	Executions      execution.WorkerConfig
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

func NewConfig(cmd *cobra.Command, v *viper.Viper) (*Config, instrumentation.ShutdownFunc, error) {
	config := new(Config)
	err := v.Unmarshal(config,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hooks.SliceTypeStringDecode[*domain.CustomMessageText],
			hooks.SliceTypeStringDecode[authz.RoleMapping],
			hooks.MapTypeStringDecode[string, *authz.SystemAPIUser],
			hooks.MapHTTPHeaderStringDecode,
			database.DecodeHook(false),
			actions.HTTPConfigDecodeHook,
			hook.EnumHookFunc(authz.MemberTypeString),
			hook.Base64ToBytesHookFunc(),
			hook.TagToLanguageHookFunc(),
			hook.StringToURLHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.TextUnmarshallerHookFunc(),
		)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read default config: %w", err)
	}

	config.Instrumentation.Metric.SetLegacyConfig(config.Metrics)
	config.Instrumentation.Log.SetLegacyConfig(config.Log)
	shutdown, err := instrumentation.Start(cmd.Context(), config.Instrumentation)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start instrumentation: %w", err)
	}
	cmd.SetContext(logging.NewCtx(cmd.Context(), logging.StreamRuntime))

	err = config.Log.SetLogger()
	if err != nil {
		err = errors.Join(err, shutdown(cmd.Context()))
		return nil, nil, fmt.Errorf("unable to set logger: %w", err)
	}

	id.Configure(config.Machine)

	// Copy the global role permissions mappings to the instance until we allow instance-level configuration over the API.
	config.DefaultInstance.RolePermissionMappings = config.InternalAuthZ.RolePermissionMappings

	return config, shutdown, nil
}

type Steps struct {
	s1ProjectionTable                       *ProjectionTable
	s2AssetsTable                           *AssetTable
	FirstInstance                           *FirstInstance
	s5LastFailed                            *LastFailed
	s6OwnerRemoveColumns                    *OwnerRemoveColumns
	s7LogstoreTables                        *LogstoreTables
	s8AuthTokens                            *AuthTokenIndexes
	CorrectCreationDate                     *CorrectCreationDate
	s12AddOTPColumns                        *AddOTPColumns
	s13FixQuotaProjection                   *FixQuotaConstraints
	s14NewEventsTable                       *NewEventsTable
	s15CurrentStates                        *CurrentProjectionState
	s16UniqueConstraintsLower               *UniqueConstraintToLower
	s17AddOffsetToUniqueConstraints         *AddOffsetToCurrentStates
	s18AddLowerFieldsToLoginNames           *AddLowerFieldsToLoginNames
	s19AddCurrentStatesIndex                *AddCurrentSequencesIndex
	s20AddByUserSessionIndex                *AddByUserIndexToSession
	s21AddBlockFieldToLimits                *AddBlockFieldToLimits
	s22ActiveInstancesIndex                 *ActiveInstanceEvents
	s23CorrectGlobalUniqueConstraints       *CorrectGlobalUniqueConstraints
	s24AddActorToAuthTokens                 *AddActorToAuthTokens
	s25User11AddLowerFieldsToVerifiedEmail  *User11AddLowerFieldsToVerifiedEmail
	s26AuthUsers3                           *AuthUsers3
	s27IDPTemplate6SAMLNameIDFormat         *IDPTemplate6SAMLNameIDFormat
	s28AddFieldTable                        *AddFieldTable
	s29FillFieldsForProjectGrant            *FillFieldsForProjectGrant
	s30FillFieldsForOrgDomainVerified       *FillFieldsForOrgDomainVerified
	s31AddAggregateIndexToFields            *AddAggregateIndexToFields
	s32AddAuthSessionID                     *AddAuthSessionID
	s33SMSConfigs3TwilioAddVerifyServiceSid *SMSConfigs3TwilioAddVerifyServiceSid
	s34AddCacheSchema                       *AddCacheSchema
	s35AddPositionToIndexEsWm               *AddPositionToIndexEsWm
	s36FillV2Milestones                     *FillV3Milestones
	s37Apps7OIDConfigsBackChannelLogoutURI  *Apps7OIDConfigsBackChannelLogoutURI
	s38BackChannelLogoutNotificationStart   *BackChannelLogoutNotificationStart
	s40InitPushFunc                         *InitPushFunc
	s42Apps7OIDCConfigsLoginVersion         *Apps7OIDCConfigsLoginVersion
	s43CreateFieldsDomainIndex              *CreateFieldsDomainIndex
	s44ReplaceCurrentSequencesIndex         *ReplaceCurrentSequencesIndex
	s45CorrectProjectOwners                 *CorrectProjectOwners
	s46InitPermissionFunctions              *InitPermissionFunctions
	s47FillMembershipFields                 *FillMembershipFields
	s48Apps7SAMLConfigsLoginVersion         *Apps7SAMLConfigsLoginVersion
	s49InitPermittedOrgsFunction            *InitPermittedOrgsFunction
	s50IDPTemplate6UsePKCE                  *IDPTemplate6UsePKCE
	s51IDPTemplate6RootCA                   *IDPTemplate6RootCA
	s52IDPTemplate6LDAP2                    *IDPTemplate6LDAP2
	s53InitPermittedOrgsFunction            *InitPermittedOrgsFunction53
	s54InstancePositionIndex                *InstancePositionIndex
	s55ExecutionHandlerStart                *ExecutionHandlerStart
	s56IDPTemplate6SAMLFederatedLogout      *IDPTemplate6SAMLFederatedLogout
	s57CreateResourceCounts                 *CreateResourceCounts
	s58ReplaceLoginNames3View               *ReplaceLoginNames3View
	s59SetupWebkeys                         *SetupWebkeys
	s60GenerateSystemID                     *GenerateSystemID
	s61IDPTemplate6SAMLSignatureAlgorithm   *IDPTemplate6SAMLSignatureAlgorithm
	s62HTTPProviderAddSigningKey            *HTTPProviderAddSigningKey
	s63AlterResourceCounts                  *AlterResourceCounts
	s64ChangePushPosition                   *ChangePushPosition
	s65FixUserMetadata5Index                *FixUserMetadata5Index
	s66SessionRecoveryCodeCheckedAt         *SessionRecoveryCodeCheckedAt
	s67SyncMemberRoleFields                 *SyncMemberRoleFields
	s68TargetAddPayloadTypeColumn           *TargetAddPayloadTypeColumn
	s69CacheTablesLogged                    *CacheTablesLogged
}

func NewSteps(ctx context.Context, v *viper.Viper) (*Steps, error) {
	v.AutomaticEnv()
	v.SetEnvPrefix("ZITADEL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(defaultSteps))
	if err != nil {
		return nil, fmt.Errorf("unable to read default steps: %w", err)
	}

	for _, file := range stepFiles {
		v.SetConfigFile(file)
		err := v.MergeInConfig()
		logging.OnError(ctx, err).Warn("unable to read setup file", "file", file)
	}

	steps := new(Steps)
	err = v.Unmarshal(steps,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hook.Base64ToBytesHookFunc(),
			hook.TagToLanguageHookFunc(),
			hook.StringToURLHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.TextUnmarshallerHookFunc(),
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to read steps: %w", err)
	}
	return steps, nil
}
