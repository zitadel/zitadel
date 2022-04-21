package start

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/hook"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	admin_es "github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	internal_authz "github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/api/oidc"
	"github.com/caos/zitadel/internal/api/ui/console"
	"github.com/caos/zitadel/internal/api/ui/login"
	auth_es "github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/authz"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/database"
	"github.com/caos/zitadel/internal/notification"
	"github.com/caos/zitadel/internal/query/projection"
	static_config "github.com/caos/zitadel/internal/static/config"
)

type Config struct {
	Log             *logging.Config
	Port            uint16
	ExternalPort    uint16
	ExternalDomain  string
	ExternalSecure  bool
	HTTP2HostHeader string
	HTTP1HostHeader string
	Database        database.Config
	Projections     projection.Config
	AuthZ           authz.Config
	Auth            auth_es.Config
	Admin           admin_es.Config
	UserAgentCookie *middleware.UserAgentCookieConfig
	OIDC            oidc.Config
	Login           login.Config
	Console         console.Config
	Notification    notification.Config
	AssetStorage    static_config.AssetStorageConfig
	InternalAuthZ   internal_authz.Config
	SystemDefaults  systemdefaults.SystemDefaults
	EncryptionKeys  *encryptionKeyConfig
	DefaultInstance command.InstanceSetup
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)

	err := v.Unmarshal(config,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hook.Base64ToBytesHookFunc(),
			hook.TagToLanguageHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		)),
	)
	err = config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	return config
}

type encryptionKeyConfig struct {
	DomainVerification   *crypto.KeyConfig
	IDPConfig            *crypto.KeyConfig
	OIDC                 *crypto.KeyConfig
	OTP                  *crypto.KeyConfig
	SMS                  *crypto.KeyConfig
	SMTP                 *crypto.KeyConfig
	User                 *crypto.KeyConfig
	CSRFCookieKeyID      string
	UserAgentCookieKeyID string
}
