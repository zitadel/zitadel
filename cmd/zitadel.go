package cmd

import (
	"bytes"
	_ "embed"
	"io"

	"github.com/caos/logging"
	"github.com/caos/zitadel/cmd/admin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Run: func(cmd *cobra.Command, args []string) {
			logging.New().Info("hello world")
		},
	}

	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	logging.New().OnError(err).Fatal("unable to read default config")

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
