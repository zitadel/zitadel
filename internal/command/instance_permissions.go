package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/permission"
)

func prepareAddRolePermissions(a *instance.Aggregate, roles []authz.RoleMapping) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, _ preparation.FilterToQueryReducer) (cmds []eventstore.Command, _ error) {
			aggregate := permission.NewAggregate(a.InstanceID)
			for _, r := range roles {
				if strings.HasPrefix(r.Role, "SYSTEM") {
					continue
				}
				for _, p := range r.Permissions {
					cmds = append(cmds, permission.NewAddedEvent(ctx, aggregate, r.Role, p))
				}
			}
			return cmds, nil
		}, nil
	}
}
