package start

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	new_logging "github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
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
- postgreSQL`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			err = tls.ModeFromFlag(cmd)
			if err != nil {
				return fmt.Errorf("invalid tlsMode: %w", err)
			}

			masterKey, err := key.MasterKey(cmd)
			if err != nil {
				return fmt.Errorf("no master key provided: %w", err)
			}

			initCtx, cancel := context.WithCancel(cmd.Context())
			initialise.InitAll(initCtx, initialise.MustNewConfig(viper.GetViper()))
			cancel()

			err = setup.BindInitProjections(cmd)
			if err != nil {
				return fmt.Errorf("unable to bind \"init-projections\" flag: %w", err)
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

	cmd.SetErr(new_logging.CommandErrorWriter("start-from-init"))

	startFlags(cmd)
	setup.Flags(cmd)

	return cmd
}
