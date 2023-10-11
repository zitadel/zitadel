package command

import (
	"context"
	errs "errors"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func (c *Commands) prepareAddOrgDomain(a *org.Aggregate, addDomain string, userIDs []string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if addDomain = strings.TrimSpace(addDomain); addDomain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "ORG-r3h4J", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			existing, err := orgDomain(ctx, filter, a.ID, addDomain)
			if err != nil && !errs.Is(err, errors.ThrowNotFound(nil, "", "")) {
				return nil, err
			}
			if existing != nil && existing.State == domain.OrgDomainStateActive {
				return nil, errors.ThrowAlreadyExists(nil, "V2-e1wse", "Errors.Already.Exists")
			}
			domainPolicy, err := domainPolicyWriteModel(ctx, filter, a.ID)
			if err != nil {
				return nil, err
			}
			events := []eventstore.Command{org.NewDomainAddedEvent(ctx, &a.Aggregate, addDomain)}
			if !domainPolicy.ValidateOrgDomains {
				events = append(events, org.NewDomainVerifiedEvent(ctx, &a.Aggregate, addDomain))
				for _, userID := range userIDs {
					claimedEvent, err := c.prepareUserDomainClaimed(ctx, filter, userID)
					if err != nil {
						logging.WithFields("userid", userID).WithError(err).Error("could not claim user")
						continue
					}
					events = append(events, claimedEvent)
				}
			}
			return events, nil
		}, nil
	}
}

func verifyOrgDomain(a *org.Aggregate, domain string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if domain = strings.TrimSpace(domain); domain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "ORG-yqlVQ", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			// no checks required because unique constraints handle it
			return []eventstore.Command{org.NewDomainVerifiedEvent(ctx, &a.Aggregate, domain)}, nil
		}, nil
	}
}

func setPrimaryOrgDomain(a *org.Aggregate, domain string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if domain = strings.TrimSpace(domain); domain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "ORG-gmNqY", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			existing, err := orgDomain(ctx, filter, a.ID, domain)
			if err != nil {
				return nil, errors.ThrowAlreadyExists(err, "V2-d0Gyw", "Errors.Already.Exists")
			}
			if existing.Primary {
				return nil, errors.ThrowPreconditionFailed(nil, "COMMA-FfoZO", "Errors.Org.DomainAlreadyPrimary")
			}
			if !existing.Verified {
				return nil, errors.ThrowPreconditionFailed(nil, "COMMA-yKA80", "Errors.Org.DomainNotVerified")
			}
			return []eventstore.Command{org.NewDomainPrimarySetEvent(ctx, &a.Aggregate, domain)}, nil
		}, nil
	}
}

func orgDomain(ctx context.Context, filter preparation.FilterToQueryReducer, orgID, domain string) (*OrgDomainWriteModel, error) {
	wm := NewOrgDomainWriteModel(orgID, domain)
	events, err := filter(ctx, wm.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "COMMA-kFHpQ", "Errors.Org.DomainNotFound")
	}
	wm.AppendEvents(events...)
	if err = wm.Reduce(); err != nil {
		return nil, err
	}

	return wm, nil
}

func (c *Commands) VerifyOrgDomain(ctx context.Context, orgID, domain string) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(orgID)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, verifyOrgDomain(orgAgg, domain))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgDomain(ctx context.Context, orgID, domain string, claimedUserIDs []string) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(orgID)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgDomain(orgAgg, domain, claimedUserIDs))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) GenerateOrgDomainValidation(ctx context.Context, orgDomain *domain.OrgDomain) (token, url string, err error) {
	if orgDomain == nil || !orgDomain.IsValid() || orgDomain.AggregateID == "" {
		return "", "", errors.ThrowInvalidArgument(nil, "ORG-R24hb", "Errors.Org.InvalidDomain")
	}
	checkType, ok := orgDomain.ValidationType.CheckType()
	if !ok {
		return "", "", errors.ThrowInvalidArgument(nil, "ORG-Gsw31", "Errors.Org.DomainVerificationTypeInvalid")
	}
	domainWriteModel, err := c.getOrgDomainWriteModel(ctx, orgDomain.AggregateID, orgDomain.Domain)
	if err != nil {
		return "", "", err
	}
	if domainWriteModel.State != domain.OrgDomainStateActive {
		return "", "", errors.ThrowNotFound(nil, "ORG-AGD31", "Errors.Org.DomainNotOnOrg")
	}
	if domainWriteModel.Verified {
		return "", "", errors.ThrowPreconditionFailed(nil, "ORG-HGw21", "Errors.Org.DomainAlreadyVerified")
	}
	token, err = orgDomain.GenerateVerificationCode(c.domainVerificationGenerator)
	if err != nil {
		return "", "", err
	}
	url, err = http_utils.TokenUrl(orgDomain.Domain, token, checkType)
	if err != nil {
		return "", "", errors.ThrowPreconditionFailed(err, "ORG-Bae21", "Errors.Org.DomainVerificationTypeInvalid")
	}

	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)

	_, err = c.eventstore.Push(
		ctx,
		org.NewDomainVerificationAddedEvent(ctx, orgAgg, orgDomain.Domain, orgDomain.ValidationType, orgDomain.ValidationCode))
	if err != nil {
		return "", "", err
	}
	return token, url, nil
}

