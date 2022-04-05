package admin

import (
	_ "embed"

	"github.com/caos/logging"
	"github.com/spf13/cobra"

	"github.com/caos/zitadel/cmd/admin/initialise"
	"github.com/caos/zitadel/cmd/admin/key"
	"github.com/caos/zitadel/cmd/admin/setup"
	"github.com/caos/zitadel/cmd/admin/start"
)

func New() *cobra.Command {
	adminCMD := &cobra.Command{
		Use:   "admin",
		Short: "The ZITADEL admin CLI let's you interact with your instance",
		Long:  `The ZITADEL admin CLI let's you interact with your instance`,
		Run: func(cmd *cobra.Command, args []string) {
			logging.New().Info("hello world")
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
