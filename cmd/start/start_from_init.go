package start

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
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

			showBasicInformation(startConfig)

			err = startZitadel(startConfig, masterKey, server)
			logging.OnError(err).Fatal("unable to start zitadel")
		},
	}

	startFlags(cmd)
	setup.Flags(cmd)

	return cmd
}

func showBasicInformation(startConfig *Config) {
	// Show basic information
	fmt.Println(color.MagentaString(figure.NewFigure("Zitadel", "", true).String()))
	http := "http"
	if startConfig.TLS.Enabled || startConfig.ExternalSecure {
		http = "https"
	}

	consoleURL := fmt.Sprintf("%s://%s:%v/ui/console\n", http, startConfig.ExternalDomain, startConfig.ExternalPort)
	healthCheckURL := fmt.Sprintf("%s://%s:%v/debug/healthz\n", http, startConfig.ExternalDomain, startConfig.ExternalPort)

	insecure := !startConfig.TLS.Enabled && !startConfig.ExternalSecure

	fmt.Printf(" ===============================================================\n\n")
	fmt.Printf(" Version          : %s\n", build.Version())
	fmt.Printf(" TLS enabled      : %v\n", startConfig.TLS.Enabled)
	fmt.Printf(" External Secure  : %v\n", startConfig.ExternalSecure)
	fmt.Printf(" Console URL      : %s", color.BlueString(consoleURL))
	fmt.Printf(" Health Check URL : %s", color.BlueString(healthCheckURL))
	if insecure {
		fmt.Printf("\n %s: you're using plain http without TLS. Be aware this is \n", color.RedString("Warning"))
		fmt.Printf(" not a secure setup and should only be used for test systems.         \n")
		fmt.Printf(" Visit: %s    \n", color.CyanString("https://zitadel.com/docs/self-hosting/manage/tls_modes"))
	}
	fmt.Printf("\n ===============================================================\n\n")
}
