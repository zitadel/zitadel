package client

import (
	"time"

	"github.com/adlerhurst/cli-client"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
	"github.com/zitadel/zitadel/internal/config/hook"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	oidc_svc "github.com/zitadel/zitadel/pkg/grpc/oidc/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/system"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "client",
		Short:            "client calls zitadel",
		PersistentPreRun: client,
	}

	cmd.AddCommand(
		admin.AdminServiceCmd,
		auth.AuthServiceCmd,
		management.ManagementServiceCmd,
		oidc_svc.OIDCServiceCmd,
		org.OrganizationServiceCmd,
		session.SessionServiceCmd,
		settings.SettingsServiceCmd,
		system.SystemServiceCmd,
		user.UserCmd,
	)

	return cmd
}

func client(cmd *cobra.Command, args []string) {
	config := mustReadConfig(viper.GetViper())

	client, err := zitadel.NewConnection(
		config.Issuer,
		config.Endpoint,
		[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
	)
	logging.OnError(err).Fatal("unable to create zitadel client")

	cli.SetConnection(client.ClientConn)
}

func mustReadConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := v.Unmarshal(config,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hook.Base64ToBytesHookFunc(),
			hook.TagToLanguageHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
			hook.EnumHookFunc(domain.FeatureString),
		)),
	)
	logging.OnError(err).Fatal("unable to read default config")

	return config
}

type Config struct {
	Issuer   string
	Endpoint string
}
