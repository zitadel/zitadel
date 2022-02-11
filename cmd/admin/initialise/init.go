package initialise

import (
	_ "embed"
	"fmt"

	"github.com/caos/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//sql import
	_ "github.com/lib/pq"
)

var (
	conn string
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize ZITADEL instance",
		Long: `init sets up the minimum requirements to start ZITADEL.
Prereqesits:
- cockroachdb`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := new(Config)
			if err := viper.Unmarshal(config); err != nil {
				return err
			}
			return initialise(config)
		},
	}

	// cmd.PersistentFlags().StringArrayVar(&configFiles, "config", nil, "path to config file to overwrite system defaults")
	//TODO(hust): simplify to multiple flags
	cmd.PersistentFlags().StringVar(&conn, "connection", "", "connection string to connect with a user which is allowed to create the database and user")

	return cmd
}

func initialise(config *Config) error {
	logging.Info("initialization started")

	if conn == "" {
		return fmt.Errorf("connection not defined")
	}

	if err := prepareDB(config.Database); err != nil {
		return err
	}

	return prepareZitadel(config.Database)
}
