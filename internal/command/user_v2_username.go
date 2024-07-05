package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) changeUsername(ctx context.Context, cmds []eventstore.Command, wm *UserV2WriteModel, userName string) ([]eventstore.Command, error) {
	if wm.UserName == userName {
		return cmds, nil
	}
	orgID := wm.ResourceOwner

	domainPolicy, err := c.domainPolicyWriteModel(ctx, orgID)
	if err != nil {
		return cmds, zerrors.ThrowPreconditionFailed(err, "COMMAND-79pv6e1q62", "Errors.Org.DomainPolicy.NotExisting")
	}
	if err = c.userValidateDomain(ctx, orgID, userName, domainPolicy.UserLoginMustBeDomain); err != nil {
		return cmds, err
	}
	return append(cmds,
		user.NewUsernameChangedEvent(ctx, &wm.Aggregate().Aggregate, wm.UserName, userName, domainPolicy.UserLoginMustBeDomain),
	), nil
}
