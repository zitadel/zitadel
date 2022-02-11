package initialise

import (
	_ "embed"

	"github.com/caos/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//sql import
	_ "github.com/lib/pq"
)

var (
	user     string
	password string
	sslCert  string
	sslKey   string
)

const (
	userFlag     = "user"
	passwordFlag = "password"
	sslCertFlag  = "ssl-cert"
	sslKeyFlag   = "ssl-key"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize ZITADEL instance",
		Long: `Sets up the minimum requirements to start ZITADEL.

Prereqesits:
- cockroachdb

The user provided by flags needs priviledge to 
- create the database if it does not exist
- see other users and create a new one if the user does not exist
- grant all rights of the ZITADEL database to the user created if not yet set
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := new(Config)
			if err := viper.Unmarshal(config); err != nil {
				return err
			}
			return initialise(config)
		},
	}

	cmd.PersistentFlags().StringVar(&password, passwordFlag, "", "password of the the provided user")
	cmd.PersistentFlags().StringVar(&sslCert, sslCertFlag, "", "ssl cert from the provided user")
	cmd.PersistentFlags().StringVar(&sslKey, sslKeyFlag, "", "ssl key from the provided user")
	cmd.PersistentFlags().StringVar(&user, userFlag, "", "(required) the user to check if the database, user and grants exists and create if not")
	cmd.MarkPersistentFlagRequired(userFlag)

	return cmd
}

func initialise(config *Config) error {
	logging.Info("initialization started")

	if err := prepareDB(config.Database, user, password, sslCert, sslKey); err != nil {
		return err
	}

	return prepareZitadel(config.Database)
}
