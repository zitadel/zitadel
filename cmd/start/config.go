package start

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

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
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/execution"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/serviceping"
	static_config "github.com/zitadel/zitadel/internal/static/config"
	metrics "github.com/zitadel/zitadel/internal/telemetry/metrics/config"
	profiler "github.com/zitadel/zitadel/internal/telemetry/profiler/config"
	tracing "github.com/zitadel/zitadel/internal/telemetry/tracing/config"
)

type Config struct {
	Log                 *logging.Config
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
	Tracing             tracing.Config
	Metrics             metrics.Config
	Profiler            profiler.Config
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

func MustNewConfig(v *viper.Viper) *Config {
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
			hooks.MapTypeStringDecode[domain.Feature, any],
			hooks.SliceTypeStringDecode[*command.SetQuota],
			hook.Base64ToBytesHookFunc(),
			hook.TagToLanguageHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.TextUnmarshallerHookFunc(),
		)),
	)
	logging.OnError(err).Fatal("unable to read config")

	err = config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	err = config.Tracing.NewTracer()
	logging.OnError(err).Fatal("unable to set tracer")

	err = config.Metrics.NewMeter()
	logging.OnError(err).Fatal("unable to set meter")

	err = config.Profiler.NewProfiler()
	logging.OnError(err).Fatal("unable to set profiler")

	id.Configure(config.Machine)
	if config.Actions != nil {
		actions.SetHTTPConfig(&config.Actions.HTTP)
	}

	// Copy the global role permissions mappings to the instance until we allow instance-level configuration over the API.
	config.DefaultInstance.RolePermissionMappings = config.InternalAuthZ.RolePermissionMappings

	return config
}
