package start

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	new_logging "github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
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
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			err = tls.ModeFromFlag(cmd)
			if err != nil {
				return err
			}

			masterKey, err := key.MasterKey(cmd)
			if err != nil {
				return err
			}

			err = setup.BindInitProjections(cmd)
			if err != nil {
				return err
			}

			setupConfig, shutdown, err := setup.NewConfig(cmd.Context(), viper.GetViper())
			if err != nil {
				return err
			}
			defer func() {
				err = errors.Join(err, shutdown(cmd.Context()))
			}()

			setupSteps := setup.MustNewSteps(viper.New())

			setupCtx, cancel := context.WithCancel(cmd.Context())
			setup.Setup(setupCtx, setupConfig, setupSteps, masterKey)
			cancel()

			startConfig, _, err := NewConfig(cmd.Context(), viper.GetViper())
			if err != nil {
				return err
			}

			return startZitadel(cmd.Context(), startConfig, masterKey, server)
		},
	}

	cmd.SetErr(new_logging.CommandErrorWriter("start-from-setup"))

	startFlags(cmd)
	setup.Flags(cmd)

	return cmd
}
