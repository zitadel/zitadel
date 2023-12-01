package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func (c *Commands) changeUsername(ctx context.Context, cmds []eventstore.Command, wm *UserHumanWriteModel, userName string) error {
	if wm.UserName == userName {
		return nil
	}
	orgID := wm.ResourceOwner

	domainPolicy, err := c.getOrgDomainPolicy(ctx, orgID)
	if err != nil {
		return errors.ThrowPreconditionFailed(err, "COMMAND-38fnu", "Errors.Org.DomainPolicy.NotExisting")
	}
	if !domainPolicy.UserLoginMustBeDomain {
		index := strings.LastIndex(userName, "@")
		if index > 1 {
			domainCheck := NewOrgDomainVerifiedWriteModel(userName[index+1:])
			if err := c.eventstore.FilterToQueryReducer(ctx, domainCheck); err != nil {
				return err
			}
			if domainCheck.Verified && domainCheck.ResourceOwner != orgID {
				return errors.ThrowInvalidArgument(nil, "COMMAND-Di2ei", "Errors.User.DomainNotAllowedAsUsername")
			}
		}
	}
	cmds = append(cmds,
		user.NewUsernameChangedEvent(ctx, UserAggregateFromWriteModel(&wm.WriteModel), wm.UserName, userName, domainPolicy.UserLoginMustBeDomain),
	)
	return nil
}
