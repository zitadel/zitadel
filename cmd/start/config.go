package start

import (
	"errors"
	"fmt"
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
	admin_es "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/api/saml"
	scim_config "github.com/zitadel/zitadel/internal/api/scim/config"
	"github.com/zitadel/zitadel/internal/api/ui/console"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	auth_es "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/cache/connector"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/hook"
	"github.com/zitadel/zitadel/internal/config/network"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/denylist"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/execution"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/serviceping"
	static_config "github.com/zitadel/zitadel/internal/static/config"
)

type Config struct {
	Instrumentation     instrumentation.Config
	Log                 *old_logging.Config
	Port                uint16
	ExternalPort        uint16
	ExternalDomain      string
	ExternalSecure      bool
	TLS                 network.TLS
	InstanceHostHeaders []string
	PublicHostHeaders   []string
	HTTP2HostHeader     string
	HTTP1HostHeader     string
	WebAuthNName        string
	Database            database.Config
	Caches              *connector.CachesConfig
	Tracing             *instrumentation.LegacyTraceConfig
	Metrics             *instrumentation.LegacyMetricConfig
	Profiler            *instrumentation.LegacyProfileConfig
	Projections         projection.Config
	Notifications       handlers.WorkerConfig
	Executions          execution.WorkerConfig
	Auth                auth_es.Config
	Admin               admin_es.Config
	UserAgentCookie     *middleware.UserAgentCookieConfig
	OIDC                oidc.Config
	SAML                saml.Config
	SCIM                scim_config.Config
	Login               login.Config
	Console             console.Config
	AssetStorage        static_config.AssetStorageConfig
	InternalAuthZ       authz.Config
	SystemAuthZ         authz.Config
	SystemDefaults      systemdefaults.SystemDefaults
	EncryptionKeys      *encryption.EncryptionKeyConfig
	DefaultInstance     command.InstanceSetup
	AuditLogRetention   time.Duration
	SystemAPIUsers      map[string]*authz.SystemAPIUser
	CustomerPortal      string
	Machine             *id.Config
	Actions             *actions.Config
	Eventstore          *eventstore.Config
	LogStore            *logstore.Configs
	Quotas              *QuotasConfig
	Telemetry           *handlers.TelemetryPusherConfig
	ServicePing         *serviceping.Config
}

type QuotasConfig struct {
	Access struct {
		logstore.EmitterConfig  `mapstructure:",squash"`
		middleware.AccessConfig `mapstructure:",squash"`
	}
	Execution *logstore.EmitterConfig
}

func NewConfig(cmd *cobra.Command, v *viper.Viper) (*Config, instrumentation.ShutdownFunc, error) {
	config, err := readConfig(v)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read config: %w", err)
	}

	config.Instrumentation.Trace.SetLegacyConfig(config.Tracing)
	config.Instrumentation.Metric.SetLegacyConfig(config.Metrics)
	config.Instrumentation.Log.SetLegacyConfig(config.Log)
	config.Instrumentation.Profile.SetLegacyConfig(config.Profiler)
	shutdown, err := instrumentation.Start(cmd.Context(), config.Instrumentation)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start instrumentation: %w", err)
	}
	cmd.SetContext(logging.NewCtx(cmd.Context(), logging.StreamRuntime))

	// Legacy logger
	err = config.Log.SetLogger()
	if err != nil {
		err = errors.Join(err, shutdown(cmd.Context()))
		return nil, nil, fmt.Errorf("unable to set logger: %w", err)
	}

	id.Configure(config.Machine)
	if config.Actions != nil {
		actions.SetHTTPConfig(&config.Actions.HTTP)
	}

	err = config.SystemDefaults.Validate()
	if err != nil {
		err = errors.Join(err, shutdown(cmd.Context()))
		return nil, nil, fmt.Errorf("system defaults config invalid: %w", err)
	}
	// Copy the global role permissions mappings to the instance until we allow instance-level configuration over the API.
	config.DefaultInstance.RolePermissionMappings = config.InternalAuthZ.RolePermissionMappings

	return config, shutdown, nil
}

func readConfig(v *viper.Viper) (*Config, error) {
	config := new(Config)

	err := v.Unmarshal(config,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hooks.SliceTypeStringDecode[*domain.CustomMessageText],
			hooks.SliceTypeStringDecode[authz.RoleMapping],
			hooks.MapTypeStringDecode[string, *authz.SystemAPIUser],
			hooks.MapHTTPHeaderStringDecode,
			database.DecodeHook(false),
			actions.HTTPConfigDecodeHook,
			denylist.DenyListDecodeHook,
			hook.EnumHookFunc(authz.MemberTypeString),
			hooks.MapTypeStringDecode[domain.Feature, any],
			hooks.SliceTypeStringDecode[*command.SetQuota],
			hook.Base64ToBytesHookFunc(),
			hook.TagToLanguageHookFunc(),
			hook.StringToURLHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.TextUnmarshallerHookFunc(),
		)),
	)
	return config, err
}
