package cmd

import (
	"bytes"
	_ "embed"
	"errors"
	"io"
	"strings"

	"github.com/caos/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/cmd/admin"
)

var (
	configFiles []string

	//go:embed defaults.yaml
	defaultConfig []byte
)

func New(out io.Writer, in io.Reader, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zitadel",
		Short: "The ZITADEL CLI let's you interact with ZITADEL",
		Long:  `The ZITADEL CLI let's you interact with ZITADEL`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("no additional command provided")
		},
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("ZITADEL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	logging.OnError(err).Fatal("unable to read default config")

	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().StringArrayVar(&configFiles, "config", nil, "path to config file to overwrite system defaults")

	cmd.AddCommand(admin.New())

	return cmd
}

func initConfig() {
	for _, file := range configFiles {
		viper.SetConfigFile(file)
		err := viper.MergeInConfig()
		logging.WithFields("file", file).OnError(err).Warn("unable to read config file")
	}
}
