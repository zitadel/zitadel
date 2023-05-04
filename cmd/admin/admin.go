package admin

import (
	_ "embed"
	"errors"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/cmd/initialise"
	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/cmd/setup"
	"github.com/zitadel/zitadel/cmd/start"
)

func New() *cobra.Command {
	adminCMD := &cobra.Command{
		Use:        "admin",
		Short:      "The ZITADEL admin CLI lets you interact with your instance",
		Long:       `The ZITADEL admin CLI lets you interact with your instance`,
		Deprecated: "please use subcommands directly, e.g. `zitadel start`",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("no additional command provided")
		},
	}

	adminCMD.AddCommand(
		initialise.New(),
		setup.New(),
		start.New(nil),
		start.NewStartFromInit(nil),
		key.New(),
	)

	return adminCMD
}
