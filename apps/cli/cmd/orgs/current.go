package orgs

import (
	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/internal/auth"
	"github.com/zitadel/zitadel/apps/cli/internal/client"
	"github.com/zitadel/zitadel/apps/cli/internal/config"
	"github.com/zitadel/zitadel/apps/cli/internal/output"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/management/managementconnect"
)

// NewCmd creates the "orgs" command group.
func NewCmd(getCfg func() *config.Config, getOutput func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orgs",
		Short: "Manage organizations",
	}

	cmd.AddCommand(newCurrentCmd(getCfg, getOutput))
	return cmd
}

func newCurrentCmd(getCfg func() *config.Config, getOutput func() string) *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the current organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := getCfg()
			actx, _, err := config.ActiveCtx(cfg)
			if err != nil {
				return err
			}

			tokenSource, err := auth.TokenSource(cmd.Context(), actx)
			if err != nil {
				return err
			}

			httpClient := client.New(tokenSource)
			baseURL := client.InstanceURL(actx.Instance)
			mgmtClient := managementconnect.NewManagementServiceClient(httpClient, baseURL)

			resp, err := mgmtClient.GetMyOrg(cmd.Context(), connect.NewRequest(&management.GetMyOrgRequest{}))
			if err != nil {
				return err
			}

			if getOutput() == "json" {
				return output.JSON(resp.Msg)
			}

			org := resp.Msg.GetOrg()
			header := []string{"ID", "NAME", "STATE", "DOMAIN"}
			rows := [][]string{
				{
					org.GetId(),
					org.GetName(),
					org.GetState().String(),
					org.GetPrimaryDomain(),
				},
			}
			output.Table(header, rows)
			return nil
		},
	}
}