func (c *Commands) ValidateOrgDomain(ctx context.Context, orgDomain *domain.OrgDomain, claimedUserIDs []string) (*domain.ObjectDetails, error) {
	if orgDomain == nil || !orgDomain.IsValid() || orgDomain.AggregateID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "ORG-R24hb", "Errors.Org.InvalidDomain")
	}
	domainWriteModel, err := c.getOrgDomainWriteModel(ctx, orgDomain.AggregateID, orgDomain.Domain)
	if err != nil {
		return nil, err
	}
	if domainWriteModel.State != domain.OrgDomainStateActive {
		return nil, errors.ThrowNotFound(nil, "ORG-Sjdi3", "Errors.Org.DomainNotOnOrg")
	}
	if domainWriteModel.Verified {
		return nil, errors.ThrowPreconditionFailed(nil, "ORG-HGw21", "Errors.Org.DomainAlreadyVerified")
	}
	if domainWriteModel.ValidationCode == nil || domainWriteModel.ValidationType == domain.OrgDomainValidationTypeUnspecified {
		return nil, errors.ThrowPreconditionFailed(nil, "ORG-SFBB3", "Errors.Org.DomainVerificationMissing")
	}

	validationCode, err := crypto.DecryptString(domainWriteModel.ValidationCode, c.domainVerificationAlg)
	if err != nil {
		return nil, err
	}
	checkType, _ := domainWriteModel.ValidationType.CheckType()
	err = c.domainVerificationValidator(domainWriteModel.Domain, validationCode, validationCode, checkType)
	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	var events []eventstore.Command
	if err == nil {
		events = append(events, org.NewDomainVerifiedEvent(ctx, orgAgg, orgDomain.Domain))

		for _, userID := range claimedUserIDs {
			userEvents, _, err := c.userDomainClaimed(ctx, userID)
			if err != nil {
				logging.WithFields("userid", userID).WithError(err).Warn("could not claim user")
				continue
			}
			events = append(events, userEvents...)
		}
		pushedEvents, err := c.eventstore.Push(ctx, events...)
		if err != nil {
			return nil, err
		}
		err = AppendAndReduce(domainWriteModel, pushedEvents...)
		if err != nil {
			return nil, err
		}
		return writeModelToObjectDetails(&domainWriteModel.WriteModel), nil
	}
	events = append(events, org.NewDomainVerificationFailedEvent(ctx, orgAgg, orgDomain.Domain))

	_, errPush := c.eventstore.Push(ctx, events...)
	logging.LogWithFields("ORG-dhTE", "orgID", orgAgg.ID, "domain", orgDomain.Domain).OnError(errPush).Error("NewDomainVerificationFailedEvent push failed")

	return nil, err
}

func (c *Commands) SetPrimaryOrgDomain(ctx context.Context, orgDomain *domain.OrgDomain) (*domain.ObjectDetails, error) {
	if orgDomain == nil || !orgDomain.IsValid() || orgDomain.AggregateID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "ORG-SsDG2", "Errors.Org.InvalidDomain")
	}
	domainWriteModel, err := c.getOrgDomainWriteModel(ctx, orgDomain.AggregateID, orgDomain.Domain)
	if err != nil {
		return nil, err
	}
	if domainWriteModel.State != domain.OrgDomainStateActive {
		return nil, errors.ThrowNotFound(nil, "ORG-GDfA3", "Errors.Org.DomainNotOnOrg")
	}
	if !domainWriteModel.Verified {
		return nil, errors.ThrowPreconditionFailed(nil, "ORG-Ggd32", "Errors.Org.DomainNotVerified")
	}
	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewDomainPrimarySetEvent(ctx, orgAgg, orgDomain.Domain))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(domainWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&domainWriteModel.WriteModel), nil
}

func (c *Commands) RemoveOrgDomain(ctx context.Context, orgDomain *domain.OrgDomain) (*domain.ObjectDetails, error) {
	if orgDomain == nil || !orgDomain.IsValid() || orgDomain.AggregateID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "ORG-SJsK3", "Errors.Org.InvalidDomain")
	}
	domainWriteModel, err := c.getOrgDomainWriteModel(ctx, orgDomain.AggregateID, orgDomain.Domain)
	if err != nil {
		return nil, err
	}
	if domainWriteModel.State != domain.OrgDomainStateActive {
		return nil, errors.ThrowNotFound(nil, "ORG-GDfA3", "Errors.Org.DomainNotOnOrg")
	}
	if domainWriteModel.Primary {
		return nil, errors.ThrowPreconditionFailed(nil, "ORG-Sjdi3", "Errors.Org.PrimaryDomainNotDeletable")
	}
	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewDomainRemovedEvent(ctx, orgAgg, orgDomain.Domain, domainWriteModel.Verified))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(domainWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&domainWriteModel.WriteModel), nil
}

