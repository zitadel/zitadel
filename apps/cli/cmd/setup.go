package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newSetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Instructions for setting up ZITADEL to work with this CLI",
		Long: `setup — Learn how to configure your ZITADEL instance for CLI access.

To use 'zitadel-cli login', you must first create a 'Native' application in your
ZITADEL project. This allows the CLI to authenticate on your behalf using the
OAuth 2.0 Device Authorization Grant flow.

Follow these steps in the ZITADEL Console:

1.  Navigate to your Project (e.g., 'Management').
2.  Click 'New Application'.
3.  Name it (e.g., 'ZITADEL CLI').
4.  Select 'Native' as the application type.
5.  On the 'Auth Method' screen, select 'PKCE'.
6.  Once the application is created, go to 'Settings' and ensure that the
    'Device Code' grant type is enabled.
7.  Copy the 'Client ID' from the application's page.

Note: No Redirect URIs are required for the Device Authorization flow.

Once configured, run:
  zitadel-cli login --instance <your-instance> --client-id <your-client-id>
`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(os.Stdout, cmd.Long)
		},
	}
}
