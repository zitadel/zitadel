package start

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/cmd/setup"
	"github.com/zitadel/zitadel/cmd/tls"
)

func NewStartFromSetup(server chan<- *Server) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-from-setup",
		Short: "cold starts zitadel",
		Long: `cold starts ZITADEL.
First the initial events are created.
Last ZITADEL starts.

Requirements:
- database
- database is initialized
`,
		Run: func(cmd *cobra.Command, args []string) {
			err := tls.ModeFromFlag(cmd)
			logging.OnError(err).Fatal("invalid tlsMode")

			masterKey, err := key.MasterKey(cmd)
			logging.OnError(err).Panic("No master key provided")

			setupConfig := setup.MustNewConfig(viper.GetViper())
			setupSteps := setup.MustNewSteps(viper.New())
			setup.Setup(setupConfig, setupSteps, masterKey)

			startConfig := MustNewConfig(viper.GetViper())

			err = startZitadel(cmd.Context(), startConfig, masterKey, server)
			logging.OnError(err).Fatal("unable to start zitadel")
		},
	}

	startFlags(cmd)
	setup.Flags(cmd)

	return cmd
}
