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
	configFiles []string

	//go:embed defaults.yaml
	defaultConfig []byte
)

func NewZitadelCMD(out io.Writer, in io.Reader, args []string) *cobra.Command {
	rootCMD := &cobra.Command{
		Use:   "zitadel",
		Short: "The ZITADEL CLI let's you interact with ZITADEL",
		Long:  `The ZITADEL CLI let's you interact with ZITADEL`,
		Run: func(cmd *cobra.Command, args []string) {
			logging.Log("CMD-t7pjR").Info("hello world")
		},
	}

	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	logging.Log("CMD-5NgGF").OnError(err).Fatal("unable to read default config")

	cobra.OnInitialize(initConfig)
	rootCMD.PersistentFlags().StringArrayVar(&configFiles, "config", nil, "path to config file to overwrite system defaults")

	rootCMD.AddCommand(admin.New())

	return rootCMD
}

func initConfig() {
	for _, file := range configFiles {
		viper.SetConfigFile(file)
		err := viper.MergeInConfig()
		logging.LogWithFields("CMD-N76SC", "file", file).OnError(err).Warn("unable to read config file")
	}
}
