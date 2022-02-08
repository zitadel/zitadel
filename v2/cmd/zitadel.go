package cmd

import (
	"bytes"
	_ "embed"
	"io"

	"github.com/caos/logging"
	"github.com/caos/zitadel/v2/cmd/admin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configPath string

	//go:embed defaults.yaml
	defaultConfig []byte
)

func NewZitadelCMD(out io.Writer, in io.ReadWriter, args []string) *cobra.Command {
	rootCMD := &cobra.Command{
		Use:   "zitadel",
		Short: "The ZITADEL CLI let's you interact with ZITADEL",
		Long:  `The ZITADEL CLI let's you interact with ZITADEL`,
		Run: func(cmd *cobra.Command, args []string) {
			logging.Log("ADMIN-t7pjR").Info("hello world")
		},
	}

	cobra.OnInitialize(initConfig)
	rootCMD.PersistentFlags().StringVar(&configPath, "config", "", "path to config file to overwrite system defaults")

	rootCMD.AddCommand(admin.New())

	return rootCMD
}

func initConfig() {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	logging.Log("ADMIN-5NgGF").OnError(err).Fatal("unable to read default config")

	if configPath != "" {
		viper.SetConfigFile(configPath)
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logging.Log("ADMIN-SX5sF").Info("using default config")
	}
}
