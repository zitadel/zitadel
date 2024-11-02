package start

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/initialise"
	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/cmd/setup"
	"github.com/zitadel/zitadel/cmd/tls"
)

func NewStartFromInit(server chan<- *Server) *cobra.Command {
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
			err := tls.ModeFromFlag(cmd)
			logging.OnError(err).Fatal("invalid tlsMode")

			masterKey, err := key.MasterKey(cmd)
			logging.OnError(err).Panic("No master key provided")

			initialise.InitAll(cmd.Context(), initialise.MustNewConfig(viper.GetViper()))

			err = setup.BindInitProjections(cmd)
			logging.OnError(err).Fatal("unable to bind \"init-projections\" flag")

			setupConfig := setup.MustNewConfig(viper.GetViper())
			setupSteps := setup.MustNewSteps(viper.New())
			setup.Setup(cmd.Context(), setupConfig, setupSteps, masterKey)

			startConfig := MustNewConfig(viper.GetViper())

			err = startZitadel(cmd.Context(), startConfig, masterKey, server)
			logging.OnError(err).Fatal("unable to start zitadel")
		},
	}

	startFlags(cmd)
	setup.Flags(cmd)

	return cmd
}
