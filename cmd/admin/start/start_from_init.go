package start

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/cmd/admin/initialise"
	"github.com/caos/zitadel/cmd/admin/setup"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewStartFromInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-from-init",
		Short: "cold starts zitadel",
		Long: `cold starts ZITADEL.
First the minimum requirements to start ZITADEL are set up.
Second the initial events are created.
Last ZITADEL starts.

Requirements:
- cockroachdb`,
		Run: func(cmd *cobra.Command, args []string) {
			initialise.InitAll(initialise.MustNewConfig(viper.GetViper()))

			setupConfig := setup.MustNewConfig(viper.GetViper())
			setupSteps := setup.MustNewSteps(viper.GetViper())
			setup.Setup(setupConfig, setupSteps)

			startConfig := MustNewConfig(viper.GetViper())
			startMasterKey, _ := cmd.Flags().GetString(flagMasterKey)

			err := startZitadel(startConfig, startMasterKey)
			logging.OnError(err).Fatal("unable to start zitadel")
		},
	}

	startFlags(cmd)

	return cmd
}
