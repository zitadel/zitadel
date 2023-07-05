package start

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/build"
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

			initialise.InitAll(initialise.MustNewConfig(viper.GetViper()))

			setupConfig := setup.MustNewConfig(viper.GetViper())
			setupSteps := setup.MustNewSteps(viper.New())
			setup.Setup(setupConfig, setupSteps, masterKey)

			startConfig := MustNewConfig(viper.GetViper())

			// Show basic information
			figure.NewFigure("Zitadel", "", true).Print()
			http := "http"
			if startConfig.TLS.Enabled || startConfig.ExternalSecure {
				http = "https"
			}
			fmt.Printf("\n ===============================================================\n\n")
			fmt.Printf(" Version          : %s\n", build.Version())
			fmt.Printf(" TLS enabled      : %v\n", startConfig.TLS.Enabled)
			fmt.Printf(" External Secure  : %v\n", startConfig.ExternalSecure)
			fmt.Printf(" Console URL      : %s://%s:%v/ui/console\n", http, startConfig.ExternalDomain, startConfig.ExternalPort)
			fmt.Printf(" Health Check URL : %s://%s:%v/debug/healthz\n", http, startConfig.ExternalDomain, startConfig.ExternalPort)
			fmt.Printf("\n ===============================================================\n\n")

			err = startZitadel(startConfig, masterKey, server)
			logging.OnError(err).Fatal("unable to start zitadel")
		},
	}

	startFlags(cmd)
	setup.Flags(cmd)

	return cmd
}
