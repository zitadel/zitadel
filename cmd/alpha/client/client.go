package client

import (
	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	oidc "github.com/zitadel/zitadel/pkg/grpc/oidc/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/system"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "client calls zitadel",
	}

	cmd.AddCommand(
		admin.AdminServiceCmd,
		auth.AuthServiceCmd,
		management.ManagementServiceCmd,
		oidc.OIDCServiceCmd,
		org.OrganizationServiceCmd,
		session.SessionServiceCmd,
		settings.SettingsServiceCmd,
		system.SystemServiceCmd,
		user.UserServiceCmd,
	)

	return cmd
}
