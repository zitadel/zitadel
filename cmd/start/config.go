package start

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/actions"
	admin_es "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing"
	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/api/saml"
	"github.com/zitadel/zitadel/internal/api/ui/console"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	auth_es "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/hook"
	"github.com/zitadel/zitadel/internal/config/network"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	"github.com/zitadel/zitadel/internal/query/projection"
	static_config "github.com/zitadel/zitadel/internal/static/config"
	metrics "github.com/zitadel/zitadel/internal/telemetry/metrics/config"
	tracing "github.com/zitadel/zitadel/internal/telemetry/tracing/config"
)

type Config struct {
	Log               *logging.Config
	Port              uint16
	ExternalPort      uint16
	ExternalDomain    string
	ExternalSecure    bool
	TLS               network.TLS
	HTTP2HostHeader   string
	HTTP1HostHeader   string
	WebAuthNName      string
	Database          database.Config
	Tracing           tracing.Config
	Metrics           metrics.Config
	Projections       projection.Config
	Auth              auth_es.Config
	Admin             admin_es.Config
	UserAgentCookie   *middleware.UserAgentCookieConfig
	OIDC              oidc.Config
	SAML              saml.Config
	Login             login.Config
	Console           console.Config
	AssetStorage      static_config.AssetStorageConfig
	InternalAuthZ     internal_authz.Config
	SystemDefaults    systemdefaults.SystemDefaults
	EncryptionKeys    *encryptionKeyConfig
	DefaultInstance   command.InstanceSetup
	AuditLogRetention time.Duration
	SystemAPIUsers    SystemAPIUsers
	CustomerPortal    string
	Machine           *id.Config
	Actions           *actions.Config
	Eventstore        *eventstore.Config
	LogStore          *logstore.Configs
	Quotas            *QuotasConfig
	Telemetry         *handlers.TelemetryPusherConfig
}

type QuotasConfig struct {
	Access *middleware.AccessConfig
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
			actions.HTTPConfigDecodeHook,
			systemAPIUsersDecodeHook,
		)),
	)
	logging.OnError(err).Fatal("unable to read config")

	err = config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	err = config.Tracing.NewTracer()
	logging.OnError(err).Fatal("unable to set tracer")

	err = config.Metrics.NewMeter()
	logging.OnError(err).Fatal("unable to set meter")

	id.Configure(config.Machine)
	actions.SetHTTPConfig(&config.Actions.HTTP)

	return config
}

type encryptionKeyConfig struct {
	DomainVerification   *crypto.KeyConfig
	IDPConfig            *crypto.KeyConfig
	OIDC                 *crypto.KeyConfig
	SAML                 *crypto.KeyConfig
	OTP                  *crypto.KeyConfig
	SMS                  *crypto.KeyConfig
	SMTP                 *crypto.KeyConfig
	User                 *crypto.KeyConfig
	CSRFCookieKeyID      string
	UserAgentCookieKeyID string
}

type SystemAPIUsers map[string]*internal_authz.SystemAPIUser

func systemAPIUsersDecodeHook(from, to reflect.Value) (any, error) {
	if to.Type() != reflect.TypeOf(SystemAPIUsers{}) {
		return from.Interface(), nil
	}

	data, ok := from.Interface().(string)
	if !ok {
		return from.Interface(), nil
	}
	users := make(SystemAPIUsers)
	err := json.Unmarshal([]byte(data), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
