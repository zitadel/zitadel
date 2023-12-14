package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func (c *Commands) changeUsername(ctx context.Context, cmds []eventstore.Command, wm *UserHumanWriteModel, userName string) ([]eventstore.Command, error) {
	if wm.UserName == userName {
		return cmds, nil
	}
	orgID := wm.ResourceOwner

	domainPolicy, err := c.domainPolicyWriteModel(ctx, orgID)
	if err != nil {
		return cmds, errors.ThrowPreconditionFailed(err, "COMMAND-38fnu", "Errors.Org.DomainPolicy.NotExisting")
	}
	if !domainPolicy.UserLoginMustBeDomain {
		index := strings.LastIndex(userName, "@")
		if index > 1 {
			domainCheck := NewOrgDomainVerifiedWriteModel(userName[index+1:])
			if err := c.eventstore.FilterToQueryReducer(ctx, domainCheck); err != nil {
				return cmds, err
			}
			if domainCheck.Verified && domainCheck.ResourceOwner != orgID {
				return cmds, errors.ThrowInvalidArgument(nil, "COMMAND-Di2ei", "Errors.User.DomainNotAllowedAsUsername")
			}
		}
	}
	return append(cmds,
		user.NewUsernameChangedEvent(ctx, UserAggregateFromWriteModel(&wm.WriteModel), wm.UserName, userName, domainPolicy.UserLoginMustBeDomain),
	), nil
}
