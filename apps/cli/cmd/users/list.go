package users

import (
	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/internal/auth"
	"github.com/zitadel/zitadel/apps/cli/internal/client"
	"github.com/zitadel/zitadel/apps/cli/internal/config"
	"github.com/zitadel/zitadel/apps/cli/internal/output"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2/userconnect"
)

// NewCmd creates the "users" command group.
func NewCmd(getCfg func() *config.Config, getOutput func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage users",
	}

	cmd.AddCommand(newListCmd(getCfg, getOutput))
	return cmd
}

func newListCmd(getCfg func() *config.Config, getOutput func() string) *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users in the current instance",
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
			userClient := userconnect.NewUserServiceClient(httpClient, baseURL)

			req := &user.ListUsersRequest{}
			if limit > 0 {
				req.Query = &object.ListQuery{
					Limit: uint32(limit),
				}
			}

			resp, err := userClient.ListUsers(cmd.Context(), connect.NewRequest(req))
			if err != nil {
				return err
			}

			if getOutput() == "json" {
				return output.JSON(resp.Msg)
			}

			header := []string{"ID", "USERNAME", "STATE"}
			var rows [][]string
			for _, u := range resp.Msg.GetResult() {
				rows = append(rows, []string{
					u.GetUserId(),
					u.GetUsername(),
					u.GetState().String(),
				})
			}
			output.Table(header, rows)
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 0, "maximum number of users to return")
	return cmd
}