func (c *Commands) addOrgDomain(ctx context.Context, orgAgg *eventstore.Aggregate, addedDomain *OrgDomainWriteModel, orgDomain *domain.OrgDomain, claimedUserIDs []string) ([]eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedDomain)
	if err != nil {
		return nil, err
	}
	if addedDomain.State == domain.OrgDomainStateActive {
		return nil, errors.ThrowAlreadyExists(nil, "COMMA-Bd2jj", "Errors.Org.Domain.AlreadyExists")
	}

	events := []eventstore.Command{
		org.NewDomainAddedEvent(ctx, orgAgg, orgDomain.Domain),
	}

	if orgDomain.Verified {
		events = append(events, org.NewDomainVerifiedEvent(ctx, orgAgg, orgDomain.Domain))
		for _, userID := range claimedUserIDs {
			userEvents, _, err := c.userDomainClaimed(ctx, userID)
			if err != nil {
				logging.WithFields("userid", userID).WithError(err).Warn("could not claim user")
				continue
			}
			events = append(events, userEvents...)
		}
	}
	if orgDomain.Primary {
		events = append(events, org.NewDomainPrimarySetEvent(ctx, orgAgg, orgDomain.Domain))
	}
	return events, nil
}

func (c *Commands) changeDefaultDomain(ctx context.Context, orgID, newName string) ([]eventstore.Command, error) {
	orgDomains := NewOrgDomainsWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, orgDomains)
	if err != nil {
		return nil, err
	}
	iamDomain := authz.GetInstance(ctx).RequestedDomain()
	defaultDomain, _ := domain.NewIAMDomainName(orgDomains.OrgName, iamDomain)
	isPrimary := defaultDomain == orgDomains.PrimaryDomain
	orgAgg := OrgAggregateFromWriteModel(&orgDomains.WriteModel)
	for _, orgDomain := range orgDomains.Domains {
		if orgDomain.State == domain.OrgDomainStateActive {
			if orgDomain.Domain == defaultDomain {
				newDefaultDomain, err := domain.NewIAMDomainName(newName, iamDomain)
				if err != nil {
					return nil, err
				}
				events := []eventstore.Command{
					org.NewDomainAddedEvent(ctx, orgAgg, newDefaultDomain),
					org.NewDomainVerifiedEvent(ctx, orgAgg, newDefaultDomain),
				}
				if isPrimary {
					events = append(events, org.NewDomainPrimarySetEvent(ctx, orgAgg, newDefaultDomain))
				}
				events = append(events, org.NewDomainRemovedEvent(ctx, orgAgg, orgDomain.Domain, orgDomain.Verified))
				return events, nil
			}
		}
	}
	return nil, nil
}

func (c *Commands) removeCustomDomains(ctx context.Context, orgID string) ([]eventstore.Command, error) {
	orgDomains := NewOrgDomainsWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, orgDomains)
	if err != nil {
		return nil, err
	}
	hasDefault := false
	defaultDomain, _ := domain.NewIAMDomainName(orgDomains.OrgName, authz.GetInstance(ctx).RequestedDomain())
	isPrimary := defaultDomain == orgDomains.PrimaryDomain
	orgAgg := OrgAggregateFromWriteModel(&orgDomains.WriteModel)
	events := make([]eventstore.Command, 0, len(orgDomains.Domains))
	for _, orgDomain := range orgDomains.Domains {
		if orgDomain.State == domain.OrgDomainStateActive {
			if orgDomain.Domain == defaultDomain {
				hasDefault = true
				continue
			}
			events = append(events, org.NewDomainRemovedEvent(ctx, orgAgg, orgDomain.Domain, orgDomain.Verified))
		}
	}
	if !hasDefault {
		return append([]eventstore.Command{
			org.NewDomainAddedEvent(ctx, orgAgg, defaultDomain),
			org.NewDomainPrimarySetEvent(ctx, orgAgg, defaultDomain),
		}, events...), nil
	}
	if !isPrimary {
		return append([]eventstore.Command{org.NewDomainPrimarySetEvent(ctx, orgAgg, defaultDomain)}, events...), nil
	}
	return events, nil
}

func (c *Commands) getOrgDomainWriteModel(ctx context.Context, orgID, domain string) (*OrgDomainWriteModel, error) {
	domainWriteModel := NewOrgDomainWriteModel(orgID, domain)
	err := c.eventstore.FilterToQueryReducer(ctx, domainWriteModel)
	if err != nil {
		return nil, err
	}
	return domainWriteModel, nil
}
