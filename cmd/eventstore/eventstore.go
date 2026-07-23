package eventstore

import (
	"errors"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	eventstoreCMD := &cobra.Command{
		Use:   "eventstore",
		Short: "eventstore maintenance commands",
		Long:  `eventstore maintenance commands`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("no additional command provided")
		},
	}

	eventstoreCMD.AddCommand(
		cleanupEvents2Cmd(),
	)

	return eventstoreCMD
}
