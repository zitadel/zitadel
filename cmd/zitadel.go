package cmd

import (
	_ "embed"
	"errors"
	"io"

	"github.com/zitadel/zitadel/internal/config/options"

	"github.com/spf13/cobra"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/admin"
	"github.com/zitadel/zitadel/cmd/initialise"
	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/cmd/setup"
	"github.com/zitadel/zitadel/cmd/start"
)

var (
	configFiles []string
)

func New(out io.Writer, in io.Reader, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zitadel",
		Short: "The ZITADEL CLI lets you interact with ZITADEL",
		Long:  `The ZITADEL CLI lets you interact with ZITADEL`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("no additional command provided")
		},
	}

	err := options.InitViper()
	logging.OnError(err).Fatalf("unable initialize config: %s", err)

	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().StringArrayVar(&configFiles, "config", nil, "path to config file to overwrite system defaults")

	cmd.AddCommand(
		admin.New(), //is now deprecated, remove later on
		initialise.New(),
		setup.New(),
		start.New(),
		start.NewStartFromInit(),
		key.New(),
	)

	return cmd
}

func initConfig() {
	options.MergeToViper(configFiles...)
}
