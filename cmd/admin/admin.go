package admin

import (
	_ "embed"
	"errors"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/cmd/admin/initialise"
	"github.com/zitadel/zitadel/cmd/admin/key"
	"github.com/zitadel/zitadel/cmd/admin/setup"
	"github.com/zitadel/zitadel/cmd/admin/start"
)

func New() *cobra.Command {
	adminCMD := &cobra.Command{
		Use:   "admin",
		Short: "The ZITADEL admin CLI lets you interact with your instance",
		Long:  `The ZITADEL admin CLI lets you interact with your instance`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("no additional command provided")
		},
	}

	adminCMD.AddCommand(
		initialise.New(),
		setup.New(),
		start.New(),
		start.NewStartFromInit(),
		key.New(),
	)

	return adminCMD
}
