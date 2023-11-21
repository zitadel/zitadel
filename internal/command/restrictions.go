package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/restrictions"
)

type SetRestrictions struct {
	PublicOrgRegistrationIsNotAllowed *bool
	AllowedLanguages                  []language.Tag
}

// SetRestrictions creates new restrictions or updates existing restrictions.
func (c *Commands) SetInstanceRestrictions(
	ctx context.Context,
	setRestrictions *SetRestrictions,
	defaultLanguage language.Tag,
) (*domain.ObjectDetails, error) {
	instanceId := authz.GetInstance(ctx).InstanceID()
	wm, err := c.getRestrictionsWriteModel(ctx, instanceId, instanceId)
	if err != nil {
		return nil, err
	}
	aggregateId := wm.AggregateID
	if aggregateId == "" {
		aggregateId, err = c.idGenerator.Next()
		if err != nil {
			return nil, err
		}
	}
	setCmd, err := c.SetRestrictionsCommand(restrictions.NewAggregate(aggregateId, instanceId, instanceId), wm, setRestrictions, defaultLanguage)()
	if err != nil {
		return nil, err
	}
	cmds, err := setCmd(ctx, nil)
	if err != nil {
		return nil, err
	}
	if len(cmds) > 0 {
		events, err := c.eventstore.Push(ctx, cmds...)
		if err != nil {
			return nil, err
		}
		err = AppendAndReduce(wm, events...)
		if err != nil {
			return nil, err
		}
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) getRestrictionsWriteModel(ctx context.Context, instanceId, resourceOwner string) (*restrictionsWriteModel, error) {
	wm := newRestrictionsWriteModel(instanceId, resourceOwner)
	return wm, c.eventstore.FilterToQueryReducer(ctx, wm)
}

func (c *Commands) SetRestrictionsCommand(a *restrictions.Aggregate, wm *restrictionsWriteModel, setRestrictions *SetRestrictions, defaultLanguage language.Tag) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if setRestrictions == nil ||
			setRestrictions.PublicOrgRegistrationIsNotAllowed == nil &&
				setRestrictions.AllowedLanguages == nil {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-oASwj", "Errors.Restrictions.NoneSpecified")
		}
		if setRestrictions.AllowedLanguages != nil {
			if err := domain.LanguagesAreSupported(setRestrictions.AllowedLanguages...); err != nil {
				return nil, err
			}
			if err := domain.LanguageIsAllowed(false, setRestrictions.AllowedLanguages, defaultLanguage); err != nil {
				return nil, err
			}
		}
		return func(ctx context.Context, _ preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			changes, languagesChanged := wm.NewChanges(setRestrictions)
			if len(changes) == 0 {
				return nil, nil
			}
			commands := []eventstore.Command{restrictions.NewSetEvent(
				eventstore.NewBaseEventForPush(
					ctx,
					&a.Aggregate,
					restrictions.SetEventType,
				),
				changes...,
			)}
			if languagesChanged {
				profiles, err := c.allProfileWriteModels(ctx)
				if err != nil {
					return nil, err
				}
				for _, profile := range profiles {
					if notAllowedErr := domain.LanguageIsAllowed(true, setRestrictions.AllowedLanguages, profile.PreferredLanguage); notAllowedErr != nil {
						changeProfile, profileChanged, profileChangedErr := profile.NewChangedEvent(
							ctx,
							UserAggregateFromWriteModel(&profile.WriteModel),
							profile.FirstName,
							profile.LastName,
							profile.NickName,
							profile.DisplayName,
							defaultLanguage,
							profile.Gender,
						)
						if profileChangedErr != nil {
							return nil, profileChangedErr
						}
						if profileChanged {
							commands = append(commands, changeProfile)
						}
					}
				}
			}
			return commands, nil
		}, nil
	}
}
